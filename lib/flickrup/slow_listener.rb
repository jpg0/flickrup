require 'rubygems'
require 'atomic'
require 'rufus/scheduler'
require "flickrup/logging"

#modification ->
#  schedule upload for 5 mins time
#
#upload run ->
#  for files modified > 5 mins ago
#    -> do upload
#  if there are files modified < 5 mins ago
#    -> schedule upload for minimum time
#
#When scheduling upload, remove all but the nearest
#When paused, upload run should be a no-op

class SlowListener
  include Logging

  IDLE = -1
  PROCESSING = 0
  CHANGES_QUEUED = 1

  def initialize(config)
    @wait = config[:timeout]
    logger.debug("Timeout set to #@wait secs")
    @config = config
    @scheduler = Rufus::Scheduler.start_new
    @processor = Processor.new(config)
    @processing = Atomic.new IDLE
  end

  def run
    logger.debug("Watching #{@config['watch_dir']}...")
    Listen.to(@config['watch_dir']) do |modified, added, removed|
      unless modified.empty? && added.empty?
        on_change
      end
    end
  end

  def do_upload
    logger.debug("Checking files for upload")
    files = Dir["#{@config['watch_dir']}/*.{jpg,JPG,jpeg,JPEG}"]

    old, new = files.partition { |x| Time.now > File.mtime(x) + @wait }

    reschedule = false

    unless old.empty?
      logger.debug("Found #{old.length} stable files")
      @processing.set PROCESSING
      @processor.run(old)
      changes = @processing.swap IDLE
      if changes > 0
        logger.debug("Changes occurred while processing")
        reschedule = true
      end
    end

    unless new.empty?
      logger.debug("Found #{new.length} new files")
      reschedule = true
    end

    if reschedule
      logger.debug("Rescheduling...")
      #todo: schedule for latest
      schedule
    end
  end

  def on_change
    logger.debug("Change detected")

    current_state = @processing.get
    was_processing = @processing.compare_and_swap(PROCESSING, CHANGES_QUEUED)

    if was_processing
      case current_state
        when IDLE
          on_change #retry
        when CHANGES_QUEUED
          on_change #retry
        when PROCESSING
          logger.debug("Scheduling queued") # already queued changes
      end
    else #wasn't processing
      case current_state
        when IDLE
          schedule
        when CHANGES_QUEUED
          logger.debug("Scheduling already queued") # already queued changes
        when PROCESSING
          on_change #retry
      end
    end
  end

  def schedule
    if @scheduler.jobs.empty?
      logger.debug("Scheduling processing for #@wait sec")
      @scheduler.in @wait do
        if @processing.get == IDLE
          do_upload
        end
      end
    else
      logger.debug("Already scheduled")
    end
  end

end
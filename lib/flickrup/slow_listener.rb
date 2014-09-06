require 'rubygems'
require 'atomic'
require 'fileutils'
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

  IDLE = 0
  SCHEDULED = 1
  PROCESSING = 2

  def initialize(config)
    @wait = config[:timeout]
    logger.info("Timeout set to #@wait secs")
    @config = config
    @scheduler = Rufus::Scheduler.start_new
    @processor = Processor.new(config)
    @processing = Atomic.new IDLE
  end

  def run
    logger.debug("Watching #{@config[:watch_dir]}...")
    @listener = Listen.to(@config[:watch_dir])
    @listener.change(&method(:on_change))
    @listener.start

    logger.info("Triggering initial run")
    schedule
  end

  def schedule
    previous_state = @processing.swap(SCHEDULED)

    case previous_state
      when SCHEDULED
        #ignore, already scheduled
        logger.info("Already scheduled")
      when PROCESSING
        #ignore, processing thread will reschedule
        logger.info("Scheduling queued")
      when IDLE
        schedule_internal @wait
    end
  end

  private
  def on_change(modified, added, removed)
    unless modified.empty? && added.empty?
      logger.info("Change detected")
      schedule
    end
  end

  def do_upload
    logger.info("Checking files for upload")

    @processing.set PROCESSING

    files = Dir["#{@config[:watch_dir]}/*.*"]
    reprocess = false

    old, new = files.partition { |x| Time.now > File.mtime(x) + @wait }

    unless new.empty?
      logger.info("Found #{new.length} new files")
      @processing.set SCHEDULED
    end

    unless old.empty?
      logger.info("Found #{old.length} stable files")
      reprocess = @processor.run(old)
    end

    @processing.compare_and_swap(PROCESSING, IDLE)

    if reprocess
      schedule_internal 1
    elsif @processing.get == SCHEDULED
      logger.info("Rescheduling...")
      schedule_internal @wait
    end
  end

  def schedule_internal(secs)
    if @scheduler.jobs.empty?
      logger.info("Scheduling processing for #{secs} sec")
      @scheduler.in secs do
        do_upload
      end
    else
      logger.info('Already scheduled')
    end
  end

end
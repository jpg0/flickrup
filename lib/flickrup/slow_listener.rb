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

  IDLE = 0
  SCHEDULED = 1
  PROCESSING = 2

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
    Listen.to(@config['watch_dir']) do |modified, added, _|
      unless modified.empty? && added.empty?
        on_change
      end
    end
  end

  def do_upload
    logger.debug("Checking files for upload")

    @processing.set PROCESSING

    files = Dir["#{@config['watch_dir']}/*.{jpg,JPG,jpeg,JPEG}"]

    old, new = files.partition { |x| Time.now > File.mtime(x) + @wait }

    unless new.empty?
      logger.debug("Found #{new.length} new files")
      @processing.set SCHEDULED
    end

    unless old.empty?
      logger.debug("Found #{old.length} stable files")
      @processor.run(old)
    end

    @processing.compare_and_swap(PROCESSING,IDLE)

    if @processing.get == SCHEDULED
      logger.debug("Rescheduling...")
      schedule
    end
  end

  def on_change
    logger.debug("Change detected")

    previous_state = @processing.swap(SCHEDULED)

    case previous_state
      when SCHEDULED
        #ignore, already scheduled
        logger.debug("Already scheduled")
      when PROCESSING
        #ignore, processing thread will reschedule
        logger.debug("Scheduling queued")
      when IDLE
        schedule
    end
  end

  def schedule
    if @scheduler.jobs.empty?
      logger.debug("Scheduling processing for #@wait sec")
      @scheduler.in @wait do
        do_upload
      end
    else
      logger.debug("Already scheduled")
    end
  end

end
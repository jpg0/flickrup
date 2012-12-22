require 'rubygems'
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

  WAIT = 5#*60 #5 min

  def initialize(config)
    @config = config
    @scheduler = Rufus::Scheduler.start_new
    @processor = Processor.new(config)
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

    old, new = files.partition { |x| Time.now > File.mtime(x) + WAIT }

    unless old.empty?
      logger.debug("Found #{old.length} stable files")
      @processor.run(old)
    end

    unless new.empty?
      logger.debug("Found #{new.length} new files")
      #todo: schedule for latest
      schedule
    end
  end

  def on_change
    logger.debug("Change detected")
    schedule
  end

  def schedule
    if @scheduler.jobs.empty?
      logger.debug("Scheduling processing for #{WAIT} sec")
      @scheduler.in WAIT do
        do_upload
      end
    end
  end

end
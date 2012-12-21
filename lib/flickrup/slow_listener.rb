require 'rubygems'
require 'rufus/scheduler'

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

  WAIT = 5#*60 #5 min

  def initialize(config)
    @config = config
    @scheduler = Rufus::Scheduler.start_new
    @processor = Processor.new(config)
  end

  def run
    Listen.to(@config['watch_dir']) do |modified, added, removed|
      unless modified.empty? && added.empty?
        on_change
      end
    end
  end

  def do_upload
    files = Dir["#{@config['watch_dir']}/*.{jpg,JPG,jpeg,JPEG}"]

    old, new = files.partition { |x| Time.now > File.mtime(x) + WAIT }

    unless old.empty?
      @processor.run(old)
    end

    unless new.empty?
      #todo: schedule for latest
      on_change
    end
  end

  def on_change
    if @scheduler.jobs.empty?
      @scheduler.in WAIT do
        do_upload
      end
    end
  end

end
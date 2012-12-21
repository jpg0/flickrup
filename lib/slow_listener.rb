require 'rubygems'
require 'rufus/scheduler'

class SlowListener

  WAIT = 5#*60 #5 min

  def initialize(config)
    @config = config
    @scheduler = Rufus::Scheduler.start_new
    @flickr = Flickrup.new(config)
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
      @flickr.run(old)
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
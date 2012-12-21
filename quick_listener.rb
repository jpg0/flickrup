require 'atomic'
require "thread"
require "listen"

class QuickListener

  def initialize(config)
    @config = config
    @changes = Atomic.new 0
    @mutex = Mutex.new
    @sync = ConditionVariable.new
    @flickup = Flickrup.new(config)
  end

  def run
    Thread.new { run_loop }
    Listen.to(@config['watch_dir']) do |modified, added, removed|
      unless modified.empty? && added.empty?
        on_change
      end
    end
  end

  private
  def run_loop

    @mutex.synchronize do
      while 1 do
        reloop = @changes.swap 0

        while reloop > 0 do
          execute
          reloop = @changes.swap 0
        end

        @sync.wait(@mutex)
      end
    end

  end

  def on_change
    @mutex.synchronize do
      @changes.set 1
      @sync.signal
    end
  end

  def execute
    @flickup.run
  end
end
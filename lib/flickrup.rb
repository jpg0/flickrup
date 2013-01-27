require "yaml"
require "rubygems"
require "trollop"
require "logger"
require "flickrup/quick_listener"
require "flickrup/slow_listener"
require "flickrup/tagged_image"
require "flickrup/uploader"
require "flickrup/processor"
require "flickrup/logging"


module Flickrup

  def self.run_cmd(args = ARGV)
    Runner.new.run(args)
  end

  class Runner
    include Logging

    def run(args = ARGV)
      opts = Trollop::options do
        opt :once, "Run single pass" # flag --once, default false
        opt :nowait, "Upload immediately on change, no delay"
        opt :loglevel, "Logging level", :type => :string, :default => "INFO"
        opt :logfile, "Log file (-- for STDOUT)", :type => :string, :default => "--"
        opt :config, "Configuration file", :type => :string, :default => './flickrup.yml'
        opt :timeout, "Timeout/Wait for file modification", :default => 5*60
        opt :dryrun, "Dry run, don't move or upload (useful for debugging)"
      end

      Logging.pre_init(Logger::Severity.const_get(opts.loglevel), opts.logfile == "--" ? STDOUT : opts.logfile)

      #load config, symbolise, merge with command line config
      config = YAML::load(File.open(opts.config, "r")).inject({}){|memo,(k,v)| memo[k.to_sym] = v; memo}.merge opts

      if opts.once
        Processor.new(config).run
      elsif opts.nowait
        QuickListener.new(config).run
      else
        SlowListener.new(config).run
      end

      # suspend thread so we don't exit immediately
      @mutex = Mutex.new
      @running = ConditionVariable.new

      @mutex.synchronize do
        @running.wait(@mutex)
      end

    end
  end
end

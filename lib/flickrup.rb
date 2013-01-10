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
        opt :config, "Configuration file", :type => :string, :default => './flickrup.yml'
        opt :timeout, "Timeout/Wait for file modification", :default => 5*60
        opt :dryrun, "Dry run, don't move or upload (useful for debugging)"
        opt :tagsetprefix, "Provide tagset support with supplied tag prefix", :type => :string, :default => 'grouping.set'
      end

      logger.level = Logger::Severity.const_get(opts.loglevel)

      config = YAML::load(File.open(opts.config, "r")).merge opts

      if opts.once
        Processor.new(config).run
      elsif opts.nowait
        QuickListener.new(config).run
      else
        SlowListener.new(config).run
      end
    end
  end
end

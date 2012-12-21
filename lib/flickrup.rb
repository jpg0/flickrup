require "yaml"
require "rubygems"
require "trollop"
require "logger"
require "flickrup/quick_listener"
require "flickrup/slow_listener"
require "flickrup/tagged_image"
require "flickrup/uploader"
require "flickrup/processor"


module Flickrup

  @log = Logger.new(STDOUT)

  def self.log
    @log
  end

  def self.run_cmd(args = ARGV)
    opts = Trollop::options do
      opt :once, "Run single pass" # flag --once, default false
      opt :nowait, "Upload immediately on change, no delay"
      opt :loglevel, "Logging level", :type => :string, :default => "INFO"
      opt :config, "Configuration file", :type => :string, :default => './flickrup.yml'
    end

    @log.sev_threshold = Logger::Severity.const_get(opts.loglevel)

    config = YAML::load(File.open(opts.config, "r"))

    if opts.once
      Processor.new(config).run
    elsif opts.nowait
      QuickListener.new(config).run
    else
      SlowListener.new(config).run
    end
  end

end

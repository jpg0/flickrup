require "yaml"
require "rubygems"
require "trollop"
require "flickr/quick_listener"
require "flickr/slow_listener"
require "flickr/tagged_image"
require "flickr/uploader"

module Flickrup
  class Processor

    def initialize(config)
      @config = config
    end

    def run(files = Dir["#{@config['watch_dir']}/*.{jpg,JPG,jpeg,JPEG}"])
      images = files.map { |x| TaggedImage.new(x) }

      tagged = images.select { |x| !(x.tags.nil? || x.tags.empty?) }

      uploader = Uploader.new(@config["api_key"], @config["shared_secret"])

      tagged.each { |x|
        x.doUpload(uploader)
        x.archive(@config["archive_dir"])
      }

      tagged.length
    end
  end

  @log = Logger.new(STDOUT)

  def self.log
    @log
  end

  def run_cmd
    opts = Trollop::options do
      opt :once, "Run single pass" # flag --once, default false
      opt :nowait, "Upload immediately on change, no delay"
    end

    config = YAML::load(File.open('./flickrup.yml', "r"))

    if opts.once
      Processor.new(config).run
    elsif opts.nowait
      QuickListener.new(config).run
    else
      SlowListener.new(config).run
    end
  end

end

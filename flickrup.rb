require "yaml"
require "rubygems"
require "trollop"
require "quick_listener"
require "slow_listener"
require "tagged_image"
require "uploader"

class Flickrup

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

opts = Trollop::options do
  opt :once, "Run single pass" # flag --once, default false
  opt :nowait, "Upload immediately on change, no delay"
end

config = YAML::load(File.open('./flickrup.yml', "r"))

if opts.once
  Flickrup.new(config).run
elsif opts.nowait
  QuickListener.new(config).run
else
  SlowListener.new(config).run
end

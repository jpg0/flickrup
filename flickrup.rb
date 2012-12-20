require "yaml"
require "rubygems"
require "trollop"
require "listener"

class Flickrup

  def initialize(config)
    @config = config
  end

  def run
    files = Dir["#{@config['watch_dir']}/*.{jpg,JPG,jpeg,JPEG}"]

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
  opt :once, "Run single pass"                      # flag --once, default false
end

config = YAML::load(File.open('./flickrup.yml', "r"))

if opts.once
  Flickrup.new(config).run
else
  Listener.new(config).run
end

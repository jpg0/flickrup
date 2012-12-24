require "flickrup/tagged_image"
require "flickrup/uploader"
require "flickrup/logging"

class Processor
  include Logging

  def initialize(config)
    @config = config
  end

  def paused
    Dir["#{@config['watch_dir']}/paused"].length > 0
  end

  def run(files = Dir["#{@config['watch_dir']}/*.{jpg,JPG,jpeg,JPEG}"])

    logger.debug("Processing #{files.length} files")

    if paused
      logger.debug("Paused detected, so skipping processing...")
      return 0
    end

    images = files.map { |x| TaggedImage.new(x) }

    tagged = images.select { |x| !(x.tags.nil? || x.tags.empty?) }

    uploader = Uploader.new(@config["api_key"], @config["shared_secret"])

    tagged.each { |x|
      x.doUpload(uploader)
      x.archive(@config["archive_dir"])
    }

    logger.debug("Processed #{tagged.length} valid files")

    tagged.length
  end
end

require "flickrup/tagged_image"
require "flickrup/uploader"
require "flickrup/logging"
require "flickrup/ext/tag_sets"
require "flickrup/ext/replacement_tags"
require "flickrup/ext/machine_tags"
require "flickrup/ext/visibility"

class Processor
  include Logging

  def initialize(config)
    @config = config
  end

  def paused
    Dir["#{@config[:watch_dir]}/paused"].length > 0
  end

  def run(files = Dir["#{@config[:watch_dir]}/*.{jpg,JPG,jpeg,JPEG}"])

    logger.debug("Processing #{files.length} files")

    if paused
      logger.debug("Paused detected, so skipping processing...")
      return 0
    end

    images = files.map { |x| TaggedImage.new(x) }

    tagged = images.select { |x| !(x.tags.nil? || x.tags.empty?) }

    logger.debug("Selected #{tagged.length} files for upload")

    uploader = get_uploader

    if @config[:dryrun]
      logger.info("Would now upload & move #{tagged.length} files, pausing for 5 secs instead")
      sleep 5
    else
      tagged.each do |x|
        begin
          logger.debug("Beginning upload for #{x.filename}...")
          x.doUpload(get_preupload_handlers, uploader)
          logger.debug("Upload complete for #{x.filename}")
          logger.debug("Beginning archive for #{x.filename}...")
          x.archive(@config[:archive_dir])
          logger.debug("Processing complete for #{x.filename}")
        rescue => e
          logger.error("Failed to process #{x.filename} due to exception: #{e}")
        end
      end
    end

    logger.debug("Processed #{tagged.length} valid files")

    tagged.length
  end

  private
  def get_uploader
    flickr_client = FlickrClient.new(@config[:api_key], @config[:shared_secret])

    context = {:flickr_client => flickr_client}

    uploader = Uploader.new(flickr_client)

    upload_handlers.reduce(uploader) do |result, clazz|
      clazz.new(@config, result, context)
    end
  end

  def get_preupload_handlers
    [MachineTags.new(@config, {}), ReplacementTags.new(@config)]
  end

  def upload_handlers()
    [TagSets, Visibility]
  end

end


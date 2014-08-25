require "flickrup/filetype/tagged_file"
require "flickrup/uploader"
require "flickrup/logging"
require "flickrup/ext/tag_sets"
require "flickrup/ext/replacement_tags"
require "flickrup/ext/machine_tags"
require "flickrup/ext/visibility"
require "flickrup/ext/sidecar"
require "flickrup/archiver"
require "flickrup/ext/associated_files"

class Processor
  include Logging

  def initialize(config)
    @MAX_DATE = Time.at(2130578470)
    @config = config
  end

  def paused
    Dir["#{@config[:watch_dir]}/paused"].length > 0
  end

  def run(all_files = Dir["#{@config[:watch_dir]}/*.*"])

    logger.info("Processing #{all_files.length} files")

    if paused
      logger.info("Paused detected, so skipping processing...")
      return 0
    end

    # some files may have been removed by now: https://bitbucket.org/jgilbert/flickrup/issue/1
    files = all_files.map { |x|
      if File.exist? x
        TaggedFile.create(x)
      else
        logger.debug("File #{x} already removed, skipping.")
      end
    }.compact

    logger.debug("Scanned #{file.length} files")

    by_date = files.sort_by { |x|
      x.date_created || @MAX_DATE
    }

    tagged = by_date.take_while { |x|
      !(x.tags.nil? || x.tags.empty?)
    }

    if tagged.empty?
      logger.info('No files selected for upload')
      unless by_date.empty?
        logger.debug("Stopped on #{by_date[0].filename}")
      end
    else
      logger.info("Selected #{tagged.length} files for upload")
    end

    uploader = get_uploader

    if @config[:dryrun]
      logger.info("Would now upload & move #{tagged.length} files, pausing for 5 secs instead")
      sleep 5
    else
      tagged.each do |x|
        begin
          logger.info("Beginning upload for #{x.filename}...")
          x.doUpload(get_preupload_handlers, uploader)
          logger.info("Upload complete for #{x.filename}")
          logger.info("Beginning archive for #{x.filename}...")
          x.archive(get_archiver, @config[:archive_dir])
          logger.info("Processing complete for #{x.filename}")
        rescue => e
          logger.error("Failed to process #{x.filename} due to exception: #{e}")
          logger.info(e.backtrace.join("\n"))
        end
      end
    end

    logger.info("Processed #{tagged.length} valid files")

    tagged.length
  end

  private
  def get_uploader

    unless @uploader

      context = get_chain_context()

      uploader = Uploader.new(context[:flickr_client])

      @uploader = upload_handlers.reduce(uploader) do |result, clazz|
        clazz.new(@config, result, context)
      end
    end

    return @uploader
  end

  def get_archiver
    unless @archiver
      @archiver = archive_handlers.reduce(Archiver.new) do |result, clazz|
        clazz.new(@config, result, get_chain_context)
      end
    end
    return @archiver
  end

  def get_chain_context
    unless @chain_context
      flickr_client = FlickrClient.new(@config[:api_key], @config[:shared_secret])
      @chain_context = {:flickr_client => flickr_client}
    end

    return @chain_context
  end

  def get_preupload_handlers
    [MachineTags.new(@config, {}), ReplacementTags.new(@config)]
  end

  def upload_handlers
    [TagSets, Visibility]
  end

  def archive_handlers
    [TagSets, Sidecar, AssociatedFiles]
  end

end


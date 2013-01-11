require 'mini_exiftool'
require 'fileutils'
require "flickrup/logging"

class TaggedImage
  include Logging

  attr_reader :file

  def initialize(file)
    @file = file
    @parsed = MiniExiftool.new @file
  end

  def tags
    keywords = @parsed.keywords
    if keywords.class == String
      [keywords]
    else
      keywords
    end
  end

  def doUpload(preupload_handlers, uploader)

    config = preupload_handlers.reduce({:tags => tags.join(" ")}) do |config, handler|
      handler.preupload(file, self, config)
    end

    #maybe save afterwards
    if @parsed.changed?
      @parsed.save!
    end

    uploader.upload(file, self, config)

    #maybe save afterwards
    if @parsed.changed?
      @parsed.save!
    end
  end

  def archive(to_dir)
    date_taken = @parsed.date_time_original || @parsed.modifydate

    target_dir = File.join(
        to_dir,
        date_taken.strftime("%y"),
        date_taken.strftime("%m")
    )

    FileUtils.makedir_p(target_dir)
    FileUtils.mv(@file, File.join(target_dir, File.basename(@file)))
    logger.info("Archived #{File.basename(@file)} to #{target_dir}")
  end

  # Returns the value of a tag.
  def [] tag
    @parsed[tag]
  end

  # Set the value of a tag.
  def []=(tag, val)
    @parsed[tag] = val
  end

end
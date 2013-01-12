require 'mini_exiftool'
require 'fileutils'
require "flickrup/logging"
require "flickrup/tag_set"
require "flickrup/processing_context"

class TaggedImage
  include Logging

  attr_reader :filename

  def initialize(filename)
    @filename = filename
    @parsed = MiniExiftool.new filename
  end

  def tags
    keywords = @parsed.keywords
    TagSet.new(if keywords == nil
      []
    elsif keywords.class == String
      [keywords]
    else
      keywords
    end)
  end

  def doUpload(preupload_handlers, uploader)

    context = preupload_handlers.reduce(ProcessingContext.new(self)) do |ctx, handler|
      handler.preupload(ctx)
    end

    #maybe save afterwards
    if @parsed.changed?
      @parsed.save!
    end

    uploader.upload(context)

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
        date_taken.strftime("%b")
    )

    FileUtils.mkdir_p(target_dir)
    FileUtils.mv(@filename, File.join(target_dir, File.basename(@filename)))
    logger.info("Archived #{File.basename(@filename)} to #{target_dir}")
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
require 'mini_exiftool'

class TaggedFile

  @@subclasses = {}

  def self.register_reader(exts)
    exts.map do |ext|
      @@subclasses[ext] = self
    end
  end

  def self.create(file)
    c = @@subclasses[File.extname(file)[1..-1]]

    if c
      c.new(file)
    end
  end

  def initialize(filename)
    @parsed = MiniExiftool.new filename
  end

  def doUpload(preupload_handlers, uploader)

    context = preupload_handlers.reduce(ProcessingContext.new(self)) do |ctx, handler|
      handler.preupload(ctx)
    end

    pre_upload
    uploader.upload(context)
    post_upload
  end

  def tags
    keywords = self['keywords']
    TagSet.new(if keywords == nil
                 []
               elsif keywords.class == String
                 [keywords]
               else
                 keywords
               end)
  end

  def archive(to_dir)
    date_taken = @parsed.date_time_original || @parsed.modifydate
    archive_by_date(to_dir, date_taken)
  end

  def pre_upload

  end

  def post_upload

  end

  private

  def archive_by_date(to_dir, date)
    archive_file_by_date(@filename, to_dir, date)
  end

  def archive_file_by_date(file, to_dir, date)
    target_dir = File.join(
        to_dir,
        date.strftime("%Y"),
        date.strftime("%b")
    )

    FileUtils.mkdir_p(target_dir)
    FileUtils.mv(file, File.join(target_dir, File.basename(file)))
    logger.info("Archived #{File.basename(file)} to #{target_dir}")
  end

end

require 'flickrup/filetype/tagged_image'
require 'flickrup/filetype/tagged_video'
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
               elsif keywords.class == Integer
                 keywords.to_s
               else
                 keywords.map do |keyword|
                   keyword.to_s
                 end
               end)
  end

  def archive(archive_handlers, to_dir)
    archive_handlers.archive(ProcessingContext.new(self, {
        :archive_dir => to_dir,
        :date_taken => @parsed.date_time_original || @parsed.modifydate || @parsed.filemodifydate,
        :filename => @filename
    }))
  end

  def pre_upload

  end

  def post_upload

  end

  def date_created
    @parsed.date_time_original
  end

end

require 'flickrup/filetype/tagged_image'
require 'flickrup/filetype/tagged_video'
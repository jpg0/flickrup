require 'mini_exiftool'
require 'ftools'
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
    @parsed.keywords.to_a
  end

  def doUpload(uploader)
    uploader.upload(file, tags)
  end

  def archive(to_dir)
    date_taken = @parsed.date_time_original || @parsed.modifydate

    target_dir = File.join(
        to_dir,
        date_taken.strftime("%y"),
        date_taken.strftime("%m")
    )

    File.makedirs(target_dir)
    FileUtils.mv(@file, File.join(target_dir, File.basename(@file)))
  end

end
require 'fileutils'
require "flickrup/logging"
require "flickrup/tag_set"
require "flickrup/processing_context"

class TaggedImage < TaggedFile
  include Logging

  TYPE = 'image'

  attr_reader :filename

  def initialize(filename)
    super
    @filename = filename
  end

  def pre_upload
    super
    #maybe save afterwards
    if @parsed.changed?
      @parsed.save!
    end
  end

  def post_upload
    super
    #maybe save afterwards
    if @parsed.changed?
      @parsed.save!
    end
  end

  # Returns the value of a tag.
  def [] tag
    @parsed[tag]
  end

  # Set the value of a tag.
  def []=(tag, val)
    @parsed[tag] = val
  end

  def type
    TYPE
  end

  register_reader(%w(jpg JPG jpeg JPEG))
end
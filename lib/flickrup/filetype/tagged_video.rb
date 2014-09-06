require 'yaml'
require 'flickrup/filetype/picasa_ini'

class TaggedVideo < TaggedFile
  include Logging

  attr_reader :filename

  TYPE = 'video'

  def initialize(filename)
    super
    @filename = filename
    @properties = PicasaIni.open(filename)
  end

  # Returns the value of a tag.
  def [] tag
    @properties[tag]
  end

  # Set the value of a tag.
  def []=(tag, val)
    @properties[tag] = val
  end

  def type
    TYPE
  end

  register_reader(%w(mp4 MP4))

end

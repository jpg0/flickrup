require 'yaml'
require 'flickrup/filetype/picasa_ini'

class TaggedVideo < TaggedFile
  include Logging

  attr_reader :filename

  def initialize(filename)
    super
    @filename = filename
    @properties = PicasaIni.open(filename)
  end

  def archive(to_dir)
    super
    write_sidecar(File.join(target_dir, "#{File.basename(file)}.meta"))
  end

  def write_sidecar(to_file)
    File.open(to_file, "w") do |f|
      f.write({:keywords => keywords}.to_yaml)
    end
  end

  # Returns the value of a tag.
  def [] tag
    @properties[tag]
  end

  # Set the value of a tag.
  def []=(tag, val)
    @properties[tag] = val
  end

  register_reader(%w(mp4 mts))

end

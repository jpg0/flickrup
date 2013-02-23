require 'flickrup/filetype/inifile'

class PicasaIni

  @@cached_mtime = DateTime.new(1900,1,1)

  def self.open(file)
    dir = File.dirname(file)
    inifile = "#{dir}/.picasa.ini"
    mtime = File.mtime(inifile)

    if @@cached_mtime == mtime then
      @@cached
    else
      @@cached_mtime = mtime
      @@cached = new(inifile)
    end

    @@cached[File.basename(file)]
  end

  def initialize(file)
    @inifile = IniFile.load(file, :escape => false)
  end

  def [](section)
    @inifile[section]
  end
end
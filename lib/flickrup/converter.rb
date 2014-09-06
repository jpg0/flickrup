require 'fileutils'
require 'tmpdir'

class Converter
  include Logging

  def initialize(mappings)
    @mappings = mappings
    @move_dir = Dir.tmpdir
  end

  def maybe_convert(file)
    mapping = @mappings[File.extname(file)[1..-1].downcase]

    if mapping then
      to_exec = mapping % {:file => file}
      logger.debug("Running: #{to_exec}")
      result = system(to_exec)

      if result
        logger.info("Removing converted source file #{file}")
        logger.debug("Moving #{file} to #{@move_dir}")
        FileUtils.mv(file, @move_dir)
      else
        logger.info("Failed to convert source file #{file}")
      end

      return result
    else
      false
    end
  end
end
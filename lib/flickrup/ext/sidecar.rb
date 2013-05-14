require 'tempfile'

class Sidecar
  def initialize(config, delegate, context)
    @delegate = delegate
  end

  def archive(ctx)
    @delegate.archive(ctx)

    if ctx.file.type == TaggedVideo::TYPE
      file = write_sidecar(ctx)
      @delegate.archive(ctx.with_properties({:filename => file}))
    end

  end

  def write_sidecar(ctx)
    #write temp file
    to_file = "#{Dir::tmpdir}/#{File.basename(ctx.properties[:filename])}.meta"
    File.open(to_file, "w") do |f|
      f.write({:keywords => ctx.file['keywords']}.to_yaml)
    end

    return to_file
  end

end
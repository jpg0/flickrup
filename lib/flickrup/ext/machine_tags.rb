class MachineTags
  def initialize(config, context)
    @config = config
  end

  def preupload(ctx)

    tags = ctx.image.tags

    remapped = tags.map do |x|
      x.sub(/^([^:]+):([^:]+)::/,"\\1:\\2=")
    end.uniq



    unless tags == remapped #write back to file
      ctx.image['Keywords'] = remapped
      ctx.image['Subject'] = remapped #also write Subject XMP tag
    end

    ctx

  end
end
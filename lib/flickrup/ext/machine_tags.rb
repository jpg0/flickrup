class MachineTags
  def initialize(config, context)
    @config = config
  end

  def preupload(image, meta, config)

    tags = config[:tags]

    remapped = tags.map do |x|
      x.sub(/^([^:]+):([^:]+)::/,"\\1:\\2=")
    end.uniq



    unless tags == remapped #write back to file
      meta['Keywords'] = remapped
      meta['Subject'] = remapped #also write Subject XMP tag
    end

    config.merge({:tags => TagSet.new(remapped)})

  end
end
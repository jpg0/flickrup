class MachineTags
  def initialize(config, delegate, context)
    @config = config
    @delegate = delegate
  end

  def upload(image, meta, config)

    tags = config[:tags].split(" ")

    remapped = tags.map do |x|
      x.sub(/^([^:]+):([^:]+)::/,"\\1:\\2=")
    end

    @delegate.upload(image, meta, config.merge({:tags => remapped.join(" ")}))

  end
end
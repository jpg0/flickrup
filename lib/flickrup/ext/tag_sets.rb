class TagSets
  def initialize(config, delegate)
    @config = config
    @delegate = delegate
  end

  def upload(image, meta, config)
    all_tags = config[:tags].split(" ")

    prefix = @config[:tagsetprefix]

    (set_tags, tags) = all_tags.partition do |x|
      x.start_with?("#{prefix}")
    end

    sets = set_tags.map do |x|
      x[(prefix.length + 1)..x.length]
    end

    @delegate.upload(image, meta, config.merge({:tags => tags.join(" ")}))

  end
end
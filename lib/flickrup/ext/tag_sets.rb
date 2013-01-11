class TagSets
  def initialize(config, delegate, context)
    @config = config
    @delegate = delegate
    @flickr = context[:flickr_client]
  end

  def upload(image, meta, config)
    all_tags = config[:tags].split(" ")

    prefix = @config[:tagsetprefix]

    (set_tags, tags) = all_tags.partition do |x|
      x.start_with?("#{prefix}")
    end

    sets = set_tags.map do |x|
      x[(prefix.length)..x.length]
    end

    uploaded = @delegate.upload(image, meta, config.merge({:tags => tags.join(" ")}))

    sets.each do |set|
      @flickr.add_to_set(uploaded, set)
    end
  end
end
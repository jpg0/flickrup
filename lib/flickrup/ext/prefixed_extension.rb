class PrefixedExtension
  def initialize(config)
    @config = config
  end

  def values(prefix_key, ctx)

    prefix = @config[prefix_key]

    if prefix
      all_tags = ctx.image.tags

      (selected_tags, tags) = all_tags.partition do |x|
        x.start_with?("#{prefix}")
      end

      selected_tags.map do |x|
        x[(prefix.length)..x.length]
      end
    else
      []
    end
  end
end
class PrefixedExtension
  def initialize(config)
    @config = config
  end

  def values(prefix_key, config)
    all_tags = config[:tags].split(" ")

    prefix = @config[prefix_key]

    (selected_tags, tags) = all_tags.partition do |x|
      x.start_with?("#{prefix}")
    end

    selected_tags.map do |x|
      x[(prefix.length)..x.length]
    end
  end
end
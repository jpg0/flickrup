class ReplacementTags
  def initialize(config, delegate)
    @config = config
    @delegate = delegate
  end

  def upload(image, meta, config)

    replacements = @config["tag_replacements"]

    if replacements != nil
      replacements.each do |tag_name, tag_value_replacements|
        existing = meta[tag_name]

        tag_value_replacements.each do |key, value|
          if existing == key
            meta[tag_name] = value
            break
          end
        end
      end
    end

    @delegate.upload(image, meta, config)
  end
end
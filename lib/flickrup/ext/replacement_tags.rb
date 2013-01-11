class ReplacementTags
  def initialize(config)
    @config = config
  end

  def preupload(image, meta, config)

    replacements = @config[:tag_replacements]
    new_config = config

    if replacements != nil
      replacements.each do |tag_name, tag_value_replacements|
        existing = meta[tag_name]

        tag_value_replacements.each do |key, value|
          if existing == key
            meta[tag_name] = value
            new_config = new_config.merge({tag_name => value})
            break
          end
        end
      end
    end

    new_config
  end
end
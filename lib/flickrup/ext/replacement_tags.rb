class ReplacementTags
  def initialize(config)
    @config = config
  end

  def preupload(ctx)

    replacements = @config[:tag_replacements]
    all_tags = ctx.file.tags

    if replacements != nil
      replacements.each do |tag_name, tag_value_replacements|
        existing = ctx.file[tag_name]

        tag_value_replacements.each do |key, value|
          if key.start_with? '$'
            if all_tags.include? key[1..-1]
              ctx.file[tag_name] = value
              break
            end
          else
            if existing == key
              ctx.file[tag_name] = value
              break
            end
          end
        end
      end
    end

    ctx
  end
end
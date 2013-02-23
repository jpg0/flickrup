class ReplacementTags
  def initialize(config)
    @config = config
  end

  def preupload(ctx)

    replacements = @config[:tag_replacements]

    if replacements != nil
      replacements.each do |tag_name, tag_value_replacements|
        existing = ctx.file[tag_name]

        tag_value_replacements.each do |key, value|
          if existing == key
            ctx.file[tag_name] = value
            break
          end
        end
      end
    end

    ctx
  end
end
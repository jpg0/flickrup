require "flickrup/ext/prefixed_extension"
require "flickrup/logging"

class TagSets < PrefixedExtension
  include Logging

  def initialize(config, delegate, context)
    super(config)
    @delegate = delegate
    @flickr = context[:flickr_client]
  end

  def upload(ctx)
    sets = values(:tagsetprefix, ctx).uniq

    uploaded = @delegate.upload(ctx)

    unless ctx.properties[UploaderConstants::NO_UPLOAD] == true
      sets.each do |set|
        logger.debug("Adding #{ctx.file.filename} to set: #{set}")
        @flickr.add_to_set(uploaded, set)
      end
    end
  end

  def archive(ctx)
    sets = values(:tagsetprefix, ctx).uniq

    if sets.length == 1
      ctx.properties[:subdir] = sets[0]
      set_date = @flickr.date_of_set(sets[0])
      ctx.properties[:date_taken] = set_date
    end
    @delegate.archive(ctx)
  end
end
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

    sets.each do |set|
      logger.debug("Adding #{ctx.file.filename} to set: #{set}")
      @flickr.add_to_set(uploaded, set)
    end
  end
end
require "flickrup/ext/prefixed_extension"

class TagSets < PrefixedExtension
  def initialize(config, delegate, context)
    super(config)
    @delegate = delegate
    @flickr = context[:flickr_client]
  end

  def upload(image, meta, config)
    sets = values(:tagsetprefix, config)

    uploaded = @delegate.upload(image, meta, config)

    sets.each do |set|
      @flickr.add_to_set(uploaded, set)
    end
  end
end
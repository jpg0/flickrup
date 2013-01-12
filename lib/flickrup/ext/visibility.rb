require "flickrup/logging"

class Visibility < PrefixedExtension
  include Logging

  def initialize(config, delegate, context)
    super(config)
    @delegate = delegate
  end

  def upload(ctx)
    visibility = values(:visibilityprefix, ctx)

    config_override = if visibility.length == 0
      logger.debug("No visibility specified for #{ctx.image.filename}.")
      {}
    elsif visibility.length == 1
      config_for_visibility(visibility[0])
    else #multiple visibilities
      logger.warn("Multiple visibilities specified. Using the first encountered: #{visibility[0]}.")
      config_for_visibility(visibility[0])
    end

    @delegate.upload(ctx.with_properties(config_override))
  end

  def config_for_visibility(visibility)
    if visibility == "offline"
      logger.debug("Offline Visibility specified. Disabling upload.")
      {UploaderConstants::NO_UPLOAD => true}
    elsif visibility == "private"
      logger.debug("Private visibility specified.")
      {"is_public" => "0", "is_friend" => "0", "is_family" => "0"}
    elsif visibility == "friends"
      logger.debug("Friends visibility specified.")
      {"is_public" => "0", "is_friend" => "1", "is_family" => "1"}
    else
      logger.warn("Unknown visibility #{visibility}. Ignoring.")
      {}
    end
  end
end

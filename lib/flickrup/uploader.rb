require "yaml"
require 'rubygems'
require "flickrup/logging"
require "flickrup/flickr_client"

module UploaderConstants
  NO_UPLOAD = :no_upload
end

class Uploader
  include Logging

  def initialize(flickr_client)
    @flickr = flickr_client
  end

  def upload(ctx)

    if ctx.upload_properties[UploaderConstants::NO_UPLOAD]
      logger.info("Skipping upload of #{ctx.image.filename} as marked as NO_UPLOAD")
    else
      logger.debug("Uploading #{ctx.image.filename} with tags #{ctx.image.tags}")

      begin
        rv = @flickr.upload_photo(ctx.image.filename, prepare_upload_properties(ctx))
        logger.debug("Uploaded #{ctx.image.filename}")
      rescue StandardError => err
        logger.error("Failed to upload #{ctx.image.filename}: #{err}")
      end
    end

    rv
  end

  def prepare_upload_properties(ctx)
    {:tags => ctx.image.tags.to_s}.merge(ctx.upload_properties)
  end

end
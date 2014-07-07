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

    if ctx.properties[UploaderConstants::NO_UPLOAD]
      logger.info("Skipping upload of #{ctx.file.filename} as marked as NO_UPLOAD")
    else
      logger.info("Uploading #{ctx.file.filename} with tags #{ctx.file.tags}")

      begin
        rv = @flickr.upload_photo(ctx.file.filename, prepare_upload_properties(ctx))
        logger.info("Uploaded #{ctx.file.filename}")
      rescue StandardError => err
        logger.error("Failed to upload #{ctx.file.filename}: #{err}")
      end
    end

    rv
  end

  def prepare_upload_properties(ctx)
    {:tags => ctx.file.tags.to_s}.merge(ctx.properties)
  end

end
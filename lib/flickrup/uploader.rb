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

  def upload(image, meta, config)

    if config[UploaderConstants::NO_UPLOAD]
      logger.info("Skipping upload of #{image} as marked as NO_UPLOAD")
    else
      logger.debug("Uploading #{image} with tags #{config[:tags]}")

      begin
        rv = @flickr.upload_photo(image, config)
        logger.debug("Uploaded #{image}")
      rescue StandardError => err
        logger.error("Failed to upload #{image}: #{err}")
      end
    end

    rv
  end

end
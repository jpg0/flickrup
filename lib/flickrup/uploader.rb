require "yaml"
require 'rubygems'
require "flickrup/logging"
require "flickrup/flickr_client"

class Uploader < FlickrClient
  include Logging

  def initialize(flickr_client)
    @flickr = flickr_client
  end

  def upload(image, meta, config)
    logger.debug("Uploading #{image} with tags #{config[:tags]}")

    begin
      rv = @flickr.upload_photo(image, config)
      logger.debug("Uploaded #{image}")
    rescue StandardError => err
      logger.error("Failed to upload #{image}: #{err}")
    end

    rv
  end

end
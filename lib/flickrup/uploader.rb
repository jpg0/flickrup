require "yaml"
require 'rubygems'
require "flickraw"
require "flickrup/logging"

class Uploader
  include Logging

  def initialize(api_key, shared_secret)
    FlickRaw.api_key= api_key
    FlickRaw.shared_secret=shared_secret

    @token_file = File.join(ENV['HOME'], ".flickrup")
    @flickr = FlickRaw::Flickr.new

    if File.exist?(@token_file)
      token = YAML.load_file(@token_file)
      @flickr.access_token = token["access"]
      @flickr.access_secret = token["secret"]
    end

    begin
      login = @flickr.test.login
    rescue FlickRaw::OAuthClient::FailedResponse => e
      get_token
    end

  end

  def get_token
    token = @flickr.get_request_token
    auth_url = @flickr.get_authorize_url(token['oauth_token'], :perms => 'write')

    puts "Open this url in your process to complete the authentication process : #{auth_url}"
    puts "Copy here the number given when you complete the process."
    verify = gets.strip

    begin
      @flickr.get_access_token(token['oauth_token'], token['oauth_token_secret'], verify)
      token = {"access" => @flickr.access_token, "secret" => @flickr.access_secret}
      login = @flickr.test.login
      puts "You are now authenticated as #{login.username}"
      File.open(@token_file, 'w') { |f| f.write(token.to_yaml) }
    rescue FlickRaw::OAuthClient::FailedResponse => e
      puts "Authentication failed : #{e.msg}"
    end

  end

  def upload(image, meta, config)
    logger.debug("Uploading #{image} with tags #{config[:tags]}")

    begin
      @flickr.upload_photo(image, config)
      logger.debug("Uploaded #{image}")
    rescue StandardError => err
      logger.error("Failed to upload #{image}: #{err}")
    end
  end

end
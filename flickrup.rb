require "yaml"
require 'rubygems'
require 'mini_exiftool'
require "flickraw"
require 'ftools'
require 'fileutils'

class Flickrup

  def initialize(config_file)
    @config_file = config_file
  end

  def run

    File.open(@config_file, "r") { |f|
      unless f.flock(File::LOCK_EX || File::LOCK_NB)
        return false
      end

      config = YAML::load(f)

      while handle_files(config) > 0
        puts "Checking for additional files"
      end
    }

    return true
  end

  private
  def handle_files(config)
    files = Dir["#{config['watch_dir']}/*.{jpg,JPG,jpeg,JPEG}"]

    images = files.map { |x| TaggedImage.new(x) }

    tagged = images.select { |x| !(x.tags.nil? || x.tags.empty?) }

    uploader = Uploader.new(config["api_key"], config["shared_secret"])

    tagged.each { |x|
      x.doUpload(uploader)
      x.archive(config["archive_dir"])
    }

    tagged.length
  end

end

class TaggedImage

  attr_reader :file

  def initialize(file)
    @file = file
    @parsed = MiniExiftool.new @file
  end

  def tags
    @parsed.keywords.to_a
  end

  def doUpload(uploader)
    uploader.upload(file, tags)
  end

  def archive(to_dir)
    date_taken = @parsed.date_time_original
    target_dir = File.join(
        to_dir,
        date_taken.strftime("%y"),
        date_taken.strftime("%m")
    )

    File.makedirs(target_dir)
    FileUtils.mv(@file, File.join(target_dir, File.basename(@file)))
  end

end

class Uploader

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

  def upload(image, tags)
    @flickr.upload_photo(image, :tags => tags.join(" "))
  end

end

Flickrup.new('./flickrup.yml').run
require "flickraw"
require 'etc'

class FlickrClient

  def initialize(api_key, shared_secret)
    @set_name_to_id = {}
    FlickRaw.api_key= api_key
    FlickRaw.shared_secret=shared_secret

    @token_file = File.join(ENV['HOME'], ".flickrup")
    set_client(FlickRaw::Flickr.new)

    if File.exist?(@token_file)
      token = YAML.load_file(@token_file)
      client.access_token = token["access"]
      client.access_secret = token["secret"]
    end

    begin
      login = client.test.login
    rescue FlickRaw::OAuthClient::FailedResponse => e
      get_token
    end
  end

  def client
    @client
  end

  def upload_photo(image, config)
    client.upload_photo(image, config)
  end

  def add_to_set(photo_id, set_name)
    set_id = get_set_id(set_name, photo_id)
    client.photosets.addPhoto(:photoset_id => set_id, :photo_id => photo_id)
  end

  def get_set_id(set_name, for_photo_id)
    unless @set_name_to_id.include? set_name
      sets_response = client.photosets.getList #if we don't have it, refresh
      @set_name_to_id = Hash[*sets_response.to_a.map{ |set|
        [set['title'], set['id']]
      }.flatten]
      unless @set_name_to_id.include? set_name #if we still don't, create it
        created = client.photosets.create(:title => set_name, :primary_photo_id => for_photo_id)
        @set_name_to_id = @set_name_to_id.merge({set_name => created.id})
      end
    end

    @set_name_to_id[set_name]
  end

  def get_token
    token = client.get_request_token
    auth_url = client.get_authorize_url(token['oauth_token'], :perms => 'write')

    puts "Open this url in your process to complete the authentication process : #{auth_url}"
    puts "Copy here the number given when you complete the process."
    verify = gets.strip

    begin
      client.get_access_token(token['oauth_token'], token['oauth_token_secret'], verify)
      token = {"access" => client.access_token, "secret" => client.access_secret}
      login = client.test.login
      puts "You are now authenticated as #{login.username}"
      File.open(@token_file, 'w') { |f| f.write(token.to_yaml) }
    rescue FlickRaw::OAuthClient::FailedResponse => e
      puts "Authentication failed : #{e.msg}"
    end

  end

  private
  def set_client(client)
    @client = client
  end

end
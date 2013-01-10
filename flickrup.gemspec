Gem::Specification.new do |s|
  s.name        = 'flickrup'
  s.version     = '0.0.8'
  s.date        = '2012-12-21'
  s.summary     = "Flickr Uploader!"
  s.description = "Flickr auto uploading script"
  s.authors     = ["Jonathan Gilbert"]
  s.email       = 'bb..jgilbert@xoxy.net'
  s.files        = Dir["{lib}/**/*.rb", "bin/*", "LICENSE", "*.md"]
  s.executables << 'flickrup'
  s.homepage    = 'http://bitbucket.org/jgilbert/flickrup'
  s.add_runtime_dependency 'atomic', '>= 1.0.0'
  s.add_runtime_dependency 'listen', '>= 0.6.0'
  s.add_runtime_dependency 'rufus-scheduler', '>= 2.0.0'
  s.add_runtime_dependency 'mini_exiftool', '>= 1.6.0'
  s.add_runtime_dependency 'flickraw', '>= 0.9.6'
  s.add_runtime_dependency 'json', '>= 1.7.0' #Flickraw seems to need this, but doesn't bring it in?
  s.add_runtime_dependency 'trollop', '>= 2.0.0'

end
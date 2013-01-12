require 'rubygems'
require 'test/unit'
require "mocha/setup"
require 'flickrup/tagged_image'

class TaggedImageTest < Test::Unit::TestCase
  def test_tags
    testee = TaggedImage.new(nil)

    #mocktool = mock("exiftool")
    #MiniExiftool.expects(:initialize).with(nil).returns(mocktool)



  end
end
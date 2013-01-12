require 'rubygems'
require 'test/unit'
require "mocha/setup"
require 'flickrup/tagged_image'

class TaggedImageTest < Test::Unit::TestCase
  def test_no_tags
    mocktool = mock("exiftool")
    MiniExiftool.stubs(:new).returns(mocktool)
    mocktool.expects(:keywords).returns(nil)

    testee = TaggedImage.new(nil)
    tags = testee.tags

    assert_equal(0, tags.length)

  end
end
require 'rubygems'
require 'test/unit'
require "mocha/setup"
require 'flickrup/filetype/tagged_image'

class TaggedImageTest < Test::Unit::TestCase
  def test_no_tags
    MiniExiftool.stubs(:new).returns({"keywords" => nil})

    testee = TaggedImage.new(nil)
    tags = testee.tags

    assert_equal(0, tags.length)
  end
end
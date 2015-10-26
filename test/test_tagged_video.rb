require 'rubygems'
require 'test/unit'
require 'mocha/test_unit'
require 'flickrup/filetype/tagged_image'

class TaggedImageTest < Test::Unit::TestCase
  def test_extract_tags
    PicasaIni.stubs(:open).returns({"keywords" => "1,2"})

    testee = TaggedVideo.new(nil)
    tags = testee.tags

    assert_equal(["1","2"], tags)

  end

  def test_extract_no_tags
    PicasaIni.stubs(:open).returns({"keywords" => nil})

    testee = TaggedVideo.new(nil)
    tags = testee.tags

    assert_equal([], tags)

  end
end
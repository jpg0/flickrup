class ProcessingContext
  attr_reader :image
  attr_reader :upload_properties

  def initialize(image, props = {})
    @image = image
    @upload_properties = props
  end

  def with_properties(props)
    ProcessingContext.new(self.image, self.upload_properties.merge(props))
  end
end
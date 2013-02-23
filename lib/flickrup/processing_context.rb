class ProcessingContext
  attr_reader :file
  attr_reader :upload_properties

  def initialize(file, props = {})
    @file = file
    @upload_properties = props
  end

  def with_properties(props)
    ProcessingContext.new(self.file, self.upload_properties.merge(props))
  end
end
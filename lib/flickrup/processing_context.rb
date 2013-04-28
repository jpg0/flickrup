class ProcessingContext
  attr_reader :file
  attr_reader :properties

  def initialize(file, props = {})
    @file = file
    @properties = props
  end

  def with_properties(props)
    ProcessingContext.new(self.file, self.properties.merge(props))
  end
end
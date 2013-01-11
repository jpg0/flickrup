class TagSet < Array
  def initialize(tags)
    super(tags)
  end

  def to_s
    reduce("") do |result, tag|
      "#{result} \"#{tag}\""
    end
  end
end
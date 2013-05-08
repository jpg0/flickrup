class AssociatedFiles

  def initialize(config, delegate, context)
    @delegate = delegate
    @replacements = config[:associated_files].map do |x|
      extract_pair(x)
    end
  end

  def archive(ctx)
    @delegate.archive(ctx)
    @replacements.map do |x|
      other = ctx.properties[:filename].gsub(*x)
      if File.exist? other
        @delegate.archive(ctx.with_properties({:filename => other}))
      end
    end
  end

  private
  def extract_pair(pair_string)
    items = pair_string[1..-1].split(pair_string[0])
    [Regexp.new("(?#{items[2]}:#{items[0]})"), items[1]]
  end
end
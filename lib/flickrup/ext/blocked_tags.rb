class BlockedTags
  include Logging

  def initialize(config)
    @config = config
  end

  def preupload(ctx)

    blockers = @config[:blocked_tags]

    if blockers != nil
      blockers.each do |tag_name, tag_value|
        if ctx.file[tag_name] == tag_value
          #block
          logger.info("Blocking upload on #{ctx.properties[:filename]} as #{tag_name}=#{tag_value}")
          return nil
        end
      end
    end

    ctx
  end
end
module Logging

  def self.pre_init(level, logdev = STDOUT)
    @level = level
    @logdev = logdev
  end

  def logger
    Logging.logger
  end

  def self.logger
    unless @logger
      @logger = Logger.new(@logdev)
      @logger.level = @level
    end

    @logger
  end
end

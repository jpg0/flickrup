require "fileutils"

class Archiver

  def archive(ctx)
    archive_file_by_date(ctx.properties[:filename], ctx.properties[:archive_dir], ctx.properties[:date_taken], ctx.properties[:subdir])
  end

  private

  def archive_file_by_date(file, to_dir, date, subdir = nil)
    target_dir = File.join(
        to_dir,
        date.strftime("%Y"),
        date.strftime("%b")
    )

    if subdir
      target_dir = "#{target_dir}/#{subdir}"
    end

    FileUtils.mkdir_p(target_dir)
    FileUtils.mv(file, File.join(target_dir, File.basename(file)))
    logger.info("Archived #{File.basename(file)} to #{target_dir}")
  end

end
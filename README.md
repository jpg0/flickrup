Flickrup
========

Purpose
-------

Simplify image uploading to Flickr, by supporting the following features:

* Automated, async upload of images, post-tagging
* Tagging via Picasa (nice interface for tagging)
* Option for tagging on another machine (i.e. via Dropbox)
* Storage of tags within image EXIF data (not solely on Flickr)
* Storage of set data with image EXIF data (not solely on Flickr)
* Archiving of photos into /year/month/ folders
* Tag rewriting
* Video support
* Allows accelerated transfer from Dropbox

Operation
---------

Flickrup is a script which will run and monitor image files in a directory. When it has detected files have changed (and then remained unchanged for a specified period), it will process them. They will be uploaded to Flickr and archived.

Additionally, machine tags are supported, but as Picasa won't allow equals (=) in tags, any tags of the form group:key::value will be converted to group:key=value and rewritten into the file.

Special tags can also direct the script to place the images in Sets, and control the visibility of the image. See the configuration file for more detail.


Usage
-----

Run flickrup, passing location of config file.

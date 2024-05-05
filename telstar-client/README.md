# Telstar 2.0


This repository includes the Telstar Viewdata server software and associated tools.
Full details of this project are available from the
[Telstar Wiki](https://github.com/johnnewcombe/telstar-2/wiki)).


For ease of deployment, support and management, it is recommended that Telstar be run from within a Docker container, see the above wiki for details.

Telstar 2.0 binaries can be downloaded from the Releases section of this repo.

## Repository Details

Some are shared modules e.g. _telstar-library_ others are standalone utilities that perform peripheral functions.

### telstar-server

The main Telstar Viewdata server.

### telstar-client

A cross platform (Linux, MacOS, Windows) client

### telstar-library

Shared module with common functions and global constants.

### telstar-emf

T.B.A

### telstar-ftse

Utility to manage the "FTSE 100 Market Overview" from Hargreaves Lansdown.

### telstar-MacViewData

A utility that can be used to create a some rawV data from binary files created with the package MacViewdata. There are versions for Linux (arm/amd64), Mac(arm/amd64),Windows (amd64) just specify the filename e.g.

    telstar-macviewdata eng2.bin

This should return something like the following. This would be the content section of the Telstar json file. This does the whole page which includes the header. Setting the frame type to "test" will prevent the Telstar header from being displayed.

    "content": {
      "data": "\u001b\u0042\u0054\u001b\u0041\u0045\u001b\u0046\u004c\u001b\u0044\u0053\u001b\u0047\u0054\u001b\u0045\u0041\u001b\u0043\u0052\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0039\u0031\u0061\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0020\u0030\u0070\u0020\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u0030\u0031\u001b\u0057\u001b\u005e\u001b\u004f\u0073\u001b\u0053\u001b\u005a\u001b\u0056\u001b\u005e\u001b\u005f\u001b\u0058\u001b\u0044\u001b\u004d\u001b\u005d\u001b\u0043\u0045\u004e\u0047\u0049\u004e\u0045\u0045\u0052\u0049\u004e\u0047\u0020\u001b\u0052\u001b\u005c\u001b\u004c\u001b\u005e\u0073\u001b\u0055\u001b\u004e\u001b\u0051\u001b\u004f\u001b\u0054\u001b\u004f\u001b\u0047\u0030\u0032\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u0030\u0033\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u0030\u0034\u001b\u0054\u001b\u005a\u001b\u005e\u0073\u001b\u0051\u001b\u0059\u001b\u0055\u001b\u0040\u001b\u0055\u001b\u0041\u001b\u004d\u0020\u001b\u0045\u001b\u005d\u001b\u0042\u0054\u0065\u0073\u0074\u0020\u0050\u0061\u0067\u0065\u0020\u0020\u001b\u005c\u001b\u004c\u001b\u005e\u001b\u0052\u0073\u001b\u0056\u001b\u0058\u001b\u0053\u001b\u0040\u001b\u0057\u001b\u0058\u001b\u0041\u0030\u0035\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u0030\u0036\u001b\u0041\u001b\u0040\u001b\u0041\u0020\u001b\u0040\u0020\u001b\u0041\u001b\u005e\u0020\u001b\u005e\u0020\u001b\u0057\u002c\u001b\u0053\u001b\u0053\u001b\u0056\u001b\u0056\u001b\u0052\u001b\u0052\u001b\u0052\u001b\u0055\u001b\u0055\u001b\u0051\u001b\u0051\u001b\u0054\u001b\u0054\u001b\u0054\u0020\u0020\u001b\u0054\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u0030\u0037\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u0030\u0038\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u0030\u0039\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u0031\u0030\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u0031\u0031\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u0031\u0032\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u0031\u0033\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u0031\u0034\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u001b\u0041\u001b\u0040\u0031\u0035\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u007e\u007f\u0031\u0036\u0057\u0068\u0069\u0074\u0065\u001b\u0043\u0059\u0065\u006c\u006c\u006f\u0077\u001b\u0046\u0043\u0079\u0061\u006e\u001b\u0042\u0047\u0072\u0065\u0065\u006e\u001b\u0045\u004d\u0061\u0067\u0065\u006e\u0074\u0061\u001b\u0041\u0052\u0065\u0064\u001b\u0044\u0042\u006c\u0075\u0065\u001b\u0057\u001b\u005a\u0021\u0022\u0023\u001b\u0053\u0024\u0025\u0026\u0027\u001b\u0056\u0028\u0029\u002a\u002b\u001b\u0052\u002c\u002d\u002e\u002f\u001b\u0059\u0030\u0031\u0032\u0033\u001b\u0055\u0034\u0035\u0036\u0037\u001b\u0051\u0038\u0039\u003a\u003b\u001b\u0054\u003c\u003d\u003e\u003f\u0020\u0020\u0021\u0022\u0023\u0020\u0024\u0025\u0026\u0027\u0020\u0028\u0029\u002a\u002b\u0020\u002c\u002d\u002e\u002f\u0020\u0030\u0031\u0032\u0033\u0020\u0034\u0035\u0036\u0037\u0020\u0038\u0039\u003a\u003b\u0020\u003c\u003d\u003e\u003f\u0020\u0040\u0041\u0042\u0043\u0020\u0044\u0045\u0046\u0047\u0020\u0048\u0049\u004a\u004b\u0020\u004c\u004d\u004e\u004f\u0020\u0050\u0051\u0052\u0053\u0020\u0054\u0055\u0056\u0057\u0020\u0058\u0059\u005a\u005b\u0020\u005c\u005d\u005e\u005f\u0020\u0060\u0061\u0062\u0063\u0020\u0064\u0065\u0066\u0067\u0020\u0068\u0069\u006a\u006b\u0020\u006c\u006d\u006e\u006f\u0020\u0070\u0071\u0072\u0073\u0020\u0074\u0075\u0076\u0077\u0020\u0078\u0079\u007a\u007b\u0020\u007c\u007d\u007e\u007f\u001b\u0054\u0060\u0061\u0062\u0063\u001b\u0051\u0064\u0065\u0066\u0067\u001b\u0055\u0068\u0069\u006a\u006b\u001b\u0052\u006c\u006d\u006e\u006f\u001b\u005a\u0070\u0071\u0072\u0073\u001b\u0056\u0074\u0075\u0076\u0077\u001b\u0053\u0078\u0079\u007a\u007b\u001b\u0057\u007c\u007d\u007e\u007f\u001b\u0043\u001b\u0058\u0043\u006f\u006e\u0063\u0065\u0061\u006c\u001b\u0048\u0046\u006c\u0061\u0073\u0068\u001b\u0043\u002a\u001b\u004b\u001b\u004b\u0042\u006f\u0078\u001b\u0049\u0053\u0074\u0065\u0061\u0064\u0079\u001b\u0058\u0047\u006f\u006e\u0065\u001b\u004a\u001b\u004a\u003f\u001b\u0056\u005e\u007f",
      "type": "rawV"
    },

### telstar-openweather

This is called from a response frame with two arguments. The first arg is the api key, the second is the town or city of interest. The utility returns a series of pages that will be placed in temporary store within telstar and can be navigated to.

To install the open weather extension all that is required is to place the binary in a location that is accessible by the Telstar Server instance e.g. /opt/Telstar/. The templates are embedded in the binary file.

If Telstar is implemented in a Docker container, the binary file can be copied from the host machine to the container using the following command.

    docker cp telstar-openweather-linux-amd64 telstar-server:/opt/telstar/volume

Where 'telstar-openweather-linux-amd64' is the platform specific binary file and 'telstar-server' is the docker Container name. Using a Volume as in the above example, is particularly useful where multiple Telstar containers exist and there is a requirement to share the binary.

A suitable response page for use within Telstar is shown in the 'response-frame' directory. There are two frames in this folder, 290a.json is simply a map that is for presentation purposes only. Frame 290b is a response frame that is used to capture the users input. This response frame is used to invoke the telstar-openweather executable and respond to the output frames it produces.

The templates are used to format the data retrieved from openweather.org into viewdata pages. The template 'weather.json' provides a template for the main Weather result page. The template 'forecast.json' provides a template for the follow on forecast pages.

### telstar-rss

This program is designed to take the rss data file (see associated 'getdatasource.sh' script) and turn it into frames based on a frame template.

Usage:

    $ telstar-rss -i ./data/rss -t ./data/template -o ./data/frames

When processing the input file 'bbc-education.xml', the software will look for a similarly named template file e.g. 'bbc-education.json' in the template directory.

To copy the templates to the remote server, use SFTP e.g.

    $ sftp root@glasstty.com

    # navigate to local telstar-rss directory using lcd
    # navigate to remote /opt/telstar directory using cd
    $ put getdatasources.sh .
    $ put createframes.sh .
    $ put cron-test.sh .
    $ put cron-live.sh .
    $ put ./telstar-rss-2.0/telstar-rss-linux-amd64 telstar-rss
    $ put ../telstar-util/telstar-util-2.0/telstar-util-linux-amd64 telstar-util

Create data/frames, data/rss and data/templates and copy files to folders.

    $ put -R ./data .

Change permissions to 755 for getdatasources.sh, createframes.sh, chron.sh and telstar-rss.

The following is a typical template for RSS data is as follows.

    "content": {
      "data": "[R][TITLE],[W][CONTENT],[R][PUBLISHDATE],[w][l.]",
      "type": "markup"
    },

The data tag above represents the template for each article, with row definitions separated with commas.

The template must consist of at least three rows with the last row being optional.

In the above example the first row represents the rss article title e.g.

    [R][TITLE]

In this example the title is displayed in red. Tthe second row represents the rss article description e.g.

    [W][CONTENT]

the description will be displayed in white. The third row represents the rss article published date e.g.

    [R][PUBLISHDATE]

In this example the published date will be displayed in white. The fourth row represents a separator e.g.

    [w][l.]

The separator, uses standard Telstar markup, which in this case is a horizontal row of dots.
The separator must be a maximum of one row and should be the last row.

### telstar-telesoftware

A utility to create Telesoftware frames from a souce file.

### telstar-upload

A utility to upload frames to Telstar using the Telstar API.

### telstar-util

A Telstar API client.



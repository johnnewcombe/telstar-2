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

## Content Template

The following is a typical template for RSS data is as follows.

    "content": {
      "data": "[R][TITLE],[W][CONTENT],[R][PUBLISHDATE],[w][l.]",
      "type": "markup"
    },

The data tag above represents the template for each article, with row definitions separated with commas.

The template must consist of at least three rows with the last row being optional.


### Title

In the above example the first row represents the rss article title e.g.

    [R][TITLE]

In this example the title is displayed in red.

### Content

In the above example the second row represents the rss article description e.g.

    [W][CONTENT]

In this example the description will be displayed in white.

### Published Date

In the above example the third row represents the rss article published date e.g.

    [R][PUBLISHDATE]

In this example the published date will be displayed in white.

### Separator

In the above example the fourth row represents a separator e.g.

    [w][l.]

In this example the separator, in standard markup, which in this case is a horizontal row of dots. The separator must be a maximum of one row and should be the last row.

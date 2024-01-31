Version 2 of the Telstar Videotex Service components. Ignore everything you have learn't !!!

# Main Differences between v1.0 and 2.0.

## Settings

* Settings should be set using environment variables 
(see https://github.com/johnnewcombe/telstar-2/wiki/Configuration-Options). 
Environment Variables can be specified in the docker-compose.yml file if appropriate 
(https://github.com/johnnewcombe/telstar-2/wiki/Orchestrating-Telstar-with-Docker-Compose). 

## Plugins

Plugins have been completely re-written and are much simpler to use (see https://github.com/johnnewcombe/telstar-2/wiki#response-frames-and-plugins).

## Content Field

The content field no longer supports a collection of strings. The collection will need to be concatenated and delimited with CR/LF.

There are now four different data formats as defined by Telstar:

* RawV (or Raw) - this is 7bit (00-7F) videotex format with control chars between 00-1F and escaped codes for alpha and graphi attributes. Ideally rows with should have any trailing spaces removed and replaced with with a CR/LF combination. This will be rendered as is.
* RawT = this is a 24 x 40 block of 7 bits chars (960 chars) in Teletext format (range 00-7F). This format is used internally when manipulating page data. It will be converted to RawV before being rendered.
* Markup - This is Telstar markup and is converted to Raw8 before being rendered.
* EditTf - This is the edit.tf editor's url format. This is converted to Raw8 before being rendered.

The following placement tags are available, if any of these tags appear in the content or header field they will be replaced as follows:

* [SERVER] e.g. CURRER
* [GREETING] e.g. GOOD EVENING, GOOD MORNING etc. 
* [DATE] e.g. TUE 17 JULY 1979
* [TIME] e.g 12:10

Both Title and Content Fields can be used to populate content, this allows two different data formats to be used, for example, 'raw' could be used for the title and 'edit.tf' could be used for the content. The only difference is that when using the 'edit.tf' format in a title, only the top four rows are taken from the Edit.tf page.

## Markup

Content ''Markup'' has been extended as follows.

    [R] ALPHA_RED
    [G] ALPHA_GREEN
    [Y] ALPHA_YELLOW
    [B] ALPHA_BLUE
    [M] ALPHA_MAGENTA
    [C] ALPHA_CYAN
    [W] ALPHA_WHITE
    [F] FLASH
    [S] STEADY
    [N] NORMAL_HEIGHT
    [D] DOUBLE_HEIGHT
    [-] BLACK_BACKGROUND
    [n] NEW_BACKGROUND
    [r] MOSAIC_RED
    [g] MOSAIC_GREEN
    [y] MOSAIC_YELLOW
    [b] MOSAIC_BLUE
    [m] MOSAIC_MAGENTA
    [c] MOSAIC_CYAN
    [w] MOSAIC_WHITE
    [h.] SEPARATOR_GRAPHIC_DOTS_HIGH
    [m.] SEPARATOR_GRAPHIC_DOTS_MID
    [l.] SEPARATOR_GRAPHIC_DOTS_LOW
    [h-] SEPARATOR_GRAPHIC_SOLID_HIGH
    [m-] SEPARATOR_GRAPHIC_SOLID_MID
    [l-] SEPARATOR_GRAPHIC_SOLID_LOW

Alpha-graphics can now be defined in markup using a special markup syntax using double square brackets e.g.

    [b[Welcome to Telstar]]

The character between the first and second bracket, i.e. the 'b' in the above exampple represents the colour to be used. Note that Alpha-graphics take up four row, therefore, sat least 4 linefeeds will be required to position the cursor to the next row.

e.g.

    [b[Welcome to Telstar]]\n\n\n\r\n




## Frame-Merge

Add a new "merge" frame type that can be passed to telstar-util adframe command. e.g.  Frame data format can be rawT, rawV or editTf, this is converted to rawT internally and stored as RawT. Any frame data format can be merged with any other.

	{
	  "pid": {
		"page-no": 101,
		"frame-id": "a"
	  },
	  "frame-type": "merge",
	  "content": {
		"data": "http://edit.tf/#0:QIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQMMKAbNw6dyCTuyZfCBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgKJEiRIkSJEiRIkSJEiRIkSJEiRIkSJEiRIkSJEiRIkSJEiRAxQR4s6LSgzEEmdUi0otOogkzo0-lNg1JM-cCg8unNYgQIASBBj67OnXllWIMu7pl5dMOndty7uixBoy4dnTQsQIECBAgBIEHfLh2dNCDDuyINmnNl59POzKuQIECBAgQIECBAgQIECBAyQTotemggzoiCvFg1JEWkCnYemnfuw7EGHdkQIECBAgQIASBBp3dMvLdh6ad-7DsQbsvfmgw7siDvlw9NGXkuQIECBAgQM0EKrTkzotOmgkzo0-lNg1JM-cCh79vDDu8oMO7IgQIECAEgQYuvPTuy8-aDdl781wM6EjVKcVBMw9MvPogoctOPLzQIEDRBHn1otKdNizqiCTOjT6U2DUkz5wKlvw5OaxByw6dixAgBIEGHdkQUMPLZpw7cu7ouQIECBAgQIECBAgQIECBAgQIECBA2QR4NSLXg2UFOLSrSYcWmCnWKkWYsQVIsWNBsLEEOHMaIASBBh3ZEHTRlQQ9-zfz54diCHh7ZUGHJ2y7unXllXIECBAgQN0EifMkxINmmCqcsPbLsQYd2RBI37NOTD5QbsvfmuQIECBA4CHQc2TDpT50WogcMGCAScgoOnLTi69MqDpvQdNGnmgQIAaBBzy8u2nHlQd9PTQgqZdmXnvzdO-HllQYd2RBt38sq5AgQOQ0yfHQT40ZYgp2adSLNQSZ0aegTIKkWnUQUIMeLTQIECAokSJEiRIkSJEiRIkSJEiRIkSJEiRIkSJEiRIkSJEiRIkSJEFeLBqSItJBGn0osODTqIBIM6EQIEFDDnyoFTJywvoECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECBAgQIECA",
		"type": "edit.tf"
	  },
	}

The above frame will be merged with frame 101a all other json fields are ignored This allows two existing pages to be merged without having to remove the unused json elements. The addframe method of telstar-util can be used to merge the frame. It is the frame type that determines the action.


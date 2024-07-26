# Telstar 2.0


This repository includes the Telstar Viewdata server software and associated tools. 
Full details of this project are available from the 
[Telstar Wiki](https://github.com/johnnewcombe/telstar-2/wiki)).

Telstar 2.0 binaries can be downloaded from the Releases section of this repo.

## Repository Details

Some are shared modules e.g. _telstar-library_ others are standalone utilities that perform peripheral functions. Full details of this project are available from the
[Telstar Wiki](https://github.com/johnnewcombe/telstar-2/wiki)).

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

A utility that can be used to create a some rawV data from binary files created with the package MacViewdata.

### telstar-openweather

This is called from a response frame with two arguments. The first arg is the api key, the second is the town or city of interest. The utility uses tTemplates to format the data retrieved from openweather.org into viewdata frames and are returns and placed in temporary store within telstar ready for viewing. 

### telstar-rss

This program is designed to take the rss data file and turn it into frames based on a frame template.

### telstar-telesoftware

A utility to create Telesoftware frames from a souce file.

### telstar-upload

A utility to upload frames to Telstar using the Telstar API.

### telstar-util

A Telstar API client.



This is called from a response frame with two arguments.

## Introduction
The first arg is the api key, the second is the town or city of interest. The utility returns a series of pages that will be placed in temporary store within telstar and can be navigated to.

## Installation

### Binary files and Templates
To install the open weather extension all that is required is to place the binary in a location that is accessible by the Telstar Server instance e.g. /opt/Telstar/. The templates are embedded in the binary file.

If Telstar is implemented in a Docker container, the binary file can be copied from the host machine to the container using the following command.

    docker cp telstar-openweather-linux-amd64 telstar-server:/opt/telstar/volume

Where 'telstar-openweather-linux-amd64' is the platform specific binary file and 'telstar-server' is the docker Container name. Using a Volume as in the above example, is particularly useful where multiple Telstar containers exist and there is a requirement to share the binary.

### Telstar Response Page
A suitable response page for use within Telstar is shown in the 'response-frame' directory. There are two frames in this folder, 290a.json is simply a map that is for presentation purposes only. Frame 290b is a response frame that is used to capture the users input. This response frame is used to invoke the telstar-openweather executable and respond to the output frames it produces.

The templates are used to format the data retrieved from openweather.org into viewdata pages. The template 'weather.json' provides a template for the main Weather result page. The template 'forecast.json' provides a template for the follow on forecast pages.


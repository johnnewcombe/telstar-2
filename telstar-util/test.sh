#-----------------------------------------------------------------------------------------------
# Utility program for interacting with the Telstar Server API.

# Usage:
#   telstar-util [command]

# Available Commands:
#   add-frame     Adds a single frame to the currently logged in system.
#   add-frames    Adds multiple frames to the currently logged in system.
#   add-page      Adds a root frame and all follow on frames to the currently logged in system.
#   add-user      Adds/updates a user in the currently logged in system.
#   completion    Generate the autocompletion script for the specified shell
#   delete-frame  Deletes a single frame from the currently logged in system.
#   delete-user   Deletes a user from the currently logged in system.
#   get-frame     Returns a single frame from the currently logged in system.
#   get-frames    Returns multiple frames from the currently logged in system.
#   get-status    Returns the status of the specified system.
#   help          Help about any command
#   login         Logs into a system.
#   publish-frame Publishes frames from the primary database to the secondary.
#   version       Returns the version of the system.

# Flags:
#   -h, --help         help for telstar-util
#       --url string   URL of the API to be accessed. (default "u")

# Use "telstar-util [command] --help" for more information about a command.
#-----------------------------------------------------------------------------------------------
PORT=25234
# API Version wrong
# Get Frames, no json output
# Publish says 0a not found

clear

# Status and Login
echo Version:
telstar-util version

echo Login:
telstar-util login --url http://localhost:$PORT --user-id 2222222222 --password 1234 -j

echo Get Status:
telstar-util get-status --url http://localhost:$PORT -j

# Frame
echo Get Frame:
telstar-util get-frame --url http://localhost:$PORT --frame-id 0a


echo Get Frames:
telstar-util get-frames --url http://localhost:$PORT -d temp  -j

exit 0

echo Delete Frame:
telstar-util delete-frame --url http://localhost:$PORT --frame-id 0a -j

echo Add Frame:
telstar-util add-frame --url http://localhost:$PORT  -s temp/0a.json -j

echo Publish Frame:
telstar-util publish-frame --url http://localhost:$PORT --frame-id 0a -j

echo Add Frames:
telstar-util add-frames --url http://localhost:$PORT  -s temp

echo Add Page:
telstar-util add-page --url http://localhost:$PORT  --page-id temp/0a.json -j

# User
echo Get User:
telstar-util add-user --url http://localhost:$PORT -j

echo Delete User:
telstar-util delete-user --url http://localhost:$PORT -j



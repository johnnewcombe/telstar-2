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

clear

echo Version:
telstar-util version
read -p "Press return any key to resume ..."

echo Login:
telstar-util login --url http://localhost:$PORT --user-id 2222222222 --password 1234 -j
read -p "Press return any key to resume ..."

echo Get Status:
telstar-util get-status --url http://localhost:$PORT -j
read -p "Press return any key to resume ..."

echo Get Frame:
telstar-util get-frame --url http://localhost:$PORT --frame-id 0a
read -p "Press return any key to resume ..."

echo Get Frames:
telstar-util get-frames --url http://localhost:$PORT -d temp  -j
read -p "Press return any key to resume ..."

echo Delete Frame:
telstar-util delete-frame --url http://localhost:$PORT --frame-id 0a -j
read -p "Press return any key to resume ..."

echo Add Frame:
telstar-util add-frame --url http://localhost:$PORT -s temp/0a.json -j
read -p "Press return any key to resume ..."

echo Publish Frame:
telstar-util publish-frame --url http://localhost:$PORT --frame-id 0a -j
read -p "Press return any key to resume ..."

echo Add Frames:
telstar-util add-frames --url http://localhost:$PORT -s temp --include-unsafe
read -p "Press return any key to resume ..."

echo Add Page:
telstar-util add-page --url http://localhost:$PORT --page-no 800171 -s temp --include-unsafe -j
read -p "Press return any key to resume ..."

exit 0

echo Add User:
telstar-util add-user --url http://localhost:$PORT -j
read -p "Press return key to resume ..."

echo Delete User:
telstar-util delete-user --url http://localhost:$PORT -j

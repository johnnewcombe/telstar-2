## Viewdata Terminal

## Creating a Viewdata Terminal using a Raspberry Pi

Install a 32 or 64 bit version of Pi OS. Use _raspi-config_ to set oit to autologon to a desktop environment.

Once logged on set the screen blanking to off (Menu/Preferences/Raspberry Pi Configuration).

Download the Viewdata Terminal from _https://bitbucket.org/johnnewcombe/telstar-2/downloads/telstar-client.zip_ and extract the files. Navigate to the _linux-arm/linux-arm64_ directory as appropriate and extract the file _telstar-client.tar.xz_. From within the _linux-arm64/linux-arm_ directory, using a terminal, enter the following command.

    $ sudo make install

If a specific connection is required create the address file and place in an appropriate place e.g. _/usr/local/bin_ alongside the binary.

    --- # Endpoint definition for the Telstar Viewdata System

    # name: Used as a display name for connection dialogues.
    name: "Redifusion Viewdata System"
    
    # address: Address details for the service to connect to.
    address:
    host: "glasstty.com"
    port: 6522
    
    init:
    # telnet: If true, sends the IAC DO-SUPPRESS_GOAHEAD some systems
    # may need this. In the case of Telstar, this will disable the
    # 1200 baud simulation from the server and run the system at full
    # internet speed.
    telnet: false
    
    # initchar: Some systems need an initial character e.g. 0x5f (Hash)
    # to detect a connection, early versions of Telstar needed this.
    initchars: []



To autostart the Viewdata Terminal, set the system to auto logon to the desktop and create the following file.

    /home/<autologon user>/.config/autostart/com.glasstty.ViewdataTerminal.desktop
    
    [Desktop Entry]
    Type=Application
    Name=Viewdata Terminal
    Exec=/usr/local/bin/telstar-client --address=/usr/local/bin/Rediffusion-CVS.yml --full-screen --startup-delay=4
    NoDisplay=false
    NotShownIn=GNOME;KDE;XFCE;

Create
Reboot.

## Video Issues

If there is no video output, consided changing the config.txt file. This can be done by mounting the SD card in a media reader etc and editing the file with a text editor.

    hdmi_force_hotplug=1
    config_hdmi_boost=4 (supports up to 9)

If the display is a computer monitor set...

    hdmi_group=1 

If the display is an older TV, try 

    hdmi_group=2.

Do not set _hdmi_safe=1_ as that overrides many of the previous options.

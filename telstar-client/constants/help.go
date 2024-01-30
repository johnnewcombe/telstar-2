package constants

const HelpMsg = `
## Introduction

This Viewdata Terminal, available for Mac, Linux and Windows allows access to viewdata systems through either a direct TCP/IP connection or via a Modem connected to a local serial port.

---

Versions are available for the following platforms/architectures.


- Linux (amd64)
- Linux (arm64)
- Linux (arm64)
- MacOS (arm64)
- Windows (386)
- Windows (amd64)

___

## Connecting to Services

When run for the first time the software will connect to the Telstar Viewdata System but from then on will attempt to connect to the last system used.

___

## Connection Files

Each supported service is defined in a .yml text file e.g. nxtel.yml. Several of these files are supplied with the software. The filename can be specified on the command line as detailed below or opened using the Open toolbar button. New TCP or Serial connection files can be created using the appropriate toolbar buttons.

___

## Command Line Parameters

The binary file can be launched from the command line (to use the command line when using MacOs, it is necessary to copy the files from the app bundle (./Contents/MacOS and ./Contents/Resources) to a convenient folder).

For example to connect to NxTel on startup:

___

		telstar-client -address=nxtel.yml

---

The full list of command line arguments is shown below:

- -address -  Endpoint definition file e.g. Telstar.yml.
- -text-size - Text size, can be used to counter display driver issues. Default = 23
- -full-screen - Full screen mode.
- -no-toolbar - Hides the toolbar and status line.
- -startup-delay - Delays startup of the application.

___

## Acknowledgements

* The font for the Telstar Client (MODE7GX2.TTF) was based on ModeSeven by Andrew Bulhak, updated by galax.xyz.
* Portions of this software are copyright (c) 2018 The FreeType Project (www.freetype.org).  All rights reserved.
* The application icon is derived from work by Dan Farrimond. This original work and the icon is licensed under the Creative Commons Attribution-Share Alike 3.0 Unported license (https://creativecommons.org/licenses/by-sa/3.0/deed.en).
* Other icons were supplied by https://www.iconsdb.com/white-icons/cloud-3-icon.html provided as CC0 1.0 Universal (CC0 1.0) Public Domain Dedication.
___

## Licence and Copyright

BSD 3-Clause License

Copyright (C) 2018 2022 John Newcombe (https://glasstty.com)
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:
* Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
* Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
* Neither the name of Fyne.io nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.


___


#
#
#
#
#
#
`

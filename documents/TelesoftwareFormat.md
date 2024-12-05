# Telesoftware


# The Telesoftware Protocol

The protocol described here is based on the _Format Recommendations for Prestel Telesoftware (CET, July 1986)._

## Introduction

In September 1980, the Council for Educational Technology (CET) invited computer manufacturers, software agencies and representatives of Prestel to discuss the formulation of a set of recommendations for the format ot Telesoftware on Prestel. As a result of these discussions, the Council fer Educational Technology published a document in January 1981, entitled __Prestel Telesoftware Format Recommendations__.

The recommendations described a method of loading a single program source file on to a viewdata system, such as Prestel, the method was particularly suited to the transfer of a program written in BASIC from Prestel to a microcomputer. At that time the recommendation was restricted to sending characters from a 7 bit character set.

At the end of 1981, as a result of work done with the Commodore PET, an extension to the recommendations to allow the transmission of files containing any character from an 8-bit character set was developed and subsequently published in February 1982.

Since that date, the experience gained from the increasing use of the format recommendations by GET during the CET Telesoftware Project, and also by other information providers on Prestel who have used the format recommendations, has resulted in the inclusion in the recommendations of some new facilities and also some clarification of the existing format.

This document details the protocol, along with these enhancements described above, as published in 1986 and is the basis for the Telesoftware protocol in use by Telstar.

__Please note that all characters described numerically in this document are shown in hexadecimal unless otherwise stated.__

## Files

The basic entity in telesoftware is the file. This is simply a quantity of data of a particular length and consists of one or more characters, where a character is an 8 bit byte having a value between 0 and FF. The file is usually given a name, called its filename.

Files are made up of blocks of data (see below).

## Pages and Frames

The basic entity on the Telstar system is the _page_, A page can consist of up to 26 _frames_ and each page is identified by a unique number (typically shown on the top row of the page). Each frame or page on Prestel consist 24 rows, each of 40 characters. The top row and the bottom row are reserved for system messages, for example, the name ef the information provider, the number of the page and the price of the page.

The remaining 22 rows are available for use in storing part of a Telesoftware file. Each of the 40 character positions on each row can contain either an alphanumeric character or one of 64 graphic symbols. It is not possible, to store characters whose ASCII value is less than 32 or greater than 127 (decimal). The Telesoftware protocol uses a combination of _escape_ sequences and character _shifts_ to mitigate this limitation.

Telstar sends a single block of data per frame. The Telesoftware protocol as published in 1986 included a mechanism to support multiple blocks on a single frame. However, that was strongly discouraged. As a result, Telstar does not support multiple blocks on a frame and therefore is not described in this document. However, full details of the protocol for multiple blocks per frame can be found in the document _Format Recommendations For Prestel Telesoftware, CET. 1986_.

## Header Frame

The first frame of a Telesoftware file is referred to as the _header frame_. The header frame contains a data block that typically contains the filename and a count of the number of frames which the file occupies.

The second and subsequent frames contain data blocks containing the file itself.

## Handling Transmission Errors

Each block of data ends with a three digit checksum. In addition, each block has a short section at the beginning which distinguishes it from the preceding and following blocks. This allows receiving software to detect missing blocks and blocks sent out of order. Frames can be re-requested by sending the _*00_ command

Prestel transmits and expects to receive characters as 7 data bits along with an even parity bit. Ideally the downloading software would check the parity bit on each character received and, if a parity error is detected, should re-request the as described above.

After requesting Prestel to send a page, the receiving software should assume that, after a Short period of inactivity, a complete frame has been received.


## Zero Page Routing

When a file will not fit on one Prestel page (consisting of 26 frames), a second, continuation, page is used. This page is reached by using the '0' route from the 'z' frame of the page. If further pages are required, then the '0' route from the frame is used to reach further continuation frames. For example a continuation frame for the frame 123z would be 1230a, a continuation frame for 1230z would be 12300a and so on.

The final frame of a telesoftware file contains an escape sequence indicating that the end of the file has been reached.

## Escape Sequences

As previously mentioned, the set of characters which can be stored within a frame is restricted to characters between 32 and 127 (decimal). The Telesoftware encoding software makes use of escape characters as defined below to allow values between 0 and 255 (decimal) to be sent.

The Telesoftware escape character has a value of 7C (hexadecimal) and this is known as the _sequence introducer_. This is displayed as a double vertical bar '||' on a viewdata terminal, however, for simplicity and to remain consistent with previous CET documentation, it is shown as a single vertical bar '|' in this document.

Whenever the character 7C occurs on a Prestel frame containing Telesoftware, it is followed by another alphabetic or numeric character which determines its meaning. For example, the character sequence 7C,41 Indicates the start of a Telesoftware block. Note that this will be displayed on the screen as;

    |A

The full list of escape sequences are shown below including a complete example of a short program.


### 7C 41 (Block Start)

Marks the start of a Telesoftware block and is followed by Block G.

### 7C 47 (|G)

This sequence is followed by the frame letter of the frame on which this telesoftware block occurs. The frame letter is a lower case alphabetic character. Where there is more than one telesoftware block on a frame, the frame letter is followed by two numeric characters (see _Sequencing of Telesoftware Blocks_ and the section on _Multiple Blocks per Frame_ below)

### 7C 5A (|Z, Block End)

Marks the end of a telesoftware block. Il is followed by three numeric characters which are ihe checksum for the block (see section on Checksum Calculation below). For example if the checksum value was 324, the sequence would be (33,32,34).

### 7C 49 (|I)

This sequence acts as a terminator for those escape sequences which were added to the original format recommendations. It is a very powerful feature and is described fully in the section _Escape Sequence Terminator_ below.

### 7C 46 (|F, End of File)

This escape sequence is inserted after the last character in the file and signifies the _end of file_ (EOF). How a downloader handles the end of file condition is up to the designer of the downloader. The following observations may be of assistance:

On a CP/M based system, the handing of and of file depends on whether the file is a text file or not. A text file is written to disc with the final 129 byte block padded to the end of the block with Ctrl+Z (1A) character.

If the file is not a text file, then the final block of the file is usually padded with nulls (00) until the last 128 byte block is filled. Other operating systems adopt different conventions.

### 7C 4C (|L, End of Line)

This sequence can be used by an Information Provider to signify the end of a line in a source file. It is usually used to replace the characters CR (carriage return) and LF (line feed), Note, however, that different microcomputers use different characters for _end of line_; some use just CR or LF, others use the more usual CR LF combination, others may use LF CR, and it is possible that some microcomputers may have a completely different system for signifying _end of line_. When you encounter this escape sequence, you should write to the output file whatever character or sequence of characters your particular microcomputer expects at the end of such source file. This escape sequence is also used on the header frame to separate the filename from the length field (see _The Header Frame_ below).

### 7C 7D (|3/4)

This sequence is used to represent the _three quarters_ character which is used for another purpose, Whenever this escape sequence in encountered, you should simply write the character to the output file.

--

### 7C 45 (|E)

This is the way in which an information Provider can insert the _double bar_ character (7C) into the output file. As with the sequence above, if you encounter this escape sequence then you should simply write the character to the output file.

--

### 7C 54 (|T)

This escape sequence is not in common usage, but it has been included to allow an Information Provider to mark the start of the header section of the file, i.e. the Filename and length (see below).

--

### 7C 44 (|D)

This escape sequence is also net commonly used - in contrast to |T (described above), It is used to allow an Information Prövider to mark the start of the data section of the file, as opposed to the header section which gives the filename and number of frames that the file occupies (see above).

--
### 7C 30 (|0) to 7C 35 (|5)

This group of escape sequences are provided to allow characters outside the range 32-127 to be Included In the Sia, See the section on _Character Code Extension_ below.

--

All other escape sequences are at present undefined. Further, all escape sequences not listed above are reserved for future use, except for of all the lower case letters. These are intended for Information Providers and others who are experimenting with extensions to this set of recommendations. During the development phase of a new feature, lower case letters should be used for a new escape sequence, When the feature has been tried and tested and is generally accepted, han oee of the reserved characters will be allocated to that function and the feature wis be incorporated into this recommendation.


### Escape Sequence Terminator

When the format recommendations were originally framed in 1980, a total of six escape sequences were defined. In 1981/82 another seven escape sequences were defined. 7C,49 (one of the new escape sequences defined) was intended to make it easier for other, as yet unforeseen, sequences to be added easily, and with as little disruption as possible to users of downloads written to comply with earlier publications of the recommendations.

Supposing we wanted to add a new escape sequence to the current list which would have the effect of reserving a certain number of blocks of disc space for our use. The escape sequence would be followed by one or more characters which would specify the amount of space to reserve.

The main problem we would encounter in introducing this facility is that any user with Telesoftware programs written before this escape sequence was announced would no longer be able to download the particular file containing it. The escape sequence which was intended to help the telesoftware program to reserve sufficient space for the file would instead be treated as part of the
file and would mean that, when users downloaded that file, they would get several seemingly spurious characters(the escape sequence and the following parameter) in the middle of their copy of the file.

To get around this problem, the following rules were devised. Escape sequences are divided into three categories:-

#### Category 1

The set of escape sequences defined in the format recommendations published in February 1902. These are

    |A, |Z, |L, |I, |0, |1, |2, |3, |4, |5, |F, |E, |3/4

i.e.

    7C 41, 7C 5A, 7C 4C, 7C 49, 7C 30, 7C 31, 7C 32, 7C 33, 7C 34, 7C 35, 7C 46, 7C 45 and 7C 7C

#### Category 2

Those escape sequences to which generally accepted meanings have subsequently been attached. At present these are

    |G, |T, and |D

i.e.

    7C 47, 7C 54 and 7C 44

#### Category 3

Those remaining escape sequences which have not yet been allocated a specific meaning.


### Processing Escape Sequences

When a telesoftware downloading program encounters an escape sequence from categories 2 and 3, it stops taking characters from the frame, converting them and storing them in the output file. Instead, if it is an escape sequence, it recognises, then it processes the escape sequence, otherwise it simply ignores all the following characters until it reaches the escape sequence 7C,4C (EOL). __How does this work for binary files with no EOL? It is therefore important not to use unsupported escape sequences which may not include EOL sequence (7C, 4C)__)

The escape sequence 7C,49 instructs the downloading program to continue taking characters from the frame, convert them and store them in the output file. One further point for Information Providers concerting escape sequences is very important. If for some reason an escape sequence will not fit completely within the current frame (leaving room for the final |Zxxx (7C,5A,xxx), then the complete escape sequence and any following parameters up to the occurrence of 7C,49 __must__ be moved to the next frame. In no circumstances can an escape sequence be split across two frames.


##Block Format

### Header Block

    |A|Gc|ITEST|L001|Z122

i.e.

    7C,41,7C,47,<frame id>,7C,47,<filename>,7C,4C,30,30,31,7C,5A,31,32,32

### Data Block

    |A|Gd|I|L1@EF|03/4|4...

## Other Special Characters

Apart the escape sequences described above, here is one other character which has a special meaning in the telesoftware file. This is the _three quarters_ character (7D). Whenever this character is encountered in a Telesoftware file, and it is NOT part of the escape sequence 7C,7D, it should be converted to the space character 20. The reason for this is a little complicated but is
related to the way in which Prestel stores frames of information. If the last, Say, six characters on a line of a Prestel frame are al spaces, then Prestel can reduce the number of characters transmitted to the microcomputer (or terminal) by stripping off the trailing spaces and transmitting a carriage return/line feed (CRLF) sequence instead. Since we do not want this feature for telesoftware, the way to ensure that Prestel does not strip trailing spaces is to convert the last character on each line when the file is loaded. If it is a space, then it is converted to the _three quarters_ character.

Over the life of the format recommendations, the instructions for this situation have varied. Originally EVERY occurrence of the space character was replaced by the character _three quarters_. Then the decision was made to allow the conversion of spaces to be purely optional. The current position for uploading files is that conversion of a final space on a line is mandatory: conversion of
spaces elsewhere within a line is optional.

The above only affects those people loading files on to Prestel. For designers of software for downloading programs, the rule is always that if you encounter the _three quarters_ on its own, then convert it to a space.

## Checksum Calculation

For each Telesoftware block, a checksum is calculated on all the characters between the escape sequences |A and |Z . The checksum ls calculated as follows:

When the escape sequence |A is encountered, the checksum is set to zero. As each character is received, its value is exclusive-OR'ed with the current value of the checksum. At the end of the block, when the escape sequence |Z is encountered, the checksum calculated above should correspond to the value of the three digits following the 7C,5A sequence. The eighth, party bit of each character __must__ be set to zero before this calculation and the resulting value will always be between 000 and 127.

For example:-

    |AThis is a checksum test|Zxxx

The value left in the checksum after all the characters between |A and |Z have been XOR'ed together is xxx in decimal. The three digits following |Z in the example above should be x,x and x. If a terminal finds that its calculated version of the checksum does not equal the value transmitted, then it can assume that there has been corruption of the frame and can ask Prestel to resend the page. Note that the checksum following the |Z is always three characters with leading zeroes it these are needed. For example if the checksum value was 24, the sequence would be (30,32,34).

The viewdata command which causes Prestel lo re-transmit a frame is the sequence *00, it is strongly recommended that designers of downloading software should keep a check on the number of times that a particular frame is re-requested and, if some particular number of retries is exceeded, should abandon the attempt to fetch the program and inform the user ef the problem. It is always possible
that an Information Provider has inadvertently put up a Telesoftware file wih
an incorrect checksum on the frame.

It is also recommended that you check for parity errors and, if a parity error occurs on a frame, you should re-request the frame in the same manner as if a checksum error had occurred. This helps to improve the error checking.

## Sequencing of Telesoftware Blocks

In the normal situation where the Information Provider is only putting one Telesoftware block on each frame, then the following sequence will be at the start of each frame, immediately after the |A sequence indicating start of block.

    |Gc|I 

Where c above is the frame letter of the current frame.

Since Prestel frames will be consecutively labeled 'a', 'b', 'c', and so on up
to 'z' and, it a continuation page is used, wil start again at 'a', it is
possible for the downloading software to detect when a frame is received out of order, or when a request for the next frame has been lost because of noise on the telephone line. A limited amount of corrective action can be taken and the user can be informed of what in occurring

## The Header Frame

As mentioned earlier, the first frame of a telesoftware file will be the header frame containing the name of the file and the number of Prestel frames (excluding the header frame itself) that the file occupies. For example:

    |AGa|ISORT.BAS|L003|Zxxx

It is recommended that the header frame should always be on the 'a' frame of a page, as in the example. However, several information Providers do also put the header frame on the 'd' frame of the frat page. Note that the _end of line_ escape sequence 7C,4C is used in the header frame to separate the filename from the length information. As with the checksum, the length information is always these decimal digits with leading zeros If these are needed. If for some reason the Information Provider is unable to provide the number of frames which the file occupies, then the number 999 should be used instead. This simply lets the downloading software that the number of frames is not known.

## Character Code Extension

In order to represent characters in the range 0 - 255 by combinations of characters from the range 32 - 127, it is necessary to use some sort of shift technique, This is similar te the use of the SHIFT key on a typewriter to give you upper case (capital) letters instead of lower case letters. The same keys are used with the shift key lo produce twice as many different symbols as there are keys on the typewriter.

As the beginning of a telesoftware file (i.e. the header frame) there is no shift
in use. The characters in the range 32 - 127 (subject only to the rules for _three quarters_ character and the escape sequences given above) represent themselves and are written to the output file exactly as they appear on the frame. In fact, at the start of every new telesoftware file the shift offset, as it is known, is always reset to zero.

The escape sequences which change this situation are the six escape sequences

    |0, |1, |2, |3, |4, |5

    i.e.

    7C 30, 7C 31, 7C 32, 7C 33, 7C 34, 7C 35

Each of these six escape sequences causes the downloader
to alter the way in which it deals with the characters until another one of
these escape sequences is encountered.

    Escape Sequence        Shift Offset set to
        |0  (7C 30)              0
        |1  (7C 31)             -64
        |2  (7C 32)             +64
        |3  (7C 33)             +96
        |4  (7C 34)             +128
        |5  (7C 35)             +160

When a character is read from a Prestel frame by the downloading software, the value of that character should be added to the current value of the shift offset. The resulting character should then be written to the output file.

Note that the shifts are what are known as locking shifts, i.e. each change to the shift offset variable is permanent until another control sequence which changes the shift offset ls encountered.

One further point for Information Providers; there is only one means to represent any particular character (0 - 255) on a Telesoftware file. The following table shows which shift offset must be selected in order to represent a particular character in the file:

    Character Value     Shift Offset which 
                        must be selected

        0-31                |1 (7C 31)
        32-127              |0 (7C 30)
        128-159             |2 (7C 32)
        160-191             |3 (7C 33)
        192-223             |4 (7C 34)
        224-255             |5 (7C 35)


## Multiple Blocks per Frame


In some fairly specialised applications, it may be necessary to put several
small Telesoftware blocks on one Prestel frame. The use of this facility is not encouraged unless it is absolutely necessary because it means that a downloader must be even more complex and is likely That only a few implementations of
the telesoftware downloader will have this capability.

In order to try and prevent the downloader from missing some of the blocks on a frame, the escape sequence used for identifying Telesoftware blocks (7C 47) is extended in this situation, Now, in addition to the frame letter which follows the (7C 47), there are also two extra characters.

The first is a number in the range 0 - 9 which identifies which block this is on the frame. The first block will be numbered block 0, the second block will be numbered block 1, and so on with a maximum of ten blocks on the frame.

The second number is the number of the last block on the frame, using the same
numbering as described above. For example, the third block on a frame which has 8 blocks in all would have at the start of the block

    |a|Gc27I___data___|Zxxx


## Characteristics of the Prestel System on which these Recommendations Rely

These recommendations assume that the Prestel database consists of a number of pages, each page being made up of a maximum ot 26 frames. To step from one frame to another, the viewdata command '#' is used. in order to reach a continuation page from the 'z' frame of the current page, the viewdata routing command 'O' is then sent. If any error is encountered in the reception of any frame, then the viewdata command '*00' will result in the frame being retransmitted.

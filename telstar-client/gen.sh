#!/bin/sh

DIR=`dirname "$0"`
FILE=constants/bundled.go
BIN=`go env GOPATH`/bin

echo $DIR
cd $DIR
rm $FILE

$BIN/fyne bundle -package constants -name AppIcon ./Icons/icon.png > $FILE
$BIN/fyne bundle -package constants -name SerialIcon -append ./Icons/phone-70-512.png >> $FILE
$BIN/fyne bundle -package constants -name CloudIcon -append ./Icons/cloud-3-512.png >> $FILE
$BIN/fyne bundle -package constants -name CloseIcon -append ./Icons/x-mark-4-512.png >> $FILE
$BIN/fyne bundle -package constants -name MODE7GX2TTF -append ./Fonts/MODE7GX2.TTF >> $FILE

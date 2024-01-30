package main

import (
	"bitbucket.org/johnnewcombe/telstar-library/file"
	"bitbucket.org/johnnewcombe/telstar-library/logger"
	"bitbucket.org/johnnewcombe/telstar-library/telesoftware"
	"bitbucket.org/johnnewcombe/telstar-library/types"
	"bitbucket.org/johnnewcombe/telstar-library/utils"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
)

const (
	MAX_CHARS_PER_BLOCK = 859 - 12 // allow for BLOCK_END (2), checksum (3) and the start sequence (7)
)

func main() {

	var (
		sourceData      []byte
		encodedData     []byte
		pageId          string
		pageNumber      int
		frameId         rune
		name            string
		flags           *flag.FlagSet
		err             error
		inputFile       string
		outputDirectory string
		deleteFiles     bool
		addRedirect     bool
		blocks          []telesoftware.Block
		frames          []types.Frame
		redirectFrame   types.Frame
		jsonFrame       []byte
	)

	flags = flag.NewFlagSet("telstar-rss", flag.ExitOnError)
	flags.StringVar(&name, "name", "", "Display name of the program or file.")
	flags.StringVar(&pageId, "pid", "", "Page ID of the start of the encoded frames. This would typically be a 'c' frame.")
	flags.StringVar(&inputFile, "i", "", "Path to the file to be encoded.")
	flags.StringVar(&outputDirectory, "o", "", "Path to the directory where the generated frames will be stored.")
	flags.BoolVar(&deleteFiles, "d", false, "Deletes the contents of the output directory before starting.")
	flags.BoolVar(&addRedirect, "add-redirect", false, "Creates a redirect 'a' frame to redirect to the first of the encoded frames.")

	if err = flags.Parse(os.Args[1:]); err != nil {
		logger.LogError.Printf("%v%", err)
		os.Exit(1)
	}

	if deleteFiles {
		// delete the output directory of json files
		file.DeleteFiles(path.Join(outputDirectory, "*.json"))
	}

	// get the data files
	if sourceData, err = file.ReadFile(inputFile); err != nil {
		log.Fatal(err)
	}

	// call encode
	if encodedData, err = telesoftware.Encode(sourceData); err != nil {
		log.Fatal(err)
	}

	// get page number and frame id
	if pageNumber, frameId, err = utils.ConvertPageIdToPID2(pageId); err != nil {
		log.Fatal(err)
	}

	// enblock
	if blocks, err = telesoftware.Enblock(encodedData, pageNumber, frameId, name, MAX_CHARS_PER_BLOCK); err != nil {
		log.Fatal(err)
	}

	if addRedirect {
		redirectFrame = types.Frame{
			PID: types.Pid{
				PageNumber: pageNumber, FrameId: "a"},
			Visible: true,
			Redirect: types.Pid{
				PageNumber: pageNumber, FrameId: string(frameId)}}

		if jsonFrame, err = redirectFrame.Dump(); err != nil {
			log.Fatal(err)
		}
		file.WriteFile(path.Join(outputDirectory, fmt.Sprintf("%da.json", pageNumber)), jsonFrame)
	}

	// create frames
	if frames, err = createFrames(blocks, pageNumber, frameId, name); err != nil {
		log.Fatal(err)

	}

	for _, frame := range frames {
		// save frames to json
		if jsonFrame, err = frame.Dump(); err != nil {
			log.Fatal(err)
		}
		if pageId, err = utils.ConvertPidToPageId(frame.PID.PageNumber, frame.PID.FrameId); err != nil {
			log.Fatal(err)
		}

		file.WriteFile(path.Join(outputDirectory, pageId+".json"), jsonFrame)
	}
}

// createFrames will create telesoftware pages starting with the header at the specified page id.
func createFrames(blocks []telesoftware.Block, pageNumber int, frameId rune, name string) ([]types.Frame, error) {

	var (
		frames []types.Frame
		//err    error
	)

	//if frameId != "a" {
	//	// add the "a" frame as a redirect to the actual frame
	//redirectFrame := frame.CreateRedirect(pageNumber, pageId)
	//	if err = dal.InsertOrReplaceFrame(settings.Database.Connection, redirectFrame, primaryDb); err != nil {
	//return err
	//}

	//}
	//if encodedData, err = telesoftware.Encode(data); err != nil {
	//	return err
	//}
	//if blocks, err = telesoftware.Enblock(encodedData, pageId, name, MAX_CHARS_PER_BLOCK); err != nil {
	//	return err
	//}

	//primaryDb = strings.ToLower(settings.Database.Collection) == "primary"

	for i := 0; i < len(blocks); i++ {

		block := blocks[i]

		var frame types.Frame
		//create a frame
		frame.PID.PageNumber = block.PageNumber
		frame.PID.FrameId = string(block.FrameId)
		frame.Content.Type = "rawV"
		frame.Content.Data = string(block.Data)
		frame.RoutingTable = utils.CreateDefaultRoutingTable(pageNumber)
		frame.Visible = true

		frames = append(frames, frame)

		//if err = dal.InsertOrReplaceFrame(settings.Database.Connection, frame, primaryDb); err != nil {
		//	return err
		//}
	}

	// all new pages created so purge any old frame with a higher id than the current frames
	// get next page id
	//if pageId, err = utils.ConvertPidToPageId(pageNumber, string(frameId)); err != nil {
	//	return frames, err
	//}

	//if pageNumber, frameId, err = utils.GetFollowOnPID(pageNumber, frameId); err != nil {
	//	return frames, err
	//}

	// purge any old frame with a higher id than the current frames
	// ignore ther deleted count
	//if _, err := dal.PurgeFrames(settings.Database.Connection, pageNumber, frameId, primaryDb); err != nil {
	//	return err
	//}

	return frames, nil
}

package main

import (
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/johnnewcombe/telstar-ftse/FtseItem"
	"github.com/johnnewcombe/telstar-library/file"
	"github.com/johnnewcombe/telstar-library/logger"
	"github.com/johnnewcombe/telstar-library/text"
	"github.com/johnnewcombe/telstar-library/types"
	"github.com/johnnewcombe/telstar-library/utils"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

const (
	// TODO make this an arg ???
	//SELECTOR = ".full-width.ftse-index-table-table"
	SELECTOR = ".stockTable"
)

type templateStruct struct {
	line   string
	length int
}

func main() {

	var (
		items             []ftseItem.FtseItem
		files             []fs.FileInfo
		err               error
		flags             *flag.FlagSet
		inputDirectory    string
		outputDirectory   string
		templateDirectory string
		deleteFiles       bool
	)

	flags = flag.NewFlagSet("telstar-ftse", flag.ExitOnError)
	flags.StringVar(&inputDirectory, "i", "./data/html", "Path to the directory containing .html files.")
	flags.StringVar(&templateDirectory, "t", "./data/templates", "Path to the directory containing template type .json files.")
	flags.StringVar(&outputDirectory, "o", "./data/frames", "Path to the directory containing where the generated frames will be stored.")
	flags.BoolVar(&deleteFiles, "d", false, "Deletes the contents of the output directory before starting.")

	if err = flags.Parse(os.Args[1:]); err != nil {
		logger.LogError.Printf("%v%", err)
		os.Exit(1)
	}

	if deleteFiles {
		// delete the output directory of json files
		file.DeleteFiles(path.Join(outputDirectory, "*.json"))
	}

	// get the data files
	if files, err = getDataFiles(inputDirectory); err != nil {
		logger.LogError.Print(err)
		return
	}
	// we need to iterate through all web pages
	for _, file := range files {

		logger.LogInfo.Printf("Processing: %s", file.Name())
		// get articles
		if !file.IsDir() {
			if items, err = parseHtml(inputDirectory, file.Name()); err != nil {
				logger.LogError.Print(err)
				continue
			}
			//fmt.Printf("%v", items)
		}

		// get the path to the template
		templatePath := getTemplatePath(file.Name(), templateDirectory)

		// create the frames for this data file
		if err = createFrames(items, templatePath, outputDirectory); err != nil {
			logger.LogError.Print(err)
			continue
		}
	}
}

func createFrames(items []ftseItem.FtseItem, templatePath string, outputDirectory string) error {

	const (
		MAXLINES        = 18
		MAXCOLS         = 39
		CODETAG         = "[CODE]"
		PRICETAG        = "[PRICE]"
		CHANGETAG       = "[CHANGE]"
		CHANGERECENTTAG = "[CAHNEPERCENT]"
	)

	var (
		err           error
		templateJson  []byte
		templateFrame types.Frame
		currentFrame  types.Frame
		linesUsed     int
		//markupText    string
	)

	// load the appropriate templateJson
	if templateJson, err = ioutil.ReadFile(templatePath); err != nil {
		return err
	}

	if err = templateFrame.Load(templateJson); err != nil {
		return err
	}

	// see how many rows on the title
	_, titleRows := utils.ParseDataType(templateFrame.Title.Type)
	if titleRows == 0 {
		titleRows = 4
	}
	maxLines := MAXLINES - titleRows

	// there are a total of 24 lines, line 0 is the header, then 4 for the title
	// text can start on line 5 to line 23. Line 24 is the sys/nav message.
	// making a total of 18 lines available

	// this gives us the details for each part of the templateJson
	// this loads the templateJson definition see the content filed of the json templates
	// each definition is comma separated, template now containes content AND
	// definitions which will be expanded e.g.
	//  "[l.],[R][TITLE],[W][CONTENT],[R][PUBLISHDATE]" represents four elements

	if currentFrame, err = getFrame(templateJson, templateFrame.PID); err != nil {
		return err
	}
	var frameContent strings.Builder

	save := func() error {
		// save old frame get new frame
		currentFrame.Content.Data = frameContent.String()
		logger.LogInfo.Printf("Saving: %d%s", currentFrame.PID.PageNumber, currentFrame.PID.FrameId)
		if currentFrame, err = saveFrame(currentFrame, outputDirectory+"/"+currentFrame.PID.String()+".json"); err != nil {
			return err
		}
		frameContent.Reset()
		linesUsed = 0
		return nil
	}

	for _, item := range items {

		frameContent.WriteString(formatContent(item, templateFrame.Content.Data))

		linesUsed++
		// all of the date must fit
		if linesUsed > maxLines {
			if err := save(); err != nil {
				return err
			}
		}

	}

	if frameContent.Len() > 0 {
		if err := save(); err != nil {
			return err
		}
	}

	//print(frameContent.String())
	//fmt.Printf("%v", maxLines)
	//fmt.Printf("%v", currentFrame)
	return nil
}

func getHtmlTable(html []byte, selector string) []string {

	var (
		result []string
	)

	// Create a goquery document from the HTTP response
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(html)))
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	doc.Find(selector).Each(func(i int, s *goquery.Selection) {

		// For each item found, get the title
		s.Find("tr").Each(func(i int, s *goquery.Selection) {
			var sb strings.Builder
			s.Find("td").Each(func(i int, s *goquery.Selection) {
				// swap commas for tags
				val := s.Text()
				if len(val) > 0 {
					sb.WriteString(fmt.Sprintf("%s|", strings.Trim(val, " ")))
				}
			})
			if s.Nodes[0].FirstChild.Data != "th" {
				result = append(result, strings.Trim(sb.String(), "|"))
			}
		})
	})
	return result
}

func exitWithHelp() {

	displayCopyright()

	fmt.Println("  Usage: get-html-table <url>")
	fmt.Println("")
	fmt.Println("  e.g.   get-html-table https://someurl.com")

	os.Exit(0)
}
func displayCopyright() {
	fmt.Println("")
	fmt.Println("Get HTML Table Utility v1.0 (c) 2021 John Newcombe")
	fmt.Println("")
}

func getDataFiles(dataDirectory string) ([]fs.FileInfo, error) {

	var (
		err   error
		files []fs.FileInfo
	)

	files, err = ioutil.ReadDir(dataDirectory)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func parseHtml(dataDirectory string, filename string) ([]ftseItem.FtseItem, error) {

	var (
		err      error
		data     []string
		items    []ftseItem.FtseItem
		htmlData []byte
	)

	if htmlData, err = ioutil.ReadFile(fmt.Sprintf("%s/%s", dataDirectory, filename)); err != nil {
		return nil, fmt.Errorf("%v\n", err)
	}

	data = getHtmlTable(htmlData, SELECTOR)
	for _, f := range data {

		fmt.Printf("%v\r\n", f)
		var item = ftseItem.FtseItem{}
		if len(f) > 0 {
			item.Load(f)
			items = append(items, item)
		}
	}

	return items, nil
}

func getTemplatePath(dataFileName string, templateDirectory string) string {

	templateName := dataFileName[:len(dataFileName)-len(path.Ext(dataFileName))]
	return fmt.Sprintf("%s/%s.json", templateDirectory, templateName)

}

func getFrame(templateJson []byte, pid types.Pid) (types.Frame, error) {

	var (
		f   types.Frame
		err error
	)

	if err = f.Load(templateJson); err != nil {
		return f, err
	}
	f.Content.Data = ""
	f.Content.Type = "markup"
	f.PID = pid

	return f, nil
}

// saveFrame Saves the specified frame and returns an empty follow on frame.
func saveFrame(frame types.Frame, outputDirectory string) (types.Frame, error) {

	frameBytes, err := frame.Dump()
	if err = ioutil.WriteFile(outputDirectory, frameBytes, 0644); err != nil {
		return frame, err
	}

	// get new pid
	pn, fid, err := utils.GetFollowOnPID(frame.PID.PageNumber, rune(frame.PID.FrameId[0]))
	if err != nil {
		return frame, err
	}

	// return an empty follow on frame
	frame.PID.PageNumber = pn
	frame.PID.FrameId = string(fid)
	frame.Content.Data = ""
	frame.Content.Type = "markup"

	return frame, nil
}

// formatContent Formats the content adding red/green colours for +/- etc. Takes the item and
// template line from the template content field
func formatContent(item ftseItem.FtseItem, templateLine string) string {

	var colour string

	if item.Change < 0 {
		colour = "[R]"
	} else {
		colour = "[G]"
	}

	content := strings.Replace(templateLine, "[CODE]", text.PadTextRight(item.Name, 20), -1)
	content = strings.Replace(content, "[PRICE]", text.PadTextLeft(fmt.Sprintf("%.2f", item.Price), 8), -1)
	content = strings.Replace(content, "[CHANGE]", text.PadTextLeft(fmt.Sprintf("%s%.2f", colour, item.Change), 8), -1)
	//	content = strings.Replace(content, "[CHANGEPERCENT]", fmt.Sprintf("%2.f", item.ChangePerCent), -1)

	return content + "\r\n"
}

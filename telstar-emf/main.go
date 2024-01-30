package main

import (
	"bitbucket.org/johnnewcombe/telstar-emf/article"
	"bitbucket.org/johnnewcombe/telstar-library/logger"
	"bitbucket.org/johnnewcombe/telstar-library/text"
	"bitbucket.org/johnnewcombe/telstar-library/types"
	"bitbucket.org/johnnewcombe/telstar-library/utils"
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/mmcdole/gofeed"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type templateStruct struct {
	line   string
	length int
}

func main() {

	var (
		files             []fs.FileInfo
		articles          []article.Article
		err               error
		fileData          []byte
		flags             *flag.FlagSet
		inputDirectory    string
		outputDirectory   string
		templateDirectory string
		dateRequired      string
		deleteFiles       bool
		fileMatches       []string
		templatePath      string
	)

	flags = flag.NewFlagSet("telstar-client", flag.ExitOnError)
	flags.StringVar(&inputDirectory, "i", "./data/json", "Path to the directory containing JSON type .json files.")
	flags.StringVar(&templateDirectory, "t", "./data/templates", "Path to the directory containing JSON frame templates files.")
	flags.StringVar(&outputDirectory, "o", "./data/frames", "Path to the directory containing output frames.")
	flags.StringVar(&dateRequired, "f", "", "Filter by start date.")
	flags.BoolVar(&deleteFiles, "d", false, "Deletes the contents of the output directory before starting.")

	if err = flags.Parse(os.Args[1:]); err != nil {
		logger.LogError.Printf("%v%", err)
		os.Exit(1)
	}

	// FIXME need to split the source to separate dates or the results
	//  could do this with separate utility or something clever here maybe?

	if deleteFiles {
		// delete the output directory of json files
		if fileMatches, err = filepath.Glob(path.Join(outputDirectory, "*.json")); err != nil {
			logger.LogError.Print(err)
			return
		}
		for _, f := range fileMatches {
			if err = os.Remove(f); err != nil {
				logger.LogError.Print(err)
				return
			}
		}
	}

	// get the files
	if files, err = getDataFiles(inputDirectory); err != nil {
		logger.LogError.Print(err)
		return
	}

	// loop though each file
	for _, file := range files {

		logger.LogInfo.Printf("Processing: %s", file.Name())

		// get articles
		if !file.IsDir() {
			if fileData, err = os.ReadFile(path.Join(inputDirectory, file.Name())); err != nil {
				logger.LogError.Print(err)
				continue
			}

			if err = json.Unmarshal(fileData, &articles); err != nil {
				logger.LogError.Print(err)
				continue
			}

			//			if articles, err = parseRss(inputDirectory, file.Name()); err != nil {
			//				logger.LogError.Print(err)
			//				continue
			//			}
			//for _, article := range articles {
			//	fmt.Printf("%s  %s\r\n", article.StartDate[:16], article.Title)
			//}
		}

		// get the path to the template
		templatePath = getTemplatePath(file.Name(), templateDirectory, dateRequired)

		// create the frames for this data file
		if err = createFrames(articles, templatePath, outputDirectory, dateRequired); err != nil {
			logger.LogError.Print(err)
			continue
		}
	}
}

func createFrames(articles []article.Article, templatePath string, outputDirectory string, dateRequired string) error {

	const (
		MAXLINES       = 18
		MAXCOLS        = 39
		TITLETAG       = "[TITLE]"
		CONTENTTAG     = "[CONTENT]"
		PUBLISHDATETAG = "[PUBLISHDATE]"
	)

	var (
		err           error
		templateJson  []byte
		templateFrame types.Frame
		currentFrame  types.Frame
		linesUsed     int
		markupText    string
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

	template := strings.Split(templateFrame.Content.Data, ",")

	// With thw above example we now have the following e.g.
	//  "[R][TITLE]"       // [TITLE] in Red
	//  "[W][CONTENT]"     // [CONTENT] in white
	//  "[R][PUBLISHDATE]" // [PUBLIHDATE] in Red
	//  "[l.]"             // Separator row of low dots
	//
	// Note that each of the [TITLE], [CONTENT] and [PUBLISHDATE] could be multiple lines of text.

	// run through the template lines and calculate widths for the TITLE, CONTENT, DATE placeholders
	var templates []templateStruct
	var length int

	for _, templateLine := range template {

		// TODO functionalise this ?
		if strings.Contains(templateLine, TITLETAG) {

			//  count how long is the markup line is  without the [TITLE] placeholder
			//  this will allow us to calculate the format textWidth needed e.g.
			length = text.GetMarkupLen(templateLine) - len(TITLETAG)

		} else if strings.Contains(templateLine, CONTENTTAG) {

			//  count how long is the markup line is  without the [CONTENT] placeholder
			//  this will allow us to calculate the format textWidth needed e.g.
			length = text.GetMarkupLen(templateLine) - len(CONTENTTAG)

		} else if strings.Contains(templateLine, PUBLISHDATETAG) {

			//  count how long is the markup line is  without the [PUBLISHDATE] placeholder
			//  this will allow us to calculate the format textWidth needed e.g.
			length = text.GetMarkupLen(templateLine) - len(PUBLISHDATETAG)
		} else {
			length = text.GetMarkupLen(templateLine)
		}

		templates = append(templates, templateStruct{templateLine, MAXCOLS - length})

	}

	//actual content for our frame
	if currentFrame, err = getFrame(templateJson, templateFrame.PID); err != nil {
		return err
	}
	var frameContent strings.Builder

	save := func(toBeContinued bool) error {
		// save old frame get new frame
		if toBeContinued {
			frameContent.WriteString("more...")
		}
		currentFrame.Content.Data = frameContent.String()
		logger.LogInfo.Printf("Saving: %d%s", currentFrame.PID.PageNumber, currentFrame.PID.FrameId)
		if currentFrame, err = saveFrame(currentFrame, outputDirectory+"/"+currentFrame.PID.String()+".json"); err != nil {
			return err
		}
		frameContent.Reset()
		linesUsed = 0
		return nil
	}

	if len(articles) > 32 {
		articles = articles[:32]
	}

	if len(articles) == 0 {
		// create no records page
		frameContent.WriteString("\r\n\r\n\r\n\r\n\r\n[D]         No Information Found")
		save(false)
		return nil
	}

	for _, article := range articles {

		if len(article.StartDate) > 0 && article.StartDate[:10] != dateRequired {
			continue
		}

		for _, templateLine := range templates {

			trows, trowCount := text.Format( article.Venue + ": " + article.Title, templateLine.length)
			crows, crowCount := text.Format(article.Description, templateLine.length)
			prows, pdrowCount := text.Format(article.StartDate, templateLine.length)

			if strings.Contains(templateLine.line, TITLETAG) {

				// get the title and apply template to each line of the title
				// get the title as a slice of strings each string representing each row
				markup := applyTemplate(trows, templateLine.line, TITLETAG)

				// all of the title must fit along with at least one row of content
				if linesUsed+trowCount+1 > maxLines {

					if err := save(false); err != nil {
						return err
					}
				}
				// add the text to the frame
				frameContent.WriteString(markup)
				linesUsed += trowCount

			} else if strings.Contains(templateLine.line, CONTENTTAG) {

				// get the content and apply template to each line of the title
				for index, row := range crows {
					// useful breakpoint = currentFrame.PID.PageNumber==2001 && currentFrame.PID.FrameId=="f"

					// at least one line should fit but we need to make sure that the publish date is not left on its
					// own at the beginning of the next page.
					//If we are on the last line of the content (index = crowCount-1) and there isn't enough
					// room for the publish date (linesUsed + 1 + pdrowCount >= MAXLINES) then put the last line of
					// the content on a new frame.
					if linesUsed > maxLines || (linesUsed+pdrowCount >= maxLines && index == crowCount-1) {
						if err := save(true); err != nil {
							return err
						}
					}
					// add the row
					markupText = applyTemplate([]string{row}, templateLine.line, CONTENTTAG)
					frameContent.WriteString(markupText)
					linesUsed++
				}

			} else if strings.Contains(templateLine.line, PUBLISHDATETAG) {

				// get the published date  and apply template to each line of the title
				//rows, rowCount = text.Format(article.Date, templateLine.length)
				markup := applyTemplate(prows, templateLine.line, PUBLISHDATETAG)

				// all of the date must fit
				if linesUsed+(pdrowCount) > maxLines {
					if err := save(false); err != nil {
						return err
					}
				}
				// add the text to the frame
				frameContent.WriteString(markup)
				linesUsed += pdrowCount

			} else { // Separator

				// this could be markup such as [l.] which is a horizontal line 39 chars long
				// or just plain text, either way no formatting is done and it is treated
				// as one line.
				frameContent.WriteString(templateLine.line)
				frameContent.WriteString("\r\n")
				linesUsed++

			}
		}
	}

	if frameContent.Len() > 0 {
		if err := save(false); err != nil {
			return err
		}
	}

	/***************** DEBUG ******************/
	//fmt.Printf("%v\n", templateFrame)
	/******************************************/

	//TODO create the json frame from templateJson
	//if jsonFrame, err = json.Marshal(templateFrame); err != nil{
	//	return err
	//}
	//TODO determine the space we have based on the templateJson
	//TODO paginate content
	/*
		for i = 0; i < len(titleLines); i++ {
			fmt.Printf("%v. Title: %v\n", i, titleLines[i])
		}
		for i = 0; i < len(descriptionLines); i++ {
			fmt.Printf("%v. Description :%v\n", i, descriptionLines[i])
		}
		for i = 0; i < len(dateLines); i++ {
			fmt.Printf("%v. Date :%v\n", i, dateLines[i])
		}
		fmt.Printf("%v. Template :%v\n\n", i, article.TemplateName)

	*/
	//}
	return nil
}

func applyTemplate(textLines []string, template string, placeHolder string) string {

	var sbText strings.Builder
	for _, line := range textLines {
		sbText.WriteString(strings.ReplaceAll(template, placeHolder, line))
		sbText.WriteString("\r\n")
	}

	return sbText.String()

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

func parseRss(dataDirectory string, filename string) ([]article.Article, error) {

	var (
		err     error
		result  []article.Article
		feed    *gofeed.Feed
		xmlData []byte
	)

	if xmlData, err = ioutil.ReadFile(fmt.Sprintf("%s/%s", dataDirectory, filename)); err != nil {
		return nil, fmt.Errorf("%v\n", err)
	}
	sData := text.CleanUtf8(string(xmlData))

	fp := gofeed.NewParser()
	if feed, err = fp.ParseString(sData); err != nil {
		return nil, fmt.Errorf("%v\n", err)
	}

	for _, item := range feed.Items {

		result = append(result, article.Article{
			Title:       item.Title,
			Description: item.Description,
			StartDate:   item.Published,
		})
	}

	return result, nil
}

func getTemplatePath(dataFileName string, templateDirectory string, dateRequired string) string {

	templateName := dataFileName[:len(dataFileName)-len(path.Ext(dataFileName))]
	return fmt.Sprintf("%s/%s-%s.json", templateDirectory, templateName, dateRequired)

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

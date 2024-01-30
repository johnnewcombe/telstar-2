package main

import (
	"bitbucket.org/johnnewcombe/telstar-library/logger"
	"bitbucket.org/johnnewcombe/telstar-library/utils"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"strings"
)

func cmdAddFrames(apiUrl string, srcDirectory string, primary bool, token string, basePage int, includeUnsafe bool) (ResponseData, error) {

	var (
		err      error
		respData ResponseData
		files    []fs.FileInfo
		count    int
		pageNo   int
		frameId  string
	)

	if files, err = ioutil.ReadDir(srcDirectory); err != nil {
		log.Fatal(err)
	}

	for _, f := range files {

		// get file extension
		ext := strings.ToLower(path.Ext(f.Name()))

		if f.IsDir() {
			continue
		}
		fullPath := filepath.Join(srcDirectory, f.Name())

		if ext == ".json" || ext == ".tf" || ext == ".edittf" || ext == ".zxnet" {

			//attempt to add the frame if a file of the correct type .json
			if respData, err = cmdAddFrame(apiUrl, fullPath, primary, token, basePage, includeUnsafe); err != nil {
				logger.LogError.Println(err)
				continue
			}

			// keep a count of frames updated
			count++

			// get the pid of the last frame uploaded
			pageNo, frameId, err = utils.ConvertPageIdToPID(getPageIdFromFilename(f.Name()))

		} else if ext == ".delete" || ext == ".deleted" {

			// get page ID
			pageId := getPageIdFromFilename(f.Name())

			// delete the frame
			if respData, err = cmdDeleteFrame(apiUrl, pageId, primary, token); err != nil {
				logger.LogError.Println(err)
				continue
			}
		}
	}

	/*  e.g.
	{"records added": 10, "last-frame-added": {"page-no": 31, "frame-id": "j"}}
	*/
	respData.Body = fmt.Sprintf("{\"records added\": %d, \"last-frame-added\": {\"page-no\": %d, \"frame-id\": \"%s\"}}\r\n", count, pageNo, frameId)
	return respData, nil
}

func getPageIdFromFilename (filename string) string{

	// get page ID
	_, filename = path.Split(filename)
	return strings.Split(filename, ".")[0]

}
package main

/*
import (
	"bitbucket.org/johnnewcombe/telstar-library/utils"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"strings"
)

func cmdAddFrames(apiUrl string, srcDirectory string, primary bool, token string) (ResponseData, error) {

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

		if strings.HasSuffix(strings.ToLower(f.Name()), ".json") ||
			strings.HasSuffix(strings.ToLower(f.Name()), ".edit.tf") {

			//load the file if a .json
			if f.IsDir() {
				continue
			}
			fullPath := filepath.Join(srcDirectory, f.Name())

			if respData, err = cmdAddFrame(apiUrl, fullPath, primary, token); err != nil {
				return respData, err
			}

			// keep a count of frames updated
			count++
			// get the pid of the last frame uploaded
			pageNo, frameId, err = utils.ConvertPageIdToPID(f.Name()[:len(path.Ext(f.Name()))-2])
		}
	}
	//  e.g.
	//	{"records added": 10, "last-frame-added": {"page-no": 31, "frame-id": "j"}}

	respData.Body = fmt.Sprintf("{\"records added\": %d, \"last-frame-added\": {\"page-no\": %d, \"frame-id\": \"%s\"}}\r\n", count, pageNo, frameId)
	return respData, nil
}


*/

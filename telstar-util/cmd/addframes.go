package cmd

import (
	"bitbucket.org/johnnewcombe/telstar-util/network"
	"github.com/spf13/cobra"
	"io/fs"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"strings"
)

var addFrames = &cobra.Command{
	Use:   "add-frames",
	Short: "Adds multiple frames to the currently logged in system.",
	Long: `
Adds multiple frame to the currently logged in system. See the login command.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		// TODO Handle GIT repositories as well as directories

		var (
			err      error
			source   string
			files    []fs.FileInfo
			count    int
			last     string
			respData network.ResponseData
		)

		respData.SetOK()

		if source, err = cmd.Flags().GetString("source"); err != nil {
			return err
		}

		if files, err = ioutil.ReadDir(source); err != nil {
			log.Fatal(err)
		}

		for _, f := range files {

			if strings.HasSuffix(strings.ToLower(f.Name()), ".json") ||
				strings.HasSuffix(strings.ToLower(f.Name()), ".edit.tf") {

				//load the file if a .json
				if f.IsDir() {
					continue
				}
				fullPath := filepath.Join(source, f.Name())

				if respData, err = addSingleFrame(cmd, fullPath); err != nil {
					return err
				}
				if respData.StatusCode < 200 || respData.StatusCode > 299 {
					//response = fmt.Sprintf(globals.Response, respData.Status)
					break
				}

				// keep a count of frames updated
				count++
				last = f.Name()
			}
		}

		result := map[string]string{
			"Updated": strconv.Itoa(count),
			"Last":    last,
		}

		stdOut(cmd, respData, result)
		return nil
	},
}

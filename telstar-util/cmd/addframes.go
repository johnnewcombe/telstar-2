package cmd

import (
	"bitbucket.org/johnnewcombe/telstar-util/network"
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"os"
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

		var (
			err       error
			source    string
			files     []os.DirEntry
			count     int
			last      string
			respData  network.ResponseData
			sourceDir string
			result    map[string]string
		)

		// init result
		result = map[string]string{}

		respData.SetOK()

		if source, err = cmd.Flags().GetString("source"); err != nil {
			return err
		}

		// if its a git repo clone the repo in temp to get the files
		// other wise assume its a local folder
		if sourceDir, err = processSource(source); err != nil {
			return err
		}

		if files, err = os.ReadDir(sourceDir); err != nil {
			return err
		}

		// FIXME Sort the range of files before this loop so that the last update reported, is of some use
		for _, f := range files {

			if strings.HasSuffix(strings.ToLower(f.Name()), ".json") ||
				strings.HasSuffix(strings.ToLower(f.Name()), ".edit.tf") {

				//load the file if a .json
				if f.IsDir() {
					continue
				}
				fullPath := filepath.Join(sourceDir, f.Name())

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

		//result = map[string]string{
		//	"Updated": strconv.Itoa(count),
		//	"Last":    last,
		//}
		result["Updated"] = strconv.Itoa(count)
		result["Last"] = last

		stdOut(cmd, respData, result)
		return nil
	},
}

func processSource(source string) (string, error) {

	const (
		maxRepoSize int64  = 2 * 1024 * 1024 //2Mb
		tempPath    string = "/tmp/foo"
	)

	var (
		r *git.Repository
		w *git.Worktree
		f os.FileInfo
	)

	if isGitUrl(source) {

		// TODO Check repo size and get files from git
		_, err := git.PlainClone(tempPath, false, &git.CloneOptions{URL: source, Progress: nil})

		if errors.Is(err, git.ErrRepositoryAlreadyExists) {

			//fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))

			// TODO Check size ? and bin if too big
			// open the repo
			if r, err = git.PlainOpen(tempPath); err != nil {
				//				fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
				return "", err
			}

			// get the working directory
			if w, err = r.Worktree(); err != nil {
				//				fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
				return "", err
			}

			// pull from origin
			err = w.Pull(&git.PullOptions{RemoteName: "origin"})
			if !errors.Is(err, git.NoErrAlreadyUpToDate) {
				//				fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
				return "", err
			}

		}

		// check repo size
		if f, err = os.Lstat(tempPath); err != nil {
			return "", errors.New("unable to determine size of repository")
		}

		if f.Size() > maxRepoSize {
			if err = os.RemoveAll(tempPath); err != nil {
				return "", fmt.Errorf("repository size too big, %s", err)
			}
			return "", fmt.Errorf("repository size too big")
		}

		// return repo folder
		return tempPath, nil // just return source for now

	} else {
		return source, nil
	}
}

func isGitUrl(url string) bool {
	if strings.HasSuffix(strings.ToLower(url), "git") {
		return true
	}
	return false
}

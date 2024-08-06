package cmd

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/johnnewcombe/telstar-library/types"
	"github.com/johnnewcombe/telstar-util/globals"
	"github.com/johnnewcombe/telstar-util/network"
	"github.com/spf13/cobra"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

var addFrames = &cobra.Command{
	Use:   "add-frames",
	Short: "Adds multiple frames to the currently logged in system.",
	Long: `
Adds multiple frame to the currently logged in system. See the login command.
If a frame already exists, it frame will be updated, if it does not exist 
then it will be created.

Setting the page-id option. restricts the update to a specific page. For 
example adding Page 222 might include 222a to 222z if they existed in the 
source. If the frames already exist, they will be updated, if they do not 
exist then they will be created.
`,
	RunE: func(cmd *cobra.Command, args []string) error {

		var (
			apiUrl        string
			primary       bool
			includeUnsafe bool
			token         string
			err           error
			source        string
			respData      network.ResponseData
			result        map[string]string
		)

		//load token - don't want to report errors here as we want an unauthorised status code to be returned
		token, _ = loadText(globals.TOKENFILE)

		if apiUrl, err = cmd.Flags().GetString("url"); err != nil {
			return err
		}

		if includeUnsafe, err = cmd.Flags().GetBool("include-unsafe"); err != nil {
			return err
		}

		if primary, err = cmd.Flags().GetBool("primary"); err != nil {
			return err
		}

		apiUrl = apiUrl + "/frame"
		if primary {
			apiUrl += "?db=primary"
		}

		respData.SetOK()

		if source, err = cmd.Flags().GetString("source"); err != nil {
			return err
		}

		// if its a git repo clone the repo in temp to get the files
		// other wise assume its a local folder
		if isGitUrl(source) {
			// note that setting page ID to "" prevents any restrictions in updating
			if respData, result, err = processFramesGit(apiUrl, source, includeUnsafe, token); err != nil {
				return err
			}
		} else {
			// note that setting page ID to "" prevents any restrictions in updating
			if respData, result, err = processFramesFs(apiUrl, source, includeUnsafe, token); err != nil {
				return err
			}
		}

		stdOut(cmd, respData, result)
		return nil
	},
}

// processFramesFs This passes the filenames of the specified folder to the AddFrame command for processing
func processFramesFs(apiUrl string, sourceDir string, includeUnsafe bool, token string) (network.ResponseData, map[string]string, error) {

	var (
		err       error
		files     []os.DirEntry
		frame     types.Frame
		count     int
		respData  network.ResponseData
		result    map[string]string
		frameData string
	)

	result = map[string]string{}

	if files, err = os.ReadDir(sourceDir); err != nil {
		return respData, result, err
	}

	for _, f := range files {

		if strings.HasSuffix(strings.ToLower(f.Name()), ".json") ||
			strings.HasSuffix(strings.ToLower(f.Name()), ".edit.tf") {

			//load the file if a .json
			if f.IsDir() {
				continue
			}
			fullPath := filepath.Join(sourceDir, f.Name())

			if frameData, err = loadText(fullPath); err != nil {
				return respData, result, err
			}
			if isEditTfFrame(frameData) {
				return respData, result, errors.New("edit.tf frames are only supported as part of a standard json frame definition")
			}

			// validate the frameData
			if frame, err = parseFrame(frameData); err != nil {
				err = fmt.Errorf("invalid frame data: %v", err)
				return respData, result, err
			}

			// validate the PID
			if !frame.IsValid() {
				err = fmt.Errorf("invalid frame id: (frame %s)", frame.PID.String())
				return respData, result, err
			}

			if respData, err = addSingleFrameJson(apiUrl, frame, includeUnsafe, token); err != nil {
				return respData, result, err
			}

			count++
			if respData.StatusCode < 200 || respData.StatusCode > 299 {
				break
			}

			// keep a count of frames updated
			//last = f.Name()
		}
	}

	result = map[string]string{
		"Updated": strconv.Itoa(count),
		//"Last":    last,
	}
	return respData, result, nil
}

// processGit The function clones/pulls the specified repo into memory and extracts each frame as json data.
// Each json frame this is passed to the helper method addSingleFrameJson used by the AddFrame command.
func processFramesGit(apiUrl string, source string, includeUnsafe bool, token string) (network.ResponseData, map[string]string, error) {

	const (
		maxFileSize int64 = 16 * 1024 // 16Kb
		maxFiles    int   = 20000     // 320Mb total e.g. 20000 * 16Kb

	)

	var (
		r        *git.Repository
		respData network.ResponseData
		ref      *plumbing.Reference
		err      error
		commit   *object.Commit
		tree     *object.Tree
		frame    types.Frame
		count    int
		result   map[string]string
		reader   io.ReadCloser
	)

	result = map[string]string{}

	// clone to memory
	r, err = git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:          source,
		SingleBranch: true,
	})

	if ref, err = r.Head(); err != nil {
		return respData, result, err
	}
	if commit, err = r.CommitObject(ref.Hash()); err != nil {
		return respData, result, err
	}
	if tree, err = commit.Tree(); err != nil {
		return respData, result, err
	}

	if err = tree.Files().ForEach(func(f *object.File) error {

		if strings.HasSuffix(strings.ToLower(f.Name), ".json") ||
			strings.HasSuffix(strings.ToLower(f.Name), ".edit.tf") {

			if f.Size > maxFileSize {
				return fmt.Errorf("file %s (%d bytes) is too big, max file size  = %d bytes", f.Name, f.Size, maxFileSize)
			}
			if count > maxFiles {
				return fmt.Errorf("the maximum number of files has ben exceeded (%d)", maxFiles)
			}

			// get the Reader
			if reader, err = f.Blob.Reader(); err != nil {
				return err
			}

			// create a suitably sized buffer
			buf := make([]byte, f.Size)

			// read the ata
			_, err = reader.Read(buf) // reader reads len(buf) bytes
			reader.Close()

			// validate the frameData
			if frame, err = parseFrame(string(buf)); err != nil {
				return fmt.Errorf("invalid frame data: %v", err)
			}

			if !frame.IsValid() {
				return fmt.Errorf("invalid frame data for frame %s", frame.GetPageId())
			}

			if respData, err = addSingleFrameJson(apiUrl, frame, includeUnsafe, token); err != nil {
				return err
			}

			// exit loop if status <200 or >299
			if respData.StatusCode < 200 || respData.StatusCode > 299 {
				// bad response so we need to cancel the iteration completely not just this file
				// we can do that with a custom error
				return &network.RequestError{}
			}

			count++

		}

		return nil

	}); err != nil {

		if !errors.Is(err, &network.RequestError{}) {
			return respData, result, err
		}
	}

	// something happened so update result
	result = map[string]string{
		"Updated": strconv.Itoa(count),
	}
	return respData, result, nil

}

func isGitUrl(url string) bool {
	if strings.HasSuffix(strings.ToLower(url), "git") {
		return true
	}
	return false
}

// this is the default sort order of golang ReadDir
func sortFileNameAscend(files []fs.DirEntry) {
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})
}

func sortFileNameDescend(files []fs.DirEntry) {
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() > files[j].Name()
	})
}

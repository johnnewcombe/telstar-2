package cmd

import (
	"bitbucket.org/johnnewcombe/telstar-library/types"
	"bitbucket.org/johnnewcombe/telstar-library/utils"
	"bitbucket.org/johnnewcombe/telstar-util/globals"
	"bitbucket.org/johnnewcombe/telstar-util/network"
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/spf13/cobra"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var addPage = &cobra.Command{
	Use:   "add-page",
	Short: "Adds a root frame and all follow on frames to the currently logged in system.",
	Long: `
Adds a collection of frames representing a page. For example if page number
222 was specified then all contiguous frames starting from 222a, that existed 
within the source, would be added to Telstar. If the frames already exist,
they will be updated, if they do not exist then they will be created. Any
contiguous frames within Telstar beyond those being added, e.g. 222g onwards
in the above example, would be deleted (purged). For example, if frames 222a,
222b, 222c, 222d and 222e already exist within Telstar and only frames 222a,
222b and 222c are being updated, the purge process would remove contiguous 
frames from 222d i.e. 222d and 222e.
`,
	RunE: func(cmd *cobra.Command, args []string) error {

		var (
			apiUrl        string
			primary       bool
			includeUnsafe bool
			pageNo        int
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
		if pageNo, err = cmd.Flags().GetInt("page-no"); err != nil {
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

		// set frame we are looking for
		if pageNo >= 0 {
			//currentFrameInPage = strconv.Itoa(pageNo) + "a"
		} else {
			//currentFrameInPage = ""
		}

		// if its a git repo clone the repo in temp to get the files
		// other wise assume its a local folder

		if isGitUrl(source) {
			if respData, result, err = processPageGit(apiUrl, source, includeUnsafe, pageNo, token); err != nil {
				return err
			}
		} else {
			if respData, result, err = processPageFs(apiUrl, source, includeUnsafe, pageNo, token); err != nil {
				return err
			}
		}

		stdOut(cmd, respData, result)
		return nil

	},
}

// processFramesFs This passes the filenames of the specified folder to the AddFrame command for processing
func processPageFs(apiUrl string, sourceDir string, includeUnsafe bool, pageNo int, token string) (network.ResponseData, map[string]string, error) {

	var (
		err          error
		files        []os.DirEntry
		frame        types.Frame
		lastFrame    types.Pid
		frames       []types.Frame
		count        int
		respData     network.ResponseData
		delResult    network.ResponseData
		result       map[string]string
		frameData    string
		currentFrame string
		urlParser    *url.URL
		apiRespose   types.ApiResponse
	)

	result = map[string]string{}

	if pageNo < 0 {
		return respData, result, errors.New("invalid page number")
	}

	// set start frame
	currentFrame = fmt.Sprintf("%da", pageNo)

	if files, err = os.ReadDir(sourceDir); err != nil {
		return respData, result, err
	}

	// get list of frames from the list of files
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

			frames = append(frames, frame)
		}
	}
	// Sort the range of files based on Page No and Frame ID
	types.SortFrames(frames)

	// check each frame for the correct page and sequence
	for _, frame = range frames {

		if frame.PID.String() == currentFrame {

			// add frame to slice

			// add frame
			if respData, err = addSingleFrameJson(apiUrl, frame, includeUnsafe, token); err != nil {
				return respData, result, err
			}

			if respData.StatusCode < 200 || respData.StatusCode > 299 {
				break
			}

			count++
			lastFrame = frame.PID

			if currentFrame, err = utils.GetFollowOnPageId(currentFrame); err != nil {
				return respData, result, err
			}
		}
	}

	// if we get here and count > 0 then an update has occurred and a purge is needed
	// the currentFrame var holds the next frame so delete that with purge=true
	if count > 0 {

		apiUrl += currentFrame + "?purge=true"

		// as we are now deleting we need to insert the frame to be deleted.
		urlParser, err = url.Parse(apiUrl)
		apiUrl = fmt.Sprintf("%s://%s%s/%s?%s", urlParser.Scheme, urlParser.Host, urlParser.Path, currentFrame, urlParser.RawQuery)

		//deleted := 0
		if delResult, err = deleteSingleFrame(apiUrl, token); err != nil {
			return delResult, result, err
		}

		//deleteResult is a local result object but the body text of that object is actually the text output of the
		// types.ApiResponse library object so we can load that and get the message element

		apiRespose.Load(delResult.Body)
	}

	result = map[string]string{
		"Updated":      strconv.Itoa(count),
		"Last Frame":   lastFrame.String(),
		"Purge Result": apiRespose.ResultText,
		"Purge From":   currentFrame,
	}
	return respData, result, nil
}

// processGit The function clones/pulls the specified repo into memory and extracts each frame as json data.
// Each json frame this is passed to the helper method addSingleFrameJson used by the AddFrame command.
func processPageGit(apiUrl string, source string, includeUnsafe bool, pageNo int, token string) (network.ResponseData, map[string]string, error) {

	const (
		maxFileSize int64 = 16 * 1024 // 16Kb
		maxFiles    int   = 20000     // 320Mb total e.g. 20000 * 16Kb

	)

	var (
		r         *git.Repository
		respData  network.ResponseData
		ref       *plumbing.Reference
		err       error
		commit    *object.Commit
		tree      *object.Tree
		frame     types.Frame
		lastFrame types.Pid
		frames    []types.Frame
		count     int
		result    map[string]string
		reader    io.ReadCloser
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
				err = fmt.Errorf("invalid frame data for file %s: %v", f.Name, err)
				return err
			}
			// all good to add to frames slice
			frames = append(frames, frame)
		}

		return nil

	}); err != nil {

		if !errors.Is(err, &network.RequestError{}) {
			return respData, result, err
		}
	}

	// sort the frames by PID
	types.SortFrames(frames)

	// send frames to telstar
	for _, frame = range frames {

		if respData, err = addSingleFrameJson(apiUrl, frame, includeUnsafe, token); err != nil {
			return respData, result, err
		}

		// exit loop if status <200 or >299
		if respData.StatusCode < 200 || respData.StatusCode > 299 {
			// bad response so we need to cancel the iteration completely not just this file
			// we can do that with a custom error
			return respData, result, &network.RequestError{}
		}

		count++
		lastFrame = frame.PID
	}

	// something happened so update result
	result = map[string]string{
		"Updated":   strconv.Itoa(count),
		"Last Fame": lastFrame.String(),
	}
	return respData, result, nil

}

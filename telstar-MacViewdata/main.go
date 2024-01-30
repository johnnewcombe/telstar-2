package main

import (
	"bitbucket.org/johnnewcombe/telstar-library/file"
	"fmt"
	"os"
	"strings"
)

func main() {

	var (
		data   []byte
		err    error
		result = strings.Builder{}
	)

	// check the number of arguments
	if len(os.Args) <2 {
		fmt.Fprint(os.Stderr, "error: filename not specified\n")
		os.Exit(1)
	}

	// get the filename
	filename := os.Args[1]

	// check it exists
	if !file.Exists(filename){
		fmt.Fprint(os.Stderr, "error: file not found\n")
		os.Exit(1)
	}

	// load the binary file
	if data, err = file.ReadFile(filename); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if len(data) != 960 {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// create suitable json from the binary file
	result.WriteString("  \"content\": {\n    \"data\": \"")
	for _, byt := range data {

		// fix for OBBS format
		if byt >= 0xC0 && byt <= 0xdf {
			byt = byt & 0x7f
			result.WriteString("\\u001b")
		}

		result.WriteString(fmt.Sprintf("\\u00%x", byt))
	}
	result.WriteString("\",\n    \"type\": \"rawV\"\n  },\n")
	fmt.Println(result.String())

	os.Exit(0)

}

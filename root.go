package editor

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

var editor = "vim"

func init() {
	if editorVar := os.Getenv("EDITOR"); editorVar != "" {
		editor = editorVar
	}
}

// Opens the default editor and returns the value.
func Read() ([]byte, error) {
	return ReadEditor(editor, "")
}

// Opens the default editor and returns the value in string format.
func ReadText() (string, error) {
	text, err := Read()
	return string(text), err
}

// Opens the editor and returns the value.
func ReadEditor(editor, programName string) ([]byte, error) {

	if programName == "" {
		programName = "editor"
	}

	// Create a temporary file.

	tempFile, tmpFileError := ioutil.TempFile("", programName)

	if tmpFileError != nil {
		return nil, tmpFileError
	}

	defer os.Remove(tempFile.Name())

	// open editor

	cmd := exec.Command("sh", "-c", fmt.Sprintf("%s %s", editor, tempFile.Name()))

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmdError := cmd.Run()

	if cmdError != nil {
		return nil, cmdError
	}

	// read tmpfile

	text, readingFileError := ioutil.ReadFile(tempFile.Name())

	if readingFileError != nil {
		return nil, readingFileError
	}

	return text, nil
}

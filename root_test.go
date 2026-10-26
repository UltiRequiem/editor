package editor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeSelf is a fake editor: it writes the temp file's own path into it,
// so tests can see which file was opened.
const writeSelf = `f() { printf '%s' "$1" > "$1"; }; f`

func TestReadEditorReturnsContent(t *testing.T) {
	got, err := ReadEditor("printf 'hello world' >", "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(got) != "hello world" {
		t.Errorf("got %q, want %q", got, "hello world")
	}
}

func TestReadEditorEmptyFile(t *testing.T) {
	got, err := ReadEditor("true", "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 0 {
		t.Errorf("got %q, want empty content", got)
	}
}

func TestReadEditorMultiline(t *testing.T) {
	want := "line one\nline two\n"

	got, err := ReadEditor(`printf 'line one\nline two\n' >`, "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(got) != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestReadEditorRemovesTempFile(t *testing.T) {
	path, err := ReadEditor(writeSelf, "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(path) == 0 {
		t.Fatal("fake editor did not write the temp file path")
	}

	if _, statErr := os.Stat(string(path)); !os.IsNotExist(statErr) {
		t.Errorf("temp file %q still exists after ReadEditor returned", path)
	}
}

func TestReadEditorUsesProgramNameAsPrefix(t *testing.T) {
	path, err := ReadEditor(writeSelf, "myprogram")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if base := filepath.Base(string(path)); !strings.HasPrefix(base, "myprogram") {
		t.Errorf("temp file %q does not start with %q", base, "myprogram")
	}
}

func TestReadEditorDefaultProgramName(t *testing.T) {
	path, err := ReadEditor(writeSelf, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if base := filepath.Base(string(path)); !strings.HasPrefix(base, "editor") {
		t.Errorf("temp file %q does not start with default prefix %q", base, "editor")
	}
}

func TestReadEditorCommandFails(t *testing.T) {
	got, err := ReadEditor("false", "test")
	if err == nil {
		t.Fatal("expected an error when the editor exits with a non-zero status")
	}

	if got != nil {
		t.Errorf("got %q, want nil content on error", got)
	}
}

func TestReadEditorMissingCommand(t *testing.T) {
	if _, err := ReadEditor("this-editor-does-not-exist-12345", "test"); err == nil {
		t.Fatal("expected an error when the editor command does not exist")
	}
}

// setEditor swaps the package-level editor for the duration of a test.
func setEditor(t *testing.T, value string) {
	t.Helper()

	previous := editor
	editor = value

	t.Cleanup(func() { editor = previous })
}

func TestRead(t *testing.T) {
	setEditor(t, "printf 'from Read' >")

	got, err := Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(got) != "from Read" {
		t.Errorf("got %q, want %q", got, "from Read")
	}
}

func TestReadText(t *testing.T) {
	setEditor(t, "printf 'from ReadText' >")

	got, err := ReadText()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "from ReadText" {
		t.Errorf("got %q, want %q", got, "from ReadText")
	}
}

func TestReadTextError(t *testing.T) {
	setEditor(t, "false")

	got, err := ReadText()
	if err == nil {
		t.Fatal("expected an error when the editor fails")
	}

	if got != "" {
		t.Errorf("got %q, want empty string on error", got)
	}
}

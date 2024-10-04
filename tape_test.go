package main

import (
	"io"
	"testing"
)

func TestTape_Write(t *testing.T) {
	file, clean := createTempFile(t, "ABCDEDF")
	defer clean()

	tape := &Tape{file}
	tape.Write([]byte("123"))

	file.Seek(0, io.SeekStart)
	newFileContents, _ := io.ReadAll(file)

	got := string(newFileContents)
	want := "123"

	if got != want {
		t.Errorf("got %s and wanted %s", got, want)
	}

}

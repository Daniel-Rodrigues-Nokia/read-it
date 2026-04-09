package internal

import (
	"fmt"
	"os"
	"strings"
)

type log struct {
	dir string
}

func createTempPath(path string) string {
	return fmt.Sprintf("%s%s%s", os.TempDir(), string(os.PathSeparator), path)
}

func NewLog() log {
	tempDir := createTempPath("read-it.log")

	return log{dir: tempDir}
}

func (l log) Log(msg string) error {
	f, err := os.OpenFile(l.dir, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Printf("\n[read-it]: error detected - check %s\n", l.dir)

	buffer := strings.Builder{}

	buffer.WriteString("\n---------------\n")
	buffer.WriteString(msg)
	buffer.WriteString("\n---------------\n")

	if _, err := f.WriteString(buffer.String()); err != nil {
		return err
	}

	return nil
}

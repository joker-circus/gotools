package file

import (
	"bytes"
	"fmt"
	"github.com/joker-circus/gotools"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Size return the file size of the specified path file
func Size(filePath string) (int64, bool, error) {
	file, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, false, nil
		}
		return 0, false, err
	}
	return file.Size(), true, nil
}

// Name Converts a string to a valid filename
func Name(name, ext string, length int) string {
	rep := strings.NewReplacer("\n", " ", "/", " ", "|", "-", ": ", "：", ":", "：", "'", "’")
	name = rep.Replace(name)
	if runtime.GOOS == "windows" {
		rep = strings.NewReplacer("\"", " ", "?", " ", "*", " ", "\\", " ", "<", " ", ">", " ")
		name = rep.Replace(name)
	}
	limitedName := gotools.LimitLength(name, length)
	if ext == "" {
		return limitedName
	}
	return fmt.Sprintf("%s.%s", limitedName, ext)
}

// Path gen valid file path
func Path(name, ext string, length int, outputPath string, escape bool) (string, error) {
	if outputPath != "" {
		if _, err := os.Stat(outputPath); err != nil {
			return "", err
		}
	}
	var fileName string
	if escape {
		fileName = Name(name, ext, length)
	} else {
		fileName = fmt.Sprintf("%s.%s", name, ext)
	}
	return filepath.Join(outputPath, fileName), nil
}

// LineCounter Counts line in file
func LineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

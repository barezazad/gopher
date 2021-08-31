package generateerrcode

import (
	"bufio"
	"fmt"
	"gopher/pkg/helper"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func GenerateNewErrCode() {
	files, _ := GetPathAllGoFiles("./")
	StartCode := 1000000
	var code int

	for _, v := range files {

		// open original file
		f, err := os.Open(v)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		if v == "main.go" {
			continue
		}

		// create temp file
		tmp, err := ioutil.TempFile("", "replace-*")
		if err != nil {
			log.Fatal(err)
		}
		defer tmp.Close()

		if code == 0 {
			code = StartCode
		}

		// replace while copying from f to tmp
		if code, err = replace(f, tmp, code); err != nil {
			log.Fatal(err)
		}

		// make sure the tmp file was successfully written to
		if err := tmp.Close(); err != nil {
			log.Fatal(err)
		}

		// close the file we're reading from
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}

		// overwrite the original file with the temp file
		if err := os.Rename(tmp.Name(), v); err != nil {
			log.Fatal(err)
		}
	}
}

func GetPathAllGoFiles(root string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match("*.go", filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func replace(r io.Reader, w io.Writer, code int) (lastCode int, err error) {
	// use scanner to read line by line
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := sc.Text()

		getErrCode := helper.RegexFindTerm(line, `"E[0-9].*[0-9]"`)

		if getErrCode != "" {
			newCode := fmt.Sprintf("\"E%v\"", code)
			line = strings.ReplaceAll(line, getErrCode, newCode)
			code++
		}
		if _, err := io.WriteString(w, line+"\n"); err != nil {
			return code, err
		}
	}
	lastCode = code
	return lastCode, sc.Err()
}

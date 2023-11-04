package common

import (
	"bytes"
	"io"
	"log"
	"regexp"
)

const chunkSize = 64000

func ShouldBeExcluded(name string, exclusions []string) bool {
	for _, exclusion := range exclusions {
		if match, _ := regexp.MatchString(exclusion, name); match {
			return true
		}
	}

	return false
}

func DeepCompare(file1, file2 SharedFile) bool {
	// Check file size ...
	f1Stat, _ := file1.Stat()
	f2Stat, _ := file2.Stat()

	if f1Stat.Size() != f2Stat.Size() {
		return false
	}

	for {
		b1 := make([]byte, chunkSize)
		_, err1 := file1.Read(b1)

		b2 := make([]byte, chunkSize)
		_, err2 := file2.Read(b2)

		if err1 != nil || err2 != nil {
			if err1 == io.EOF && err2 == io.EOF {
				return true
			} else if err1 == io.EOF || err2 == io.EOF {
				return false
			} else {
				log.Fatal(err1, err2)
			}
		}

		if !bytes.Equal(b1, b2) {
			return false
		}
	}
}

package misc

import (
	"io"
	"os"

	"github.com/anotherhadi/eleakxir/leak-utils/settings"
)

func MergeFiles(lu settings.LeakUtils, outputFile string, inputFiles ...string) error {
	out, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer out.Close()

	for _, inputFile := range inputFiles {
		file, err := os.Open(inputFile)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			return err
		}
	}

	return nil
}

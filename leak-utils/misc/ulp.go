package misc

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/anotherhadi/eleakxir/leak-utils/settings"
)

func RemoveUrlSchemeFromUlp(lu settings.LeakUtils, inputFile string) error {
	file, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	outputFile := inputFile + ".clean"
	out, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer out.Close()

	reader := bufio.NewReader(file)
	writer := bufio.NewWriter(out)

	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}

		firstColumn := strings.Index(line, ":")
		firstScheme := strings.Index(line, "://")
		if firstScheme != -1 && firstColumn == firstScheme {
			line = line[firstScheme+3:]
		}

		_, werr := writer.WriteString(line)
		if werr != nil {
			return err
		}

		if err == io.EOF {
			break
		}
	}

	err = writer.Flush()
	if err != nil {
		return err
	}

	err = os.Remove(inputFile)
	if err != nil {
		return err
	}

	err = os.Rename(outputFile, inputFile)
	if err != nil {
		return err
	}

	return nil
}

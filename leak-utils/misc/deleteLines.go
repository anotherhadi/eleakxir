package misc

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/anotherhadi/eleakxir/leak-utils/settings"
)

func deleteXFirstLines(lu settings.LeakUtils, inputFile, outputFile string, x int) error {
	in, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer out.Close()

	reader := bufio.NewReader(in)
	writer := bufio.NewWriter(out)
	defer writer.Flush()

	lineNum := 0
	for {
		line, err := reader.ReadString('\n')

		if err != nil && err != io.EOF {
			return err
		}

		lineNum++
		if lineNum > x {
			if _, writeErr := writer.WriteString(line); writeErr != nil {
				return writeErr
			}
		}

		if err == io.EOF {
			break
		}
	}

	return nil
}

// Delete the last X lines of a file
func deleteXLastLines(lu settings.LeakUtils, inputFile, outputFile string, x int) error {
	in, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer in.Close()

	lines := []string{}
	reader := bufio.NewReader(in)

	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}

		lines = append(lines, strings.TrimSuffix(line, "\n"))

		if err == io.EOF {
			break
		}
	}
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	if x > len(lines) {
		x = len(lines)
	}

	out, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer out.Close()
	writer := bufio.NewWriter(out)
	defer writer.Flush()

	for _, line := range lines[:len(lines)-x] {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

// Print the first X lines of a file with index starting from 1
func printFirstXLines(lu settings.LeakUtils, inputFile string, x int) error {
	in, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer in.Close()

	reader := bufio.NewReader(in)

	lineNum := 0
	for lineNum < x {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}

		lineNum++
		text := strings.TrimSuffix(line, "\n")

		if len(text) > 100 {
			text = text[:100] + "..."
		}
		fmt.Println(settings.Muted.Render(fmt.Sprintf("%d: %s", lineNum, text)))

		if err == io.EOF {
			break
		}
	}

	return nil
}

// Print the last X lines of a file with index starting from 1 from the bottom
func printLastXLines(lu settings.LeakUtils, inputFile string, x int) error {
	in, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer in.Close()

	lines := []string{}
	reader := bufio.NewReader(in)

	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}

		lines = append(lines, strings.TrimSuffix(line, "\n"))

		if err == io.EOF {
			break
		}
	}
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	if x > len(lines) {
		x = len(lines)
	}

	start := len(lines) - x
	subLines := lines[start:]

	for i := 0; i < len(subLines); i++ {
		index := len(subLines) - i
		text := subLines[i]
		if len(text) > 100 {
			text = text[:100] + "..."
		}
		fmt.Println(settings.Muted.Render(fmt.Sprintf("%d: %s", index, text)))
	}

	return nil
}

func DeleteFirstLines(lu settings.LeakUtils, inputFile, outputFile string, show int) error {
	err := printFirstXLines(lu, inputFile, show)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println(settings.Muted.Render("How many lines do you want to delete from the start? "))
	input, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	input = strings.TrimSpace(input)
	x, err := strconv.Atoi(input)
	if err != nil {
		return fmt.Errorf("bad input: %s", err)
	}

	err = deleteXFirstLines(lu, inputFile, outputFile, x)
	if err != nil {
		return err
	}
	return nil
}

func DeleteLastLines(lu settings.LeakUtils, inputFile, outputFile string, show int) error {
	err := printLastXLines(lu, inputFile, show)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println(settings.Muted.Render("How many lines do you want to delete from the end? "))
	input, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	input = strings.TrimSpace(input)
	x, err := strconv.Atoi(input)
	if err != nil {
		return fmt.Errorf("bad input: %s", err)
	}

	err = deleteXLastLines(lu, inputFile, outputFile, x)
	if err != nil {
		return err
	}

	fmt.Println(settings.Muted.Render(fmt.Sprintf("%d lines deleted from the end. Result saved to %s", x, outputFile)))
	return nil
}

package misc

import (
	"bufio"
	"fmt"
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

	scanner := bufio.NewScanner(in)
	writer := bufio.NewWriter(out)
	defer writer.Flush()

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		if lineNum <= x {
			continue
		}
		_, err := writer.WriteString(scanner.Text() + "\n")
		if err != nil {
			return err
		}
	}

	return scanner.Err()
}

// Delete the last X lines of a file
func deleteXLastLines(lu settings.LeakUtils, inputFile, outputFile string, x int) error {
	in, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer in.Close()

	lines := []string{}
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
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

	scanner := bufio.NewScanner(in)
	lineNum := 0
	for scanner.Scan() && lineNum < x {
		lineNum++
		fmt.Println(settings.Muted.Render(fmt.Sprintf("%d: %s", lineNum, scanner.Text())))
	}

	return scanner.Err()
}

// Print the last X lines of a file with index starting from 1 from the bottom
func printLastXLines(lu settings.LeakUtils, inputFile string, x int) error {
	in, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer in.Close()

	lines := []string{}
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	// Si x est supérieur au nombre de lignes, on affiche tout
	if x > len(lines) {
		x = len(lines)
	}

	start := len(lines) - x
	subLines := lines[start:]

	// Affiche les lignes avec 1 en bas
	for i := 0; i < len(subLines); i++ {
		// index = nombre total de lignes affichées - position dans le slice
		index := len(subLines) - i
		fmt.Println(settings.Muted.Render(fmt.Sprintf("%d: %s", index, subLines[i])))
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

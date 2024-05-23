package utils

import (
	"bufio"
	"io/fs"
	"os"
	"strings"
)

const (
	DefaultDirPerms  fs.FileMode = 0755
	DefaultFilePerms fs.FileMode = 0644
)

func ReadLinesMap(path string) (map[string]bool, error) {
	if !FileExists(path) {
		return map[string]bool{}, nil
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lines := map[string]bool{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		entry := scanner.Text()
		if entry != "" {
			lines[entry] = true
		}
	}
	return lines, scanner.Err()
}

func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func ReadLineByLine(path string, action func(line string)) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		action(scanner.Text())
	}
	return nil
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func DirExists(dir string) bool {
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func Mkdir(dirs ...string) []error {
	var errors []error
	for _, dir := range dirs {
		err := os.MkdirAll(dir, DefaultDirPerms)
		if err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

func WriteLines(path string, lines []string) error {
	return os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), DefaultFilePerms)
}

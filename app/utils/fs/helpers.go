package fs

import (
	"iac/utils/errors"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

type DirEntry = fs.DirEntry

// Resolved home directory, converts to absolute path and then creates the directory if it doesn't exist. Returns the absolute path on success.
func MkdirIfNotExists(path string) (string, error) {
	slog.Debug("Enter MkdirIfNotExists", "path", path)

	path, err := AbsPath(path)
	if err != nil {
		return "", err
	}

	stat, err := os.Stat(path)
	if err != nil || stat == nil {
		slog.Debug("Directory does not exist, attempting to create it")
		slog.Info("Creating directory", "path", path)
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return "", errors.New("Failed to create directory", "path", path, "error", err)
		}
	} else if !stat.IsDir() {
		return "", errors.New("Path exists but is not a directory", "path", path)
	}

	slog.Debug("Directory already exists", "path", path)
	return path, nil
}

func AbsPath(path string) (string, error) {
	slog.Debug("ExpandPath", "path", path)

	if path == "" {
		return "", errors.New("Path is empty")
	}

	expandedPath, err := homedir.Expand(path)
	if err != nil {
		return "", errors.New("Could not expand path", "path", path)
	}

	absPath, err := filepath.Abs(expandedPath)
	if err != nil {
		return "", errors.New("Could not turn path into absolute path", "path", path)
	}

	slog.Debug("Converted to absolute path.", "path", path, "absolute path", absPath)
	return absPath, nil
}

type StatResult string

const (
	StatResultStatError StatResult = ""
	StatResultNotExist  StatResult = "StatResultNotExist"
	StatResultFile      StatResult = "StatResultFile"
	StatResultDirectory StatResult = "StatResultDirectory"
)

func Stat(path string) (StatResult, error) {
	slog.Debug("Stat (Enter)", "path", path)

	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			slog.Debug("Stat (Exit)", "path", path, "result", StatResultNotExist)
			return StatResultNotExist, nil
		}
		return StatResultStatError, errors.New("Could not stat path", "path", path, "error", err)
	}

	if stat.IsDir() {
		slog.Debug("Stat (Exit)", "path", path, "result", StatResultDirectory)
		return StatResultDirectory, nil
	} else {
		slog.Debug("Stat (Exit)", "path", path, "result", StatResultFile)
		return StatResultFile, nil
	}
}

func WriteFileFromString(path string, content string) error {
	slog.Debug("WriteFile (Enter)", "path", path, "content", content)

	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return errors.New("Could not write file", "path", path, "error", err)
	}

	slog.Debug("WriteFile (Exit)", "path", path)
	return nil
}

func WriteFileFromBytes(path string, content []byte) error {
	slog.Debug("WriteFile (Enter)", "path", path, "content length", len(content))

	err := os.WriteFile(path, content, 0644)
	if err != nil {
		return errors.New("Could not write file", "path", path, "error", err)
	}
	slog.Debug("WriteFile (Exit)", "path", path)
	return nil
}

func ReadFileAsString(path string) (string, error) {
	slog.Debug("ReadFileAsString (Enter)", "path", path)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", errors.New("Could not read file", "path", path, "error", err)
	}
	content := string(data)
	slog.Debug("ReadFileAsString (Exit)", "path", path, "content", content)
	return content, nil
}

func ReadFileAsBytes(path string) ([]byte, error) {
	slog.Debug("ReadFileAsBytes (Enter)", "path", path)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.New("Could not read file", "path", path, "error", err)
	}
	slog.Debug("ReadFileAsBytes (Exit)", "path", path, "date length", len(data))
	return data, nil
}

func ReadDir(path string) ([]DirEntry, error) {
	slog.Debug("ReadDir (Enter)", "path", path)
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, errors.New("Could not read directory", "path", path, "error", err)
	}
	slog.Debug("ReadDir (Exit)", "path", path, "entries count", len(entries))
	return entries, nil
}

func DeleteFile(path string) error {
	slog.Debug("DeleteFile (Enter)", "path", path)
	err := os.Remove(path)
	if err != nil {
		return errors.New("Could not delete file", "path", path, "error", err)
	}
	slog.Debug("DeleteFile (Exit)", "path", path)
	return nil
}

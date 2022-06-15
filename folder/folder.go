package folder

import "os"

// If the folder path exists
func Exists(folder string) bool {
	if info, err := os.Stat(folder); err == nil && info.IsDir() {
		return true
	}
	return false
}

// Create folder recursively
func Create(path string) error {
	err := os.MkdirAll(path, 0750)
	if err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}

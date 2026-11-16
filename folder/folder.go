package folder

import "os"

// Exists reports whether path exists and is a directory.
func Exists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// Create creates the folder and any missing parents. It is a no-op if the
// folder already exists.
func Create(path string) error {
	return os.MkdirAll(path, 0o750)
}

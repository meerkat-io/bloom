package folder

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExists(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "f.txt")
	if err := os.WriteFile(file, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name string
		path string
		want bool
	}{
		{"dir", dir, true},
		{"file is not a folder", file, false},
		{"missing", filepath.Join(dir, "nope"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Exists(tt.path); got != tt.want {
				t.Errorf("Exists(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	dir := t.TempDir()
	nested := filepath.Join(dir, "a", "b", "c")
	if err := Create(nested); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if !Exists(nested) {
		t.Fatal("folder was not created")
	}
	if err := Create(nested); err != nil {
		t.Fatalf("Create on existing folder: %v", err)
	}
}

func TestCreateOverFile(t *testing.T) {
	file := filepath.Join(t.TempDir(), "f")
	if err := os.WriteFile(file, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := Create(file); err == nil {
		t.Fatal("expected error when path is an existing file")
	}
}

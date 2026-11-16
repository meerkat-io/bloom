package flag

import (
	"os"
	"strings"
	"testing"
)

type cfg struct {
	Help  bool    `flag:"h" usage:"help"`
	Name  string  `flag:"n" usage:"name" tip:"name" default:"anon"`
	Count int     `flag:"c" usage:"count" default:"3"`
	Rate  float64 `flag:"r" usage:"rate" default:"1.5"`
}

type reqCfg struct {
	Lang string `flag:"t" required:"true"`
}

func setArgs(t *testing.T, args ...string) {
	t.Helper()
	old := os.Args
	os.Args = append([]string{"prog"}, args...)
	t.Cleanup(func() { os.Args = old })
}

func TestParseCommandLine(t *testing.T) {
	setArgs(t, "-h", "-n", "bob", "-c", "0x10", "-r", "2.5")
	var c cfg
	if err := ParseCommandLine(&c); err != nil {
		t.Fatal(err)
	}
	if !c.Help || c.Name != "bob" || c.Count != 16 || c.Rate != 2.5 {
		t.Errorf("got %+v", c)
	}
}

func TestParseCommandLineDefaults(t *testing.T) {
	setArgs(t)
	var c cfg
	if err := ParseCommandLine(&c); err != nil {
		t.Fatal(err)
	}
	if c.Help || c.Name != "anon" || c.Count != 3 || c.Rate != 1.5 {
		t.Errorf("got %+v", c)
	}
}

func TestParseCommandLineErrors(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{"not a flag", []string{"x"}, "invalid flag"},
		{"double dash", []string{"--n", "x"}, "bad flag syntax"},
		{"unknown", []string{"-z"}, "not defined"},
		{"missing value", []string{"-n"}, "no value"},
		{"empty value", []string{"-n", ""}, ""},
		{"value is flag", []string{"-n", "-h"}, "no value"},
		{"bad int", []string{"-c", "abc"}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setArgs(t, tt.args...)
			var c cfg
			err := ParseCommandLine(&c)
			if tt.name == "empty value" {
				if err != nil {
					t.Errorf("unexpected error %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Errorf("err = %v, want containing %q", err, tt.want)
			}
		})
	}
}

func TestRequired(t *testing.T) {
	setArgs(t)
	var c reqCfg
	if err := ParseCommandLine(&c); err == nil {
		t.Error("expected required error")
	}
	setArgs(t, "-t", "go")
	if err := ParseCommandLine(&c); err != nil || c.Lang != "go" {
		t.Errorf("err=%v cfg=%+v", err, c)
	}
}

func TestInvalidConfig(t *testing.T) {
	setArgs(t)
	type noTag struct {
		A string
	}
	type badType struct {
		A []string `flag:"a"`
	}
	type badDefault struct {
		A int `flag:"a" default:"x"`
	}
	type badFloat struct {
		A float64 `flag:"a" default:"x"`
	}
	type unexported struct {
		a string `flag:"a"`
	}
	for name, v := range map[string]any{
		"non-pointer":           cfg{},
		"nil":                   nil,
		"pointer to non-struct": new(int),
		"no tag":                &noTag{},
		"bad type":              &badType{},
		"bad int default":       &badDefault{},
		"bad float default":     &badFloat{},
		"unexported":            &unexported{},
	} {
		if err := ParseCommandLine(v); err == nil {
			t.Errorf("%s: ParseCommandLine expected error", name)
		}
		if err := ParseEnvironment(v); err == nil {
			t.Errorf("%s: ParseEnvironment expected error", name)
		}
	}
}

func TestParseEnvironment(t *testing.T) {
	t.Setenv("h", "true")
	t.Setenv("n", "alice")
	t.Setenv("c", "7")
	t.Setenv("r", "0.25")
	var c cfg
	if err := ParseEnvironment(&c); err != nil {
		t.Fatal(err)
	}
	if !c.Help || c.Name != "alice" || c.Count != 7 || c.Rate != 0.25 {
		t.Errorf("got %+v", c)
	}
}

func TestParseEnvironmentDefaultsAndBoolFalse(t *testing.T) {
	for _, k := range []string{"h", "n", "c", "r"} {
		if v, ok := os.LookupEnv(k); ok {
			os.Unsetenv(k)
			t.Cleanup(func() { os.Setenv(k, v) })
		}
	}
	var c cfg
	if err := ParseEnvironment(&c); err != nil {
		t.Fatal(err)
	}
	if c.Help || c.Name != "anon" || c.Count != 3 {
		t.Errorf("defaults not kept: %+v", c)
	}
	t.Setenv("h", "false")
	if err := ParseEnvironment(&c); err != nil || c.Help {
		t.Errorf("err=%v help=%v", err, c.Help)
	}
}

func TestParseEnvironmentErrors(t *testing.T) {
	t.Setenv("c", "notanint")
	var c cfg
	if err := ParseEnvironment(&c); err == nil {
		t.Error("expected parse error")
	}
	t.Setenv("c", "")
	if err := ParseEnvironment(&c); err == nil {
		t.Error("expected empty value error")
	}
	var r reqCfg
	os.Unsetenv("t")
	if err := ParseEnvironment(&r); err == nil {
		t.Error("expected required error")
	}
}

func TestUsage(t *testing.T) {
	setArgs(t)
	var c cfg
	if err := ParseCommandLine(&c); err != nil {
		t.Fatal(err)
	}
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	Usage()
	w.Close()
	os.Stdout = old
	buf := make([]byte, 4096)
	n, _ := r.Read(buf)
	if !strings.Contains(string(buf[:n]), "-n:    name") {
		t.Errorf("usage output: %q", buf[:n])
	}
}

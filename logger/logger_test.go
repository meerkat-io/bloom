package logger

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func capture(t *testing.T) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	oldW, oldF := log.Writer(), log.Flags()
	log.SetOutput(&buf)
	log.SetFlags(0)
	oldLevel := int(logLevel.Load())
	oldExit := exit
	t.Cleanup(func() {
		log.SetOutput(oldW)
		log.SetFlags(oldF)
		logLevel.Store(int32(oldLevel))
		exit = oldExit
	})
	return &buf
}

func TestSetLevel(t *testing.T) {
	capture(t)
	for _, bad := range []int{0, -1, CRITICAL + 1} {
		if err := SetLevel(bad); err == nil {
			t.Errorf("SetLevel(%d) expected error", bad)
		}
	}
	for l := DEBUG; l <= CRITICAL; l++ {
		if err := SetLevel(l); err != nil {
			t.Errorf("SetLevel(%d): %v", l, err)
		}
	}
}

func TestLevelFiltering(t *testing.T) {
	buf := capture(t)
	calls := []struct {
		level int
		tag   string
		plain func(...any)
		fmtd  func(string, ...any)
	}{
		{DEBUG, "[DEBUG]", Debug, Debugf},
		{INFO, "[INFO]", Info, Infof},
		{NOTICE, "[NOTICE]", Notice, Noticef},
		{WARN, "[WARN]", Warn, Warnf},
		{ERROR, "[ERROR]", Error, Errorf},
	}
	for cur := DEBUG; cur <= CRITICAL; cur++ {
		if err := SetLevel(cur); err != nil {
			t.Fatal(err)
		}
		for _, c := range calls {
			buf.Reset()
			c.plain("msg")
			c.fmtd("f%d", 1)
			out := buf.String()
			want := c.level >= cur
			if got := strings.Contains(out, c.tag); got != want {
				t.Errorf("level=%d call=%s logged=%v want %v", cur, c.tag, got, want)
			}
			if want && (!strings.Contains(out, "msg") || !strings.Contains(out, "f1")) {
				t.Errorf("missing content in %q", out)
			}
		}
	}
}

func TestPosition(t *testing.T) {
	buf := capture(t)
	_ = SetLevel(DEBUG)
	Info("where")
	if !strings.Contains(buf.String(), "logger_test.go:") {
		t.Errorf("expected caller position in %q", buf.String())
	}
	buf.Reset()
	Infof("where %s", "f")
	if !strings.Contains(buf.String(), "logger_test.go:") {
		t.Errorf("expected caller position in %q", buf.String())
	}
}

func TestCritical(t *testing.T) {
	buf := capture(t)
	code := -1
	exit = func(c int) { code = c }
	_ = SetLevel(CRITICAL)
	Critical("boom")
	if code != 1 || !strings.Contains(buf.String(), "boom") {
		t.Errorf("code=%d out=%q", code, buf.String())
	}
	code = -1
	Criticalf("boom %d", 2)
	if code != 1 || !strings.Contains(buf.String(), "boom 2") {
		t.Errorf("code=%d out=%q", code, buf.String())
	}
}

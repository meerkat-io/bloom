package time

import (
	"testing"
	"time"
)

func TestTimestamp(t *testing.T) {
	before := time.Now().Unix()
	got := Timestamp()
	after := time.Now().Unix()
	if got < before || got > after {
		t.Errorf("Timestamp() = %d, want within [%d, %d]", got, before, after)
	}
}

func TestUnix(t *testing.T) {
	got := Unix(1700000000)
	if got.Unix() != 1700000000 || got.Nanosecond() != 0 {
		t.Errorf("Unix = %v", got)
	}
}

func TestElapse(t *testing.T) {
	got := Elapse(Timestamp() - 10)
	if got < 10*time.Second || got > 12*time.Second {
		t.Errorf("Elapse = %v, want ~10s", got)
	}
	if got%time.Second != 0 {
		t.Errorf("Elapse = %v, want whole seconds", got)
	}
	if got := Elapse(Timestamp() + 100); got >= 0 {
		t.Errorf("future timestamp should give negative duration, got %v", got)
	}
}

func TestExpired(t *testing.T) {
	now := Timestamp()
	if !Expired(now - 10) {
		t.Error("past timestamp should be expired")
	}
	if Expired(now + 100) {
		t.Error("future timestamp should not be expired")
	}
}

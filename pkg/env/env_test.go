package env

import (
	"github.com/jpfourny/papaya/pkg/pair"
	"github.com/jpfourny/papaya/pkg/stream"
	"os"
	"testing"
	"time"
)

func TestSet(t *testing.T) {
	Set("foo", "bar")

	if os.Getenv("foo") != "bar" {
		t.Errorf("expected os.Getenv(%q) to return %q; got %q", "foo", "bar", os.Getenv("foo"))
	}
}

func TestUnset(t *testing.T) {
	_ = os.Setenv("foo", "bar")
	Unset("foo")

	if os.Getenv("foo") != "" {
		t.Errorf("expected os.Getenv(%q) to return %q; got %q", "foo", "", os.Getenv("foo"))
	}
}

func TestSetAllPairs(t *testing.T) {
	_ = os.Setenv("foo", "bar")

	undo := SetAllPairs(
		pair.Of("foo", "bar2"),
		pair.Of("baz", "qux"),
	)

	if os.Getenv("foo") != "bar2" {
		t.Errorf("expected os.Getenv(%q) to return %q; got %q", "foo", "bar2", os.Getenv("foo"))
	}
	if os.Getenv("baz") != "qux" {
		t.Errorf("expected os.Getenv(%q) to return %q; got %q", "baz", "qux", os.Getenv("baz"))
	}

	undo()

	if os.Getenv("foo") != "bar" {
		t.Errorf("expected os.Getenv(%q) to return %q; got %q", "foo", "bar", os.Getenv("foo"))
	}
	if os.Getenv("baz") != "" {
		t.Errorf("expected os.Getenv(%q) to return %q; got %q", "baz", "", os.Getenv("baz"))
	}
}

func TestSetAllMap(t *testing.T) {
	_ = os.Setenv("foo", "bar")

	undo := SetAllMap(map[string]string{
		"foo": "bar2",
		"baz": "qux",
	})

	if os.Getenv("foo") != "bar2" {
		t.Errorf("expected os.Getenv(%q) to return %q; got %q", "foo", "bar2", os.Getenv("foo"))
	}
	if os.Getenv("baz") != "qux" {
		t.Errorf("expected os.Getenv(%q) to return %q; got %q", "baz", "qux", os.Getenv("baz"))
	}

	undo()

	if os.Getenv("foo") != "bar" {
		t.Errorf("expected os.Getenv(%q) to return %q; got %q", "foo", "bar", os.Getenv("foo"))
	}
	if os.Getenv("baz") != "" {
		t.Errorf("expected os.Getenv(%q) to return %q; got %q", "baz", "", os.Getenv("baz"))
	}
}

func TestGet(t *testing.T) {
	_ = os.Setenv("foo", "bar")

	got := Get("foo")
	if !got.Present() || got.GetOrZero() != "bar" {
		t.Errorf("expected Get(%q) to return %q; got %v", "foo", "bar", got)
	}

	got = Get("baz")
	if got.Present() {
		t.Errorf("expected Get(%q) to return empty opt; got %v", "baz", got)
	}
}

func TestGetBool(t *testing.T) {
	_ = os.Setenv("foo", "true")

	got := GetBool("foo")
	if !got.Present() || !got.GetOrZero() {
		t.Errorf("expected GetBool(%q) to return true; got %v", "foo", got)
	}

	_ = os.Setenv("foo", "false")

	got = GetBool("foo")
	if !got.Present() || got.GetOrZero() {
		t.Errorf("expected GetBool(%q) to return false; got %v", "foo", got)
	}

	got = GetBool("baz")
	if got.Present() {
		t.Errorf("expected GetBool(%q) to return empty opt; got %v", "baz", got)
	}

	_ = os.Setenv("foo", "not-a-bool")

	got = GetBool("foo")
	if got.Present() {
		t.Errorf("expected GetBool(%q) to return empty opt; got %v", "foo", got)
	}
}

func TestGetInt(t *testing.T) {
	_ = os.Setenv("foo", "42")

	got := GetInt[int]("foo")
	if !got.Present() || got.GetOrZero() != 42 {
		t.Errorf("expected GetInt(%q) to return 42; got %v", "foo", got)
	}

	got = GetInt[int]("baz")
	if got.Present() {
		t.Errorf("expected GetInt(%q) to return empty opt; got %v", "baz", got)
	}

	_ = os.Setenv("foo", "not-an-int")

	got = GetInt[int]("foo")
	if got.Present() {
		t.Errorf("expected GetInt(%q) to return empty opt; got %v", "foo", got)
	}
}

func TestGetUInt(t *testing.T) {
	_ = os.Setenv("foo", "42")

	got := GetUInt[uint]("foo")
	if !got.Present() || got.GetOrZero() != 42 {
		t.Errorf("expected GetUInt(%q) to return 42; got %v", "foo", got)
	}

	got = GetUInt[uint]("baz")
	if got.Present() {
		t.Errorf("expected GetUInt(%q) to return empty opt; got %v", "baz", got)
	}

	_ = os.Setenv("foo", "not-an-uint")

	got = GetUInt[uint]("foo")
	if got.Present() {
		t.Errorf("expected GetUInt(%q) to return empty opt; got %v", "foo", got)
	}
}

func TestGetFloat(t *testing.T) {
	_ = os.Setenv("foo", "42.42")

	got := GetFloat[float64]("foo")
	if !got.Present() || got.GetOrZero() != 42.42 {
		t.Errorf("expected GetFloat(%q) to return 42.42; got %v", "foo", got)
	}

	got = GetFloat[float64]("baz")
	if got.Present() {
		t.Errorf("expected GetFloat(%q) to return empty opt; got %v", "baz", got)
	}

	_ = os.Setenv("foo", "not-a-float")

	got = GetFloat[float64]("foo")
	if got.Present() {
		t.Errorf("expected GetFloat(%q) to return empty opt; got %v", "foo", got)
	}
}

func TestGetDuration(t *testing.T) {
	_ = os.Setenv("foo", "1s")

	got := GetDuration("foo")
	if !got.Present() || got.GetOrZero() != time.Second {
		t.Errorf("expected GetDuration(%q) to return 1s; got %v", "foo", got)
	}

	got = GetDuration("baz")
	if got.Present() {
		t.Errorf("expected GetDuration(%q) to return empty opt; got %v", "baz", got)
	}

	_ = os.Setenv("foo", "not-a-duration")

	got = GetDuration("foo")
	if got.Present() {
		t.Errorf("expected GetDuration(%q) to return empty opt; got %v", "foo", got)
	}
}

func TestGetTime(t *testing.T) {
	_ = os.Setenv("foo", "2006-01-02T15:04:05Z")

	got := GetTime("foo", time.RFC3339)
	if !got.Present() || got.GetOrZero().Format(time.RFC3339) != "2006-01-02T15:04:05Z" {
		t.Errorf("expected GetTime(%q, %q) to return 2006-01-02T15:04:05Z; got %v", "foo", time.RFC3339, got)
	}

	got = GetTime("baz", time.RFC3339)
	if got.Present() {
		t.Errorf("expected GetTime(%q, %q) to return empty opt; got %v", "baz", time.RFC3339, got)
	}

	_ = os.Setenv("foo", "not-a-time")

	got = GetTime("foo", time.RFC3339)
	if got.Present() {
		t.Errorf("expected GetTime(%q, %q) to return empty opt; got %v", "foo", time.RFC3339, got)
	}
}

func TestGetTimeInLocation(t *testing.T) {
	_ = os.Setenv("foo", "2006-01-02T15:04:05Z")

	got := GetTimeInLocation("foo", time.RFC3339, time.UTC)
	if !got.Present() || got.GetOrZero().Format(time.RFC3339) != "2006-01-02T15:04:05Z" {
		t.Errorf("expected GetTimeInLocation(%q, %q, %v) to return 2006-01-02T15:04:05Z; got %v", "foo", time.RFC3339, time.UTC, got)
	}

	got = GetTimeInLocation("baz", time.RFC3339, time.UTC)
	if got.Present() {
		t.Errorf("expected GetTimeInLocation(%q, %q, %v) to return empty opt; got %v", "baz", time.RFC3339, time.UTC, got)
	}

	_ = os.Setenv("foo", "not-a-time")

	got = GetTimeInLocation("foo", time.RFC3339, time.UTC)
	if got.Present() {
		t.Errorf("expected GetTimeInLocation(%q, %q, %v) to return empty opt; got %v", "foo", time.RFC3339, time.UTC, got)
	}
}

func TestToStream(t *testing.T) {
	_ = os.Setenv("foo", "bar")
	_ = os.Setenv("baz", "qux")

	got := ToStream()
	m := stream.CollectMap(got)

	if m["foo"] != "bar" {
		t.Errorf("expected ToStream() to contain %q; got %v", "foo", m["foo"])
	}
	if m["baz"] != "qux" {
		t.Errorf("expected ToStream() to contain %q; got %v", "baz", m["baz"])
	}
}

func TestToMap(t *testing.T) {
	_ = os.Setenv("foo", "bar")
	_ = os.Setenv("baz", "qux")

	got := ToMap()

	if got["foo"] != "bar" {
		t.Errorf("expected ToMap() to contain %q; got %v", "foo", got["foo"])
	}
	if got["baz"] != "qux" {
		t.Errorf("expected ToMap() to contain %q; got %v", "baz", got["baz"])
	}
}

package pred

import (
	"testing"
)

func TestStringContains(t *testing.T) {
	p := StringContains("foo")
	if !p("foo") {
		t.Errorf("StringContains(\"foo\")(\"foo\") = false; want true")
	}
	if !p("foobar") {
		t.Errorf("StringContains(\"foo\")(\"foobar\") = false; want true")
	}
	if !p("barfoo") {
		t.Errorf("StringContains(\"foo\")(\"barfoo\") = false; want true")
	}
	if !p("barfoobar") {
		t.Errorf("StringContains(\"foo\")(\"barfoobar\") = false; want true")
	}
	if p("bar") {
		t.Errorf("StringContains(\"foo\")(\"bar\") = true; want false")
	}
}

func TestStringContainsAny(t *testing.T) {
	p := StringContainsAny("abc")
	if !p("a") {
		t.Errorf("StringContainsAny(\"abc\")(\"a\") = false; want true")
	}
	if !p("b") {
		t.Errorf("StringContainsAny(\"abc\")(\"b\") = false; want true")
	}
	if !p("c") {
		t.Errorf("StringContainsAny(\"abc\")(\"c\") = false; want true")
	}
	if !p("ab") {
		t.Errorf("StringContainsAny(\"abc\")(\"ab\") = false; want true")
	}
	if !p("ac") {
		t.Errorf("StringContainsAny(\"abc\")(\"ac\") = false; want true")
	}
	if p("d") {
		t.Errorf("StringContainsAny(\"abc\")(\"d\") = true; want false")
	}
}

func TestStringContainsRune(t *testing.T) {
	p := StringContainsRune('a')
	if !p("a") {
		t.Errorf("StringContainsRune('a')(\"a\") = false; want true")
	}
	if p("b") {
		t.Errorf("StringContainsRune('a')(\"b\") = true; want false")
	}
	if !p("ab") {
		t.Errorf("StringContainsRune('a')(\"ab\") = false; want true")
	}
	if p("c") {
		t.Errorf("StringContainsRune('a')(\"c\") = true; want false")
	}
}

func TestStringContainsFunc(t *testing.T) {
	p := StringContainsFunc(func(r rune) bool {
		return r == 'a'
	})
	if !p("a") {
		t.Errorf("StringContainsFunc(func(r rune) bool { return r == 'a' })(\"a\") = false; want true")
	}
	if p("b") {
		t.Errorf("StringContainsFunc(func(r rune) bool { return r == 'a' })(\"b\") = true; want false")
	}
	if !p("ab") {
		t.Errorf("StringContainsFunc(func(r rune) bool { return r == 'a' })(\"ab\") = false; want true")
	}
	if p("c") {
		t.Errorf("StringContainsFunc(func(r rune) bool { return r == 'a' })(\"c\") = true; want false")
	}
}

func TestStringHasPrefix(t *testing.T) {
	p := StringHasPrefix("foo")
	if !p("foo") {
		t.Errorf("StringHasPrefix(\"foo\")(\"foo\") = false; want true")
	}
	if !p("foobar") {
		t.Errorf("StringHasPrefix(\"foo\")(\"foobar\") = false; want true")
	}
	if p("barfoo") {
		t.Errorf("StringHasPrefix(\"foo\")(\"barfoo\") = true; want false")
	}
	if p("barfoobar") {
		t.Errorf("StringHasPrefix(\"foo\")(\"barfoobar\") = true; want false")
	}
	if p("bar") {
		t.Errorf("StringHasPrefix(\"foo\")(\"bar\") = true; want false")
	}
}

func TestStringHasSuffix(t *testing.T) {
	p := StringHasSuffix("foo")
	if !p("foo") {
		t.Errorf("StringHasSuffix(\"foo\")(\"foo\") = false; want true")
	}
	if p("foobar") {
		t.Errorf("StringHasSuffix(\"foo\")(\"foobar\") = true; want false")
	}
	if !p("barfoo") {
		t.Errorf("StringHasSuffix(\"foo\")(\"barfoo\") = false; want true")
	}
	if p("barfoobar") {
		t.Errorf("StringHasSuffix(\"foo\")(\"barfoobar\") = true; want false")
	}
	if p("bar") {
		t.Errorf("StringHasSuffix(\"foo\")(\"bar\") = true; want false")
	}
}

func TestStringEqualFold(t *testing.T) {
	p := StringEqualFold("foo")
	if !p("foo") {
		t.Errorf("StringEqualFold(\"foo\")(\"foo\") = false; want true")
	}
	if !p("FOO") {
		t.Errorf("StringEqualFold(\"foo\")(\"FOO\") = false; want true")
	}
	if p("bar") {
		t.Errorf("StringEqualFold(\"foo\")(\"bar\") = true; want false")
	}
}

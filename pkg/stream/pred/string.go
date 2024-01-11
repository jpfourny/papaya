package pred

import (
	"strings"
)

// StringContains returns a function that returns true if the provided string contains the provided substring.
// It uses the strings.Contains function to compare strings.
//
// Examples:
//
//	p := pred.StringContains("foobar")
//	p("foo") // true
//	p("bar") // true
//	p("baz") // false
func StringContains(substr string) func(string) bool {
	return func(e string) bool {
		return strings.Contains(e, substr)
	}
}

// StringContainsAny returns a function that returns true if the provided string contains any of the provided chars.
// It uses the strings.ContainsAny function to compare strings.
//
// Examples:
//
//	p := pred.StringContainsAny("abc")
//	p("foo") // false
//	p("bar") // true
//	p("baz") // true
func StringContainsAny(chars string) func(string) bool {
	return func(e string) bool {
		return strings.ContainsAny(e, chars)
	}
}

// StringContainsRune returns a function that returns true if the provided string contains the provided rune.
// It uses the strings.ContainsRune function to compare strings.
//
// Examples:
//
//	p := pred.StringContainsRune('a')
//	p("foo") // false
//	p("bar") // true
//	p("baz") // true
func StringContainsRune(r rune) func(string) bool {
	return func(e string) bool {
		return strings.ContainsRune(e, r)
	}
}

// StringContainsFunc returns a function that returns true if the provided string contains a rune that satisfies the provided function.
// It uses the strings.ContainsFunc function to compare strings.
//
// Examples:
//
//	p := pred.StringContainsFunc(func(r rune) bool { return r == 'a' })
//	p("foo") // false
//	p("bar") // true
//	p("baz") // true
func StringContainsFunc(f func(rune) bool) func(string) bool {
	return func(e string) bool {
		return strings.ContainsFunc(e, f)
	}
}

// StringHasPrefix returns a function that returns true if the provided string has the provided prefix.
// It uses the strings.HasPrefix function to compare strings.
//
// Examples:
//
//	p := pred.StringHasPrefix("foo")
//	p("foobar") // true
//	p("fuzbar") // false
func StringHasPrefix(prefix string) func(string) bool {
	return func(e string) bool {
		return strings.HasPrefix(e, prefix)
	}
}

// StringHasSuffix returns a function that returns true if the provided string has the provided suffix.
// It uses the strings.HasSuffix function to compare strings.
//
// Examples:
//
//	p := pred.StringHasSuffix("bar")
//	p("foobar") // true
//	p("foobaz") // false
func StringHasSuffix(suffix string) func(string) bool {
	return func(e string) bool {
		return strings.HasSuffix(e, suffix)
	}
}

// StringEqualFold returns a function that returns true if the provided string is equal to the provided want string, ignoring case.
// It uses the strings.EqualFold function to compare strings.
//
// Examples:
//
//	p := pred.StringEqualFold("foo")
//	p("foo") // true
//	p("FOO") // true
//	p("bar") // false
func StringEqualFold(s string) func(string) bool {
	return func(e string) bool {
		return strings.EqualFold(e, s)
	}
}

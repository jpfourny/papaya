package env

import (
	"github.com/jpfourny/papaya/pkg/constraint"
	"github.com/jpfourny/papaya/pkg/opt"
	"github.com/jpfourny/papaya/pkg/pair"
	"github.com/jpfourny/papaya/pkg/stream"
	"github.com/jpfourny/papaya/pkg/stream/mapper"
	"os"
	"strings"
	"time"
)

// Set sets the environment variable with the given key to the given value.
func Set(key, value string) {
	_ = os.Setenv(key, value)
}

// Unset unsets the environment variable with the given key.
func Unset(key string) {
	_ = os.Unsetenv(key)
}

// SetAllPairs sets the environment variables to the given pairs.
// Returns a function that can be called to revert the changes.
func SetAllPairs(pairs ...pair.Pair[string, string]) (revert func()) {
	return setAll(stream.Of(pairs...))
}

// SetAllMap sets the environment variables to the given map.
// Returns a function that can be called to revert the changes.
func SetAllMap(m map[string]string) (revert func()) {
	return setAll(stream.FromMap(m))
}

func setAll(vars stream.Stream[pair.Pair[string, string]]) (revert func()) {
	backupVars := make(map[string]string)
	var addedKeys []string

	stream.ForEach(vars, func(p pair.Pair[string, string]) {
		if prev, ok := os.LookupEnv(p.First()); ok {
			backupVars[p.First()] = prev // Backup existing key before overwriting.
		} else {
			addedKeys = append(addedKeys, p.First()) // Keep track of added keys.
		}
		Set(p.First(), p.Second())
	})

	return func() {
		for k, v := range backupVars {
			Set(k, v)
		}
		for _, k := range addedKeys {
			Unset(k)
		}
	}
}

// Get returns the value of the environment variable with the given key, if it exists.
func Get(key string) opt.Optional[string] {
	return opt.Maybe(os.LookupEnv(key))
}

// GetBool returns the value of the environment variable with the given key, if it exists and can be parsed as a boolean.
// An empty Optional is returned if the variable is unset or if value cannot be parsed as a boolean.
func GetBool(key string) opt.Optional[bool] {
	return opt.OptionalMap(
		Get(key),
		mapper.TryParseBool[string](),
	)
}

// GetInt returns the value of the environment variable with the given key, if it exists and can be parsed as an integer of the desired type I.
// An empty Optional is returned if the variable is unset or if value cannot be parsed as an integer.
func GetInt[I constraint.SignedInteger](key string) opt.Optional[I] {
	return opt.OptionalMap(
		Get(key),
		mapper.TryParseInt[string, I](10, 64),
	)
}

// GetUInt returns the value of the environment variable with the given key, if it exists and can be parsed as an unsigned integer of the desired type I.
// An empty Optional is returned if the variable is unset or if value cannot be parsed as an unsigned integer.
func GetUInt[I constraint.UnsignedInteger](key string) opt.Optional[I] {
	return opt.OptionalMap(
		Get(key),
		mapper.TryParseUint[string, I](10, 64),
	)
}

// GetFloat returns the value of the environment variable with the given key, if it exists and can be parsed as a float of the desired type F.
// An empty Optional is returned if the variable is unset or if value cannot be parsed as a float.
func GetFloat[F constraint.Float](key string) opt.Optional[F] {
	return opt.OptionalMap(
		Get(key),
		mapper.TryParseFloat[string, F](64),
	)
}

// GetDuration returns the value of the environment variable with the given key, if it exists and can be parsed as a duration.
// An empty Optional is returned if the variable is unset or if value cannot be parsed as a duration.
// See time.ParseDuration for details on the expected format.
func GetDuration(key string) opt.Optional[time.Duration] {
	return opt.OptionalMap(
		Get(key),
		mapper.TryParseDuration[string](),
	)
}

// GetTime returns the value of the environment variable with the given key, if it exists and can be parsed as a time.Time with the given layout.
// An empty Optional is returned if the variable is unset or if value cannot be parsed as a time.Time.
// See time.Parse for details on the expected format.
func GetTime(key string, layout string) opt.Optional[time.Time] {
	return opt.OptionalMap(
		Get(key),
		mapper.TryParseTime[string](layout),
	)
}

// GetTimeInLocation returns the value of the environment variable with the given key, if it exists and can be parsed as a time.Time with the given layout and location.
// An empty Optional is returned if the variable is unset or if value cannot be parsed as a time.Time.
// See time.ParseInLocation for details on the expected format.
func GetTimeInLocation(key string, layout string, loc *time.Location) opt.Optional[time.Time] {
	return opt.OptionalMap(
		Get(key),
		mapper.TryParseTimeInLocation[string](layout, loc),
	)
}

// ToStream returns a stream of pairs representing the environment variables.
func ToStream() stream.Stream[pair.Pair[string, string]] {
	return stream.Map(
		stream.FromSlice(os.Environ()),
		func(s string) pair.Pair[string, string] {
			parts := strings.SplitN(s, "=", 2)
			return pair.Of(parts[0], parts[1])
		},
	)
}

// ToMap returns a map representing the environment variables.
func ToMap() map[string]string {
	return stream.CollectMap(ToStream())
}

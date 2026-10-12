package rule

import (
	"go/token"
	"math"
	"regexp"
	"strings"

	"github.com/mgechev/revive/internal/config"
	"github.com/mgechev/revive/lint"
)

func srcLine(src []byte, p token.Position) string {
	// Run to end of line in both directions if not at line start/end.
	lo, hi := p.Offset, p.Offset+1
	for lo > 0 && src[lo-1] != '\n' {
		lo--
	}
	for hi < len(src) && src[hi-1] != '\n' {
		hi++
	}
	return string(src[lo:hi])
}

// isRuleOption returns true if arg and name are the same after normalization.
func isRuleOption(arg, name string) bool {
	return config.NormalizeOption(arg) == config.NormalizeOption(name)
}

var normalizePathReplacer = strings.NewReplacer("-", "", "_", "", ".", "")

// normalizePath removes hyphens, underscores, and dots from the name
//
// Example: normalizePath("foo.bar-_buz") -> "foobarbuz".
func normalizePath(name string) string {
	return normalizePathReplacer.Replace(name)
}

// isVersionPath checks if a directory name is a version directory (v1, V2, etc.)
func isVersionPath(name string) bool {
	if len(name) < 2 || (name[0] != 'v' && name[0] != 'V') {
		return false
	}

	for i := 1; i < len(name); i++ {
		if name[i] < '0' || name[i] > '9' {
			return false
		}
	}

	return true
}

var directiveCommentRE = regexp.MustCompile("^//(line |extern |export |[a-z0-9]+:[a-z0-9])") // see https://go-review.googlesource.com/c/website/+/442516/1..2/_content/doc/comment.md#494

func isDirectiveComment(line string) bool {
	return directiveCommentRE.MatchString(line)
}

// newInternalFailureError returns a slice of Failure with a single internal failure in it.
func newInternalFailureError(e error) []lint.Failure {
	return []lint.Failure{lint.NewInternalFailure(e.Error())}
}

// integer is the set of integer types a rule argument can hold.
type integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// float is the set of floating-point types a rule argument can hold.
type float interface {
	~float32 | ~float64
}

// toInt converts a rule argument to an int.
//
// A rule cannot expect a single Go type for a numeric argument, because the type depends on
// the decoder that produced the configuration:
//   - a TOML decoder, the one revive itself uses, yields int64 for every integer;
//   - a YAML decoder yields int;
//   - a JSON decoder yields float64, as JSON has no integer type;
//   - a tool embedding revive through revivelib can pass any numeric type it holds.
//
// golangci-lint, for instance, accepts its configuration in all three formats and passes the
// decoded arguments straight to the rules. So a numeric argument is accepted whatever its Go
// type, as long as it holds an integral value that fits in an int.
func toInt(arg any) (int, bool) {
	switch v := arg.(type) {
	case int:
		return integerToInt(v)
	case int8:
		return integerToInt(v)
	case int16:
		return integerToInt(v)
	case int32:
		return integerToInt(v)
	case int64:
		return integerToInt(v)
	case uint:
		return integerToInt(v)
	case uint8:
		return integerToInt(v)
	case uint16:
		return integerToInt(v)
	case uint32:
		return integerToInt(v)
	case uint64:
		return integerToInt(v)
	case float32:
		return floatToInt(v)
	case float64:
		return floatToInt(v)
	default:
		return 0, false
	}
}

// integerToInt converts an integer to an int, reporting whether it is representable as such.
func integerToInt[T integer](v T) (int, bool) {
	i := int(v)

	truncated := T(i) != v
	signFlipped := (i < 0) != (v < 0) // an unsigned value above math.MaxInt round-trips, but changes sign
	if truncated || signFlipped {
		return 0, false
	}

	return i, true
}

// floatToInt converts a float holding an integral value to an int, reporting whether it is representable as such.
func floatToInt[T float](v T) (int, bool) {
	// math.MaxInt is not exactly representable as a float, so the upper bound is -math.MinInt,
	// the smallest power of two above it.
	const maxIntExclusive = -math.MinInt

	f := float64(v)
	if f != math.Trunc(f) || f < math.MinInt || f >= maxIntExclusive {
		return 0, false
	}

	return int(f), true
}

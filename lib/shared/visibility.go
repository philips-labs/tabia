package shared

import "strings"

//go:generate stringer -type=Visibility

// Visibility indicates repository visibility
type Visibility int

const (
	// Public repositories are publicly visible
	Public Visibility = iota
	// Internal repositories are only visible to organization members
	Internal
	// Private repositories are only visible to authorized users
	Private
)

func (v *Visibility) UnmarshalText(text []byte) error {
	*v = VisibilityFromText(string(text))
	return nil
}

func (v Visibility) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

func VisibilityFromText(text string) Visibility {
	switch strings.ToLower(text) {
	default:
		return Public
	case "public":
		return Public
	case "internal":
		return Internal
	case "private":
		return Private
	}
}

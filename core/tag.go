// shared interfaces and utility functions
package core

// interface for an tag markup
type Tag interface {
	AddCode(code Code)
	AddTag(tag Tag) (bool, bool)
	Name() string
	IsEnd() bool
	Code
}

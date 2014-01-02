package core

type Container interface {
	AddToken(token Token)
	AddTag(tag Tag) (bool, bool)
}

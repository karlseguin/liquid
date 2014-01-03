package core

// interface for something that can render itself
type Code interface {
	Render(data map[string]interface{}) []byte
}

package core

type Token interface {
	Render(data map[string]interface{}) []byte
}

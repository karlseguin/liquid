package core

import (
	"github.com/karlseguin/bytepool"
	"io"
)

var BytePool = bytepool.New(128, 4096)

// cache interface
type Cache interface {
	Get(key string) Code
	Set(key string, template Code)
	Clear()
}

// The callback to execute to resolve include tags
// If you're going to use name to read from the filesystem, beware of directory
// traversal.
type IncludeHandler func(name string, writer io.Writer, data map[string]interface{})

// Configuration used for generating a template
type Configuration struct {
	cache              Cache
	includeHandler     IncludeHandler
	preserveWhitespace bool
}

// Set the caching engine, or nil for no caching
func (c *Configuration) Cache(cache Cache) *Configuration {
	c.cache = cache
	return c
}

// Gets the configured cache
func (c *Configuration) GetCache() Cache {
	return c.cache
}

// Set the include handler
func (c *Configuration) IncludeHandler(handler IncludeHandler) *Configuration {
	c.includeHandler = handler
	return c
}

// Gets the configured include handler
func (c *Configuration) GetIncludeHandler() IncludeHandler {
	return c.includeHandler
}

// Preserves whitespace
func (c *Configuration) PreserveWhitespace() *Configuration {
	c.preserveWhitespace = true
	return c
}

// Gets the preserves whitespace value
func (c *Configuration) GetPreserveWhitespace() bool {
	return c.preserveWhitespace
}

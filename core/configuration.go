package core

import (
	"github.com/karlseguin/bytepool"
	"io"
)

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
	bytepool           *bytepool.Pool
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

// Occasionally, Liquid needs to create temporary buffers (supporting the
// capture tag, for example). It uses a fixed-length byte pool. You can control
// the number of buffers to keep in the pool as well as the maximum size of each
// item. By default, 512 items are kept with a maximum of 4KB. If you expect
// large captures, you should increase the size parameter
func (c *Configuration) SetInternalBuffer(count, size int) *Configuration {
	c.bytepool = bytepool.New(count, size)
	return c
}

// Gets the writer provider
func (c *Configuration) GetWriter() *bytepool.Item {
	return c.bytepool.Checkout()
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

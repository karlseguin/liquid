package core

// cache interface
type Cache interface {
	Get(key string) Code
	Set(key string, template Code)
	Clear()
}

type Configuration struct {
	cache          Cache
	includeHandler IncludeHandler
}

// The callback to execute to resolve include tags
// If you're going to use name to read from the filesystem, beware of directory
// traversal.
type IncludeHandler func(name string, data map[string]interface{}) []byte

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

//gets the configured include handler
func (c *Configuration) GetIncludeHandler() IncludeHandler {
	return c.includeHandler
}

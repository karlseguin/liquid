package liquid

type Configuration struct {
	cache Cache
}

var (
	defaultConfig = Configure()
	//A Configuration with caching disabled
	NoCache = Configure().Cache(nil)
)

// Entry into the fluent-configuration
func Configure() *Configuration {
	return &Configuration{
		cache: TemplateCache,
	}
}

// Set the caching engine, or nil for no caching
func (c *Configuration) Cache(cache Cache) *Configuration {
	c.cache = cache
	return c
}

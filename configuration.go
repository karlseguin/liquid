package liquid

type Configuration struct {
	cache Cache
}

var defaultConfig = Configure()

func Configure() *Configuration {
	return &Configuration{
		cache: TemplateCache,
	}
}

func (c *Configuration) Cache(cache Cache) *Configuration {
	c.cache = cache
	return c
}

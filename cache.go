package liquid

import(
	"sync"
)

type Cache interface {
	Get(key string) *Template
	Set(key string, template *Template)
	Clear()
}
var TemplateCache = &SimpleCache{lookup: make(map[string]*Template)}

type SimpleCache struct {
	sync.RWMutex
	lookup map[string]*Template
}

func (c *SimpleCache) Get(key string) *Template {
	c.RLock()
	defer c.RUnlock()
	return c.lookup[key]
}

func (c *SimpleCache) Set(key string, template *Template) {
	c.Lock()
	defer c.Unlock()
	c.lookup[key] = template
}

func (c *SimpleCache) Clear() {
	c.Lock()
	defer c.Unlock()
	c.lookup = make(map[string]*Template)
}

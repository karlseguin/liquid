package liquid

import (
	"github.com/karlseguin/liquid/core"
	"sync"
)

var TemplateCache = &SimpleCache{lookup: make(map[string]core.Code)}

type SimpleCache struct {
	sync.RWMutex
	lookup map[string]core.Code
}

func (c *SimpleCache) Get(key string) core.Code {
	c.RLock()
	defer c.RUnlock()
	return c.lookup[key]
}

func (c *SimpleCache) Set(key string, template core.Code) {
	c.Lock()
	defer c.Unlock()
	c.lookup[key] = template
}

func (c *SimpleCache) Clear() {
	c.Lock()
	defer c.Unlock()
	c.lookup = make(map[string]core.Code)
}

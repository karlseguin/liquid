package liquid

import (
	"github.com/karlseguin/bytepool"
	"github.com/karlseguin/liquid/core"
)

var (
	defaultConfig = Configure()
	//A Configuration with caching disabled
	NoCache = Configure().Cache(nil)
)

// Entry into the fluent-configuration
func Configure() *core.Configuration {
	c := new(core.Configuration)
	return c.Cache(TemplateCache)
}

// Set's the count and size of the internal bytepool
func SetInternalBuffer(count, size int) {
	core.BytePool = bytepool.New(count, size)
}

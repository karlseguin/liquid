package core

import (
	"time"
)

var (
	Now = func() time.Time { return time.Now() }
)

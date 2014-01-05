package filters

import (
	"github.com/karlseguin/gspec"
	"github.com/karlseguin/liquid/core"
	"testing"
	"time"
)

func init() {
	core.Now = func() time.Time {
		t, _ := time.Parse("Mon Jan 02 15:04:05 2006", "Mon Jan 02 15:04:05 2006")
		return t
	}
}

func TestDateNowWithBasicFormat(t *testing.T) {
	spec := gspec.New(t)
	filter := DateFactory([]core.Value{stringValue("%Y %m %d")})
	spec.Expect(filter("now", nil).(string)).ToEqual("2006 01 02")
}

func TestDateTodayWithBasicFormat(t *testing.T) {
	spec := gspec.New(t)
	filter := DateFactory([]core.Value{stringValue("%H:%M:%S%%")})
	spec.Expect(filter("today", nil).(string)).ToEqual("15:04:05%")
}

func TestDateWithSillyFormat(t *testing.T) {
	spec := gspec.New(t)
	filter := DateFactory([]core.Value{stringValue("%w  %U  %j")})
	spec.Expect(filter("2014-01-10 21:31:28 +0800", nil).(string)).ToEqual("05  02  10")
}

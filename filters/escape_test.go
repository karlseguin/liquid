package filters

import (
	"github.com/karlseguin/gspec"
	"testing"
)

func TestEscapesAString(t *testing.T) {
	spec := gspec.New(t)
	filter := EscapeFactory(nil)
	spec.Expect(filter("<script>hack</script>", nil).(string)).ToEqual("&lt;script&gt;hack&lt;/script&gt;")
}

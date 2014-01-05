package filters

import (
	"fmt"
	"github.com/karlseguin/liquid/core"
	"time"
)

var (
	zeroTime = time.Time{}
)

// Creates an date filter
func DateFactory(parameters []core.Value) core.Filter {
	if len(parameters) == 0 {
		return Noop
	}
	return (&DateFilter{parameters[0]}).ToString
}

type DateFilter struct {
	format core.Value
}

func (d *DateFilter) ToString(input interface{}, data map[string]interface{}) interface{} {
	time, ok := inputToTime(input)
	if ok == false {
		return input
	}
	return formatTime(time, core.ToString(d.format.Resolve(data)))
}

func inputToTime(input interface{}) (time.Time, bool) {
	switch typed := input.(type) {
	case time.Time:
		return typed, true
	case string:
		return timeFromString(typed)
	case []byte:
		return timeFromString(string(typed))
	}
	if n, ok := core.ToInt(input); ok {
		return core.Now().Add(time.Minute * time.Duration(n)), true
	}
	return zeroTime, false
}

func timeFromString(s string) (time.Time, bool) {
	if s == "now" || s == "today" {
		return core.Now(), true
	}
	t, err := time.Parse("2006-01-02 15:04:05 -0700", s)
	if err != nil {
		return zeroTime, false
	}
	return t, true
}

func formatTime(t time.Time, ruby string) string {
	l := len(ruby) - 1
	format := make([]byte, 0, l*2)
	for i := 0; i < l; i++ {
		if ruby[i] != '%' {
			format = append(format, ruby[i])
			continue
		}
		n := ruby[i+1]
		if n == '%' {
			format = append(format, '%')
		} else {
			format = append(format, convertTimeFormat(t, n)...)
		}
		i += 1
	}
	return string(format)
}

func convertTimeFormat(t time.Time, n byte) string {
	switch n {
	case 'a':
		return t.Weekday().String()[:3]
	case 'A':
		return t.Weekday().String()
	case 'b':
		return t.Month().String()[:3]
	case 'B':
		return t.Month().String()
	case 'c':
		return t.Format("ANSIC")
	case 'd':
		return fmt.Sprintf("%02d", t.Day())
	case 'H':
		return fmt.Sprintf("%02d", t.Hour())
	case 'I':
		hr := t.Hour() % 12
		if hr == 0 {
			hr = 12
		}
		return fmt.Sprintf("%02d", hr)
	case 'm':
		return fmt.Sprintf("%02d", t.Month())
	case 'M':
		return fmt.Sprintf("%02d", t.Minute())
	case 'p':
		if t.Hour() > 11 {
			return "PM"
		}
		return "AM"
	case 'S':
		return fmt.Sprintf("%02d", t.Second())
	case 'x':
		return t.Format("Mon Jan 02")
	case 'X':
		return t.Format("15:04:05")
	case 'y':
		year := fmt.Sprintf("%04d", t.Year())
		if l := len(year); l > 2 {
			return year[l-2 : l]
		}
		return year
	case 'Y':
		return fmt.Sprintf("%04d", t.Year())
	case 'j':
		return fmt.Sprintf("%02d", t.YearDay())
	case 'w':
		return fmt.Sprintf("%02d", t.Weekday())
	case 'U':
		_, w := t.ISOWeek()
		return fmt.Sprintf("%02d", w)
	case 'W':
		_, w := t.ISOWeek()
		return fmt.Sprintf("%02d", w)
	default:
		return string(n)
	}
}

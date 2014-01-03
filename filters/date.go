package filters

import (
	"github.com/karlseguin/liquid/core"
	"time"
	"sync"
	"strconv"
)

var (
	now = func() time.Time { return time.Now() }
	zeroTime = time.Time{}
	timeFormatCache = make(map[string]string)
	timeFormatCacheLock sync.RWMutex
	timePartLookup = map[byte]string{
		'a': "Mon",
		'A': "Monday",
		'b': "Jan",
		'B': "January",
		'c': "ANSIC", //fail
		'd': "02",
		'H': "15",
		'I': "03",
		'm': "01",
		'M': "04",
		'p': "PM",
		'S': "05",
		'x': "Mon Jan 02",
		'X': "15:04:05",
		'y': "06",
		'Y': "2006",
		'Z': "MST",
	}
)

// Creates an date filter
func DateFactory(parameters []core.Value) Filter {
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
	format := alignFormat(time, core.ToString(d.format.Resolve(data)))
	return time.Format(format)
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
		return now().Add(time.Minute * time.Duration(n)), true
	}
	return zeroTime, false
}

func timeFromString(s string) (time.Time, bool) {
	if s == "now" || s == "today" {
		return now(), true
	}
	t, err := time.Parse("2006-01-02 15:04:05 -0700", s)
	if err != nil {
		return zeroTime, false
	}
	return t, true
}

func alignFormat(t time.Time, ruby string) string {
	timeFormatCacheLock.RLock()
	format, exists := timeFormatCache[ruby]
	timeFormatCacheLock.RUnlock()
	if exists == false {
		format = buildTimeFormat(t, ruby)
		timeFormatCacheLock.Lock()
		defer timeFormatCacheLock.Unlock()
		if len(timeFormatCache) > 500 {
			timeFormatCache = make(map[string]string)
		}
		timeFormatCache[ruby] = format
	}
	return format
}

func buildTimeFormat(t time.Time, ruby string) string {
	l := len(ruby) - 1
	format := make([]byte, 0, l * 5)
	for i := 0; i < l; i++ {
		if ruby[i] != '%' {
			format = append(format, ruby[i])
			continue
		}
		n := ruby[i+1]
		if n == '%' {
			format = append(format, '%')
		} else {
			alt, ok := timePartLookup[n]
			if ok == false {
				alt, ok = timePartCalculate(t, n)
			}
			if ok {
				b := []byte(alt)
				for j := 0; j < len(b); j++ {
					format = append(format, b[j])
				}
			}
		}
		i += 1
	}
	return string(format)
}

func timePartCalculate(t time.Time, n byte) (string, bool) {
	switch n {
	case 'j':
		return strconv.Itoa(t.YearDay()), true
	case 'w':
		return strconv.Itoa(int(t.Weekday())), true
	}
	return "", false
}

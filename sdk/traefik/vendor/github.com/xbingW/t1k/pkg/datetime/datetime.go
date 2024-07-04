package datetime

import "time"

func Now() int64 {
	return time.Now().UnixNano() / 1e3
}

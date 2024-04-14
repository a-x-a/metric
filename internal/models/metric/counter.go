package metric

import (
	"strconv"
)

type (
	Counter int64
)

func (c Counter) Kind() string {
	return string(KindCounter)
}

func (c Counter) String() string {
	return strconv.FormatInt(int64(c), 10)
}

func (c Counter) IsCounter() bool {
	return true
}

func (c Counter) IsGauge() bool {
	return false
}

func ToCounter(value string) (Counter, error) {
	val, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}

	return Counter(val), nil
}

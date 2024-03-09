package metric

import "fmt"

type (
	// Counter interface {
	// 	Get() int64
	// 	Set(v int64) error

	// 	Inc() error
	// 	Dec() error
	// }

	Counter int64
)

func (c Counter) Kind() string {
	return string(KindCounter)
}

func (c Counter) String() string {
	return fmt.Sprintf("%d", c)
}

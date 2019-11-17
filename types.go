package influx

import (
	"errors"

	influx "github.com/influxdata/influxdb1-client/v2"
)

type (
	// Tags alias of string/string map for tsdb tags
	Tags = map[string]string
	// Fields alias of string/interface{} map for tsdb fields
	Fields = map[string]interface{}
	// Point alias of influx.Point
	Point = influx.Point
	// Result alias of influx.Result
	Result = influx.Result
)

const (
	rfc3339     = "rfc3339"
	hour        = "h"
	minute      = "m"
	second      = "s"
	millisecond = "ms"
	microsecond = "us"
	nanosecond  = "ns"
)

// application defined errors
var (
	ErrNotImplemented = errors.New("not implemented")
	ErrEmptyResults   = errors.New("empty results")
)

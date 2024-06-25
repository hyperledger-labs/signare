package tracer

import "fmt"

const (
	traceContextVersion contextVersion = 0

	traceParentLength = 55
	traceIDLength     = 32
	spanIDLength      = 16

	traceParentHeader = "traceparent"
	traceStateHeader  = "tracestate"
)

type contextVersion uint

func (v contextVersion) String() string {
	return fmt.Sprintf("%02d", v)
}

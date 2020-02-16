package jobs

import (
	"strings"

	"github.com/phomer/scheduler/log"
)

type TimeScale int

var (
	SEC  TimeScale = 1
	MIN  TimeScale = 2
	HOUR TimeScale = 4
	DAY  TimeScale = 8
)

var lookup = map[string]TimeScale{
	"sec":   SEC,
	"secs":  SEC,
	"min":   MIN,
	"mins":  MIN,
	"hour":  HOUR,
	"hours": HOUR,
	"day":   HOUR,
	"days":  HOUR,
}

var ConvertScale = map[TimeScale]int{
	SEC:  1,
	MIN:  60,
	HOUR: 60 * 60,
	DAY:  24 * 60 * 60,
}

func LookupTimeScale(value string) *TimeScale {
	scale, ok := lookup[strings.ToLower(value)]
	if !ok {
		return nil
	}
	return &scale
}

func RelativeUnixTime(number int, scale *TimeScale) int64 {
	return int64(number * ConvertScale[*scale])
}

func AbsoluteUnixTime(base int64, number int, scale *TimeScale) int64 {
	log.Dump(number, scale)
	return base + int64(number*ConvertScale[*scale])
}

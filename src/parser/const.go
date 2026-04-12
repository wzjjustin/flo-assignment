package parser

import "time"

const (
	HEADER_RECORD_TYPE  string = "100"
	NMI_DATA_TYPE       string = "200"
	INTERVAL_DATA_TYPE  string = "300"
	INTERVAL_EVENT_TYPE string = "400" //currently unused
	B2B_DETAILS_TYPE    string = "500" //currently unused
	END_OF_DATA         string = "900"

	IntervalLengthPT5M  = time.Duration(5 * time.Minute)
	IntervalLengthPT15M = time.Duration(15 * time.Minute)
	IntervalLengthPT30M = time.Duration(30 * time.Minute)

	minutesIn1Day = 60 * 24 // 1440
)

var IntervalLengths = map[int]time.Duration{ //nolint:gochecknoglobals
	5:  IntervalLengthPT5M,
	15: IntervalLengthPT15M,
	30: IntervalLengthPT30M,
}

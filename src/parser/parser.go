package parser

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseRecord(recordType, data string) Record {
	// Remove trailing newline and split into fields
	cleanData := strings.TrimSpace(data)
	fields := strings.Split(cleanData, ",")
	return Record{Type: recordType, Data: fields}
}

func ParseNMIRecord(data []string) (nmi NMIRecord, err error) {
	// Expected format: 200,NMI,E1E2,1,E1,N1,01009,kWh,30,20050610
	// Field 1: NMI
	// Field 8: Interval length
	// Field 9: Next scheduled read
	var possibleIntervals = map[int]struct{}{
		5:  {},
		15: {},
		30: {},
	}

	if len(data) < 10 {
		return nmi, fmt.Errorf("insufficient data - NMI record")
	}

	nmi.NMI = data[1]
	intervalValue, err := strconv.Atoi(data[8])
	if err != nil {
		return nmi, err
	}
	if _, found := possibleIntervals[intervalValue]; !found {
		return nmi, fmt.Errorf("invalid interval value %d", intervalValue)
	}

	nmi.IntervalLength = intervalValue
	nmi.NextScheduledRead = data[9]
	// validation

	return nmi, err
}

func ParseIntervalRecord(intervalValue int, data []string) (IntervalRecord, error) {
	// Expected format: 300,YYYYMMDD,value1,value2,...,valueN,status,...
	rec := IntervalRecord{}
	statusIdx := 1440/intervalValue + 2
	if len(data) < statusIdx {
		return rec, fmt.Errorf("invalid number of interval values")
	}
	if len(data) > 1 {
		rec.Date = data[1]
	}

	// Collect interval values (fields 2 through~N)
	if len(data) > 2 {
		if len(data) > statusIdx {
			rec.Values = data[2:statusIdx]
			if statusIdx < len(data) {
				rec.Status = data[statusIdx]
			}
		} else {
			if len(data) > 2 {
				rec.Values = data[2:]
			}
		}
	}

	return rec, nil
}

func IsRecordStart(field []string) bool {
	if len(field) > 0 {
		switch field[0] {
		case END_OF_DATA,
			HEADER_RECORD_TYPE,
			NMI_DATA_TYPE,
			INTERVAL_DATA_TYPE:
			/* currently unused
			INTERVAL_EVENT_TYPE,
			B2B_DETAILS_TYPE:
			*/
			return true
		}
	}
	return false
}

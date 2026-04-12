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

var possibleIntervals = map[int]struct{}{
	5:  {},
	15: {},
	30: {},
}

func ParseNMIRecord(data []string) (nmi NMIRecord, err error) {
	// Expected format: 200,NMI,E1E2,1,E1,N1,01009,kWh,30,20050610
	// Field 1: NMI
	// Field 8: Interval length

	// data length validation
	if len(data) < 10 {
		return nmi, fmt.Errorf("insufficient data - NMI record")
	}

	// nmi field validation
	if len(data[1]) > 10 || data[1] == "" {
		return nmi, fmt.Errorf("invalid NMI field data")
	}

	// interval length validation
	intervalValue, err := strconv.Atoi(data[8])
	if err != nil {
		return nmi, err
	}

	if _, found := possibleIntervals[intervalValue]; !found {
		return nmi, fmt.Errorf("invalid interval value %d", intervalValue)
	}

	return NMIRecord{
		IntervalLength: intervalValue,
		NMI:            data[1],
	}, err
}

func calculateStatusIdx(intervalValue int) int {
	return minutesIn1Day/intervalValue + 2
}

func ParseIntervalRecord(intervalValue int, data []string) (IntervalRecord, error) {
	// Expected format: 300,YYYYMMDD,value1,value2,...,valueN,status,...
	rec := IntervalRecord{}
	if _, found := possibleIntervals[intervalValue]; !found {
		return IntervalRecord{}, fmt.Errorf("invalid interval value %d", intervalValue)
	}

	statusIdx := calculateStatusIdx(intervalValue)
	if len(data) < statusIdx {
		return IntervalRecord{}, fmt.Errorf("invalid number of interval values")
	}

	// Collect interval values (fields 2 through~N)
	if len(data) > statusIdx {
		rec.Values = data[2:statusIdx]
		if statusIdx < len(data) {
			rec.QualityMethod = data[statusIdx]
		}
	} else {
		rec.Values = data[2:]
	}

	rec.Date = data[1]

	return rec, nil
}

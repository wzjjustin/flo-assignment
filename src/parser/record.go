package parser

type Record struct {
	Type string // 100, 200, 300, 400, 500, 900
	Data []string
}

type HeaderRecord struct {
	Version   string
	Timestamp string
	FromParty string
}

type NMIRecord struct {
	NMI          string   // Field 1: NMI identifier
	UnitOfMeasure string   // Field 7: Unit of measure (e.g., kWh)
	IntervalLength int    // Field 8: Interval length in minutes
	NextScheduledRead string // Field 9: Next scheduled read date
}

type IntervalRecord struct {
	Date  string   // Field 1: Date of reading (YYYYMMDD)
	Values []string // Fields 2-49: Consumption values for each interval
	Status string   // Field 50: Status flag (A, M, etc.)
}

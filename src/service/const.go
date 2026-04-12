package service

import "flo-assignment/src/parser"

const (
	YYYYMMDDformat = "20060102"
)

type Void struct{}

// using struct{} instead of boolean as struct{} is zero-byte
var validBlockCycleOrder = map[string]map[string]Void{
	parser.HEADER_RECORD_TYPE: { //100
		"": {},
	},
	parser.NMI_DATA_TYPE: { //200
		parser.HEADER_RECORD_TYPE: {}, //100
		parser.B2B_DETAILS_TYPE:   {}, //500
	},
	parser.INTERVAL_DATA_TYPE: { //300
		parser.NMI_DATA_TYPE:      {}, //200
		parser.INTERVAL_DATA_TYPE: {}, //300
	},
	parser.INTERVAL_EVENT_TYPE: { //400
		parser.INTERVAL_DATA_TYPE: {}, //300
	},
	parser.B2B_DETAILS_TYPE: { //500
		parser.INTERVAL_DATA_TYPE:  {}, //300
		parser.INTERVAL_EVENT_TYPE: {}, //400
	},
	parser.END_OF_DATA: { //900
		parser.B2B_DETAILS_TYPE: {}, //500
	},
}

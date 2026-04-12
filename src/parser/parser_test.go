package parser

import (
	mockdata "flo-assignment/src/mock_data"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseNMIRecord(t *testing.T) {
	type testCases struct {
		name        string
		input       []string
		expected    NMIRecord
		expectedErr error
	}

	testCategories := []struct {
		category  string
		testCases []testCases
	}{
		{
			category: "happy flow",
			testCases: []testCases{
				{
					name:     "happy flow - interval 5",
					input:    []string{"200", "NEM1201009", "E1E2", "1", "E1", "N1", "01009", "kWh", "5", "20050610"},
					expected: NMIRecord{NMI: "NEM1201009", IntervalLength: 5},
				},
				{
					name:     "happy flow - interval 15",
					input:    []string{"200", "NEM1201010", "E1E2", "1", "E1", "N1", "01009", "kWh", "15", "20050610"},
					expected: NMIRecord{NMI: "NEM1201010", IntervalLength: 15},
				},
				{
					name:     "happy flow - interval 30",
					input:    []string{"200", "NEM1201011", "E1E2", "1", "E1", "N1", "01009", "kWh", "30", "20050610"},
					expected: NMIRecord{NMI: "NEM1201011", IntervalLength: 30},
				},
			},
		},
		{
			category: "insufficient data",
			testCases: []testCases{
				{
					name:        "missing NextScheduledReadDate data",
					input:       []string{"200", "NEM1201009", "E1E2", "1", "E1", "N1", "01009", "kWh", "5"},
					expected:    NMIRecord{},
					expectedErr: fmt.Errorf("insufficient data - NMI record"),
				},
				{
					name:        "missing multiple data",
					input:       []string{"200"},
					expected:    NMIRecord{},
					expectedErr: fmt.Errorf("insufficient data - NMI record"),
				},
			},
		},
		{
			category: "invalid data",
			testCases: []testCases{
				{
					name:        "nmi field too long",
					input:       []string{"200", "NEM1234567890", "E1E2", "1", "E1", "N1", "01009", "kWh", "5", "20050610"},
					expected:    NMIRecord{},
					expectedErr: fmt.Errorf("invalid NMI field data"),
				},
				{
					name:        "empty nmi field",
					input:       []string{"200", "", "E1E2", "1", "E1", "N1", "01009", "kWh", "5", "20050610"},
					expected:    NMIRecord{},
					expectedErr: fmt.Errorf("invalid NMI field data"),
				},
				{
					name:        "invalid interval value",
					input:       []string{"200", "NEM1201010", "E1E2", "1", "E1", "N1", "01009", "kWh", "1", "20050610"},
					expected:    NMIRecord{},
					expectedErr: fmt.Errorf("invalid interval value 1"),
				},
			},
		},
	}

	for _, tc := range testCategories {
		for _, tt := range tc.testCases {
			t.Run("", func(t *testing.T) {
				result, err := ParseNMIRecord(tt.input)
				if result != tt.expected {
					t.Errorf("got %v, want %v", result, tt.expected)
				}
				if tt.expectedErr != nil && err.Error() != tt.expectedErr.Error() {
					t.Errorf("got %v, want %v", err, tt.expectedErr)
				}
			})
		}
	}
}

func TestParseIntervalRecord(t *testing.T) {
	type intervalRecordInput struct {
		intervalValue int
		data          []string
	}
	type testCases struct {
		name        string
		input       intervalRecordInput
		expected    IntervalRecord
		expectedErr error
	}

	testCategories := []struct {
		category  string
		testCases []testCases
	}{
		{
			category: "happy flow",
			testCases: []testCases{
				{
					name: "interval 5",
					input: intervalRecordInput{
						intervalValue: 5,
						data:          mockdata.MockParseIntervalRecordInput5Data,
					},
					expected: IntervalRecord{
						Date:          "20050301",
						Values:        mockdata.MockParseIntervalRecordInput5Data[2:calculateStatusIdx(5)],
						QualityMethod: mockdata.MockParseIntervalRecordInput5Data[calculateStatusIdx(5)],
					},
				},
				{
					name: "interval 15",
					input: intervalRecordInput{
						intervalValue: 15,
						data:          mockdata.MockParseIntervalRecordInput15Data,
					},
					expected: IntervalRecord{
						Date:          "20050301",
						Values:        mockdata.MockParseIntervalRecordInput15Data[2:calculateStatusIdx(15)],
						QualityMethod: mockdata.MockParseIntervalRecordInput15Data[calculateStatusIdx(15)],
					},
				},
				{
					name: "interval 30",
					input: intervalRecordInput{
						intervalValue: 30,
						data:          mockdata.MockParseIntervalRecordInput30Data,
					},
					expected: IntervalRecord{
						Date:          "20050301",
						Values:        mockdata.MockParseIntervalRecordInput30Data[2:calculateStatusIdx(30)],
						QualityMethod: mockdata.MockParseIntervalRecordInput30Data[calculateStatusIdx(30)],
					},
				},
			},
		},
		{
			category: "invalid data",
			testCases: []testCases{
				{
					name: "invalid interval value",
					input: intervalRecordInput{
						intervalValue: 0,
						data:          mockdata.MockParseIntervalRecordInput5Data,
					},
					expected:    IntervalRecord{},
					expectedErr: fmt.Errorf("invalid interval value 0"),
				},
				{
					name: "invalid interval value",
					input: intervalRecordInput{
						intervalValue: 15,
						data:          mockdata.MockParseIntervalRecordInputInvalidData,
					},
					expected:    IntervalRecord{},
					expectedErr: fmt.Errorf("invalid number of interval values"),
				},
			},
		},
	}

	for _, tc := range testCategories {
		for _, tt := range tc.testCases {
			t.Run("", func(t *testing.T) {
				result, err := ParseIntervalRecord(tt.input.intervalValue, tt.input.data)
				if tt.expectedErr != nil && err.Error() != tt.expectedErr.Error() {
					t.Errorf("got %v, want %v", err, tt.expectedErr)
				}
				assert.New(t)
				assert.Equal(t, result, tt.expected)
			})
		}
	}
}

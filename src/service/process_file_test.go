package service

import (
	mockdata "flo-assignment/src/mock_data"
	"flo-assignment/src/parser"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
)

var checksum = "123"

func (t *TestSuite) TestLoadDataToDB_Happy_Flow() {
	results := make(chan parseResult, 100)
	input := [][]string{mockdata.MockParseNMIDataSplit, mockdata.MockParseIntervalRecordInput30Data}

	for i, rec := range input {
		results <- parseResult{index: i, rec: parser.Record{Type: rec[0], Data: rec}}
	}

	close(results)

	t.mockDB.ExpectBegin()
	t.mockDB.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "processed_files" ("created_at","updated_at","deleted_at","checksum") VALUES ($1,$2,$3,$4)`)).WillReturnRows(&sqlmock.Rows{})
	t.mockDB.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "meter_readings" ("created_at","updated_at","deleted_at","nmi","timestamp","consumption") VALUES ($1,$2,$3,$4,$5,$6)`)).WillReturnRows(&sqlmock.Rows{})
	t.mockDB.ExpectCommit()

	err := t.service.LoadDataToDB(checksum, results)
	t.NoError(err)
}

func (t *TestSuite) TestLoadDataToDB_Invalid_Data() {
	checksum := "123"
	results := make(chan parseResult, 100)

	input := [][]string{mockdata.MockParseNMIDataSplit, mockdata.MockParseIntervalRecordInputInvalidData}

	for i, rec := range input {
		results <- parseResult{index: i, rec: parser.Record{Type: rec[0], Data: rec}}
	}

	close(results)

	t.mockDB.ExpectBegin()
	t.mockDB.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "processed_files" ("created_at","updated_at","deleted_at","checksum") VALUES ($1,$2,$3,$4)`)).WillReturnRows(&sqlmock.Rows{})
	t.mockDB.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "meter_readings" ("created_at","updated_at","deleted_at","nmi","timestamp","consumption") VALUES ($1,$2,$3,$4,$5,$6)`)).WillReturnRows(&sqlmock.Rows{})

	err := t.service.LoadDataToDB(checksum, results)
	t.ErrorContains(err, "invalid number of interval values")
}

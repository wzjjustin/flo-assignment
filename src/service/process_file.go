package service

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"flo-assignment/src/model"
	"flo-assignment/src/parser"
	"fmt"
	"hash"
	"os"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"
)

type parseTask struct {
	index      int
	recordType string
	raw        string
}

type parseResult struct {
	index int
	rec   parser.Record
	err   error
}

// worker's method to handle parsing of records
func (s *Service) parseWorker(tasks <-chan parseTask, results chan<- parseResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range tasks {
		rec := parser.ParseRecord(task.recordType, task.raw)
		results <- parseResult{index: task.index, rec: rec}
	}
}

func (s *Service) processIntervalRecord(tx *gorm.DB, nmi string, intervalLength int, rec parser.IntervalRecord) error {
	t, err := time.Parse(YYYYMMDDformat, rec.Date)
	if err != nil {
		return fmt.Errorf("failed to parse date in %s format: %v", YYYYMMDDformat, err)
	}

	intervalIncrement := parser.IntervalLengths[intervalLength]

	toCreate := make([]model.MeterReading, len(rec.Values))
	for i := range rec.Values {
		toCreate[i] = model.MeterReading{
			NMI:         nmi,
			Timestamp:   t,
			Consumption: rec.Values[i],
		}
		t = t.Add(intervalIncrement)
	}

	tx = tx.CreateInBatches(toCreate, len(rec.Values))
	if tx.Error != nil {
		return fmt.Errorf("failed to create in batch: %v", tx.Error)
	}

	return nil
}

func addToWriter(currentRecord *strings.Builder, line string, h hash.Hash) error {
	if _, err := currentRecord.WriteString(line); err != nil {
		return fmt.Errorf("failed to write to string builder: %v", err)
	}

	if _, err := h.Write([]byte(line)); err != nil {
		return fmt.Errorf("failed to write to hash: %v", err)
	}

	return nil
}

func (s *Service) ProcessFileWithWorkers(ctx context.Context, path string) error {
	var (
		numWorkers = s.serviceCfg.NumWorkers
		tasks      = make(chan parseTask, 100)
		results    = make(chan parseResult, 100)
		wg         sync.WaitGroup
	)

	// initiate worker pool
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go s.parseWorker(tasks, results, &wg)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	// open file
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}

	defer file.Close() //nolint:errcheck

	var (
		scanner       = bufio.NewScanner(file)
		h             = sha256.New()
		currentRecord strings.Builder
		recordType    string
		recordIndex   int
	)

	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	// read each line till EOF
	for scanner.Scan() { // defaulted to ScanLine()
		line := scanner.Text()
		fields := strings.Split(line, ",")
		if len(fields) > 0 {
			// validate block cycle order
			prevValidTypes, isRecord := validBlockCycleOrder[fields[0]]
			if !isRecord { // if new token is not a record, add to current record and continue
				if err := addToWriter(&currentRecord, line, h); err != nil {
					return err
				}

				continue
			}

			if _, isPrevValid := prevValidTypes[recordType]; !isPrevValid { // if new old record type does not path into new record type
				return fmt.Errorf("blocking order is incorrect, new type = %s, want = %v", fields[0], prevValidTypes)
			}

			// submit data for processing if any
			if currentRecord.Len() > 0 { // send preprocessed record for processing via worker
				tasks <- parseTask{index: recordIndex, recordType: recordType, raw: currentRecord.String()}
				recordIndex++
			}

			// clean up current record to be added in addToWriter() outside of if-statment
			recordType = fields[0]
			currentRecord.Reset()
		}

		// append current record
		if err := addToWriter(&currentRecord, line, h); err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan error: %v", err)
	}
	// send 'end' data for processing (900 type)
	if currentRecord.Len() > 0 {
		tasks <- parseTask{index: recordIndex, recordType: recordType, raw: currentRecord.String()}
	}

	close(tasks) // closing tasks channel, indicating no new task data will be sent.

	return s.LoadDataToDB(hex.EncodeToString(h.Sum(nil)), results)
}

func (s *Service) LoadDataToDB(checkSum string, results <-chan parseResult) error {
	pending := make(map[int]parser.Record)
	nextIndex := 0
	currentNMI := ""
	currentIntervalLength := 0

	err := s.db.Transaction(func(tx *gorm.DB) error {
		// create a file-processed record to ensure no duplicated processing of same file
		createResult := tx.Create(&model.ProcessedFile{
			Checksum: checkSum,
		})
		if errors.Is(createResult.Error, gorm.ErrDuplicatedKey) {
			return fmt.Errorf("file was processed before: %v", createResult.Error)
		}

		// check for channel buffer for results from processing
		for res := range results {
			if res.err != nil {
				return res.err
			}
			// perform index to maintain block of 200-500 order
			pending[res.index] = res.rec

			for { // process in order of index to ensure nmi is correct for records
				rec, ok := pending[nextIndex]
				if !ok {
					break
				}

				switch rec.Type {
				case parser.NMI_DATA_TYPE:
					nmi, err := parser.ParseNMIRecord(rec.Data)
					if err != nil {
						return err
					}
					currentNMI = nmi.NMI
					currentIntervalLength = nmi.IntervalLength
				case parser.INTERVAL_DATA_TYPE:
					intervalRecord, err := parser.ParseIntervalRecord(currentIntervalLength, rec.Data)
					if err != nil {
						return err
					}
					fmt.Printf("Ordered: Processing interval for NMI %s on date %s with %d values\n",
						currentNMI, intervalRecord.Date, len(intervalRecord.Values))
					if err := s.processIntervalRecord(tx, currentNMI, currentIntervalLength, intervalRecord); err != nil {
						return err
					}
				}

				delete(pending, nextIndex)
				nextIndex++
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("rollback: %v", err)
	}

	return nil
}

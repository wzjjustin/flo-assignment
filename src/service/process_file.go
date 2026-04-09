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
	"os"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"
)

func (s *Service) processIntervalRecord(tx *gorm.DB, nmi string, intervalLength int, rec parser.IntervalRecord) error {
	t, err := time.Parse("20060102", rec.Date)
	if err != nil {
		return fmt.Errorf("failed to parse date: %v", err)
	}

	toCreate := make([]model.MeterReading, len(rec.Values))
	for i := range rec.Values {
		toCreate[i] = model.MeterReading{
			NMI:         nmi,
			Timestamp:   t,
			Consumption: rec.Values[i],
		}
		t = t.Add(time.Minute * time.Duration(intervalLength))
	}
	tx = tx.CreateInBatches(toCreate, len(rec.Values))
	if tx.Error != nil {
		return fmt.Errorf("failed to create in batch: %v", tx.Error)
	}

	return nil
}

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

func (s *Service) parseWorker(tasks <-chan parseTask, results chan<- parseResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range tasks {
		rec := parser.ParseRecord(task.recordType, task.raw)
		results <- parseResult{index: task.index, rec: rec}
	}
}

func (s *Service) ProcessFileWithWorkers(ctx context.Context, path string) error {
	var (
		numWorkers = s.serviceCfg.NumWorkers
		tasks      = make(chan parseTask, 100)
		results    = make(chan parseResult, 100)
		wg         sync.WaitGroup
	)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go s.parseWorker(tasks, results, &wg)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	var (
		scanner       = bufio.NewScanner(file)
		h             = sha256.New()
		currentRecord strings.Builder
		recordType    string
		recordIndex   int
		isHeaderFound bool
		isEndFound    bool
	)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, ",")
		if parser.IsRecordStart(fields) {
			switch fields[0] {
			case parser.HEADER_RECORD_TYPE:
				if isEndFound {
					return fmt.Errorf("header not at top of file")
				}
				isHeaderFound = true
			case parser.END_OF_DATA:
				if !isHeaderFound {
					return fmt.Errorf("header not at top of file")
				}
				isEndFound = true
			default:
				if !isHeaderFound {
					return fmt.Errorf("header not at top of file")
				}
			}
			if currentRecord.Len() > 0 {
				tasks <- parseTask{index: recordIndex, recordType: recordType, raw: currentRecord.String()}
				recordIndex++
			}
			recordType = fields[0]
			currentRecord.Reset()
		}
		currentRecord.WriteString(line)
		h.Write(scanner.Bytes())
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan error: %v", err)
	}

	if currentRecord.Len() > 0 {
		tasks <- parseTask{index: recordIndex, recordType: recordType, raw: currentRecord.String()}
		recordIndex++
	}

	close(tasks)

	pending := make(map[int]parser.Record)
	nextIndex := 0
	currentNMI := ""
	currentIntervalLength := 0
	checkSum := hex.EncodeToString(h.Sum(nil))
	err = s.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Create(&model.ProcessedFile{
			Checksum: checkSum,
		})
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return fmt.Errorf("file was processed before: %v", result.Error)
		}

		for res := range results {
			if res.err != nil {
				return res.err
			}
			pending[res.index] = res.rec

			for {
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
				default:
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

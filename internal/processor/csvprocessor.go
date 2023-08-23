package processor

import (
	"encoding/csv"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"io"
	"promotions-service/internal"
	"runtime"
	"sync"
)

type Repository interface {
	Insert(prs []internal.PromotionRecord) error
}

type CsvFile struct {
	repo       Repository
	linesPool  *sync.Pool
	NumWorkers int
	BatchSize  int
}

func NewCsvFile(repo Repository) *CsvFile {
	return &CsvFile{
		repo:       repo,
		NumWorkers: runtime.NumCPU(),
		BatchSize:  1000,
		linesPool: &sync.Pool{New: func() interface{} {
			lines := internal.PromotionRecord{}
			return lines
		}},
	}
}

// Process processes the csv file
// it work by batching line in BatchSize and processing
// each batch in a worker in worker pool
func (cf *CsvFile) Process(f io.Reader) error {
	r := csv.NewReader(f)

	var wg sync.WaitGroup

	cChan := make(chan [][]string, cf.NumWorkers)
	errChan := make(chan error, cf.NumWorkers)

	for i := 0; i < cf.NumWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cf.processChunk(cChan, errChan)
		}()
	}

	var records [][]string

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			break
		}
		records = append(records, record)
		if len(records) == cf.BatchSize {
			cChan <- records
			records = nil
		}
	}
	cChan <- records

	close(cChan)
	close(errChan)

	var errs error
	for err := range errChan {
		errs = multierror.Append(errs, err)
	}

	wg.Wait()

	return errs
}

// processChunk take one batch of lines
// and process and store that in storage
func (cf *CsvFile) processChunk(cChan <-chan [][]string, errChan chan error) {
	var err error

	for chunk := range cChan {
		promotions := make([]internal.PromotionRecord, len(chunk))
		for i, row := range chunk {
			if len(row) != 3 {
				continue
			}
			promotions[i] = internal.PromotionRecord{
				ID:             row[0],
				Price:          row[1],
				ExpirationDate: row[2],
			}
		}
		err = cf.repo.Insert(promotions)
		if err != nil {
			errChan <- err
		}
	}
}

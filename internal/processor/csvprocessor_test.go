package processor

import (
	"os"
	"promotions-service/internal/database"
	"testing"
)

func TestCsvFile_Process(t *testing.T) {

	db := database.NewInMemory()
	_ = db

	f, err := os.Open("../../data/promotions.csv")
	if err != nil {
		t.Errorf("error should be nil %s", err)
	}
	pr := NewCsvFile(db)
	pr.NumWorkers = 4
	err = pr.Process(f)

	db.SwitchPrimaryTable()
	if err != nil {
		t.Errorf("error should be nil %s", err)
	}

	if db.Len() != 199999 {
		t.Errorf("len should be %d, is %d", 199999, db.Len())
	}

	val := db.Get("6ecc8087-74ad-4e8c-a63d-7ea42cd4b811")

	if val == nil {
		t.Errorf("val should not be nil")
	}

}

func BenchmarkCsvFile_Process(t *testing.B) {

	for i := 0; i < t.N; i++ {
		db := database.NewInMemory()
		_ = db

		f, err := os.Open("../../data/promotions.csv")
		if err != nil {
			t.Errorf("error should be nil %s", err)
		}
		pr := NewCsvFile(db)
		pr.NumWorkers = 4
		err = pr.Process(f)
		db.SwitchPrimaryTable()
		if err != nil {
			t.Errorf("error should be nil %s", err)
		}

	}
}

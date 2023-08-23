package database

import (
	"promotions-service/internal"
	"strconv"
	"testing"
	"time"
)

func TestInMemory(t *testing.T) {

	newDB := NewInMemory()

	var input []internal.PromotionRecord
	for i := 0; i < 100; i++ {
		input = append(input, internal.PromotionRecord{
			ID:             strconv.Itoa(i),
			Price:          "",
			ExpirationDate: "",
		})
	}

	// will update secondary table
	newDB.Insert(input)

	// will read from primary table
	if newDB.Get("99") != nil {
		t.Errorf("Get() should return nil")
	}

	// will switch primary with secondary
	newDB.SwitchPrimaryTable()

	// will read from new primary table
	if newDB.Get("99") == nil {
		t.Errorf("Get() should not return nil")
	}

	// will delete secondary table
	newDB.PurgeSecondaryTable()

	// will switch primary with secondary again
	newDB.SwitchPrimaryTable()

	// will read from old primary table which is purged
	if newDB.Get("99") != nil {
		t.Errorf("Get() should return nil")
	}
}

func TestInMemory_Concurrent(t *testing.T) {

	newDB := NewInMemory()

	var input []internal.PromotionRecord
	for i := 0; i < 100; i++ {
		input = append(input, internal.PromotionRecord{
			ID:             strconv.Itoa(i),
			Price:          "",
			ExpirationDate: "",
		})
	}

	newDB.Insert(input)
	newDB.SwitchPrimaryTable()

	go func() {
		for i := 0; i < 100; i++ {
			tRecord := internal.PromotionRecord{
				ID:             strconv.Itoa(i),
				Price:          "",
				ExpirationDate: "",
			}
			time.Sleep(1 * time.Millisecond)
			newDB.Insert([]internal.PromotionRecord{tRecord})
		}
	}()

	for i := 0; i < 100; i++ {
		if newDB.Get(strconv.Itoa(i)) == nil {
			t.Errorf("Get() should not return nil")
		}
	}

}

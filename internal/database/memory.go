package database

import (
	"promotions-service/internal"
	"sync"
)

const (
	TABLE_A = "dataTableA"
	TABLE_B = "dataTableB"
)

type InMemory struct {
	mu          *sync.RWMutex
	activeTable string
	dataTableA  *sync.Map
	dataTableB  *sync.Map
}

type record struct {
	Price          string
	ExpirationDate string
}

func NewInMemory() *InMemory {
	mu := sync.RWMutex{}
	return &InMemory{
		mu:          &mu,
		activeTable: TABLE_A,
		dataTableA:  &sync.Map{},
		dataTableB:  &sync.Map{},
	}
}

// SwitchPrimaryTable makes the secondary table
// primary that will then serve the read request
func (db *InMemory) SwitchPrimaryTable() {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.activeTable == TABLE_A {
		db.activeTable = TABLE_B
	} else {
		db.activeTable = TABLE_A
	}
}

func (db *InMemory) getPrimaryTable() *sync.Map {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.activeTable == TABLE_A {
		return db.dataTableA
	} else {
		return db.dataTableB
	}
}

func (db *InMemory) getSecondaryTable() *sync.Map {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.activeTable == TABLE_A {
		return db.dataTableB
	} else {
		return db.dataTableA
	}
}

// Insert takes slice of PromotionRecord and insert
// them in storage
func (db *InMemory) Insert(prs []internal.PromotionRecord) error {
	data := db.getSecondaryTable()

	for _, pr := range prs {
		data.Store(pr.ID, record{
			Price:          pr.Price,
			ExpirationDate: pr.ExpirationDate,
		})
	}

	return nil
}

// Get takes an ID and return PromotionRecord
// return nil if record not found
func (db *InMemory) Get(id string) *internal.PromotionRecord {
	data := db.getPrimaryTable()
	val, ok := data.Load(id)
	if !ok {
		return nil
	}

	return &internal.PromotionRecord{
		ID:             id,
		Price:          (val.(record)).Price,
		ExpirationDate: (val.(record)).ExpirationDate,
	}
}

// PurgeSecondaryTable remove all data from table
// that is not in use
func (db *InMemory) PurgeSecondaryTable() {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.activeTable == TABLE_A {
		db.dataTableB = &sync.Map{}
	} else {
		db.dataTableA = &sync.Map{}
	}
}

// Len is helper function to get length of
// the primary table
func (db *InMemory) Len() int {
	t := db.getPrimaryTable()
	var i int
	t.Range(func(k, v interface{}) bool {
		i++
		return true
	})
	return i
}

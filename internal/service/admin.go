package service

import (
	"fmt"
	"os"
	"promotions-service/internal/processor"
)

type AdminRepository interface {
	SwitchPrimaryTable()
	PurgeSecondaryTable()
}
type Admin struct {
	repo         AdminRepository
	csvProcessor *processor.CsvFile
}

func NewAdmin(proc *processor.CsvFile, repo AdminRepository) Admin {
	return Admin{
		repo:         repo,
		csvProcessor: proc,
	}
}
func (ad *Admin) Process(file string) error {
	if isDirectory(file) {
		return fmt.Errorf("open file %s", "invalid file path")
	}
	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("open file %s", err)
	}
	return ad.csvProcessor.Process(f)
}
func (ad *Admin) SwitchStorage() error {
	ad.repo.SwitchPrimaryTable()
	ad.repo.PurgeSecondaryTable()
	return nil
}

func isDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fileInfo.IsDir()
}

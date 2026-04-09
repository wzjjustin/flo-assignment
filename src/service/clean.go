package service

import "fmt"

func (s *Service) CleanDB() error {
	result := s.db.Exec("DROP TABLE processed_files")
	if result.Error != nil {
		return result.Error
	}
	result = s.db.Exec("DROP TABLE meter_readings")
	if result.Error != nil {
		return result.Error
	}

	fmt.Println("Clean up successful!")
	return nil
}

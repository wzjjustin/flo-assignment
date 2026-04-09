package service

import (
	"flo-assignment/src/model"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Service struct {
	serviceCfg *Config
	db         *gorm.DB
}

func NewService(cfgPath string) (*Service, error) {
	cfgFile, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil, err
	}

	serviceConfig := &Config{}
	if err = yaml.Unmarshal(cfgFile, serviceConfig); err != nil {
		return nil, err
	}

	db, err := gorm.Open(postgres.Open(buildURI(serviceConfig.DBConfig)), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	result := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`) // to fix uuid_generate_v4() does not exist on postgresql
	if result.Error != nil {
		return nil, fmt.Errorf("failed to create extention: %v", result.Error)
	}

	if err = db.AutoMigrate(model.MeterReading{}, model.ProcessedFile{}); err != nil {
		return nil, err
	}

	return &Service{
		serviceCfg: serviceConfig,
		db:         db,
	}, nil
}

func buildURI(dbCfg DBConfig) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		dbCfg.Host, dbCfg.User, dbCfg.Password, dbCfg.DBname, dbCfg.Port)
}

package service

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TestSuite struct {
	suite.Suite
	service *Service
	mockDB  sqlmock.Sqlmock
}

func TestService(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (t *TestSuite) SetupSuite() {
	mockDb, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})
	t.service = &Service{db: db}
	t.mockDB = mock

}

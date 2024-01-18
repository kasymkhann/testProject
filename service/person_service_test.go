package service

import (
	"testProject/internal/model"
	"testProject/pkg/logging"
	"testProject/repository"
	"testing"

	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

type MockLogger struct {
	mock.Mock
}

func (m *MockRepository) CreatePerson(person *model.Person) error {
	args := m.Called(person)
	return args.Error(0)
}

func (m *MockLogger) Logger(message string) {
	m.Called(message)
}

func TestCreatePerson(t *testing.T) {
	repo := new(MockRepository)
	logger := new(MockLogger)

	repoWrapper := struct {
		repository.Repository
	}{
		Repository: repository.Repository{},
	}

	loggerWrapper := struct {
		logging.Logger
	}{
		Logger: logging.Logger{},
	}

	service := NewService(&repoWrapper.Repository, &loggerWrapper.Logger)

	testPerson := &model.Person{
		ID:          1,
		Name:        "TestName",
		Surname:     "TestSurname",
		Patronymic:  "TestPatronymic",
		Age:         22,
		Gender:      "ManGender",
		Nationality: "TestNationality",
	}

	repo.On("CreatePerson", testPerson).Return(nil)

	logger.On("Debug", mock.Anything)

	err := service.CreatePerson(testPerson)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	repo.AssertExpectations(t)
	logger.AssertExpectations(t)
}

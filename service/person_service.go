package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"testProject/internal/model"
	"testProject/pkg/logging"
	"testProject/repository"
)

// Service представляет собой сервис для работы с данными о людях.
type Service struct {
	repo   *repository.Repository
	logger *logging.Logger
}

// NewService создает новый экземпляр сервиса с переданным репозиторием и логгером в конструкторе.
func NewService(repo *repository.Repository, logger *logging.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// CreatePerson создает новую запись о человеке в базе данных.
// Обогащает данные о возрасте, поле и национальности с использованием внешних сервисов Agify, Genderize и Nationalize.
func (s *Service) CreatePerson(person *model.Person) error {
	s.logger.Debug("Service: Handling CreatePerson request")

	age, err := s.enrichWithAge(person.Name)
	if err != nil {
		s.logger.Errorf("Failed to enrich with age: %v", err)
		return err //errors.New("failed to enrich with age")
	}
	gender, err := s.enrichWithGender(person.Name)
	if err != nil {
		s.logger.Errorf("Failed to enrich with gender: %v", err)
		return err //errors.New("failed to enrich with gender")
	}
	nationality, err := s.enrichWithNationality(person.Name)
	if err != nil {
		s.logger.Errorf("Failed to enrich with nationality: %v", err)
		return err //errors.New("failed to enrich with nationality")
	}

	person.Age = age
	person.Gender = gender
	person.Nationality = nationality

	return s.repo.CreatePerson(person)

}

// GetPeople возвращает список людей с учетом переданных фильтров, смещения и лимита.
// Возрашаеть ошибку если не удолась
func (s *Service) GetPeople(filter map[string]interface{}, offset, limit int) ([]model.Person, error) {
	s.logger.Debug("Service: Handling GetPeople request")

	people, err := s.repo.GetPeople(filter, offset, limit)
	if err != nil {
		return nil, err
	}
	return people, nil
}

// GetPersonById возвращает информацию о человеке по его идентификатору.
// Возвращает ошибку, если человек не найден или при возникновении других проблем.
func (s *Service) GetPersonById(id int) (*model.Person, error) {
	s.logger.Debug("Service: Handling GetPersonById request")
	person, err := s.repo.GetPersonById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Warn("Person not found:", err)
			return nil, errors.New("person not found")
		}
		s.logger.Error("Failed to get person:", err)
		return nil, errors.New("failed to get person")

	}
	return person, nil
}

// UpdatePerson обновляет информацию о человеке в базе данных.
// Возвращает ошибку, если не удалось обновить информацию или при возникновении других проблем
func (s *Service) UpdatePerson(person *model.Person) error {
	s.logger.Debug("Service: Handling UpdatePerson request")
	err := s.repo.UpdatePerson(person)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Warn("Person not found:", err)
			return errors.New("person not found")
		}
		s.logger.Error("Failed to update person:", err)
		return errors.New("failed to update person")
	}
	return nil
}

// DeletePerson удаляет запись о человеке из базы данных по его id.
func (s *Service) DeletePerson(id int) error {
	s.logger.Debug("Service: Handling DeletePerson request")
	err := s.repo.DeletePerson(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Warn("person not found:", err)
			return errors.New("person not found")
		}
		s.logger.Error("Failed to delete person:", err)
		return errors.New("failed to delete person")
	}
	return nil
}

// enrichWithAge обогащает данные возрастом,
// подробнее: сразу не сдается при проблемах с внешним сервисом,
// а предпринимает попытки восстановления это делают код более устойчивым к временным проблемам с внешним сервисом.
func (s *Service) enrichWithAge(name string) (int, error) {
	s.logger.Debug("Service: Enriching with age")

	const maxRetries = 3
	var age int

	for retry := 0; retry < maxRetries; retry++ {
		url := fmt.Sprintf("https://api.agify.io/?name=%s", name)
		resp, err := http.Get(url)
		if err != nil {
			s.logger.Errorf("Failed to get age from Agify: %v", err)
			time.Sleep(time.Second) // Пауза перед повторной попыткой
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			s.logger.Errorf("Failed to get age from Agify. Status code: %d", resp.StatusCode)
			time.Sleep(time.Second) // Пауза перед повторной попыткой
			continue
		}

		if resp.StatusCode != http.StatusOK {
			s.logger.Errorf("Attempt %d: Unexpected status code from Agify: %d", retry+1, resp.StatusCode)
			time.Sleep(time.Second) // Пауза перед повторной попыткой
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			s.logger.Errorf("Attempt %d: Failed to read Agify response body: %v", retry+1, err)
			time.Sleep(time.Second) // Пауза перед повторной попыткой
			continue
		}
		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			s.logger.Errorf("Attempt %d: Failed to parse Agify response: %v", retry+1, err)
			time.Sleep(time.Second) // Пауза перед повторной попыткой
			continue
		}
		age, ok := result["age"].(float64)
		if !ok {
			s.logger.Errorf("Attempt %d: Failed to parse age from Agify response", retry+1)
			return 0, fmt.Errorf("failed to parse age from Agify response")
		}
		return int(age), nil
	}
	return age, fmt.Errorf("failed to get age from Agify after multiple retries")
}

// enrichWithGender обогащает данные полом с использованием внешнего сервиса Genderize.
// Возвращает пол и ошибку, если запрос к сервису не удался.
func (s *Service) enrichWithGender(name string) (string, error) {
	s.logger.Debug("Service: Enriching with gender")

	url := fmt.Sprintf("https://api.genderize.io/?name=%s", name)
	resp, err := http.Get(url)
	if err != nil {
		s.logger.Errorf("Failed to get gender from Genderize: %v", err)
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.logger.Errorf("Failed to read Genderize response body: %v", err)
		return "", err
	}

	var result map[string]interface{}

	if err := json.Unmarshal(body, &result); err != nil {
		s.logger.Errorf("Failed to read Genderize response body: %v", err)
		return "", err
	}

	gender, ok := result["gender"].(string)
	if !ok {
		s.logger.Errorf("Failed to parse gender from Genderize response")
		return "", err
	}
	return gender, nil

}

// enrichWithNationality обогащает данные национальностью с использованием внешнего сервиса Nationalize.
// Возвращает национальность и ошибку, если запрос к сервису не удался.
func (s *Service) enrichWithNationality(name string) (string, error) {
	s.logger.Debug("Service: Enriching with nationality")
	url := fmt.Sprintf("https://api.nationalize.io/?name=%s", name)
	resp, err := http.Get(url)
	if err != nil {
		s.logger.Errorf("Failed to get nationality from Nationalize: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.logger.Errorf("Unexpected status code from Nationalize: %d", resp.StatusCode)
		return "", fmt.Errorf("unexpected status code from Nationalize: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.logger.Errorf("Failed to read Nationalize response body: %v", err)
		return "", err
	}

	s.logger.Debugf("Nationalize response body: %s", string(body))

	var result map[string]interface{}

	if err := json.Unmarshal(body, &result); err != nil {
		s.logger.Errorf("Failed to parse Nationalize response: %v", err)
		return "", err
	}

	countryArray, ok := result["country"].([]interface{})
	if !ok || len(countryArray) == 0 {
		s.logger.Errorf("Failed to parse nationality from Nationalize response")
		return "", fmt.Errorf("failed to parse nationality from Nationalize response")

	}
	nationalityMap, ok := countryArray[0].(map[string]interface{})
	if !ok {
		s.logger.Errorf("Failed to parse nationality from Nationalize response")
		return "", fmt.Errorf("failed to parse nationality from Nationalize response")
	}

	nationality, ok := nationalityMap["country_id"].(string)
	if !ok {
		s.logger.Errorf("Failed to parse nationality from Nationalize response")
		return "", fmt.Errorf("failed to parse nationality from Nationalize response")
	}
	return nationality, nil

}

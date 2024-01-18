package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"testProject/internal/model"
	"testProject/pkg/helpers"
	"testProject/pkg/logging"

	"github.com/jmoiron/sqlx"
)

var ErrRowsAffected = errors.New("Error getting RowsAffected. This may indicate a problem with the underlying database or an issue with the query execution. Please check the database connection and the correctness of the query.")
var ErrNamedExec = errors.New("Error in NamedExec. This may indicate a problem with the underlying database or an issue with the query execution. Please check the database connection and the correctness of the query.")

type Repository struct {
	db     *sqlx.DB
	logger *logging.Logger
}

// спросить у макса покрывать код логгами или нет и еще как вытащить все ошибки в логе в одну чтобы было стурктурировано
func NewRepository(databaseURL string, logging *logging.Logger) (*Repository, error) {
	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}
	return &Repository{db: db, logger: logging}, nil
}

// Создания
func (r *Repository) CreatePerson(person *model.Person) error {
	r.logger.Debug("Repository: Handling CreatePerson request")

	query := `
        INSERT INTO people(name, surname, patronymic, age, gender, nationality) 
        VALUES($1, $2, $3, $4, $5, $6) 
        RETURNING id
    `

	err := r.db.QueryRow(query, person.Name, person.Surname, person.Patronymic, person.Age, person.Gender, person.Nationality).Scan(&person.ID)
	if err != nil {
		return err
	}

	return nil

}

// получения
func (r *Repository) GetPeople(filters map[string]interface{}, offset, limit int) ([]model.Person, error) {
	query := `SELECT * FROM people WHERE`
	args := make([]interface{}, 0)
	for key, value := range filters {
		query += key + "=? AND "
		args = append(args, value)
	}
	query += "1=1 LIMIT $1 OFFSET $2"

	args = append(args, limit, offset)

	var people []model.Person
	if err := r.db.Select(&people, query, args...); err != nil {
		helpers.LogAndReturnError(r.logger, "error when querying the database:", err)
		return nil, err
	}

	return people, nil
}

// получения по id
func (r *Repository) GetPersonById(id int) (*model.Person, error) {
	var person model.Person
	err := r.db.Get(&person, "SELECT * FROM people WHERE id = $1", id)
	if err != nil {
		helpers.LogAndReturnError(r.logger, "error when querying the database:", err)
		return nil, err
	}
	return &person, nil
}

// обноваления человека в бд
func (r *Repository) UpdatePerson(person *model.Person) error {
	query := `UPDATE people SET name=:name, surname=:surname, patronymic=:patronymic, 
	age=:age, gender=:gender, nationality=:nationality WHERE id=:id`

	result, err := r.db.NamedExec(query, person)
	if err != nil {
		return ErrNamedExec
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRowsAffected, err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// удалеения person
func (r *Repository) DeletePerson(id int) error {
	result, err := r.db.Exec("DELETE FROM people WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRowsAffected, err)

	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type Course struct {
	db          *sql.DB
	Name        string `json:"name"`
	Description string `json:"description"`
	ID          string `json:"id"`
	CategoryID  string `json:"category_id"`
}

func NewCourse(
	db *sql.DB,
) *Course {
	return &Course{
		db: db,
	}
}

func (c *Course) Create(name, description, categoryId string) (*Course, error) {

	cat := &Course{
		Name:        name,
		Description: description,
		ID:          uuid.NewString(),
		CategoryID:  categoryId,
	}

	_, err := c.db.Exec(
		"INSERT INTO courses (id,name,description, category_id) VALUES ($1,$2,$3,$4)",
		cat.ID,
		cat.Name,
		cat.Description,
		cat.CategoryID,
	)

	if !errors.Is(err, nil) {
		return nil, err
	}

	return cat, nil
}

func (c *Course) FindById(id string) (*Course, error) {

	row, err := c.db.Query("SELECT * FROM courses WHERE id = $1", id)

	if !errors.Is(err, nil) {
		return nil, err
	}

	var course *Course

	row.Scan(&course.ID, &course.Name, &course.Description)

	return course, nil
}

func (c *Course) FindAll() (*[]Course, error) {

	rows, err := c.db.Query("SELECT * FROM courses")

	if !errors.Is(err, nil) {
		return nil, err
	}

	var courses []Course

	for rows.Next() {

		var course Course

		rows.Scan(
			&course.ID,
			&course.Name,
			&course.Description,
			&course.CategoryID,
		)

		courses = append(courses, course)
	}

	return &courses, nil
}

func (c *Course) FindByCategoryId(categoryId string) (*[]Course, error) {
	rows, err := c.db.Query("SELECT * FROM courses WHERE category_id = ?", categoryId)

	if !errors.Is(err, nil) {
		return nil, err
	}

	var courses []Course

	for rows.Next() {
		var course Course

		rows.Scan(
			&course.ID,
			&course.Name,
			&course.Description,
			&course.CategoryID,
		)

		courses = append(courses, course)
	}
	fmt.Println(courses)
	return &courses, nil
}

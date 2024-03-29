// Package db ...
package db

import (
	"fmt"
	"log"
	"strconv"
)

type User struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// GetUsers ...
func GetUsers(page, limit string) ([]User, error) {
	var offset int = 0
	var users []User
	var total int

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return users, fmt.Errorf("page parse error: %s", err)
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return users, fmt.Errorf("limit parse error: %s", err)
	}

	if pageInt > 1 {
		offset = (pageInt - 1) * limitInt
	}

	queryRowTotal := fmt.Sprintf("SELECT count(*) FROM %s.users;", Postgres.schema)

	err = Postgres.DB.QueryRow(queryRowTotal).Scan(&total)
	if err != nil {
		return users, err
	}

	queryRow := fmt.Sprintf("SELECT * FROM %s.users ORDER BY id OFFSET %d LIMIT %d;", Postgres.schema, offset, limitInt)

	rows, err := Postgres.DB.Query(queryRow)
	if err != nil {
		return users, err
	}

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.Name, &user.Email); err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return users, nil
}

// GetUserById ...
func GetUserById(id string) (User, error) {
	var user User

	// SQL injection check
	if _, err := strconv.Atoi(id); err != nil {
		return user, fmt.Errorf("id must be a integer")
	}

	queryRow := fmt.Sprintf("SELECT * FROM %s.users WHERE id = %s;", Postgres.schema, id)

	row := Postgres.DB.QueryRow(queryRow)

	err := row.Scan(&user.Id, &user.Name, &user.Email)
	if err != nil {
		return user, err
	}

	return user, nil
}

// CreateUser ...
func CreateUser(name, email string) error {
	// TODO name and email validate
	_, err := Postgres.DB.Exec("INSERT INTO test.users (name, email) VALUES ($1, $2) RETURNING id", name, email)

	if err != nil {
		return err
	}
	// TODO must be return user error
	return nil
}

// UpdateUser ...
func UpdateUser(id, name, email string) error {
	setParamsString := ""

	if name != "" {
		setParamsString += fmt.Sprintf("name='%s'", name)
	}
	if email != "" {
		setParamsString += fmt.Sprintf(", email='%s'", email)
	}
	queryString := fmt.Sprintf("UPDATE test.users SET %s WHERE id = %s;", setParamsString, id)

	_, err := Postgres.DB.Exec(queryString)
	if err != nil {
		return err
	}
	// TODO must be return user error
	return nil
}

// DeleteUserById ...
func DeleteUserById(id string) error {
	// SQL injection check
	if _, err := strconv.Atoi(id); err != nil {
		return fmt.Errorf("id must be a integer")
	}

	execRow := fmt.Sprintf("DELETE FROM %s.users WHERE id = %s;", Postgres.schema, id)

	_, err := Postgres.DB.Exec(execRow)
	if err != nil {
		return err
	}
	// TODO  must be return user_id error
	return nil
}

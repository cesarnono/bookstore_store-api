package users

import (
	"fmt"
	"strings"

	"github.com/cesarnono/bookstore_users-api/datasources/mysql/users_db"
	"github.com/cesarnono/bookstore_users-api/logger"
	"github.com/cesarnono/bookstore_users-api/utils/mysql_utils"
	"github.com/cesarnono/bookstore_users-api/utils/rest_errors"
)

const (
	queryInsertUser             = "INSERT INTO users(first_name,last_name,email,date_created,status,password) VALUES(?,?,?,?,?,?);"
	emailUNIQUE                 = "email_UNIQUE"
	errorNoRows                 = "no rows in result set"
	queryGetUser                = "SELECT id, first_name, last_name,email,date_created,status FROM users WHERE id =?"
	queryUpdateUser             = "UPDATE users SET first_name=?, last_name=?, email=? WHERE id=?;"
	queryDeleteUser             = "DELETE FROM users WHERE id=?;"
	queryFindByStatus           = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE status=?;"
	queryFindByEmailAndPassword = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE email=? AND password=? AND status=?"
)

var (
	userDB = make(map[int64]*User)
)

func (user *User) Get() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		logger.Error("error when trying to prepare get user statement", err)
		return rest_errors.NewInternalServerError("database error")
	}
	defer stmt.Close()
	result := stmt.QueryRow(user.Id)
	if err := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
		return mysql_utils.ParseError(err)
	}
	return nil
}

func (user *User) Save() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		return rest_errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()
	insertResult, er := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Status, user.Password)
	if er != nil {
		return mysql_utils.ParseError(er)
	}
	userID, er := insertResult.LastInsertId()
	if er != nil {
		return mysql_utils.ParseError(er)
	}
	user.Id = userID
	return nil
}

func (user *User) Update() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		return rest_errors.NewInternalServerError("error when trying to update user1")
	}
	defer stmt.Close()
	fmt.Println(user)
	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id)
	if err != nil {
		return rest_errors.NewInternalServerError("error when trying to update user2")
	}
	return nil
}

func (user *User) Delete() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		return rest_errors.NewInternalServerError("error when tying to update user")
	}
	defer stmt.Close()

	if _, err = stmt.Exec(user.Id); err != nil {
		return rest_errors.NewInternalServerError("error when tying to save user")
	}
	return nil
}

func (user *User) FindByStatus(status string) ([]User, *rest_errors.RestErr) {
	stmt, err := users_db.Client.Prepare(queryFindByStatus)
	if err != nil {
		return nil, rest_errors.NewInternalServerError("error when trying to get user")
	}
	defer stmt.Close()

	rows, err := stmt.Query(status)
	if err != nil {
		return nil, rest_errors.NewInternalServerError("error when trying to get user")
	}
	defer rows.Close()

	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
			return nil, rest_errors.NewInternalServerError("error when trying to get user")
		}
		results = append(results, user)
	}
	if len(results) == 0 {
		return nil, rest_errors.NewNotFoundError(fmt.Sprintf("no users matching status %s", status))
	}
	return results, nil
}

func (user *User) FindByEmailAndPassword() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryFindByEmailAndPassword)
	if err != nil {
		logger.Error("error when trying to prepare get user by email and password statement", err)
		return rest_errors.NewInternalServerError("error when tying to find user")
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Email, user.Password, StatusActive)
	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); getErr != nil {
		if strings.Contains(getErr.Error(), mysql_utils.ErrorNoRows) {
			return rest_errors.NewNotFoundError("invalid user credentials")
		}
		logger.Error("error when trying to get user by email and password", getErr)
		return rest_errors.NewInternalServerError("error when tying to find user")
	}
	return nil
}

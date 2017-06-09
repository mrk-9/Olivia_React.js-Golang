package models

import (
	"database/sql"
	"errors"
	_ "github.com/davecgh/go-spew/spew"
	"strings"
	"regexp"
)

type SupportUser struct {
	Id         int    `json:"id"`
	FirstName  string `json:"fname"`
	SecondName string `json:"sname"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	CreatedAt  string `json:"created_at"`
}

func CreateSupportUsersTable() {
	stmt, queryErr := DB.Prepare("CREATE TABLE IF NOT EXISTS `support_users` (`id` INT UNSIGNED NOT NULL AUTO_INCREMENT, `fname` VARCHAR(128) NOT NULL, `sname` VARCHAR(128) NOT NULL, `email` VARCHAR(256) NOT NULL, `password` VARCHAR(40) NOT NULL, `created_at` TIMESTAMP NOT NULL DEFAULT NOW(), PRIMARY KEY (id))")
	defer stmt.Close()

	HandleErr(queryErr)

	_, execErr := stmt.Exec()

	HandleErr(execErr)
}

func (u *SupportUser) Save() error {
	u.Password = encodePassword(u.Password)

	var stmt *sql.Stmt
	var queryErr, execErr error

	if u.Id == 0 {
		stmt, queryErr = DB.Prepare("INSERT INTO `support_users` (`email`, `password`, `fname`, `sname`) VALUES (?,?,?,?)")
		defer stmt.Close()

		HandleErr(queryErr)

		reg := regexp.MustCompile("[a-zA-Z]+")
		var formattedFirstName string = strings.Join(reg.FindAllString(u.FirstName, -1), "")
		formattedFirstName = strings.ToLower(formattedFirstName)

		var result sql.Result
		result, execErr = stmt.Exec(u.Email, u.Password, u.FirstName, u.SecondName)
		id, _ := result.LastInsertId()
		u.Id = int(id)
	} else {
		stmt, queryErr = DB.Prepare("UPDATE `support_users` SET `email` = ?, `password` = ?, `fname` = ?, `sname` = ? WHERE `id` = ?")
		defer stmt.Close()

		HandleErr(queryErr)
		_, execErr = stmt.Exec(u.Id, u.Email, u.Password, u.FirstName, u.SecondName)
	}

	return execErr
}

func (u *SupportUser) GetById() error {
	if u.Id == 0 {
		return errors.New("User ID should be present")
	}

	stmt, queryErr := DB.Prepare("SELECT `email`, `fname`, `sname` FROM `support_users` WHERE `id` = ?")
	defer stmt.Close()

	HandleErr(queryErr)

	execErr := stmt.QueryRow(u.Id).Scan(&u.Email, &u.FirstName, &u.SecondName)

	return execErr
}

func (SupportUser) CheckExist(email string) (bool, error) {
	stmt, queryErr := DB.Prepare("SELECT `id` FROM `support_users` WHERE `email` = ?")
	defer stmt.Close()

	if queryErr != nil {
		return false, queryErr
	}

	var FoundUser User
	execErr := stmt.QueryRow(email).Scan(&FoundUser.Id)

	return (execErr == nil), execErr
}

func (u *SupportUser) GetByCredentials(email string, password string) error {
	password = encodePassword(password)

	stmt, queryErr := DB.Prepare("SELECT `id`, `email`, `fname`, `sname`, `created_at` FROM `support_users` WHERE `email` = ? AND `password` = ?")
	defer stmt.Close()

	if queryErr != nil {
		return queryErr
	}

	scanErr := stmt.QueryRow(email, password).Scan(&u.Id, &u.Email, &u.FirstName, &u.SecondName, &u.CreatedAt)

	return scanErr
}

func (u *SupportUser) GetByToken(token string) error {
	stmt, queryErr := DB.Prepare("SELECT `id`, `email`, `fname`, `sname`, `created_at` FROM `support_users` WHERE `id` = (SELECT `user_id` FROM `support_tokens` WHERE `token` = ?)")
	defer stmt.Close()

	if queryErr != nil {
		return queryErr
	}

	scanErr := stmt.QueryRow(token).Scan(&u.Id, &u.Email, &u.FirstName, &u.SecondName, &u.CreatedAt)

	return scanErr
}

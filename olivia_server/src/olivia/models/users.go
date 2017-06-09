package models

import (
	"database/sql"
	"errors"
	_ "github.com/davecgh/go-spew/spew"
	"log"
	"regexp"
	"strings"
)

type User struct {
	Id         int64  `json:"id"`
	FirstName  string `json:"fname"`
	SecondName string `json:"sname"`
	OEmail     string `json:"o_email"`
	Email      string `json:"email"`
	Code       string `json:"code"`
	Phone      string `json:"phone"`
	Password   string `json:"password"`
	CreatedAt  string `json:"created_at"`
}

func CreateUsersTable() {
	stmt, queryErr := DB.Prepare("CREATE TABLE IF NOT EXISTS `users` (`id` INT UNSIGNED NOT NULL AUTO_INCREMENT, `o_email` VARCHAR(256) NOT NULL, `fname` VARCHAR(128) NOT NULL, `sname` VARCHAR(128) NOT NULL, `email` VARCHAR(256) NOT NULL, `phone` VARCHAR(64) NOT NULL, `code` VARCHAR(8) NOT NULL, `password` VARCHAR(40) NOT NULL, `created_at` TIMESTAMP NOT NULL DEFAULT NOW(), PRIMARY KEY (id))")
	defer stmt.Close()

	HandleErr(queryErr)

	_, execErr := stmt.Exec()

	HandleErr(execErr)
}

func (u *User) Save() error {
	u.Password = encodePassword(u.Password)

	var stmt *sql.Stmt
	var queryErr, execErr error

	if u.Id == 0 {
		stmt, queryErr = DB.Prepare("INSERT INTO `users` (`o_email`, `email`, `phone`, `code`, `password`, `fname`, `sname`) VALUES (?,?,?,?,?,?,?)")
		defer stmt.Close()

		HandleErr(queryErr)

		reg := regexp.MustCompile("[a-zA-Z]+")
		var formattedFirstName string = strings.Join(reg.FindAllString(u.FirstName, -1), "")
		formattedFirstName = strings.ToLower(formattedFirstName)

		u.OEmail = strings.Join([]string{formattedFirstName, "-", CreateRandomString(16)}, "")

		var result sql.Result
		result, execErr = stmt.Exec(u.OEmail, u.Email, u.Phone, u.Code, u.Password, u.FirstName, u.SecondName)
		id, _ := result.LastInsertId()
		u.Id = int64(id)
	} else {
		stmt, queryErr = DB.Prepare("UPDATE `users` SET `email` = ?, `phone` = ?, `code` = ?, `password` = ?, `fname` = ?, `sname` = ? WHERE `id` = ?")
		defer stmt.Close()

		HandleErr(queryErr)
		_, execErr = stmt.Exec(u.Id, u.Email, u.Phone, u.Code, u.Password, u.FirstName, u.SecondName)
	}

	return execErr
}

func (u *User) GetById() error {
	if u.Id == 0 {
		return errors.New("User ID should be present")
	}

	stmt, queryErr := DB.Prepare("SELECT `email`, `fname`, `sname` FROM `users` WHERE `id` = ?")
	defer stmt.Close()

	HandleErr(queryErr)

	execErr := stmt.QueryRow(u.Id).Scan(&u.Email, &u.FirstName, &u.SecondName)

	return execErr
}

func GetAllUsers() []User {
	stmt, queryErr := DB.Prepare("SELECT `id`, `email`, `fname`, `sname` FROM `users`")
	defer stmt.Close()

	HandleErr(queryErr)

	rows, execErr := stmt.Query()
	if execErr != nil {
		log.Println(execErr)
	}

	var users = make([]User, 0)

	for rows.Next() {
		var u User

		err := rows.Scan(&u.Id, &u.Email, &u.FirstName, &u.SecondName)
		HandleErr(err)

		users = append(users, u)
	}

	return users
}

func (User) CheckExist(email string) (bool, error) {
	stmt, queryErr := DB.Prepare("SELECT `id` FROM `users` WHERE `email` = ?")
	defer stmt.Close()

	if queryErr != nil {
		return false, queryErr
	}

	var FoundUser User
	execErr := stmt.QueryRow(email).Scan(&FoundUser.Id)

	return (execErr == nil), execErr
}

func (u *User) GetByCredentials(email string, password string) error {
	password = encodePassword(password)

	stmt, queryErr := DB.Prepare("SELECT `id`, `o_email`, `email`, `phone`, `code`, `fname`, `sname`, `created_at` FROM `users` WHERE `email` = ? AND `password` = ?")
	defer stmt.Close()

	if queryErr != nil {
		return queryErr
	}

	scanErr := stmt.QueryRow(email, password).Scan(&u.Id, &u.OEmail, &u.Email, &u.Phone, &u.Code, &u.FirstName, &u.SecondName, &u.CreatedAt)

	return scanErr
}

func (u *User) GetByToken(token string) error {
	stmt, queryErr := DB.Prepare("SELECT `id`, `o_email`, `email`, `phone`, `code`, `fname`, `sname`, `created_at` FROM `users` WHERE `id` = (SELECT `user_id` FROM `tokens` WHERE `token` = ?)")
	defer stmt.Close()

	if queryErr != nil {
		return queryErr
	}

	scanErr := stmt.QueryRow(token).Scan(&u.Id, &u.OEmail, &u.Email, &u.Phone, &u.Code, &u.FirstName, &u.SecondName, &u.CreatedAt)

	return scanErr
}

func (u *User) GetByOEmail() error {
	stmt, queryErr := DB.Prepare("SELECT `id`, `email`, `code`, `phone`, `fname` FROM `users` WHERE `o_email` = ?")
	defer stmt.Close()

	if queryErr != nil {
		return queryErr
	}

	execErr := stmt.QueryRow(u.OEmail).Scan(&u.Id, &u.Email, &u.Code, &u.Phone, &u.FirstName)

	return execErr
}

func (u *User) GetByInternalPhoneNumber(phoneNumber string) error {
	var pn PhoneNumber
	pn.GetByPhoneNumber(phoneNumber)

	if pn.Id == 0 {
		log.Printf("Number not found: %s", phoneNumber)
		return nil
	}

	u.Id = pn.Owner.Id
	u.Email = pn.Owner.Email
	u.FirstName = pn.Owner.FirstName
	u.SecondName = pn.Owner.SecondName
	u.OEmail = pn.Owner.OEmail
	u.Phone = pn.Owner.Phone
	u.CreatedAt = pn.Owner.CreatedAt

	return nil
}

func (u *User) GetPhoneNumber() (PhoneNumber, error) {
	var pn PhoneNumber

	stmt, queryErr := DB.Prepare(
		"SELECT `id`, `sid`, `number`, `created_at` " +
			"FROM `phone_numbers` " +
			"WHERE `owner` = ? " +
			"LIMIT 1")
	defer stmt.Close()

	if queryErr != nil {
		return pn, queryErr
	}

	stmt.QueryRow(u.Id).Scan(&pn.Id, &pn.Sid, &pn.Number, &pn.CreatedAt)

	return pn, nil
}

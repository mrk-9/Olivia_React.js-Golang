package models

import (
	"database/sql"
	"log"
	"olivia/parsers"
)

type Email struct {
	Id          int64  `json:"id"`
	UserId      int64  `json:"user_id"`
	From        string `json:"from"`
	Subject     string `json:"subject"`
	Status      string `json:"status"`
	EnquiryType string `json:"enquiry_type"`
	Lead        int    `json:"lead"`
	Text        string `json:""`
	HTML        string `json:""`
	Raw         string `json:""`
	CreatedAt   string `json:"created_at"`
}

func CreateEmailsTable() {
	stmt, queryErr := DB.Prepare("" +
		"CREATE TABLE IF NOT EXISTS `emails` (" +
		"`id` INT UNSIGNED NOT NULL AUTO_INCREMENT, " +
		"`user_id` INT UNSIGNED NOT NULL, " +
		"`from` VARCHAR(512) NOT NULL, " +
		"`subject` VARCHAR(512) NOT NULL, " +
		"`status` VARCHAR(512), " +
		"`enquiry_type` VARCHAR(512), " +
		"`lead` INT, " +
		"`text` LONGTEXT NOT NULL, " +
		"`html` LONGTEXT NOT NULL, " +
		"`raw` LONGTEXT NOT NULL, " +
		"`created_at` TIMESTAMP, PRIMARY KEY (id)" +
		")")
	defer stmt.Close()

	HandleErr(queryErr)

	_, execErr := stmt.Exec()

	HandleErr(execErr)
}

func (e *Email) Create() (sql.Result, error) {
	stmt, queryErr := DB.Prepare("INSERT INTO `emails` (`user_id`, `from`, `subject`, `status`, `enquiry_type`, `text`, `html`, `raw`) VALUES (?,?,?,?,?,?,?,?)")
	defer stmt.Close()

	HandleErr(queryErr)

	result, execErr := stmt.Exec(e.UserId, e.From, e.Subject, e.Status, e.EnquiryType, e.Text, e.HTML, e.Raw)
	if execErr != nil {
		log.Println(execErr.Error())
	}

	lastId, _ := result.LastInsertId()
	e.Id = lastId

	return result, execErr
}

func (e *Email) Save() error {
	var stmt *sql.Stmt
	var queryErr, execErr error

	if e.Id == 0 {
		stmt, queryErr = DB.Prepare("INSERT INTO `emails` (`user_id`, `from`, `text`, `html`) VALUES (?,?,?,?)")
		defer stmt.Close()

		HandleErr(queryErr)

		_, execErr = stmt.Exec(e.UserId, e.From, e.Text, e.HTML)
	} else {
		stmt, queryErr = DB.Prepare("UPDATE `emails` SET `lead` = ? WHERE `id` = ? ")
		defer stmt.Close()

		HandleErr(queryErr)
		_, execErr = stmt.Exec(e.Lead, e.Id)
	}

	return execErr
}

func GetEmailsByUserId(id int64) []parsers.EmailData {
	stmt, queryErr := DB.Prepare("SELECT `id`, `from`, `subject` `created_at` FROM `emails` WHERE `user_id` = ?")
	defer stmt.Close()

	HandleErr(queryErr)

	rows, err := stmt.Query(id)
	HandleErr(err)

	defer rows.Close()

	var emails = make([]parsers.EmailData, 0)

	for rows.Next() {
		var e Email

		err := rows.Scan(&e.Id, &e.From, &e.Subject, &e.CreatedAt)
		HandleErr(err)

		var emailData parsers.EmailData

		emailData.Id = e.Id

		parsers.Parse(e.From, e.HTML, &emailData)

		emailData.CreatedAt = e.CreatedAt

		emails = append(emails, emailData)
	}

	return emails
}

func GetEmailByLead(lead Lead) (Email, error) {
	stmt, queryErr := DB.Prepare("SELECT `id`, `from`, `subject`, `status`, `enquiry_type`, `lead`, `created_at` FROM `emails` WHERE `lead` = ?")
	defer stmt.Close()

	HandleErr(queryErr)

	var e Email
	scanErr := stmt.QueryRow(lead.Id).Scan(&e.Id, &e.From, &e.Subject, &e.Status, &e.EnquiryType, &e.Lead, &e.CreatedAt)

	return e, scanErr
}

func GetEmailsByLeadAndUser(lead *Lead, user *User) (Email, error) {
	stmt, queryErr := DB.Prepare(
		"SELECT `id`, `from`, `subject`, `status`, `enquiry_type`, `lead`, `created_at` " +
			"FROM `emails` " +
			"WHERE `user_id` = ? AND `lead` = ? ")
	defer stmt.Close()

	if queryErr != nil {
		log.Println(queryErr.Error())
	}

	var e Email
	scanErr := stmt.QueryRow(user.Id, lead.Id).Scan(&e.Id, &e.From, &e.Subject, &e.Status, &e.EnquiryType, &e.Lead, &e.CreatedAt)

	return e, scanErr
}

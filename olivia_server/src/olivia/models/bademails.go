package models

import (
	"database/sql"
	"log"
)

type BadEmail struct {
	Id          int    `json:"id"`
	From        string `json:"from"`
	Subject     string `json:"subject"`
	Text        string `json:""`
	HTML        string `json:""`
	Raw         string `json:""`
	CreatedAt   string `json:"created_at"`
}

func CreateBadEmailsTable() {
	stmt, queryErr := DB.Prepare("" +
		"CREATE TABLE IF NOT EXISTS `bademails` (" +
		"`id` INT UNSIGNED NOT NULL AUTO_INCREMENT, " +
		"`from` VARCHAR(512) NOT NULL, " +
		"`subject` VARCHAR(512) NOT NULL, " +
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

func (e *BadEmail) Create() (sql.Result, error) {
	stmt, queryErr := DB.Prepare("INSERT INTO `bademails` (`from`, `subject`, `text`, `html`, `raw`) VALUES (?,?,?,?,?)")
	defer stmt.Close()

	HandleErr(queryErr)

	result, execErr := stmt.Exec(e.From, e.Subject, e.Text, e.HTML, e.Raw)
	if execErr != nil {
		log.Println(execErr.Error())
	}

	lastId, _ := result.LastInsertId()
	e.Id = int(lastId)

	return result, execErr
}
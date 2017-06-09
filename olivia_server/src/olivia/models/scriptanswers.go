package models

import (
	"log"
)

type ScriptAnswer struct {
	Id          int64  `json:"id"`
	Lead        int64  `json:"owner"`
	Script      int64  `json:"title"`
	ScriptTitle string `json:"script_title"`
	Text        string `json:"text"`
	CreatedAt   string `json:"created_at"`
}

func CreateScriptAnswersTable() {
	stmt, queryErr := DB.Prepare("" +
		"CREATE TABLE IF NOT EXISTS `scriptanswers` (" +
		"`id` INT UNSIGNED NOT NULL AUTO_INCREMENT, " +
		"`lead` INT UNSIGNED NOT NULL, " +
		"`script` INT UNSIGNED NOT NULL, " +
		"`text` TEXT NOT NULL, " +
		"`created_at` TIMESTAMP, PRIMARY KEY (id)" +
		")")
	defer stmt.Close()

	HandleErr(queryErr)

	_, execErr := stmt.Exec()

	HandleErr(execErr)
}

func GetAnswerByScriptAndLead(scriptId int64, leadId int64) (error, ScriptAnswer) {
	stmt, queryErr := DB.Prepare("SELECT `id`, `lead`, `script`, `text` " +
		"FROM `scriptanswers` WHERE `script` = ? AND `lead` = ?")
	defer stmt.Close()

	if queryErr != nil {
		return queryErr, ScriptAnswer{}
	}

	var sa ScriptAnswer
	stmt.QueryRow(scriptId, leadId).Scan(&sa.Id, &sa.Lead, &sa.Script, &sa.Text)

	return nil, sa
}

func (sa *ScriptAnswer) Create() error {
	stmt, queryErr := DB.Prepare("INSERT INTO `scriptanswers` " +
		"(`lead`, `script`, `text`) " +
		"VALUES (?,?,?)")
	defer stmt.Close()

	if queryErr != nil {
		return queryErr
	}

	result, execErr := stmt.Exec(&sa.Lead, &sa.Script, &sa.Text)
	lastId, _ := result.LastInsertId()
	sa.Id = lastId

	return execErr
}

func (sa *ScriptAnswer) Save() error {
	stmt, queryErr := DB.Prepare("UPDATE `scriptanswers` " +
		"SET `text` = ? WHERE `id` = ?")
	defer stmt.Close()

	if queryErr != nil {
		return queryErr
	}

	result, execErr := stmt.Exec(sa.Text, sa.Id)

	log.Println(result)

	return execErr
}

func GetLeadAnswers(leadId int64, userId int64) ([]ScriptAnswer, error) {
	stmt, queryErr := DB.Prepare("SELECT s.id, s.title, IFNULL(sa.text, \"\") as `text` " +
		"FROM `scripts` `s` " +
		"LEFT JOIN `scriptanswers` AS `sa` ON `s`.`id` = `sa`.`script` AND `sa`.`lead` = ? " +
		"WHERE `s`.`owner` = ?")
	defer stmt.Close()

	if queryErr != nil {
		return nil, queryErr
	}

	rows, scanErr := stmt.Query(leadId, userId)
	if scanErr != nil {
		log.Println(scanErr.Error())
	}

	defer rows.Close()

	var scriptAnswers = make([]ScriptAnswer, 0)

	for rows.Next() {
		var sa ScriptAnswer

		err := rows.Scan(&sa.Id, &sa.ScriptTitle, &sa.Text)
		HandleErr(err)

		scriptAnswers = append(scriptAnswers, sa)
	}

	return scriptAnswers, scanErr
}

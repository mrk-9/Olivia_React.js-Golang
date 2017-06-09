package models

import (
	"log"
	"strings"
)

type Script struct {
	Id        int64  `json:"id"`
	Owner     int64  `json:"owner"`
	Title     string `json:"title"`
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
}

type ScriptAndAnswer struct {
	Id     int64  `json:"id"`
	Title  string `json:"title"`
	Script string `json:"script"`
	Answer string `json:"answer"`
}

func CreateScriptsTable() {
	stmt, queryErr := DB.Prepare("" +
		"CREATE TABLE IF NOT EXISTS `scripts` (" +
		"`id` INT UNSIGNED NOT NULL AUTO_INCREMENT, " +
		"`owner` INT UNSIGNED NOT NULL, " +
		"`title` VARCHAR(32) NOT NULL, " +
		"`text` VARCHAR(200) NOT NULL, " +
		"`created_at` TIMESTAMP, PRIMARY KEY (id)" +
		")")
	defer stmt.Close()

	HandleErr(queryErr)

	_, execErr := stmt.Exec()

	HandleErr(execErr)
}

func CreateDefaultScripts(userId int64) {
	var s1 Script
	var s2 Script
	var s3 Script
	var s4 Script
	var s5 Script

	s1.Title = "General Questions"
	s1.Text = "Hi {first_name}, I’m responding to your request about the property at {address}. Do you have any questions about the property?"

	s2.Title = "Budget & Loan Approval"
	s2.Text = "Have you already been approved for a home loan? If yes, what is your approximate budget?"

	s3.Title = "Property Viewing Time"
	s3.Text = "What’s the best time for you to see the property {first_name}?"

	s4.Title = "Interested in Other Properties?"
	s4.Text = "Would you be interested in other properties in this area or price range?"

	s5.Title = "Conclusion"
	s5.Text = "Thanks for your answers {first_name}, I’ll have the realtor contact you with more information within the next day. Have a great day!"

	s1.Owner = userId
	s2.Owner = userId
	s3.Owner = userId
	s4.Owner = userId
	s5.Owner = userId

	s1.Create()
	s2.Create()
	s3.Create()
	s4.Create()
	s5.Create()
}

func (s *Script) Create() error {
	stmt, queryErr := DB.Prepare("INSERT INTO `scripts` " +
		"(`owner`, `title`, `text`) " +
		"VALUES (?,?,?)")
	defer stmt.Close()

	if queryErr != nil {
		return queryErr
	}

	result, execErr := stmt.Exec(&s.Owner, &s.Title, &s.Text)
	lastId, _ := result.LastInsertId()
	s.Id = lastId

	return execErr
}

func (s *Script) Save() error {
	stmt, queryErr := DB.Prepare("UPDATE `scripts` " +
		"SET `title` = ?, `text` = ? WHERE `id` = ?")
	defer stmt.Close()

	if queryErr != nil {
		return queryErr
	}

	result, execErr := stmt.Exec(s.Title, s.Text, s.Id)

	log.Println(result)

	return execErr
}

func GetUserScripts(user User) ([]Script, error) {
	stmt, queryErr := DB.Prepare("SELECT `id`, `owner`, `title`, `text` FROM `scripts` WHERE " +
		"`owner` = ? ORDER BY `id` ASC")
	defer stmt.Close()

	if queryErr != nil {
		return nil, queryErr
	}

	rows, scanErr := stmt.Query(user.Id)
	if scanErr != nil {
		log.Println(scanErr.Error())
	}

	defer rows.Close()

	var scripts = make([]Script, 0)

	for rows.Next() {
		var s Script

		err := rows.Scan(&s.Id, &s.Owner, &s.Title, &s.Text)
		HandleErr(err)

		scripts = append(scripts, s)
	}

	return scripts, scanErr
}

func GetScriptById(id int64) Script {
	var s Script

	stmt, queryErr := DB.Prepare("SELECT `owner`, `title`, `text`, `created_at` FROM `scripts` WHERE " +
		"`id` = ?")
	defer stmt.Close()

	if queryErr != nil {
		return s
	}

	stmt.QueryRow(id).Scan(&s.Owner, &s.Title, &s.Text, &s.CreatedAt)
	s.Id = id

	return s
}

func RemoveScript(id int64, ownerId int64) {
	stmt, queryErr := DB.Prepare("DELETE FROM `scripts` WHERE `id` = ? AND `owner` = ?")
	defer stmt.Close()

	if queryErr != nil {
		log.Println(queryErr)
	}

	log.Println(id, ownerId)

	stmt.Exec(id, ownerId)
}

func FormatScriptMessage(text string, lead *Lead, property *Property) string {
	leadName := strings.SplitN(lead.Name, " ", 2)
	firstName := lead.Name
	lastName := lead.Name
	if len(leadName) == 2 {
		firstName = leadName[0]
		lastName = leadName[1]
	}

	text = strings.Replace(text, "{first_name}", firstName, -1)
	text = strings.Replace(text, "(first_name}", firstName, -1)
	text = strings.Replace(text, "{last_name}", lastName, -1)
	text = strings.Replace(text, "{sale_price}", property.Price, -1)
	text = strings.Replace(text, "{address}", property.Address, -1)
	text = strings.Replace(text, "{description}", property.Description, -1)

	return text
}

func GetScriptsAndAnswers(user User, lead Lead, property Property) []ScriptAndAnswer {
	stmt, queryErr := DB.Prepare(
		"SELECT `s`.`id`, `s`.`title`, `s`.`text`, " +
			"CASE WHEN `sa`.`text` IS NULL THEN \"\" ELSE `sa`.`text` END AS `answer` " +
			"FROM `scripts` `s` " +
			"LEFT JOIN `scriptanswers` `sa` ON `s`.`id` = `sa`.`script` AND `sa`.`lead` = ? " +
			"WHERE `s`.`owner` = ? " +
			"ORDER BY `s`.`id` ")
	defer stmt.Close()

	if queryErr != nil {
		log.Println(queryErr)
	}

	rows, scanErr := stmt.Query(lead.Id, user.Id)
	if scanErr != nil {
		log.Println(scanErr.Error())
	}

	defer rows.Close()

	var scripts = make([]ScriptAndAnswer, 0)

	for rows.Next() {
		var s ScriptAndAnswer

		err := rows.Scan(&s.Id, &s.Title, &s.Script, &s.Answer)
		HandleErr(err)

		s.Script = FormatScriptMessage(s.Script, &lead, &property)

		scripts = append(scripts, s)
	}

	return scripts
}

package models

import "log"

type Message struct {
	Id        int64  `json:"id"`
	Lead      int64  `json:"user"`
	User      int64  `json:"lead"`
	Message   string `json:"message"`
	ForLead   bool   `json:"for_lead"`
	CreatedAt string `json:"created_at"`
}

func CreateMessagesTable() {
	stmt, queryErr := DB.Prepare("" +
		"CREATE TABLE IF NOT EXISTS `messages` (" +
		"`id` INT UNSIGNED NOT NULL AUTO_INCREMENT, " +
		"`lead` INT UNSIGNED NOT NULL, " +
		"`user` INT UNSIGNED NOT NULL, " +
		"`message` TEXT NOT NULL, " +
		"`for_lead` TINYINT(1) NOT NULL, " +
		"`created_at` TIMESTAMP, PRIMARY KEY (id)" +
		")")
	defer stmt.Close()

	HandleErr(queryErr)

	_, execErr := stmt.Exec()

	HandleErr(execErr)
}

func (m *Message) Create() error {
	stmt, queryErr := DB.Prepare("INSERT INTO `messages` " +
		"(`lead`, `user`, `message`, `for_lead`) " +
		"VALUES (?,?,?,?)")
	defer stmt.Close()

	if queryErr != nil {
		return queryErr
	}

	result, execErr := stmt.Exec(m.Lead, m.User, m.Message, m.ForLead)
	lastId, _ := result.LastInsertId()
	m.Id = lastId

	m.GetById(m.Id)

	return execErr
}

func (m *Message) GetById(id int64) error {
	stmt, queryErr := DB.Prepare("SELECT `lead`, `user`, `message`, `for_lead`, `created_at` " +
		"FROM `messages` WHERE `id` = ?")
	defer stmt.Close()

	HandleErr(queryErr)

	execErr := stmt.QueryRow(id).Scan(&m.Lead, &m.User, &m.Message, &m.ForLead, &m.CreatedAt)
	m.Id = id

	return execErr
}

func GetChat(user User, lead Lead) ([]Message, error) {
	stmt, queryErr := DB.Prepare("SELECT `id`, `message`, `for_lead`, `created_at` " +
		"FROM `messages` " +
		"WHERE `user` = ? AND `lead` = ? " +
		"ORDER BY `created_at` ASC")
	defer stmt.Close()

	if queryErr != nil {
		return nil, queryErr
	}

	rows, scanErr := stmt.Query(user.Id, lead.Id)
	if scanErr != nil {
		log.Println(scanErr.Error())
	}

	defer rows.Close()

	var messages = make([]Message, 0)

	for rows.Next() {
		var m Message

		err := rows.Scan(&m.Id, &m.Message, &m.ForLead, &m.CreatedAt)
		if err != nil {
			log.Printf("Can't scan chat query: %s", err.Error())
		}

		m.Lead = lead.Id
		m.User = user.Id

		messages = append(messages, m)
	}

	return messages, scanErr
}

func GetLastMessageFromChat(user User, lead Lead) (Message, error) {
	stmt, queryErr := DB.Prepare("SELECT `id`, `message`, `for_lead`, `created_at` " +
		"FROM `messages` " +
		"WHERE `user` = ? AND `lead` = ? " +
		"ORDER BY `id` DESC " +
		"LIMIT 1")
	defer stmt.Close()

	if queryErr != nil {
		return Message{}, queryErr
	}

	var m Message

	scanErr := stmt.QueryRow(user.Id, lead.Id).Scan(&m.Id, &m.Message, &m.ForLead, &m.CreatedAt)
	if scanErr != nil {
		return Message{}, scanErr
	}

	return m, nil
}

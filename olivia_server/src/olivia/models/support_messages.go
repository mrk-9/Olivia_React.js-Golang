package models

type SupportMessage struct {
	Id        int    `json:"id"`
	Lead      int    `json:"user"`
	User      int    `json:"lead"`
	Message   string `json:"message"`
	ForLead   bool   `json:"for_lead"`
	CreatedAt string `json:"created_at"`
}

func GetLastSupportMessageFromChat(user User, lead Lead) (SupportMessage, error) {
	stmt, queryErr := DB.Prepare("SELECT `id`, `message`, `for_lead`, `created_at` " +
		"FROM `messages` " +
		"WHERE `user` = ? AND `lead` = ? " +
		"ORDER BY `id` DESC " +
		"LIMIT 1")
	defer stmt.Close()

	if queryErr != nil {
		return SupportMessage{}, queryErr
	}

	var m SupportMessage

	scanErr := stmt.QueryRow(user.Id, lead.Id).Scan(&m.Id, &m.Message, &m.ForLead, &m.CreatedAt)
	if scanErr != nil {
		return SupportMessage{}, scanErr
	}

	return m, nil
}

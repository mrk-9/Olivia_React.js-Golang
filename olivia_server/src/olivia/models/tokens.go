package models

type Token struct {
	Id     int    `json:"-"`
	Token  string `json:"token"`
	UserId int    `json:"-"`
}

func CreateTokensTable() {
	stmt, queryErr := DB.Prepare("CREATE TABLE IF NOT EXISTS `tokens` (`id` INT UNSIGNED NOT NULL AUTO_INCREMENT, `token` VARCHAR(256) NOT NULL, `user_id` INT UNSIGNED NOT NULL, PRIMARY KEY (id))")
	defer stmt.Close()

	HandleErr(queryErr)

	_, execErr := stmt.Exec()

	HandleErr(execErr)
}

func (t *Token) Create(user User) {
	stmt, queryErr := DB.Prepare("INSERT INTO `tokens` (`token`, `user_id`) VALUES (?,?)")
	defer stmt.Close()

	HandleErr(queryErr)

	t.Token = CreateRandomString(32)
	t.UserId = int(user.Id)

	stmt.Exec(t.Token, t.UserId)
}

func (t *Token) GetTokenForUser(user User) {
	stmt, queryErr := DB.Prepare("SELECT `id`, `token`, `user_id` FROM `tokens` WHERE `user_id` = ?")
	defer stmt.Close()

	HandleErr(queryErr)

	stmt.QueryRow(user.Id).Scan(&t.Id, &t.Token, &t.UserId)
}

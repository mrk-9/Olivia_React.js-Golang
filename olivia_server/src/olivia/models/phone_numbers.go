package models

import ()

type PhoneNumber struct {
	Id        int    `json:"id"`
	Sid       string `json:"sid"`
	Number    string `json:"number"`
	Owner     User   `json:"owner"`
	CreatedAt string `json:"created_at"`
}

func CreatePhoneNumbersTable() {
	stmt, queryErr := DB.Prepare("" +
		"CREATE TABLE IF NOT EXISTS `phone_numbers` (" +
		"`id` int(11) NOT NULL AUTO_INCREMENT, " +
		"`sid` varchar(64) NOT NULL, " +
		"`number` varchar(128) NOT NULL, " +
		"`created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, " +
		"`owner` int(11) NOT NULL, " +
		"PRIMARY KEY (`id`)" +
		")")
	defer stmt.Close()

	HandleErr(queryErr)

	_, execErr := stmt.Exec()

	HandleErr(execErr)
}

func (pn *PhoneNumber) Add() (error) {
	stmt, queryErr := DB.Prepare("INSERT INTO `phone_numbers` (`number`, `sid`, `owner`) VALUES (?,?,?)")
	defer stmt.Close()

	if queryErr != nil {
		return queryErr
	}

	result, execErr := stmt.Exec(pn.Number, pn.Sid, pn.Owner.Id)
	if execErr != nil {
		return execErr
	}

	lastId, lastIdError := result.LastInsertId()
	if lastIdError != nil {
		return lastIdError
	}

	pn.Id = int(lastId)

	return nil
}

func (pn *PhoneNumber) GetNumberById() (error) {
	stmt, queryErr := DB.Prepare("SELECT `id`, `number`, `owner`, `created_at` FROM `phone_numbers` WHERE `id` = ?")
	defer stmt.Close()

	if queryErr != nil {
		return queryErr
	}

	var ownerId int

	stmt.QueryRow(pn.Id).Scan(&pn.Id, &pn.Number, &ownerId, &pn.CreatedAt)

	var user User
	user.Id = int64(ownerId)

	userError := user.GetById()
	if userError != nil {
		return userError
	}

	pn.Owner = user

	return nil
}

func (pn *PhoneNumber) GetByPhoneNumber(number string) (error) {
	stmt, queryErr := DB.Prepare(
		"SELECT `id`, `sid`, `owner`, `created_at` " +
			"FROM `phone_numbers` " +
			"WHERE `number` = ?")
	defer stmt.Close()

	if queryErr != nil {
		return queryErr
	}

	var ownerId int

	stmt.QueryRow(number).Scan(&pn.Id, &pn.Sid, &ownerId, &pn.CreatedAt)
	pn.Number = number

	pn.Owner.Id = int64(ownerId)
	pn.Owner.GetById()

	return nil
}
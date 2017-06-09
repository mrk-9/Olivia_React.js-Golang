package models

import (
	"database/sql"
	"errors"
	"log"
	"strconv"
	"strings"
)

type Lead struct {
	Id              int64      `json:"id"`
	Name            string     `json:"name"`
	Phone           string     `json:"phone"`
	Email           string     `json:"email"`
	Realtor         int64      `json:"realtor"`
	InterestedInIds string     `json:"interested_in_ids"`
	InterestedIn    []Property `json:"interested_in"`
	NeedsAssistance bool       `json:"needs_assistance"`
	CreatedAt       string     `json:"created_at"`
	ZendeskTicketId int64      `json:"zendesk_ticket_id"`
}

func CreateLeadsTable() {
	stmt, queryErr := DB.Prepare("" +
		"CREATE TABLE IF NOT EXISTS `leads` (" +
		"`id` INT UNSIGNED NOT NULL AUTO_INCREMENT, " +
		"`name` VARCHAR(256) NOT NULL, " +
		"`phone` VARCHAR(256) NOT NULL, " +
		"`email` VARCHAR(256) NOT NULL, " +
		"`realtor` INT UNSIGNED, " +
		"`interested_in` TEXT, " +
		"`needs_assistance` TINYINT(1) NOT NULL, " +
		"`created_at` TIMESTAMP, PRIMARY KEY (id), " +
		"`zendesk_ticket_id` INT UNSIGNED NOT NULL" +
		")")
	defer stmt.Close()

	HandleErr(queryErr)

	_, execErr := stmt.Exec()

	HandleErr(execErr)
}

func (l *Lead) GetById(id int64) error {
	stmt, queryErr := DB.Prepare(
		"SELECT `name`, `phone`, `email`, " +
			"CASE WHEN `realtor` IS NULL THEN 0 ELSE `realtor` END AS `realtor`, " +
			"`interested_in`, `needs_assistance`, `created_at`, `zendesk_ticket_id` " +
			"FROM `leads` WHERE `id` = ?")
	defer stmt.Close()

	HandleErr(queryErr)

	execErr := stmt.QueryRow(id).Scan(&l.Name, &l.Phone, &l.Email, &l.Realtor, &l.InterestedInIds, &l.NeedsAssistance, &l.CreatedAt, &l.ZendeskTicketId)
	l.Id = int64(id)

	return execErr
}

func (l *Lead) SetZendeskTicketId(id int64) error {
	stmt, queryErr := DB.Prepare("UPDATE `leads` SET `zendesk_ticket_id` = ? WHERE `id` = ?")
	defer stmt.Close()

	if queryErr != nil {
		return queryErr
	}

	_, execErr := stmt.Exec(id, l.Id)

	return execErr
}

func (l *Lead) Refresh() error {
	return l.GetById(l.Id)
}

func (l *Lead) Create() (sql.Result, error) {
	stmt, queryErr := DB.Prepare("INSERT INTO `leads` " +
		"(`name`, `phone`, `email`, `needs_assistance`) " +
		"VALUES (?,?,?,?)")
	defer stmt.Close()

	if queryErr != nil {
		log.Println(queryErr)
		return nil, queryErr
	}

	log.Println("lead:")
	log.Println(l)

	result, execErr := stmt.Exec(l.Name, l.Phone, l.Email, l.NeedsAssistance)
	if execErr != nil {
		log.Println(execErr)
	}

	lastId, lastIdErr := result.LastInsertId()
	log.Println(lastIdErr)
	l.Id = lastId

	return result, execErr
}

func (l *Lead) AddInterest(property Property) error {
	stmt, queryErr := DB.Prepare("UPDATE `leads`" +
		"SET `interested_in` = CASE " +
		"WHEN (`interested_in` IS NULL OR `interested_in` = '') " +
		"THEN ? " +
		"ELSE CONCAT(`interested_in`, ',', ?) END " +
		"WHERE `id` = ?")
	defer stmt.Close()

	if queryErr != nil {
		return queryErr
	}

	_, execErr := stmt.Exec(property.Id, property.Id, l.Id)

	return execErr
}

func (l *Lead) Search() error {
	if l.Phone == "" {
		return errors.New("Searching lead should contain name, email and phone number")
	}

	stmt, queryErr := DB.Prepare("SELECT `id` FROM `leads` WHERE `phone` = ?")
	defer stmt.Close()

	if queryErr != nil {
		return queryErr
	}

	scanErr := stmt.QueryRow(l.Phone).Scan(&l.Id)

	return scanErr
}

func (l *Lead) FindByPhone() error {
	if l.Phone == "" {
		return errors.New("Phone number should be filled")
	}

	stmt, queryErr := DB.Prepare("SELECT `id`, `name`, `email`, `interested_in`, `zendesk_ticket_id` FROM `leads` WHERE `phone` = ?")
	defer stmt.Close()

	if queryErr != nil {
		return queryErr
	}

	scanErr := stmt.QueryRow(l.Phone).Scan(&l.Id, &l.Name, &l.Email, &l.InterestedInIds, &l.ZendeskTicketId)

	return scanErr
}

func GetLeadsByRealtorOrFree(realtor int64) ([]Lead, error) {
	stmt, queryErr := DB.Prepare(
		"SELECT `l`.`id` AS `id`, `l`.`name` AS `name`, `l`.`phone` AS `phone`, `l`.`email` AS `email`, " +
			"CASE WHEN `l`.`realtor` IS NULL THEN 0 ELSE `l`.`realtor` END AS `realtor`, " +
			"CASE WHEN `l`.`interested_in` IS NOT NULL THEN `l`.`interested_in` ELSE \"\" END AS `interested_in`, " +
			"CASE WHEN `m2`.`created_at` IS NOT NULL THEN `m2`.`created_at` ELSE CASE WHEN MAX(`e`.`created_at`) IS NOT NULL THEN MAX(`e`.`created_at`) ELSE `l`.`created_at` END END AS `created_at` " +
			"FROM `leads` `l` " +
			"LEFT JOIN " +
			"( " +
			"SELECT `m`.`lead` AS `lead`, MAX(`m`.`id`) AS `id` " +
			"FROM `messages` `m` " +
			"WHERE `m`.`user` = ? " +
			"GROUP BY `m`.`lead` " +
			") `m` ON `m`.`lead` = `l`.`id` " +
			"LEFT JOIN " +
			"( " +
			"SELECT * " +
			"FROM `messages` " +
			") `m2` ON `m`.`id` = `m2`.`id` " +
			"LEFT JOIN `emails` `e` ON `e`.`lead` = `l`.`id` " +
			"WHERE `l`.`realtor` = ? OR `l`.`realtor` IS NULL AND (`e`.`user_id` = ? OR `l`.`id` = 20)" +
			"GROUP BY `l`.`id` " +
			"ORDER BY `created_at` DESC")

	defer stmt.Close()

	if queryErr != nil {
		log.Println(queryErr.Error())
		return nil, queryErr
	}

	rows, scanErr := stmt.Query(realtor, realtor, realtor)
	if scanErr != nil {
		log.Println(scanErr.Error())
	}

	defer rows.Close()

	var leads = make([]Lead, 0)

	for rows.Next() {
		var l Lead

		err := rows.Scan(&l.Id, &l.Name, &l.Phone, &l.Email, &l.Realtor, &l.InterestedInIds, &l.CreatedAt)
		HandleErr(err)

		leads = append(leads, l)
	}

	return leads, scanErr
}

func GetFreeLeads() ([]Lead, error) {
	stmt, queryErr := DB.Prepare(
		"SELECT `l`.`id` AS `id`, `l`.`name` AS `name`, `l`.`phone` AS `phone`, `l`.`email` AS `email`, " +
			"CASE WHEN `l`.`interested_in` IS NOT NULL THEN `l`.`interested_in` ELSE \"\" END AS `interested_in`, " +
			"`e`.`user_id`, `l`.`needs_assistance`, " +
			"CASE WHEN `m2`.`created_at` IS NOT NULL THEN `m2`.`created_at` ELSE CASE WHEN MAX(`e`.`created_at`) IS NOT NULL THEN MAX(`e`.`created_at`) ELSE `l`.`created_at` END END AS `created_at` " +
			"FROM `leads` `l` " +
			"LEFT JOIN " +
			"( " +
			"SELECT `m`.`lead` AS `lead`, MAX(`m`.`id`) AS `id` " +
			"FROM `messages` `m` " +
			"GROUP BY `m`.`lead` " +
			") `m` ON `m`.`lead` = `l`.`id` " +
			"LEFT JOIN " +
			"( " +
			"SELECT * " +
			"FROM `messages` " +
			") `m2` ON `m`.`id` = `m2`.`id` " +
			"LEFT JOIN `emails` `e` ON `e`.`lead` = `l`.`id` " +
			"WHERE `l`.`realtor` IS NULL AND `l`.`id` != 20 " +
			"GROUP BY `l`.`id`, `e`.`user_id` " +
			"ORDER BY `created_at` DESC")

	defer stmt.Close()

	if queryErr != nil {
		log.Println(queryErr.Error())
		return nil, queryErr
	}

	rows, scanErr := stmt.Query()
	if scanErr != nil {
		log.Println(scanErr.Error())
	}

	defer rows.Close()

	var leads = make([]Lead, 0)

	for rows.Next() {
		var l Lead

		err := rows.Scan(&l.Id, &l.Name, &l.Phone, &l.Email, &l.InterestedInIds, &l.Realtor, &l.NeedsAssistance, &l.CreatedAt)
		HandleErr(err)

		leads = append(leads, l)
	}

	return leads, scanErr
}

func (l *Lead) GetProperties() error {
	if l.InterestedInIds == "" {
		return errors.New("Lead should contain properties IDs")
	}

	ids := strings.Split(l.InterestedInIds, ",")

	l.InterestedIn = make([]Property, 0)

	for _, idString := range ids {
		var property Property

		id, _ := strconv.Atoi(idString)
		property.GetById(id)

		l.InterestedIn = append(l.InterestedIn, property)
	}

	return nil
}

func (l *Lead) EnsureLeadCameToUser(user *User) bool {
	email, _ := GetEmailsByLeadAndUser(l, user)

	if email.Id == 0 {
		return false
	}

	return true
}

func (l *Lead) SetOwnership(user *User) error {
	stmt, queryErr := DB.Prepare("UPDATE `leads` SET `realtor` = ? WHERE `id` = ?")
	defer stmt.Close()

	if queryErr != nil {
		return queryErr
	}

	_, execErr := stmt.Exec(user.Id, l.Id)

	return execErr
}

func (l *Lead) SetAssistance(flag bool) error {
	stmt, queryErr := DB.Prepare("UPDATE `leads` SET `needs_assistance` = ? WHERE `id` = ?")
	defer stmt.Close()

	if queryErr != nil {
		return queryErr
	}

	_, execErr := stmt.Exec(flag, l.Id)

	return execErr
}

func (l *Lead) ReleaseOwnership() error {
	stmt, queryErr := DB.Prepare("UPDATE `leads` SET `realtor` = NULL WHERE `id` = ?")
	defer stmt.Close()

	if queryErr != nil {
		return queryErr
	}

	_, execErr := stmt.Exec(l.Id)

	return execErr
}

func (l *Lead) GetMostRecentProperty() (error, Property) {
	ids := strings.Split(l.InterestedInIds, ",")
	lastIdString := ids[len(ids)-1]

	lastId, _ := strconv.Atoi(lastIdString)

	var property Property
	property.GetById(lastId)

	return nil, property
}

func (l *Lead) GetRealtorFromEmail() User {
	stmt, queryErr := DB.Prepare(
		"SELECT `u`.`id`, `u`.`o_email`, `u`.`email`, `u`.`phone`, `u`.`code`, `u`.`fname`, `u`.`sname`, `u`.`created_at` " +
			"FROM `leads` `l` " +
			"LEFT JOIN `emails` `e` ON `e`.`lead` = `l`.`id` " +
			"LEFT JOIN `users` `u` ON `u`.`id` = `e`.`user_id` " +
			"WHERE `l`.`id` = ? " +
			"GROUP BY `u`.`id`")
	defer stmt.Close()

	if queryErr != nil {
		log.Println(queryErr)
	}

	var u User
	stmt.QueryRow(l.Id).Scan(&u.Id, &u.OEmail, &u.Email, &u.Phone, &u.Code, &u.FirstName, &u.SecondName, &u.CreatedAt)

	return u
}

func GetLeadsWhoNeedAssistanceByRealtor(userId int64) []Lead {
	stmt, queryErr := DB.Prepare(
		"SELECT `id`, `name`, `interested_in` " +
			"FROM `leads` " +
			"WHERE `id` IN (SELECT `lead` FROM `emails` WHERE `user_id` = ?) AND `needs_assistance` = 1")
	defer stmt.Close()

	if queryErr != nil {
		log.Println(queryErr)
	}

	rows, _ := stmt.Query(userId)
	defer rows.Close()

	var leads = make([]Lead, 0)

	for rows.Next() {
		var l Lead

		err := rows.Scan(&l.Id, &l.Name, &l.InterestedInIds)
		HandleErr(err)

		leads = append(leads, l)
	}

	return leads
}

func GetSlippingAwayLeadsByRealtor(userId int64) []Lead {
	stmt, queryErr := DB.Prepare(
		"SELECT `l`.`id`, `l`.`name`, `l`.`interested_in` " +
			"FROM `leads` `l` " +
			"LEFT JOIN `messages` `m` ON `m`.`lead` = `l`.`id` AND `m`.`for_lead` = 0 " +
			"WHERE `l`.`id` IN (SELECT `lead` FROM `emails` WHERE `user_id` = ?) " +
			"GROUP BY `l`.`id` " +
			"HAVING MAX(`m`.`created_at`) < DATE_SUB(NOW(), INTERVAL 1 WEEK)")
	defer stmt.Close()

	if queryErr != nil {
		log.Println(queryErr)
	}

	rows, _ := stmt.Query(userId)
	defer rows.Close()

	var leads = make([]Lead, 0)

	for rows.Next() {
		var l Lead

		err := rows.Scan(&l.Id, &l.Name, &l.InterestedInIds)
		HandleErr(err)

		leads = append(leads, l)
	}

	return leads
}

func GetNewLeadsByRealtor(userId int64) []Lead {
	stmt, queryErr := DB.Prepare(
		"SELECT `id`, `name`, `interested_in` " +
			"FROM `leads` " +
			"WHERE `id` IN (SELECT `lead` FROM `emails` WHERE `user_id` = ?) AND `created_at` > DATE_SUB(NOW(), INTERVAL 3 DAY)")
	defer stmt.Close()

	if queryErr != nil {
		log.Println(queryErr)
	}

	rows, _ := stmt.Query(userId)
	defer rows.Close()

	var leads = make([]Lead, 0)

	for rows.Next() {
		var l Lead

		err := rows.Scan(&l.Id, &l.Name, &l.InterestedInIds)
		HandleErr(err)

		leads = append(leads, l)
	}

	return leads
}

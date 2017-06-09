package models

import (
	"database/sql"
	"errors"
)

type Property struct {
	Id          int    `json:"id"`
	Price       string `json:"price"`
	Address     string `json:"address"`
	Description string `json:"description"`
	ImageUrl    string `json:"image_url"`
	Link        string `json:"link"`
	CreatedAt   string `json:"created_at"`
}

func CreatePropertyTable() {
	stmt, queryErr := DB.Prepare("" +
		"CREATE TABLE IF NOT EXISTS `property` (" +
		"`id` INT UNSIGNED NOT NULL AUTO_INCREMENT, " +
		"`price` VARCHAR(32) NOT NULL, " +
		"`address` VARCHAR(512) NOT NULL, " +
		"`description` TEXT NOT NULL, " +
		"`image_url` VARCHAR(1024) NOT NULL, " +
		"`link` VARCHAR(1024) NOT NULL, " +
		"`created_at` TIMESTAMP, PRIMARY KEY (id)" +
		")")
	defer stmt.Close()

	HandleErr(queryErr)

	_, execErr := stmt.Exec()

	HandleErr(execErr)
}

func (p *Property) Create() (sql.Result, error) {
	stmt, queryErr := DB.Prepare("INSERT INTO `property` " +
		"(`address`, `description`, `image_url`, `link`, `price`) " +
		"VALUES (?,?,?,?,?)")
	defer stmt.Close()

	if queryErr != nil {
		return nil, queryErr
	}

	result, execErr := stmt.Exec(p.Address, p.Description, p.ImageUrl, p.Link, p.Price)
	lastId, _ := result.LastInsertId()
	p.Id = int(lastId)

	return result, execErr
}

func (p *Property) Search() error {
	if p.Address == "" || p.Description == "" || p.ImageUrl == "" || p.Link == "" {
		return errors.New("Searching property should contain address, description, image URL and link")
	}

	stmt, queryErr := DB.Prepare("SELECT `id`, `price` FROM `property` WHERE " +
		"`address` = ? AND `description` = ? AND `image_url` = ? AND `link` = ?")
	defer stmt.Close()

	if queryErr != nil {
		return queryErr
	}

	scanErr := stmt.QueryRow(p.Address, p.Description, p.ImageUrl, p.Link).Scan(&p.Id, &p.Price)

	return scanErr
}

func (p *Property) GetById(id int) error {
	stmt, queryErr := DB.Prepare("SELECT `price`, `address`, `description`, `image_url`, `link`, `created_at` " +
		"FROM `property` WHERE `id` = ?")
	defer stmt.Close()

	HandleErr(queryErr)

	execErr := stmt.QueryRow(id).Scan(&p.Price, &p.Address, &p.Description, &p.ImageUrl, &p.Link, &p.CreatedAt)
	p.Id = id

	return execErr
}

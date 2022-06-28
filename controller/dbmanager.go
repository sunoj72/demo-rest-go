package controllers

import (
	"database/sql"
	"errors"
	"log"
	"suno/demo-rest/model"
)

var db sql.DB

func InitDB(file string) (*sql.DB, error) {
	db_internal, err := sql.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}

	db = *db_internal

	createMemberTableQuery := `
	CREATE TABLE IF NOT EXISTS members ( 
		id text NOT NULL PRIMARY KEY,
		name text,
		email text
		)
	`

	_, err = db_internal.Exec(createMemberTableQuery)
	if err != nil {
		return nil, err
	}

	createFavoritesTableQuery := `
		CREATE TABLE IF NOT EXISTS favorites (
		id text NOT NULL,
		favorite text,
		UNIQUE (id, favorite)
		)
	`

	_, err = db_internal.Exec(createFavoritesTableQuery)
	if err != nil {
		return nil, err
	}

	err = AddMember("sunoj", "Soonho Kim", "sunoj@lgcns.com")
	if err != nil {
		return nil, err
	}

	myFavorites := make([]string, 3)
	myFavorites[0] = "eating"
	myFavorites[1] = "sleeping"
	myFavorites[2] = "watching"

	err = AddFavorites("sunoj", myFavorites)
	if err != nil {
		return nil, err
	}

	return db_internal, nil
}

func AddMember(id string, name string, email string) error {
	if id == "" {
		return errors.New("InvalidArguments Error in AddMember")
	}

	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("INSERT INTO members (id, name, email) VALUES (?,?,?)")
	_, err := stmt.Exec(id, name, email)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	tx.Commit()
	return nil
}

func AddFavorites(id string, favorites []string) error {
	if id == "" || favorites == nil {
		return errors.New("InvalidArguments Error in AddFavorites")
	}

	tx, _ := db.Begin()

	stmt, _ := tx.Prepare("INSERT INTO favorites (id, favorite) VALUES (?,?)")

	for i := 0; i < len(favorites); i++ {
		_, err := stmt.Exec(id, favorites[i])
		if err != nil {
			log.Println(err.Error())
			return err
		}
	}

	tx.Commit()

	return nil
}

func GetMembers() ([]model.ShortMember, error) {
	var member model.ShortMember
	var members []model.ShortMember

	rows, err := db.Query("SELECT id, name, email FROM members")
	if err != nil {
		return nil, err
	}

	var cnt int = 0
	for rows.Next() {
		cnt++
		if err := rows.Scan(&member.Id, &member.Name, &member.Email); err != nil {
			return nil, err
		}

		members = append(members, member)
	}

	if cnt == 0 {
		return nil, errors.New("404")
	}

	defer rows.Close()

	return members, nil
}

func GetMember(id string) (model.Member, error) {
	var member model.Member
	var favorite model.Favorite
	var favorites model.Favorites

	row := db.QueryRow("select id, name, email from members where id = $1", id)
	if err := row.Scan(&member.Id, &member.Name, &member.Email); err != nil {
		if err == sql.ErrNoRows {
			return model.Member{}, errors.New("404")
		} else {
			return model.Member{}, err
		}
	}

	rows, err := db.Query("select id, favorite from favorites where id = $1", id)
	if err != nil {
		return member, err
	}

	for rows.Next() {
		if err := rows.Scan(&favorite.Id, &favorite.Favorite); err != nil {
			return member, err
		}

		favorites = append(favorites, favorite.Favorite)
	}

	member.Favorites = favorites

	defer rows.Close()

	return member, nil
}

func UpdateMember(id string, name string, email string, favorites model.Favorites) error {
	var member model.Member

	row := db.QueryRow("SELECT id, name, email FROM members WHERE id = $1", id)
	err := row.Scan(&member.Id, &member.Name, &member.Email)
	if err != nil && err == sql.ErrNoRows {
		return errors.New("204")
	}

	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("UPDATE members SET name = ?, email = ? WHERE id = ?")
	_, err = stmt.Exec(name, email, id)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	stmt, _ = tx.Prepare("DELETE FROM favorites WHERE id = ?")
	_, err = stmt.Exec(id)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	stmt, _ = tx.Prepare("INSERT INTO favorites (id, favorite) VALUES (?, ?)")
	for i := 0; i < len(favorites); i++ {
		_, err = stmt.Exec(id, favorites[i])
		if err != nil {
			log.Println(err.Error())
			return err
		}
	}

	tx.Commit()

	return nil
}

func DeleteMember(id string) error {
	var member model.Member

	row := db.QueryRow("SELECT id, name, email FROM members WHERE id = $1", id)
	err := row.Scan(&member.Id, &member.Name, &member.Email)
	if err != nil && err == sql.ErrNoRows {
		return errors.New("204")
	}

	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("DELETE FROM members WHERE id = ?")
	_, err = stmt.Exec(id)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	stmt, _ = tx.Prepare("DELETE FROM favorites WHERE id = ?")
	_, err = stmt.Exec(id)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	tx.Commit()

	return nil
}

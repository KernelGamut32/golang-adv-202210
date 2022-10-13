package service

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/KernelGamut32/gameserver/internal/users"
	database "github.com/KernelGamut32/gameserver/internal/users/db"

	"golang.org/x/crypto/bcrypt"
)

type UsersDB struct {
	*sql.DB
}

func GetUsersDataStore() users.UserDatastore {
	return &UsersDB{database.Get()}
}

func (db *UsersDB) CreateUser(user *users.User) error {
	if user.Email == "" || user.Password == "" || user.Name == "" {
		return errors.New("need values for all fields")
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return errors.New("password encryption failed")
	}
	user.Password = string(pass)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	result, err := db.ExecContext(ctx, "insert into users (name, email, password) values (?, ?, ?)",
		user.Name, user.Email, user.Password)

	if err != nil {
		return err
	}

	id, e := result.LastInsertId()
	if e != nil {
		return e
	}

	user.ID = uint(id)

	return nil
}

func (db *UsersDB) GetAllUsers() ([]users.User, error) {
	var theUsers []users.User

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	rows, err := db.QueryContext(ctx, "select id, name, email, password from users")
	if err != nil {
		log.Print("error occurred in GetAllUsers ", err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user users.User
		rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
		theUsers = append(theUsers, user)
	}
	return theUsers, nil
}

func (db *UsersDB) FindUser(email, password string) (*users.User, error) {
	user := &users.User{}

	if email == "" || password == "" {
		return nil, errors.New("email and password are required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	row := db.QueryRowContext(ctx, "select id, name, email, password from users where email = ?",
		email)

	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	if err == sql.ErrNoRows {
		return nil, err
	}

	errf := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errf != nil { //Password does not match!
		return nil, errors.New("invalid login credentials")
	}

	return user, nil
}

func (db *UsersDB) UpdateUser(id string, user *users.User) error {
	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return errors.New("password encryption failed")
	}
	user.Password = string(pass)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	result, err := db.ExecContext(ctx, "update users set name = ?, email = ?, password = ? where id = ?",
		user.Name, user.Email, user.Password, id)
	if err != nil {
		log.Print("error occurred in UpdateUser ", err.Error())
		return err
	}
	num, err := result.RowsAffected()
	if err != nil {
		log.Fatal("could not update database ", err.Error())
		return err
	}

	log.Println("number of rows affected is ", num)
	return nil
}

func (db *UsersDB) DeleteUser(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	result, err := db.ExecContext(ctx, "delete from users where id = ?", id)
	if err != nil {
		log.Print("error occurred in DeleteUser ", err.Error())
		return err
	}
	_, err = result.RowsAffected()
	if err != nil {
		log.Fatal("could not update database ", err.Error())
		return err
	}
	return nil
}

func (db *UsersDB) GetUser(id string) (users.User, error) {
	var user users.User

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	row := db.QueryRowContext(ctx, "select id, name, email, password from users where id = ?", id)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	if err == sql.ErrNoRows {
		return users.User{}, err
	}
	return user, nil
}

package db

import (
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username            string `gorm:"unique"`
	Email               string `gorm:"unique"`
	Password            string
	Offices             []Office `gorm:"many2many:user_offices;"`
	Rankings            []Ranking
	MatchParticipations []MatchParticipant
	Approvals           []MatchApproval
}

type CreateUserErrors struct {
	Username error
	Email    error
	Error    error
}

func (d *database) CreateUser(username, email, password string) *CreateUserErrors {
	user := &User{Username: username, Email: email, Password: password}
	err := d.c.Create(user).Error
	if err != nil {
		postgresErr, ok := err.(*pgconn.PgError)
		if !ok {
			return &CreateUserErrors{Error: err}
		}

		if postgresErr.SQLState() == "23505" {
			constaintArray := strings.Split(postgresErr.ConstraintName, "_")
			columnName := constaintArray[len(constaintArray)-1]
			if columnName == "username" {
				return &CreateUserErrors{Username: errors.New("Username is taken")}
			}
			if columnName == "email" {
				return &CreateUserErrors{Email: errors.New("Email is taken")}
			}
		}

		return &CreateUserErrors{Error: err}
	}

	return nil
}

func (d *database) GetUserByUsername(username string) (*User, error) {
	var user User
	err := d.c.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

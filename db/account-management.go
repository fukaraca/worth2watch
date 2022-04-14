package db

import (
	"fmt"
	"github.com/fukaraca/worth2watch/model"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"golang.org/x/net/context"
	"log"
)

//QueryLogin queries the password for given username and returns hashed-password or error depending on the result
func QueryLogin(c *gin.Context, username string) (string, error) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), TIMEOUT)
	defer cancel()
	result, err := Conn.Query(ctx, "SELECT password FROM users WHERE username LIKE $1;", username)
	defer result.Close()
	if err != nil {
		log.Println("login query for password failed:", err)
	}

	password := ""
	for result.Next() {
		if err := result.Scan(&password); err == pgx.ErrNoRows {
			return "", fmt.Errorf("username not found")
		} else if err == nil {
			return password, nil
		}
	}
	return "", err
}

//IsAdmin checks DB for given user whether he/she is admin or not
func IsAdmin(c *gin.Context, username string) (bool, error) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), TIMEOUT)
	defer cancel()
	isAdmin := false
	result := Conn.QueryRow(ctx, "SELECT isadmin FROM users WHERE username = $1;", username)
	err := result.Scan(&isAdmin)
	if err != nil {
		log.Println("login query for password failed:", err)
		return false, err
	}

	if isAdmin {
		return true, nil
	}
	return false, nil
}

//CreateNewUser simply inserts new user contents to DB
func CreateNewUser(c *gin.Context, newUser *model.User) error {
	ctx, cancel := context.WithTimeout(c.Request.Context(), model.TIMEOUT)
	defer cancel()

	_, err := Conn.Exec(ctx, "INSERT INTO users (user_id,username,password,email,name,lastname,createdon,lastlogin,isadmin)  VALUES (nextval('users_user_id_seq'),$1,$2,$3,$4,$5,$6,$7,$8);", newUser.Username, newUser.Password, newUser.Email, newUser.Name, newUser.Lastname, newUser.CreatedOn, newUser.LastLogin, newUser.Isadmin)

	if err != nil {
		return fmt.Errorf("user infos for register was failed to insert to DB:%v", err)
	}
	return nil
}

//QueryUserInfo returns user info from db except password
func QueryUserInfo(c *gin.Context, username string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), TIMEOUT)
	defer cancel()

	tempUser := new(model.User)
	row := Conn.QueryRow(ctx, "SELECT * FROM users WHERE username = $1;", username)
	err := row.Scan(&tempUser.UserID, &tempUser.Username, &tempUser.Password, &tempUser.Email, &tempUser.Name, &tempUser.Lastname, &tempUser.CreatedOn, &tempUser.LastLogin, &tempUser.Isadmin)
	if err != nil {
		return nil, fmt.Errorf("scanning the user infos from DB was failed:%v", err)
	}
	tempUser.Password = ""
	return tempUser, nil
}

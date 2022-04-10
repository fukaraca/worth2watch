package api

import (
	"context"
	"github.com/fukaraca/worth2watch/auth"
	"github.com/fukaraca/worth2watch/db"
	"github.com/fukaraca/worth2watch/model"
	"github.com/fukaraca/worth2watch/util"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

//Auth is the authentication middleware
func Auth(fn gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("Req ID:", requestid.Get(c))
		if !auth.CheckSession(c) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"notification": "you must login",
			})
			c.Abort()
			return
		}

		fn(c)
	}
}

//CheckRegistration is for registeration of new user or admin.
//Data must be POSTed as "form data".
//Leading or trailing whitespaces will be handled by frontend
func CheckRegistration(c *gin.Context) {
	if auth.CheckSession(c) {
		c.JSON(http.StatusBadRequest, gin.H{
			"notification": "already logged in",
		})
		return
	}
	newUser := new(model.User)
	err := c.BindJSON(newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"notification": err.Error(),
		})
		log.Println("new user info as JSON couldn't be binded:", err)
		return
	}
	if newUser.Username == "" || newUser.Email.String == "" || newUser.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"notification": "username or email or password cannot be empty",
		})
		log.Println("failed due to 'username or email or password cannot be empty':")
		return
	}
	newUser.Username = *util.Striper(newUser.Username)
	newUser.Email.String = *util.Striper(newUser.Email.String)
	pass, err := util.HashPassword(*util.Striper(newUser.Password))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"notification": err.Error(),
		})
		log.Println("new user's password couldn't be hashed:", err)
		return
	}
	newUser.Password = pass
	newUser.CreatedOn.Time = time.Now()
	newUser.LastLogin.Time = time.Now()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"notification": err.Error(),
		})
		log.Println("new user info creation time couldn't be assigned:", err)
		return
	}
	//insert new users if not exist
	ctx, cancel := context.WithTimeout(c.Request.Context(), model.TIMEOUT)
	defer cancel()
	_, err = db.Conn.Exec(ctx, "INSERT INTO users (user_id,username,password,email,name,lastname,createdon,lastlogin,isadmin) VALUES (nextval('users_user_id_seq'),$1,$2,$3,$4,$5,$6,$7,$8);", newUser.Username, newUser.Password, newUser.Email, newUser.Name, newUser.Lastname, newUser.CreatedOn.Time, newUser.LastLogin.Time, newUser.Isadmin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"notification": err.Error(),
		})
		log.Println("user infos for register was failed to insert to DB:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"notification": "account created successfully",
	})
}

func ForgotPassword(c *gin.Context) {
}

//Login is handler function for login process.
//It requires form-data for 'logUsername' and 'logPassword' keys.
func Login(c *gin.Context) {
	if auth.CheckSession(c) {
		c.JSON(http.StatusBadRequest, gin.H{
			"notification": "already logged in",
		})
		return
	}
	logUsername := c.PostForm("logUsername")
	logPassword := c.PostForm("logPassword")

	hashedPass, err := db.QueryLogin(c, logUsername)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"notification": err.Error(),
		})
		log.Println("password couldn't be fetched from DB:", err)
		return
	}
	if !util.CheckPasswordHash(logPassword, hashedPass) {
		c.JSON(http.StatusBadRequest, gin.H{
			"notification": "password or username is incorrect"})
		log.Println("password or username is incorrect")
		return
	}
	lastLoginTime := time.Now()

	ctx, cancel := context.WithTimeout(c.Request.Context(), model.TIMEOUT)
	defer cancel()
	_, err = db.Conn.Exec(ctx, "UPDATE users SET lastlogin = $1 WHERE username = $2;", lastLoginTime, logUsername)
	if err != nil {
		log.Println("last login time update has error:", err)
	}
	auth.CreateSession(logUsername, c)
	log.Printf("%s has logged in:\n", logUsername)
	c.JSON(http.StatusOK, gin.H{
		"notification": "user " + logUsername + " successfully logged in",
	})
}

//Logout handler
func Logout(c *gin.Context) {
	ok, err := auth.DeleteSession(c)
	if err != nil || !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"notification": err.Error()})
		log.Println("session couldn't be deleted:", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"notification": "logged out successfully",
	})
}

//UpdateUser is handler for updating user/admins informations. Changing admin/user role was not implemented.
func UpdateUser(c *gin.Context) {
	firstname := *util.Striper(c.PostForm("firstname"))
	lastname := *util.Striper(c.PostForm("lastname"))

	username, _ := c.Cookie("uid")

	ctx, cancel := context.WithTimeout(c.Request.Context(), model.TIMEOUT)
	defer cancel()
	_, err := db.Conn.Exec(ctx, "UPDATE users SET name = $1,lastname=$2 WHERE username = $3;", firstname, lastname, username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"notification": "user information update failed",
		})
		log.Println("user information update has failed:", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"notification": "user informations updated succesfully",
	})
}

//GetUserInfo handles GET request for user infos. Only admin and the user can get the info
func GetUserInfo(c *gin.Context) {
	username, _ := c.Cookie("uid")
	usernameP := c.Param("username")
	//only user and admin may peek user info
	if usernameP != username && !auth.CheckAdminForLoggedIn(c, username) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"notification": "you are not allowed to see another users info",
		})
		return
	}
	user, err := db.QueryUserInfo(c, usernameP)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"notification": "user information get failed",
		})
		log.Println("user information get has failed:", err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func Movies(c *gin.Context) {

}

func AddMovie(c *gin.Context) {

}

func Series(c *gin.Context) {

}

func Seasons(c *gin.Context) {

}

func Episodes(c *gin.Context) {

}

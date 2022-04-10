package api

import (
	"github.com/fukaraca/worth2watch/model"
)

func Endpoints() {
	//todo    vvvv
	model.R.POST("/addMovie", AddMovie)
	model.R.GET("/movies", Movies)
	model.R.GET("/series", Series)
	model.R.GET("/seasons", Seasons)
	model.R.GET("episodes", Episodes)
	///todo  ^^^^
	model.R.GET("/user/:username", Auth(GetUserInfo))

	///
	model.R.POST("/register", CheckRegistration)
	model.R.POST("/login", Login)
	model.R.PATCH("/updateUser", Auth(UpdateUser))
	model.R.POST("/logout", Auth(Logout))
}

package api

import (
	"github.com/fukaraca/worth2watch/model"
)

func Endpoints() {
	//todo    vvvv

	model.R.GET("/movies", Movies)
	model.R.GET("/series", Series)
	model.R.GET("/seasons", Seasons)
	model.R.GET("episodes", Episodes)
	///todo  ^^^^

	//content management
	model.R.POST("/addContentByID", Auth(AddContentByID))
	model.R.POST("/addContentWithJSON", Auth(AddContentWithJSON))
	model.R.DELETE("/deleteMovieByID", Auth(DeleteContentByID))
	//account management
	model.R.GET("/user/:username", Auth(GetUserInfo))
	model.R.POST("/register", CheckRegistration)
	model.R.POST("/login", Login)
	model.R.PATCH("/updateUser", Auth(UpdateUser))
	model.R.POST("/logout", Auth(Logout))
}

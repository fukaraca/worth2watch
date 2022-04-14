package api

import (
	"github.com/fukaraca/worth2watch/model"
)

func Endpoints() {
	//todo    vvvv

	///todo  ^^^^
	//public
	model.R.GET("/movies/:id", GetThisMovie)
	model.R.GET("/movies/list", GetMoviesWithPage)
	model.R.GET("/searchContent", SearchContent)
	model.R.GET("/series/:seriesid", GetThisSeries)
	model.R.GET("/series/list", GetSeriesWithPage)
	model.R.GET("/series/:seriesid/:season", GetEpisodesForaSeason)
	model.R.GET("/getSimilarContent", GetSimilarContent)

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

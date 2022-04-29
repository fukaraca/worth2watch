package api

import (
	"github.com/gin-gonic/gin"
)

var R *gin.Engine

func Endpoints() {

	//public
	R.GET("/movies/:id", GetThisMovie)
	R.GET("/movies/list", GetMoviesWithPage)
	R.GET("/searchContent", SearchContent)
	R.GET("/series/:seriesid", GetThisSeries)
	R.GET("/series/list", GetSeriesWithPage)
	R.GET("/series/:seriesid/:season", GetEpisodesForaSeason)
	R.GET("/getSimilarContent", GetSimilarContent)
	//user accessed
	R.POST("/addFavorites", Auth(AddToFavorites))
	R.GET("/getFavorites", Auth(GetFavorites))
	R.GET("/searchFavorites", Auth(SearchFavorites))
	//content management
	R.POST("/addContentByID", Auth(AddContentByID))
	R.POST("/addContentWithJSON", Auth(AddContentWithJSON))
	R.DELETE("/deleteMovieByID", Auth(DeleteContentByID))
	//account management
	R.GET("/user/:username", Auth(GetUserInfo))
	R.POST("/register", CheckRegistration)
	R.POST("/login", Login)
	R.PATCH("/updateUser", Auth(UpdateUser))
	R.POST("/logout", Auth(Logout))
}

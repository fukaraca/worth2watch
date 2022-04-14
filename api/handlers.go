package api

import (
	"context"
	"encoding/json"
	"github.com/fukaraca/worth2watch/auth"
	"github.com/fukaraca/worth2watch/db"
	"github.com/fukaraca/worth2watch/model"
	"github.com/fukaraca/worth2watch/util"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"log"
	"net/http"
	"strconv"
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

//CheckRegistration is a func for registering a new user or admin.
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
	if newUser.Username == "" || *newUser.Email == "" || newUser.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"notification": "username or email or password cannot be empty",
		})
		log.Println("failed due to 'username or email or password cannot be empty':")
		return
	}

	newUser.Username = *util.Striper(newUser.Username)
	newUser.Email = util.Striper(*newUser.Email)
	pass, err := util.HashPassword(*util.Striper(newUser.Password))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"notification": err.Error(),
		})
		log.Println("new user's password couldn't be hashed:", err)
		return
	}
	newUser.Password = pass
	err = newUser.CreatedOn.Set(time.Now())
	newUser.LastLogin.Set(time.Now())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"notification": err.Error(),
		})
		log.Println("new user info creation time couldn't be assigned:", err)
		return
	}

	err = db.CreateNewUser(c, newUser)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"notification": err.Error(),
		})
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

//AddContentByID is handler func for content adding by IMDB ID and contents type.
//It requires "movie" or "series" for content-type key.
//Insertion is being maintained asynchronously due to expensive amount of time that was required to be got all data related to series.
func AddContentByID(c *gin.Context) {
	username, _ := c.Cookie("uid")
	if !auth.CheckAdminForLoggedIn(c, username) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"notification": "unauthorized attempt",
		})
		log.Println("unauthorized attempt by user:", username)
		return
	}
	IMDBID := c.PostForm("imdb-id")
	contentType := c.PostForm("content-type")
	switch contentType {
	case "movie":
		go db.AddMovieContentWithID(IMDBID)

	case "series":
		go db.AddSeriesContentWithID(IMDBID)

	default:
		log.Println("invalid content-type:", contentType)
		c.JSON(http.StatusBadRequest, gin.H{
			"notification": "invalid content-type: " + contentType,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"notification": "content will be inserted shortly after",
	})
}

//AddContentWithJSON handles adding new content with JSON format.
//Content-type "movie" or "series" must be provided.
//This is not practical at all.
//However, for an internal movie database that includes content that is not contained in IMDB, it can be useful with an additional struct field exposes e.g. BetterIMDB_ID..
func AddContentWithJSON(c *gin.Context) {
	username, _ := c.Cookie("uid")
	if !auth.CheckAdminForLoggedIn(c, username) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"notification": "unauthorized attempt",
		})
		log.Println("unauthorized attempt by user:", username)
		return
	}
	/*	inputJSON := struct {
		ContentType string `json:"content-type"`
		RawData     *model.Movie `json:"content-raw-data"`
	}{}*/

	contentType := c.PostForm("content-type")
	inputJSON := c.PostForm("content-raw-data")

	switch contentType {
	case "movie":
		movie := new(model.Movie)
		err := json.Unmarshal([]byte(inputJSON), movie)
		if err != nil {
			log.Println("movie content couldn't be added: ", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"notification": "movie content couldn't be added. error: " + err.Error(),
			})
			return
		}
		err = db.AddMovieContentWithStruct(c, movie)
		if err != nil {
			log.Println("movie content couldn't be added: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"notification": "movie content couldn't be added. error: " + err.Error(),
			})
			return
		}
	case "series":
		series := new(model.Series)
		err := json.Unmarshal([]byte(inputJSON), series)
		if err != nil {
			log.Println("series content couldn't be added: ", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"notification": "series content couldn't be added. error: " + err.Error(),
			})
			return
		}
		err = db.AddSeriesContentWithStruct(c, series)
		if err != nil {
			log.Println("series content couldn't be added: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"notification": "series content couldn't be added. error: " + err.Error(),
			})
			return
		}
	default:
		log.Println("invalid content-type:", contentType)
		c.JSON(http.StatusBadRequest, gin.H{
			"notification": "invalid content-type: " + contentType,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"notification": "content has been created successfully",
	})
}

//DeleteContentByID deletes content for given IMDB id. Content-type must be provided
func DeleteContentByID(c *gin.Context) {
	username, _ := c.Cookie("uid")
	if !auth.CheckAdminForLoggedIn(c, username) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"notification": "unauthorized attempt",
		})
		log.Println("unauthorized attempt by user:", username)
		return
	}
	IMDBID := c.PostForm("imdb-id")
	contentType := c.PostForm("content-type")

	if !(contentType == "movie" || contentType == "series") {
		log.Println("invalid content-type:", contentType)
		c.JSON(http.StatusBadRequest, gin.H{
			"notification": "invalid content-type: " + contentType,
		})
		return
	}

	err := db.DeleteContent(c, IMDBID, contentType)
	if err != nil {
		log.Println("content ", IMDBID, " couldn't be deleted: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"notification": "content " + IMDBID + " couldn't be deleted. error: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"notification": "content has been deleted successfully",
	})
}

//GetThisMovie is a handler function for responsing a specific movie details
func GetThisMovie(c *gin.Context) {
	id := c.Param("id")
	movie, err := db.GetThisMovieFromDB(c, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"notification": "no such movie",
			})
			return
		} else {
			log.Println("get this movie failed for id: ", id, " err: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"notification": err,
			})
			return
		}
	}

	c.JSON(http.StatusOK, movie)
}

//GetThisSeries is a handler function for responsing a specific serie with its seasons
func GetThisSeries(c *gin.Context) {
	id := c.Param("seriesid")
	series, seasons, err := db.GetThisSeriesFromDB(c, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"notification": "no such serie",
			})
			return
		} else if err.Error() == "there is no season for given series" {
			c.JSON(http.StatusOK, gin.H{
				"series":       series,
				"notification": err.Error(),
			})
		} else {
			log.Println("get this movie failed for id: ", id, " err: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"notification": err,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"series":  series,
		"seasons": seasons,
	})
}

//GetEpisodesForaSeason is a handle function for responsing episodes for a certain season of a series
func GetEpisodesForaSeason(c *gin.Context) {
	seriesid := c.Param("seriesid")
	seasonNumber := c.Param("season")

	season, err := db.GetEpisodesForaSeasonFromDB(c, seriesid, seasonNumber)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"notification": "no such season",
			})
			return
		} else {
			log.Println("get this season failed for serie: ", seriesid, " season:", seasonNumber, " err: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"notification": err,
			})
			return
		}
	}

	c.JSON(http.StatusOK, season)
}

//GetMoviesWithPage handles request for given amount of item and page
func GetMoviesWithPage(c *gin.Context) {
	var err error
	q := c.Request.URL.Query()
	page := 1
	items := 10

	if q.Has("page") {
		page, err = strconv.Atoi(q.Get("page"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"notification": "invalid page format",
			})
			return
		}
	}
	if q.Has("items") {
		items, err = strconv.Atoi(q.Get("items"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"notification": "invalid items format",
			})
			return
		}
	}

	movies, err := db.GetMoviesListWithPage(c, page, items)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusOK, gin.H{
				"notification": "end of the list",
			})
			return
		} else {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"notification": err,
			})
			return
		}
	}
	c.JSON(http.StatusOK, movies)

}

//GetSeriesWithPage handles request for given amount of item and page
func GetSeriesWithPage(c *gin.Context) {
	var err error
	q := c.Request.URL.Query()
	page := 1
	items := 10

	if q.Has("page") {
		page, err = strconv.Atoi(q.Get("page"))
		if err != nil || page < 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"notification": "invalid page format",
			})
			return
		}
	}
	if q.Has("items") {
		items, err = strconv.Atoi(q.Get("items"))
		if err != nil || items < 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"notification": "invalid items format",
			})
			return
		}
	}

	series, err := db.GetSeriesListWithPage(c, page, items)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusOK, gin.H{
				"notification": "end of the list",
			})
			return
		} else {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"notification": err,
			})
			return
		}
	}
	c.JSON(http.StatusOK, series)

}

//SearchContent is handler function for searching movies/series by name and genres.
//Page and item amount for movie or series on a page must be provided
func SearchContent(c *gin.Context) {
	var err error
	q := c.Request.URL.Query()
	name := ""
	if q.Has("name") {
		if len(q["name"]) > 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"notification": "only one name can be accepted",
			})
			return
		}
		name = q.Get("name")
	}
	genres := []string{}
	if q.Has("genre") {

		genres = q["genre"]
	}
	page := 1
	items := 10

	if q.Has("page") {
		page, err = strconv.Atoi(q.Get("page"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"notification": "invalid page format",
			})
			return
		}
	}
	if q.Has("items") {
		items, err = strconv.Atoi(q.Get("items"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"notification": "invalid items format",
			})
			return
		}
	}

	movies, series, err := db.SearchContent(c, name, genres, page, items)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"notification": "en error occurred",
		})
		return
	}
	if movies == nil && series == nil {
		c.JSON(http.StatusOK, gin.H{
			"notification": "no content found for given filter",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"movies": movies,
		"series": series,
	})
}

//GetSimilarContent handles for similar content request. Similarity evaluation is made by genre tags and sorted descending of rating
func GetSimilarContent(c *gin.Context) {
	IMDBID := c.PostForm("imdb-id")
	contentType := c.PostForm("content-type")

	if !(contentType == "movie" || contentType == "series") {
		log.Println("invalid content-type:", contentType)
		c.JSON(http.StatusBadRequest, gin.H{
			"notification": "invalid content-type: " + contentType,
		})
		return
	}
	movies, series, err := db.FindSimilarContent(c, IMDBID, contentType)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusOK, gin.H{
				"notification": "no match of content",
			})
			return
		} else {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"notification": err,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"movies": movies,
		"series": series,
	})
}

//AddToFavorites adds content to users favorites.
func AddToFavorites(c *gin.Context) {
	id := c.PostForm("imdb-id")
	contentType := c.PostForm("content-type")

	if !(contentType == "movie" || contentType == "series") {
		log.Println("invalid content-type:", contentType)
		c.JSON(http.StatusBadRequest, gin.H{
			"notification": "invalid content-type: " + contentType,
		})
		return
	}
	err := db.AddContentToFavorites(c, id, contentType)
	if err != nil {
		log.Println("content add failed: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"notification": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"notification": "content has been added successfully",
	})
}

//GetFavorites ...
func GetFavorites(c *gin.Context) {
	var err error
	q := c.Request.URL.Query()
	page := 1
	items := 10

	if q.Has("page") {
		page, err = strconv.Atoi(q.Get("page"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"notification": "invalid page format",
			})
			return
		}
	}
	if q.Has("items") {
		items, err = strconv.Atoi(q.Get("items"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"notification": "invalid items format",
			})
			return
		}
	}

	movies, series, err := db.GetFavoriteContents(c, page, items)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"notification": "en error occurred",
		})
		return
	}
	if movies == nil && series == nil {
		c.JSON(http.StatusOK, gin.H{
			"notification": "there is no favorite item",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"movies": movies,
		"series": series,
	})
}

//SearchFavorites is handler function for searching on favorites
func SearchFavorites(c *gin.Context) {
	var err error
	q := c.Request.URL.Query()
	name := ""
	if q.Has("name") {
		if len(q["name"]) > 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"notification": "only one name can be accepted",
			})
			return
		}
		name = q.Get("name")
	}
	genres := []string{}
	if q.Has("genre") {

		genres = q["genre"]
	}
	page := 1
	items := 10

	if q.Has("page") {
		page, err = strconv.Atoi(q.Get("page"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"notification": "invalid page format",
			})
			return
		}
	}
	if q.Has("items") {
		items, err = strconv.Atoi(q.Get("items"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"notification": "invalid items format",
			})
			return
		}
	}

	movies, series, err := db.SearchFavorites(c, name, genres, page, items)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"notification": "en error occurred",
		})
		return
	}
	if movies == nil && series == nil {
		c.JSON(http.StatusOK, gin.H{
			"notification": "no content found for given filter",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"movies": movies,
		"series": series,
	})
}

package model

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgtype"
	"time"
)

var R *gin.Engine

//os.Getenv("key_string")
var TIMEOUT = 5 * time.Second
var ServerHost = "localhost"
var ServerPort = ":8080"

type User struct {
	UserID    int                `db.Conn:"user_id" json:"userID,omitempty"`
	Username  string             `db.Conn:"username" json:"username"`
	Password  string             `db.Conn:"password" json:"password"`
	Email     pgtype.Text        `db.Conn:"email" json:"email"`
	Name      pgtype.Text        `db.Conn:"name" json:"name"`
	Lastname  pgtype.Text        `db.Conn:"name" json:"lastname"`
	CreatedOn pgtype.Timestamptz `db.Conn:"createdon" json:"createdOn"`
	LastLogin pgtype.Timestamptz `db.Conn:"lastlogin" json:"lastLogin"`
	Isadmin   bool               `db.Conn:"isadmin" json:"isAdmin"`
}

type Movie struct {
	MovieID     int              `db.Conn:"movie_id" json:"movieID,omitempty"`
	Title       pgtype.Text      `db.Conn:"title" json:"title"`
	Description pgtype.Text      `db.Conn:"description" json:"description"`
	Rating      pgtype.Numeric   `db.Conn:"rating" json:"rating"`
	ReleaseDate pgtype.Timestamp `db.Conn:"release_date" json:"releaseDate"`
	Director    []pgtype.Text    `db.Conn:"director" json:"director,omitempty"`
	Writer      []pgtype.Text    `db.Conn:"writer" json:"writer,omitempty"`
	Stars       []pgtype.Text    `db.Conn:"stars" json:"stars,omitempty"`
	Duration    int              `db.Conn:"duration_min" json:"duration,omitempty"`
	IMDBid      pgtype.Text      `db.Conn:"imdb_id" json:"IMDBid"`
	Year        int              `db.Conn:"year" json:"year,omitempty"`
	Genre       pgtype.Text      `db.Conn:"genre" json:"genre"`
	Audio       []pgtype.Text    `db.Conn:"audio" json:"audio,omitempty"`
	Subtitles   []pgtype.Text    `db.Conn:"subtitles" json:"subtitles,omitempty"`
}

type Series struct {
	SerieID     int              `db.Conn:"serie_id" json:"serieID,omitempty"`
	Title       pgtype.Text      `db.Conn:"title" json:"title"`
	Description pgtype.Text      `db.Conn:"description" json:"description"`
	Rating      pgtype.Numeric   `db.Conn:"rating" json:"rating"`
	ReleaseDate pgtype.Timestamp `db.Conn:"release_date" json:"releaseDate"`
	Director    []pgtype.Text    `db.Conn:"director" json:"director,omitempty"`
	Writer      []pgtype.Text    `db.Conn:"writer" json:"writer,omitempty"`
	Stars       []pgtype.Text    `db.Conn:"stars" json:"stars,omitempty"`
	Duration    int              `db.Conn:"duration_min" json:"duration,omitempty"`
	IMDBid      pgtype.Text      `db.Conn:"imdb_id" json:"IMDBid"`
	Year        int              `db.Conn:"year" json:"year,omitempty"`
	Genre       pgtype.Text      `db.Conn:"genre" json:"genre"`
	Seasons     int              `db.Conn:"seasons" json:"seasons"`
}

type Seasons struct {
	SeasonID     int `db.Conn:"season_id" json:"seasonID,omitempty"`
	SeasonNumber int `db.Conn:"season_number" json:"seasonNumber,omitempty"`
	Episodes     int `db.Conn:"episodes" json:"episodes,omitempty"`
	SerieID      int `db.Conn:"serie_id" json:"serieID,omitempty" json:"serieID,omitempty"`
}

type Episodes struct {
	EpisodeID   int              `db.Conn:"episode_id" json:"episodeID,omitempty"`
	Title       pgtype.Text      `db.Conn:"title" json:"title"`
	Description pgtype.Text      `db.Conn:"description" json:"description"`
	Rating      pgtype.Numeric   `db.Conn:"rating" json:"rating"`
	ReleaseDate pgtype.Timestamp `db.Conn:"release_date" json:"releaseDate"`
	Director    []pgtype.Text    `db.Conn:"director" json:"director,omitempty"`
	Writer      []pgtype.Text    `db.Conn:"writer" json:"writer,omitempty"`
	Stars       []pgtype.Text    `db.Conn:"stars" json:"stars,omitempty"`
	Duration    int              `db.Conn:"duration_min" json:"duration,omitempty"`
	IMDBid      pgtype.Text      `db.Conn:"imdb_id" json:"IMDBid"`
	Year        int              `db.Conn:"year" json:"year,omitempty"`
	Audio       []pgtype.Text    `db.Conn:"audio" json:"audio,omitempty"`
	Subtitles   []pgtype.Text    `db.Conn:"subtitles" json:"subtitles,omitempty"`
	SeasonID    int              `db.Conn:"season_id" json:"seasonID"`
}

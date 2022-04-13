package model

import (
	"github.com/fukaraca/worth2watch/config"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgtype"
)

var R *gin.Engine

var TIMEOUT = config.GetEnv.GetDuration("TIMEOUT")
var ServerHost = config.GetEnv.GetString("SERVER_HOST")
var ServerPort = config.GetEnv.GetString("SERVER_PORT")

type User struct {
	UserID    int                `db.Conn:"user_id" json:"userID,omitempty"`
	Username  string             `db.Conn:"username" json:"username"`
	Password  string             `db.Conn:"password" json:"password"`
	Email     pgtype.Text        `db.Conn:"email" json:"email"`
	Name      pgtype.Text        `db.Conn:"name" json:"name"`
	Lastname  pgtype.Text        `db.Conn:"lastname" json:"lastname"`
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
	Directors   []pgtype.Text    `db.Conn:"directors" json:"director,omitempty"`
	Writers     []pgtype.Text    `db.Conn:"writers" json:"writer,omitempty"`
	Stars       []pgtype.Text    `db.Conn:"stars" json:"stars,omitempty"`
	Duration    int              `db.Conn:"duration_min" json:"duration,omitempty"`
	IMDBid      pgtype.Text      `db.Conn:"imdb_id" json:"IMDBid"`
	Year        int              `db.Conn:"year" json:"year,omitempty"`
	Genres      []pgtype.Text    `db.Conn:"genres" json:"genre,omitempty"`
	Audio       []pgtype.Text    `db.Conn:"audios" json:"audio,omitempty"`
	Subtitles   []pgtype.Text    `db.Conn:"subtitles" json:"subtitles,omitempty"`
}

type Series struct {
	SerieID     int              `db.Conn:"serie_id" json:"serieID,omitempty"`
	Title       pgtype.Text      `db.Conn:"title" json:"title"`
	Description pgtype.Text      `db.Conn:"description" json:"description"`
	Rating      pgtype.Numeric   `db.Conn:"rating" json:"rating"`
	ReleaseDate pgtype.Timestamp `db.Conn:"release_date" json:"releaseDate"`
	Directors   []pgtype.Text    `db.Conn:"directors" json:"director,omitempty"`
	Writers     []pgtype.Text    `db.Conn:"writers" json:"writer,omitempty"`
	Stars       []pgtype.Text    `db.Conn:"stars" json:"stars,omitempty"`
	Duration    int              `db.Conn:"duration_min" json:"duration,omitempty"`
	IMDBid      pgtype.Text      `db.Conn:"imdb_id" json:"IMDBid"`
	Year        int              `db.Conn:"year" json:"year,omitempty"`
	Genres      []pgtype.Text    `db.Conn:"genres" json:"genre,omitempty"`
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
	Directors   []pgtype.Text    `db.Conn:"directors" json:"director,omitempty"`
	Writers     []pgtype.Text    `db.Conn:"writers" json:"writer,omitempty"`
	Stars       []pgtype.Text    `db.Conn:"stars" json:"stars,omitempty"`
	Duration    int              `db.Conn:"duration_min" json:"duration,omitempty"`
	IMDBid      pgtype.Text      `db.Conn:"imdb_id" json:"IMDBid"`
	Year        int              `db.Conn:"year" json:"year,omitempty"`
	Audios      []pgtype.Text    `db.Conn:"audios" json:"audio,omitempty"`
	Subtitles   []pgtype.Text    `db.Conn:"subtitles" json:"subtitles,omitempty"`
	SeasonID    int              `db.Conn:"season_id" json:"seasonID"`
}

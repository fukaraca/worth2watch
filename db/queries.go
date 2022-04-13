package db

import (
	"fmt"
	"github.com/fukaraca/worth2watch/api/admin"
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

//AddMovieContentWithID inserts movie to DB..
func AddMovieContentWithID(imdb string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	id, err := admin.FindIDWithIMDB(imdb)
	if err != nil {
		return
	}
	movie := admin.GetMovie(id)
	err = AddMovieContentWithStruct(ctx, movie)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("movie: ", movie.Title, " succesfully added")
	}
}

//AddSeriesContentWithID adds series to DB with its seasons.
func AddSeriesContentWithID(imdb string) {

	id, err := admin.FindIDWithIMDB(imdb)
	if err != nil {
		return
	}

	series := admin.GetSeries(id)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = AddSeriesContentWithStruct(ctx, series)
	if err != nil {
		log.Println(err)
	} else {
		log.Println(*series.Title, " saved to DB successfully")
	}

}

func AddMovieContentWithStruct(ctx context.Context, movie *model.Movie) error {
	ctx1, cancel1 := context.WithTimeout(ctx, TIMEOUT)
	defer cancel1()
	_, err = Conn.Exec(ctx1, "INSERT INTO movies (movie_id,title,description,rating,release_date,directors,writers,stars,duration_min,imdb_id,year,genres,audios,subtitles) VALUES (nextval('movies_movie_id_seq'),$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13);", movie.Title, movie.Description, movie.Rating, movie.ReleaseDate, movie.Directors, movie.Writers, movie.Stars, movie.Duration, movie.IMDBid, movie.Year, movie.Genres, movie.Audios, movie.Subtitles)
	if err != nil {
		return fmt.Errorf("insert for movie %s failed: %v", *movie.Title, err)
	}

	return nil
}

func AddSeriesContentWithStruct(ctx context.Context, series *model.Series) error {
	ctx1, cancel1 := context.WithCancel(ctx)
	defer cancel1()
	//insert series
	_, err = Conn.Exec(ctx1, "INSERT INTO series (serie_id,title,description,rating,release_date,directors,writers,stars,duration_min,imdb_id,year,genres,seasons) VALUES (nextval('series_serie_id_seq'),$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12);", series.Title, series.Description, series.Rating, series.ReleaseDate, series.Directors, series.Writers, series.Stars, series.Duration, series.IMDBid, series.Year, series.Genres, series.Seasons)
	if err != nil {
		return err
	}

	row := Conn.QueryRow(ctx1, "SELECT serie_id FROM series WHERE imdb_id=$1;", series.IMDBid)
	seriesId := 0
	err = row.Scan(&seriesId)
	if err != nil {
		return fmt.Errorf("series_id couldn't be scanned from db: %v", err)
	}

	//insert seasons
	for i := 1; i < series.Seasons+1; i++ {
		season, episodes := admin.GetSeason(series, i)

		ctx2, cancel2 := context.WithTimeout(ctx, TIMEOUT)
		defer cancel2()
		_, err = Conn.Exec(ctx2, "INSERT INTO seasons (season_id,season_number,episodes,imdb_id) VALUES (nextval('seasons_season_id_seq'),$1,$2,$3);", season.SeasonNumber, season.Episodes, season.IMDBid)
		if err != nil {
			return fmt.Errorf("insert season %d for %s failed: %v", season.SeasonNumber, *series.Title, err)
		}
		//foreign key assignment for seasons
		_, err = Conn.Exec(ctx2, "UPDATE seasons SET serie_id=$1 WHERE imdb_id=$2;", seriesId, season.IMDBid)
		if err != nil {
			return fmt.Errorf("update  season %d for %s failed when FK assignment: %v", season.SeasonNumber, *series.Title, err)
		}

		row := Conn.QueryRow(ctx2, "SELECT season_id FROM seasons WHERE serie_id=$1 AND season_number=$2;", seriesId, season.SeasonNumber)
		seasonId := 0
		err = row.Scan(&seasonId)
		if err != nil {
			return fmt.Errorf("season_id couldn't be scanned from db : %v", err)
		}

		//insert episodes
		for _, episode := range episodes {

			ctx3, cancel3 := context.WithTimeout(ctx, TIMEOUT)
			defer cancel3()
			_, err = Conn.Exec(ctx3, "INSERT INTO episodes (episode_id,title,description,rating,release_date,directors,writers,stars,duration_min,imdb_id,year,audios,subtitles,episode_number) VALUES (nextval('episodes_episode_id_seq'),$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13);", episode.Title, episode.Description, episode.Rating, episode.ReleaseDate, episode.Directors, episode.Writers, episode.Stars, episode.Duration, episode.IMDBid, episode.Year, episode.Audios, episode.Subtitles, episode.EpisodeNumber)
			if err != nil {
				return fmt.Errorf("insert episode %d for %s failed: %v", episode.EpisodeNumber, *series.Title, err)
			}

			//foreign key assignment for episodes
			_, err = Conn.Exec(ctx3, "UPDATE episodes SET season_id=$1 WHERE imdb_id=$2;", seasonId, episode.IMDBid)
			if err != nil {
				return fmt.Errorf("update  season %d for %s failed when fk assignment: %v", episode.EpisodeNumber, *series.Title, err)
			}
		}
	}
	return nil
}

//DeleteContent deletes given content from DB
func DeleteContent(c *gin.Context, username, id, contentType string) error {
	ctx, cancel := context.WithTimeout(c.Request.Context(), TIMEOUT)
	defer cancel()

	switch contentType {
	case "movie":
		_, err := Conn.Exec(ctx, "DELETE FROM movies WHERE imdb_id=$1;", id)
		//_, err := Conn.Exec(ctx, "DELETE FROM favorite_movies WHERE =$1;",id)
		if err != nil {
			return err
		}
	case "series":
		_, err := Conn.Exec(ctx, "DELETE FROM series WHERE imdb_id=$1;", id)
		if err != nil {
			return err
		}
	}
	return nil
}

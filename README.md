# a Better IMDB is possible

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white) ![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white) ![Postgres](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white) ![Redis](https://img.shields.io/badge/redis-%23DD0031.svg?style=for-the-badge&logo=redis&logoColor=white)
 
In this project, a functional application that manages back-end needs of a  movie/series database.

Features:
- You can make account management additionally with administration role
- You can manage contents with admin role by adding with IMDB-ID, raw-JSON and delete content by IMDB-ID
- Users can add/delete movie/series to their Favorites and search by genre and content name
- On public access, any guest can request movies list , a specific movie, series list, a specific series along with its seasons and episodes. Additionally, the guest can search by genre and content name.
- Dockerized PostgreSQL and Redis by docker-compose. (application will be Dockerized)

## Get Started

```
git clone https://github.com/fukaraca/worth2watch.git
```

TLDR; 

- Insert API key to env file
- If you will use provided docker-compose file, follow instructions. (if you have a solution with the issue, any kind of feedback or PR is accepted gratefully)

In order to initialize the application, PostgreSQL and Redis must be running with required environment values. (see config/config.env)
You can use your Psql and Redis, also if you want, you can use given Dockerized version. 
 
Note: If you will use provided docker-compose file, on "docker-compose up" step, at first Psql will not be accepting request due to user previlige-and persistent volume permissions conflict.
To handle it following post can help. [Stackoverflow](https://stackoverflow.com/a/71827179/12664011)

Now, Psql and Redis running, we must create a database as per in config.env . Application doesn't create database, 
so we must create it beforehand. Also API key that was provided by [TMDB](https://www.themoviedb.org) must be inserted to config.env file.

And we can start:

` go run .`


## Endpoints


```go
package api

func Endpoints() {
	//public
	model.R.GET("/movies/:id", GetThisMovie)
	model.R.GET("/movies/list", GetMoviesWithPage)
	model.R.GET("/searchContent", SearchContent)
	model.R.GET("/series/:seriesid", GetThisSeries)
	model.R.GET("/series/list", GetSeriesWithPage)
	model.R.GET("/series/:seriesid/:season", GetEpisodesForaSeason)
	model.R.GET("/getSimilarContent", GetSimilarContent)
	//user accessed
	model.R.POST("/addFavorites", Auth(AddToFavorites))
	model.R.GET("/getFavorites", Auth(GetFavorites))
	model.R.GET("/searchFavorites", Auth(SearchFavorites))
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
```


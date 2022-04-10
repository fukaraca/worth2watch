-- noinspection SqlCurrentSchemaInspectionForFile

CREATE TABLE IF NOT EXISTS users (
    user_id serial primary key not null ,
    username varchar(50) not null unique ,
    password char(64) not null ,
    email varchar(255) not null unique ,
    name varchar(50),
    lastname varchar(50),
    createdon timestamptz not null ,
    lastlogin timestamp not null ,
    isAdmin boolean not null
);
CREATE TABLE IF NOT EXISTS movies (
    movie_id serial primary key not null ,
    title varchar(100) not null ,
    description text,
    rating NUMERIC(3,1) check ( rating>=0 ),
    release_date timestamp,
    director varchar(70) [],
    writer varchar(70) [],
    stars varchar(70) [],
    duration_min int,
    IMDB_id varchar(20),
    year int,
    genre varchar(30),
    audio varchar(20) [],
    subtitles varchar(20) []
);
CREATE TABLE IF NOT EXISTS series (
    serie_id serial primary key not null ,
    title varchar(100) not null ,
    description text,
    rating NUMERIC(3,1) check ( rating>=0 ),
    release_date timestamp,
    director varchar(70) [],
    writer varchar(70) [],
    stars varchar(70) [],
    duration_min int,
    IMDB_id varchar(20),
    year int,
    genre varchar(30),
    seasons int
);
CREATE TABLE IF NOT EXISTS seasons (
    season_id serial primary key not null ,
    season_number int not null ,
    episodes int not null ,
    serie_id int references series(serie_id) not null
);
CREATE TABLE IF NOT EXISTS episodes (
    episode_id serial unique not null,
    title varchar(100),
    description text,
    rating NUMERIC(3,1) check ( rating>=0 ),
    release_date timestamp,
    director varchar(70) [],
    writer varchar(70) [],
    stars varchar(70) [],
    duration_min int,
    IMDB_id varchar(20),
    year int,
    audio varchar(20) [],
    subtitles varchar(20) [],
    season_id int references seasons(season_id) not null
);
CREATE TABLE IF NOT EXISTS favorite_movies (
    favorite_id serial unique not null,
    user_id int references users(user_id) not null,
    movie_id int references movies(movie_id) not null
);
CREATE TABLE IF NOT EXISTS favorite_series (
    favorite_id serial unique not null ,
    user_id int references users(user_id) not null ,
    serie_id int references series(serie_id) not null
);
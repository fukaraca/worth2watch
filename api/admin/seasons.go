package admin

import (
	"encoding/json"
	"fmt"
	"github.com/jackc/pgtype"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

//GetSeason function constructs and returns the given season and its episodes
func (series *Series) GetSeason(season int) (*Seasons, *[]Episodes) {
	id, err := FindIDWithIMDB(series.IMDBid.String)
	if err != nil {
		log.Println("serie id couldn't be gotten", err)
		return nil, nil
	}
	getURL := fmt.Sprintf("https://api.themoviedb.org/3/tv/%d/season/%d?api_key=%s", id, season, API_KEY)
	resp, err := http.Get(getURL)
	if err != nil {
		log.Println(err)
		return nil, nil
	}
	defer resp.Body.Close()

	read, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, nil
	}
	seasonsFromAPI := new(SeasonsAPI)
	err = json.Unmarshal(read, &seasonsFromAPI)
	if err != nil {
		log.Println("unmarshall failed:", err)
		return nil, nil
	}

	//construct seasons and episodes
	var retSeason Seasons
	var retEpisodes []Episodes
	//seasons
	retSeason.SeasonNumber = season
	retSeason.Episodes = len(seasonsFromAPI.Episodes)
	//episodes
	for i, ep := range seasonsFromAPI.Episodes {
		temp := Episodes{}
		//title
		temp.Title.String = ep.Name
		//Description
		temp.Description.String = ep.Overview
		//Rating
		if err = temp.Rating.Set(ep.VoteAverage); err != nil {
			log.Println("rating value couldn't be assigned for pgtype", err)
		}
		//Release Date
		parsed, err := time.Parse("2006-01-02", ep.AirDate)
		if err != nil {
			{
				parsed, err = time.Parse("2006/01/02", ep.AirDate)
				if err != nil {
					parsed, err = time.Parse("2006.01.02", ep.AirDate)
					if err != nil {
						log.Println("release date couldn't be parsed for episode:", ep.EpisodeNumber, "err:", err, "time format:", ep.AirDate)
					}
				}
			}
		}
		temp.ReleaseDate.Time = parsed
		//Directors
		directors := seasonsFromAPI.getDirectors(i, "Director")
		for director := range directors {
			directorToBeAppended := pgtype.Text{String: director}
			temp.Directors = append(temp.Directors, directorToBeAppended)
		}
		//Writers
		writers := seasonsFromAPI.getWriters(i, "Writer")
		for writer := range writers {
			writerToBeAppended := pgtype.Text{String: writer}
			temp.Writers = append(temp.Writers, writerToBeAppended)
		}
		//Stars
		stars := seasonsFromAPI.getStars(id, season, ep.EpisodeNumber, 6, 5)
		for star := range stars {
			starToBeAppended := pgtype.Text{String: star}
			temp.Stars = append(temp.Stars, starToBeAppended)
		}
		//Duration
		temp.Duration = series.Duration
		//IMDB id
		if imdb, err := findIMDBIDForEpisode(id, season, ep.EpisodeNumber); err != nil {
			log.Println("imdb id for season ", ep.SeasonNumber, " episode ", ep.EpisodeNumber, " couldn't be get")
		} else {
			temp.IMDBid.String = imdb
		}
		//Year
		temp.Year = temp.ReleaseDate.Time.Year()
		//Audios and Subtitles
		if tr, err := getTranslationDataOfEpisode(id, season, ep.EpisodeNumber); err != nil {
			log.Println(err)
		} else {
			for _, translation := range tr.Translations {
				temp.Audios = append(temp.Audios, pgtype.Text{String: translation.EnglishName})
				temp.Subtitles = append(temp.Subtitles, pgtype.Text{String: translation.EnglishName})
			}
		}
		retEpisodes = append(retEpisodes, temp)
		log.Println(series.Title.String, " season ", season, " episode ", ep.EpisodeNumber, " has been fetched")
	}
	return &retSeason, &retEpisodes
}

//getDirectors is a helper func for GetSeason
func (crew *SeasonsAPI) getDirectors(idx int, jobTitle string) map[string]struct{} {
	ret := make(map[string]struct{})
	var empty struct{}
	for _, c := range crew.Episodes[idx].Crew {
		if c.Job == jobTitle {
			ret[c.Name] = empty
		}
	}
	return ret
}

//getWriters is a helper func for GetSeason
func (crew *SeasonsAPI) getWriters(idx int, jobTitle string) map[string]struct{} {
	ret := make(map[string]struct{})
	var empty struct{}
	for _, c := range crew.Episodes[idx].Crew {
		if c.Job == jobTitle {
			ret[c.Name] = empty
		}
	}
	return ret
}

//getStars is helper func for GetSeason.For given popularity and amount, it looks up for it among the cast and retuns that match
func (s *SeasonsAPI) getStars(serieID, season, episode, count int, popularity float64) map[string]struct{} {
	getURL := fmt.Sprintf("https://api.themoviedb.org/3/tv/%d/season/%d/episode/%d/credits?api_key=%s", serieID, season, episode, API_KEY)
	resp, err := http.Get(getURL)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer resp.Body.Close()

	read, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil
	}
	castFromAPI := new(EpisodeCastAPI)
	err = json.Unmarshal(read, &castFromAPI)
	if err != nil {
		log.Println("unmarshall failed:", err)
		return nil
	}
	ret := make(map[string]struct{})
	type temp struct {
		name   string
		rating float64
	}
	populers := []temp{}
	var empty struct{}
	for _, c := range castFromAPI.Cast {
		if c.Popularity > popularity && c.KnownForDepartment == "Acting" {
			populers = append(populers, temp{
				name:   c.Name,
				rating: c.Popularity,
			})
		}
	}
	sort.Slice(populers, func(i, j int) bool {
		return populers[i].rating > populers[j].rating
	})
	for i := 0; i < count && i < len(populers); i++ {
		ret[populers[i].name] = empty
	}
	return ret
}

//findIMDBWithID finds and returns episodes IMDB ID for given TMDB ID
func findIMDBIDForEpisode(serieID, season, episode int) (string, error) {
	getURL := fmt.Sprintf("https://api.themoviedb.org/3/tv/%d/season/%d/episode/%d/external_ids?api_key=%s", serieID, season, episode, API_KEY)
	resp, err := http.Get(getURL)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer resp.Body.Close()

	read, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return "", err
	}
	_, after, ok := strings.Cut(string(read), "\"imdb_id\":\"")
	if !ok {
		return "", fmt.Errorf("imdb_id not found")
	}
	before, _, ok := strings.Cut(after, "\",\"")
	if !ok {
		return "", fmt.Errorf("imdb_id not found")
	}

	return before, nil
}

func getTranslationDataOfEpisode(serieID, season, episode int) (*TranslationDataOfEpisode, error) {
	getURL := fmt.Sprintf("https://api.themoviedb.org/3/tv/%d/season/%d/episode/%d/translations?api_key=%s", serieID, season, episode, API_KEY)
	resp, err := http.Get(getURL)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	read, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	tr := TranslationDataOfEpisode{}
	err = json.Unmarshal(read, &tr)
	if err != nil {
		log.Println("translation data couldn't be unmarshalled for ", serieID, " season ", season, " episode ", episode)
		return nil, err
	}
	return &tr, nil
}

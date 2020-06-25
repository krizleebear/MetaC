package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"time"

	tmdb "github.com/cyruzin/golang-tmdb"
)

var (
	apiKey = "" //defined via build flag
)

func initClient() *tmdb.Client {
	tmdbClient, err := tmdb.Init(apiKey)
	if err != nil {
		panic(err)
	}

	return tmdbClient
}

// Search TheMovieDB online for a movie of the given string
// Returns all matched movies.
func Search(movieName string) (*tmdb.SearchMovies, error) {
	tmdbClient := initClient()

	options := map[string]string{
		"language": "de-DE",
		//"year": "2008",
	}

	return tmdbClient.GetSearchMovies(movieName, options)
}

func getPosterSizes() []string {
	tmdbClient := initClient()
	api, _ := tmdbClient.GetConfigurationAPI()
	return api.Images.PosterSizes
}

type singleMovie struct {
	VoteCount        int64   `json:"vote_count"`
	ID               int64   `json:"id"`
	Video            bool    `json:"video"`
	VoteAverage      float32 `json:"vote_average"`
	Title            string  `json:"title"`
	Popularity       float32 `json:"popularity"`
	PosterPath       string  `json:"poster_path"`
	OriginalLanguage string  `json:"original_language"`
	OriginalTitle    string  `json:"original_title"`
	GenreIDs         []int64 `json:"genre_ids"`
	BackdropPath     string  `json:"backdrop_path"`
	Adult            bool    `json:"adult"`
	Overview         string  `json:"overview"`
	ReleaseDate      string  `json:"release_date"`
}

func downloadPoster(imageURI string) string {
	fileName := path.Base(imageURI)

	tmpfile, err := ioutil.TempFile("", "MetaY.*"+fileName)
	if err != nil {
		log.Fatal(err)
	}

	//https://image.tmdb.org/t/p/original/687NOelgrgtsKEFsotLCH0YZn6H.jpg
	resp, err := http.Get("https://image.tmdb.org/t/p/original" + imageURI)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return ""
	}

	b, _ := io.Copy(tmpfile, resp.Body)
	fmt.Println("Poster file size: ", b)
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}

	return tmpfile.Name()
}

// GetYearFromReleaseDate or empty string
func GetYearFromReleaseDate(releaseDate string) string {
	parsedDate, err := time.Parse("2006-01-02", releaseDate)
	if err != nil {
		return ""
	}
	return parsedDate.Format("2006")
}

func getMovieCredits(movieID int64) (*tmdb.MovieCredits, error) {
	tmdbClient := initClient()

	return tmdbClient.GetMovieCredits(int(movieID), nil)
}

// getCast extracts the first names of the full cast
func getCast(credits *tmdb.MovieCredits) []string {
	maxCount := 5
	members := make([]string, 0, maxCount)
	for i, member := range credits.Cast {
		if i > maxCount-1 {
			break
		}
		members = append(members, member.Name)
	}
	return members
}

func getDirectors(credits *tmdb.MovieCredits) []string {
	maxCount := 5
	members := make([]string, 0, maxCount)
	expectedDepartment := "Directing"
	for _, member := range credits.Crew {
		if len(members) > maxCount-1 {
			break
		}
		if member.Department == expectedDepartment {
			members = append(members, member.Name)
		}
	}

	return members
}

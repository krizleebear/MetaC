package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"regexp"
	"strings"
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

	name, year := splitNameAndYear(movieName)

	options := map[string]string{
		"language": "de-DE",
		"year":     year,
	}

	return tmdbClient.GetSearchMovies(name, options)
}

// SearchMulti performs a combined search for both movies and tv shows
func SearchMulti(title string) (*tmdb.SearchMulti, error) {
	tmdbClient := initClient()

	name, year := splitNameAndYear(title)

	options := map[string]string{
		"language": "de-DE",
		"year":     year,
	}

	// normalize fields between movies and tv show results
	response, error := tmdbClient.GetSearchMulti(name, options)
	for i, result := range response.Results {
		result.ReleaseDate = FirstNonEmpty(result.ReleaseDate, result.FirstAirDate)
		result.Title = FirstNonEmpty(result.Title, result.Name, result.OriginalName)
		response.Results[i] = result
	}

	return response, error
}

// FirstNonEmpty string of all given strings is being returned.
// If all strings are empty or if none is given, an empty string will be returned.
func FirstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func getPosterSizes() []string {
	tmdbClient := initClient()
	api, _ := tmdbClient.GetConfigurationAPI()
	return api.Images.PosterSizes
}

type singleMovie struct {
	PosterPath       string   `json:"poster_path,omitempty"`
	Popularity       float32  `json:"popularity"`
	ID               int64    `json:"id"`
	Overview         string   `json:"overview,omitempty"`
	BackdropPath     string   `json:"backdrop_path,omitempty"`
	VoteAverage      float32  `json:"vote_average,omitempty"`
	MediaType        string   `json:"media_type"`
	FirstAirDate     string   `json:"first_air_date,omitempty"`
	OriginCountry    []string `json:"origin_country,omitempty"`
	GenreIDs         []int64  `json:"genre_ids,omitempty"`
	OriginalLanguage string   `json:"original_language,omitempty"`
	VoteCount        int64    `json:"vote_count,omitempty"`
	Name             string   `json:"name,omitempty"`
	OriginalName     string   `json:"original_name,omitempty"`
	Adult            bool     `json:"adult,omitempty"`
	ReleaseDate      string   `json:"release_date,omitempty"`
	OriginalTitle    string   `json:"original_title,omitempty"`
	Title            string   `json:"title,omitempty"`
	Video            bool     `json:"video,omitempty"`
	ProfilePath      string   `json:"profile_path,omitempty"`
	KnownFor         []struct {
		PosterPath       string  `json:"poster_path"`
		Adult            bool    `json:"adult"`
		Overview         string  `json:"overview"`
		ReleaseDate      string  `json:"release_date"`
		OriginalTitle    string  `json:"original_title"`
		GenreIDs         []int64 `json:"genre_ids"`
		ID               int64   `json:"id"`
		MediaType        string  `json:"media_type"`
		OriginalLanguage string  `json:"original_language"`
		Title            string  `json:"title"`
		BackdropPath     string  `json:"backdrop_path"`
		Popularity       float32 `json:"popularity"`
		VoteCount        int64   `json:"vote_count"`
		Video            bool    `json:"video"`
		VoteAverage      float32 `json:"vote_average"`
	} `json:"known_for,omitempty"`
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

func splitNameAndYear(movieName string) (string, string) {
	year := ""

	re := regexp.MustCompile(`\([0-9]{4}\)`)
	loc := re.FindStringIndex(movieName)
	if loc != nil {
		year = movieName[loc[0]+1 : loc[1]-1]
		movieName = strings.TrimSpace(movieName[0:loc[0]])
	}

	return movieName, year
}

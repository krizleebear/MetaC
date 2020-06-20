package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"
	"time"

	tmdb "github.com/cyruzin/golang-tmdb"
)

var (
	apiKey = "" //defined via build flag
)

// Search TheMovieDB online for a movie of the given string
// Returns all matched movies.
func Search(movieName string) (*tmdb.SearchMovies, error) {
	tmdbClient, err := tmdb.Init(apiKey)
	if err != nil {
		return nil, err
	}

	options := map[string]string{
		"language": "de-DE",
		//"year": "2008",
	}

	return tmdbClient.GetSearchMovies(movieName, options)
}

func downloadPoster(imageURI string) string {
	fileName := path.Base(imageURI)

	tmpfile, err := ioutil.TempFile("", "MetaY.*"+fileName)
	if err != nil {
		log.Fatal(err)
	}

	//https://image.tmdb.org/t/p/w500/687NOelgrgtsKEFsotLCH0YZn6H.jpg
	resp, err := http.Get("https://image.tmdb.org/t/p/w500" + imageURI)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	b, _ := io.Copy(tmpfile, resp.Body)
	fmt.Println("Poster file size: ", b)
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}

	return tmpfile.Name()
}

// ToAtomicParsleyArguments returns the command line arguments for AtomicParsley tool
// It's your job to delete the poster file after you've used it.
func ToAtomicParsleyArguments(movieFile string, searchResult *tmdb.SearchMovies, resultIndex int) ([]string, string) {
	movie := searchResult.Results[resultIndex]
	year := GetYearFromReleaseDate(movie.ReleaseDate)
	description := strings.ReplaceAll(movie.Overview, "\"", "\\\"")
	posterFile := downloadPoster(movie.PosterPath)

	arguments := make([]string, 0, 10)
	arguments = append(arguments, "AtomicParsley")
	arguments = append(arguments, movieFile)

	arguments = append(arguments, "--overWrite")

	arguments = append(arguments, "--stik")
	arguments = append(arguments, "Movie")

	arguments = append(arguments, "--title")
	arguments = append(arguments, movie.Title)

	arguments = append(arguments, "--year")
	arguments = append(arguments, year)

	arguments = append(arguments, "--longdesc")
	arguments = append(arguments, description)

	arguments = append(arguments, "--artwork")
	arguments = append(arguments, "REMOVE_ALL")

	arguments = append(arguments, "--artwork")
	arguments = append(arguments, posterFile)

	return arguments, posterFile
}

// GetYearFromReleaseDate or empty string
func GetYearFromReleaseDate(releaseDate string) string {
	parsedDate, err := time.Parse("2006-01-02", releaseDate)
	if err != nil {
		return ""
	}
	return parsedDate.Format("2006")
}

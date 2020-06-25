package main

import (
	"regexp"

	tmdb "github.com/cyruzin/golang-tmdb"
)

// ToAtomicParsleyArguments returns the command line arguments for AtomicParsley tool
// It's your job to delete the poster file after you've used it.
//searchResult *tmdb.SearchMovies, resultIndex int
//movie := searchResult.Results[resultIndex]
func ToAtomicParsleyArguments(movieFile string, movie singleMovie, credits *tmdb.MovieCredits) ([]string, string) {

	year := GetYearFromReleaseDate(movie.ReleaseDate)
	description := movie.Overview
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
	description = removeDuplicateWhitespace(description)
	arguments = append(arguments, description)

	if posterFile != "" {
		arguments = append(arguments, "--artwork")
		arguments = append(arguments, "REMOVE_ALL")

		arguments = append(arguments, "--artwork")
		arguments = append(arguments, posterFile)
	}

	return arguments, posterFile
}

func removeDuplicateWhitespace(s string) string {
	space := regexp.MustCompile(`\s+`)
	return space.ReplaceAllString(s, " ")
}

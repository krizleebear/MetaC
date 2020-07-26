package main

import (
	"reflect"
	"strings"
	"testing"
)

func Test_removeDuplicateWhitespace(t *testing.T) {
	assertEquals(t, "a b", removeDuplicateWhitespace("a   b"))
}

func TestToAtomicParsleyArguments(t *testing.T) {
	movie := singleMovie{
		ID:          83542,
		Title:       "SomeTitle",
		PosterPath:  "/687NOelgrgtsKEFsotLCH0YZn6H.jpg",
		GenreIDs:    []int64{18, 878},
		Overview:    "Detailled description",
		ReleaseDate: "2012-10-26",
	}

	gotArgs, gotPosterFile := ToAtomicParsleyArguments("myMovie.mp4", movie, nil)
	wantedArgs := []string{"AtomicParsley",
		"myMovie.mp4", "--overWrite", "--stik", "Movie", "--title", "SomeTitle", "--year", "2012",
		"--longdesc", "Detailled description", "--artwork", "REMOVE_ALL", "--artwork",
		"/var/folders/8h/5k3zp7yd1rg8mtpyg4mnl9hr0000gn/T/MetaY.371467386somePoster.jpg"}

	if !strings.Contains(gotPosterFile, "MetaY") {
		t.Errorf("got %q, expected other filename", gotPosterFile)
	}

	// adapt postfile because it's a temporary file with random name
	wantedArgs[len(wantedArgs)-1] = gotPosterFile

	if !reflect.DeepEqual(gotArgs, wantedArgs) {
		t.Errorf("got %q, want %q", gotArgs, wantedArgs)
	}
}

func TestTVShowArguments(t *testing.T) {
	movie := singleMovie{
		ID:           46849,
		Title:        "Unsere Mütter, unsere Väter",
		Name:         "Unsere Mütter, unsere Väter",
		PosterPath:   "/loXl5jZ5X0Q5sxZwqjkjQQS0Jll.jpg",
		GenreIDs:     []int64{18, 10768},
		Overview:     "Detailled description",
		ReleaseDate:  "2012-10-26",
		MediaType:    "tv",
		FirstAirDate: "2012-10-26",
	}

	gotArgs, gotPosterFile := ToAtomicParsleyArguments("myMovie.mp4", movie, nil)
	wantedArgs := []string{"AtomicParsley",
		"myMovie.mp4", "--overWrite", "--stik", "TV Show", "--title", "Unsere Mütter, unsere Väter", "--year", "2012",
		"--longdesc", "Detailled description", "--artwork", "REMOVE_ALL", "--artwork",
		"/var/folders/8h/5k3zp7yd1rg8mtpyg4mnl9hr0000gn/T/MetaY.371467386somePoster.jpg"}

	if !strings.Contains(gotPosterFile, "MetaY") {
		t.Errorf("got %q, expected other filename", gotPosterFile)
	}

	// adapt postfile because it's a temporary file with random name
	wantedArgs[len(wantedArgs)-1] = gotPosterFile

	if !reflect.DeepEqual(gotArgs, wantedArgs) {
		t.Errorf("got %q, want %q", gotArgs, wantedArgs)
	}
}

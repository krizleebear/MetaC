package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

func init() {
	apiKey = os.Getenv("apiKey")
}

func TestSearch(t *testing.T) {
	result, err := Search("Cloud Atlas")
	if err != nil {
		t.Error(err)
	}

	if len(result.Results) < 1 {
		t.Error("found no movie for search term")
	}

	for _, result := range result.Results {
		fmt.Printf("%+v\n", result)
	}
}

func TestDownload(t *testing.T) {
	filename := downloadPoster("/687NOelgrgtsKEFsotLCH0YZn6H.jpg")
	defer os.Remove(filename)

	fileinfo, err := os.Stat(filename)
	if err != nil {
		t.Error(err)
	}

	if fileinfo.Size() == 0 {
		t.Errorf("Poster file must not be empty")
	}
}

func Test_getPosterSizes(t *testing.T) {

	fmt.Printf("%+v\n", getPosterSizes())
}

func Test_getMovieCredits(t *testing.T) {
	credits, err := getMovieCredits(83542)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(credits)
}

func TestToAtomicParsleyArguments(t *testing.T) {
	movie := singleMovie{
		ID:          83542,
		Title:       "SomeTitle",
		PosterPath:  "/somePoster.jpg",
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

func Test_getCast(t *testing.T) {
	credits, err := getMovieCredits(83542)
	if err != nil {
		t.Error(err)
	}
	cast := getCast(credits)
	fmt.Println(cast)
}

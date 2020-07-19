package main

import (
	"fmt"
	"os"
	"testing"
)

func init() {
	apiKey = os.Getenv("apiKey")
}

func assertEquals(t *testing.T, wanted interface{}, got interface{}) {
	if wanted != got {
		t.Fatalf("wanted %v but got %v", wanted, got)
	}
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

func Test_getCast(t *testing.T) {
	credits, _ := getMovieCredits(83542)
	cast := getCast(credits)

	assertEquals(t, 5, len(cast))
	assertEquals(t, "Tom Hanks", cast[0])
}

func Test_getDirectors(t *testing.T) {
	credits, _ := getMovieCredits(83542)
	directors := getDirectors(credits)

	assertEquals(t, 5, len(directors))
	assertEquals(t, "Tom Tykwer", directors[0])
}

func Test_splitNameAndYear(t *testing.T) {
	n, y := splitNameAndYear("")
	assertEquals(t, "", n)
	assertEquals(t, "", y)

	n, y = splitNameAndYear("Blabla")
	assertEquals(t, "Blabla", n)
	assertEquals(t, "", y)

	n, y = splitNameAndYear("Blabla (2002)")
	assertEquals(t, "Blabla", n)
	assertEquals(t, "2002", y)
}

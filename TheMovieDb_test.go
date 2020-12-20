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

func assertNotEquals(t *testing.T, notWanted interface{}, got interface{}) {
	if notWanted == got {
		t.Fatalf("didn't expect it, but got %v", got)
	}
}

func assertNil(t *testing.T, got interface{}) {
	if got != nil {
		t.Fatalf("expected nil, but got %v", got)
	}
}

func assertNotNil(t *testing.T, got interface{}) {
	if got == nil {
		t.Fatalf("expected a value, but got nil")
	}
}

func assertSliceContains(t *testing.T, wanted interface{}, got []string) {
	assertNotNil(t, got)
	for _, element := range got {
		if wanted == element {
			return
		}
	}
	t.Fatalf("wanted %v but got %v", wanted, got)
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

func TestSearchMulti(t *testing.T) {

	expectedTitle := "Unsere Mütter, unsere Väter"
	response, err := SearchMulti(expectedTitle)
	if err != nil {
		t.Error(err)
	}

	result := response.Results[0]
	fmt.Printf("%+v\n", result)

	// Respective attributes of TV Show and Movie must be equal
	assertEquals(t, result.ReleaseDate, result.FirstAirDate)
	assertEquals(t, result.Title, result.Name)
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

	credits := getCastMembers(83542, "movie")

	fmt.Println(credits)
}

func Test_getCast(t *testing.T) {
	cast := getCastMembers(83542, "movie")

	assertEquals(t, 5, len(cast))
	assertEquals(t, "Tom Hanks", cast[0])
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

func Test_extractEpisodeID(t *testing.T) {
	ep, err := extractEpisodeID("")
	assertEquals(t, episodeID{}, ep)
	assertNotNil(t, err)

	ep, err = extractEpisodeID("S01E02")
	assertEquals(t, episodeID{1, 2}, ep)

	ep, err = extractEpisodeID("bla - S03E04 - abc")
	assertEquals(t, episodeID{3, 4}, ep)
}

func Test_GetMovieDetails(t *testing.T) {
	detail, err := GetMovieDetails(715025)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%v", detail)

	assertEquals(t, "Heaven", detail.Title)
	assertNotNil(t, detail.Overview)
	assertNotEquals(t, "", detail.Overview)
	//detail.Overview
}

func Test_GetShowDetails(t *testing.T) {

	cast := getCastMembers(94495, "tv")
	assertNotNil(t, cast)
	assertSliceContains(t, "Kristofer Hivju", cast)
}

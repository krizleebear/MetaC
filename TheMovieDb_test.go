package main

import (
	"fmt"
	"os"
	"testing"
)

func TestSearch(t *testing.T) {
	result, err := Search("Cloud Atlas")
	if err != nil {
		panic(err)
	}
	//fmt.Println(result)
	for _, result := range result.Results {
		//fmt.Println(result)
		//fmt.Println(result.Overview)
		//fmt.Println(result.PosterPath)
		fmt.Printf("%+v\n", result)
		//result.Title
		//result.ReleaseDate

		//https://image.tmdb.org/t/p/w440_and_h660_face/687NOelgrgtsKEFsotLCH0YZn6H.jpg
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

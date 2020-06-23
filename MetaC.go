package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

func main() {
	movieFile, _ := filepath.Abs(os.Args[1])
	fmt.Println(movieFile)
	_, err := os.Stat(movieFile)
	if err != nil {
		panic(err)
	}

	fileName := filepath.Base(movieFile)
	fileExt := filepath.Ext(movieFile)
	title := strings.TrimSuffix(fileName, fileExt)

	movieResults, err := Search(path.Base(title))
	if err != nil {
		panic(err)
	}

	if movieResults.TotalResults == 0 {
		panic("No movie result found for title " + title)
	}

	for i, result := range movieResults.Results {
		year := GetYearFromReleaseDate(result.ReleaseDate)
		fmt.Printf("%v) %v (%v)\n", i+1, result.Title, year)
	}

	movieIndex := 0
	fmt.Printf("%v results found. Please select the correct movie.\n", movieResults.TotalResults)
	if movieResults.TotalResults > 1 {
		movieIndex = readNumberFromTerminal() - 1
	}

	selectedMovie := movieResults.Results[movieIndex]

	fmt.Printf("%+v\n", selectedMovie)

	movieCredits, err := getMovieCredits(selectedMovie.ID)
	//fmt.Printf("%+v\n", movieCredits)

	args, posterFile := ToAtomicParsleyArguments(movieFile, selectedMovie, movieCredits)
	defer os.Remove(posterFile)

	fmt.Printf("AtomicParsley %v\n", strings.Join(args, " "))
	localExec("AtomicParsley", args)
}

func localExec(localBinary string, args []string) {
	binary, lookErr := exec.LookPath(localBinary)
	if lookErr != nil {
		panic(lookErr)
	}

	fmt.Println(binary)
	fmt.Println(args)

	env := os.Environ()
	execErr := syscall.Exec(binary, args, env)
	if execErr != nil {
		panic(execErr)
	}
}

func readNumberFromTerminal() int {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	token := scanner.Text()

	number, err := strconv.Atoi(token)
	if err != nil {
		panic(err)
	}
	return number
}

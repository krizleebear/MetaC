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

	tmdb "github.com/cyruzin/golang-tmdb"
)

func main() {

	movieFile := getFileFromArgs()
	title := getTitleFromFile(movieFile)

	movieResults := search(title)

	if movieResults.TotalResults == 0 {
		title, movieResults = askForTitle(title, movieResults)
	}

	for i, result := range movieResults.Results {
		year := GetYearFromReleaseDate(result.ReleaseDate)
		fmt.Printf("%v) %v (%v)\n", i+1, result.Title, year)
	}

	movieIndex := 0
	if movieResults.TotalResults > 1 {
		fmt.Printf("%v results found. Please select the correct movie.\n", movieResults.TotalResults)
		movieIndex = readNumberFromTerminal() - 1
	}

	selectedMovie := movieResults.Results[movieIndex]

	fmt.Printf("%+v\n", selectedMovie)

	movieCredits, _ := getMovieCredits(selectedMovie.ID)
	cast := getCast(movieCredits)
	fmt.Printf("Cast: %v\n", cast)

	args, posterFile := ToAtomicParsleyArguments(movieFile, selectedMovie, movieCredits)
	defer os.Remove(posterFile)

	fmt.Printf("AtomicParsley %v\n", strings.Join(args, " "))
	fmt.Println("Press enter to continue (or ctrl+c to abort)...")
	readTokenFromTerminal()

	localExec("AtomicParsley", args)
}

func askForTitle(title string, movieResults *tmdb.SearchMulti) (string, *tmdb.SearchMulti) {
	for movieResults.TotalResults == 0 {
		fmt.Printf("No movie result found for '%v'. Please enter an alternative title\n> ", title)
		title = readTokenFromTerminal()
		movieResults = search(title)
	}
	return title, movieResults
}

func getFileFromArgs() string {
	movieFile, _ := filepath.Abs(os.Args[1])
	fmt.Println(movieFile)
	_, err := os.Stat(movieFile)
	if err != nil {
		panic(err)
	}

	return movieFile
}

func getTitleFromFile(movieFile string) string {
	fileName := filepath.Base(movieFile)
	fileExt := filepath.Ext(movieFile)
	title := strings.TrimSuffix(fileName, fileExt)
	return path.Base(title)
}

func search(title string) *tmdb.SearchMulti {
	movieResults, err := SearchMulti(title)
	if err != nil {
		panic(err)
	}
	return movieResults
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

func readTokenFromTerminal() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func readNumberFromTerminal() int {
	token := readTokenFromTerminal()

	number, err := strconv.Atoi(token)
	if err != nil {
		panic(err)
	}
	return number
}

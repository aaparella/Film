package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/fatih/color"
)

var REQ = `http://www.omdbapi.com/?t={{.Title}}&y=&plot=short&r=json`

type Movie struct {
	Title      string
	Year       string
	ImdbRating string
	Metascore  string
	Runtime    string
	Genre      string
	Director   string
	Writer     string
	Actors     string
	Plot       string
}

func colorsForRating(rating float64) (*color.Color, *color.Color) {
	var c *color.Color
	if rating < 50 {
		c = color.New(color.FgRed)
	} else if rating < 70 {
		c = color.New(color.FgYellow)
	} else {
		c = color.New(color.FgGreen)
	}
	return c, c.Add(color.Bold)
}

func printMetascore(score string) {
	rating, _ := strconv.Atoi(score)
	color, boldColor := colorsForRating(float64(rating))
	boldColor.Printf("Metascore   : %.0d%% ", rating)
	printRatingBar(float64(rating), color)
	fmt.Println("")
}

func printIMDBRating(r string) {
	rating, _ := strconv.ParseFloat(r, 64)
	rating = rating * 10
	color, boldColor := colorsForRating(rating)
	boldColor.Printf("IMDB Rating : %.0f%% ", rating)
	printRatingBar(rating, color)
	fmt.Println("")
}

func printRatingBar(rating float64, color *color.Color) {
	color.Printf("[")
	color.Printf("%s", strings.Repeat("=", int(math.Floor(rating))))
	color.Printf("%s", strings.Repeat(" ", int(100-math.Floor(rating))))
	color.Printf("]")
}

func printValue(title, value string) {
	boldCyan := color.New(color.FgCyan).Add(color.Bold)
	boldCyan.Printf("%11s : ", title)
	color.Cyan(value)
}

func getMovie(title string) (*Movie, error) {
	data := struct {
		Title string
	}{
		title,
	}

	var URL bytes.Buffer
	tpl, _ := template.New("req").Parse(REQ)
	_ = tpl.Execute(&URL, data)

	resp, err := http.Get(URL.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var movie Movie
	if err = json.NewDecoder(resp.Body).Decode(&movie); err != nil {
		return nil, err
	}

	return &movie, nil
}

func printMovieInformation(movie *Movie) {
	printValue("Title", movie.Title)
	printValue("Director", movie.Director)
	printValue("Year", movie.Year)
	printValue("Genre", movie.Genre)
	printValue("Actors", movie.Actors)
	printValue("Writer(s)", movie.Writer)

	printIMDBRating(movie.ImdbRating)
	printMetascore(movie.Metascore)

	fmt.Println("\n", movie.Plot)
}

func main() {
	title := strings.Join(os.Args[1:], "+")
	movie, err := getMovie(title)
	if err != nil {
		log.Fatal(err)
	}

	if movie.Title == "" {
		color.Red("Could not find movie titled : %s", title)
		return
	}

	printMovieInformation(movie)
}

package main

import (
	"net/http"

	"errors"
	"math/rand"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Movie struct {
	ID    string `json:"id"`
	Isbn  string `json:"isbn"`
	Title string `json:"title"`
}

var movies []Movie

//==============================================================================================

func MoviebyID(id string) (*Movie, error) {
	for k, v := range movies {
		if v.ID == id {
			return &movies[k], nil
		}
	}
	return nil, errors.New("Invalid ID")
}

func removeFromList(id string) ([]Movie, error) {
	for k, v := range movies {
		if v.ID == id {
			movies = append(movies[:k], movies[k+1:]...)
			return movies, nil
		}
	}
	return movies, errors.New("Wrong ID")
}

func getMovies(c *gin.Context) {
	sort.Slice(movies, func(i, j int) bool {
		return movies[i].ID < movies[j].ID
	})
	c.IndentedJSON(http.StatusOK, movies)
}

func getMovie(c *gin.Context) {
	id := c.Param("id")
	movie, err := MoviebyID(id)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}
	c.IndentedJSON(http.StatusOK, movie)
}

func createMovie(c *gin.Context) {
	var newmovie Movie

	if err := c.BindJSON(&newmovie); err != nil {
		return
	}

	newmovie.ID = strconv.Itoa(rand.Intn(10000))

	movies = append(movies, newmovie)
	c.IndentedJSON(http.StatusOK, newmovie)
}

func updateMovie(c *gin.Context) {
	var newmovie Movie
	var err error

	id, ok := c.GetQuery("id")
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Please Give ID"})
		return
	}

	movies, err = removeFromList(id)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}

	if err1 := c.BindJSON(&newmovie); err1 != nil {
		return
	}

	newmovie.ID = id
	movies = append(movies, newmovie)
	c.IndentedJSON(http.StatusOK, newmovie)
}

func deleteMovie(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Please Give ID"})
		return
	}

	movies, err := removeFromList(id)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}

	c.IndentedJSON(http.StatusOK, movies)
}

//==============================================================================================

func main() {
	r := gin.Default()

	movies = append(movies, Movie{
		ID:    "1",
		Isbn:  "418227",
		Title: "Movie One",
	})
	movies = append(movies, Movie{
		ID:    "2",
		Isbn:  "45455",
		Title: "Movie two",
	})

	r.GET("/movies", getMovies)
	r.GET("/movies/:id", getMovie)
	r.POST("/movies", createMovie)
	r.PUT("/movies", updateMovie)
	r.DELETE("/movies", deleteMovie)

	r.Run("localhost:8000")
}

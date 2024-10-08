package movies

import (
	"movie-reservation-system/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

const DEFAULT_ID = "0"

type Movie struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Year        int    `json:"year"`
	Description string `json:"description"`
	ImageUrl    string `json:"image_url"`
	Genres      string `json:"genres"`
	Cast        string `json:"cast"`
}

func GetMovies(c *gin.Context) {
	lastIdParam := c.Query("last_id")
	lastId := DEFAULT_ID

	if lastIdParam != "" {
		lastId = lastIdParam
	}

	query := `
		SELECT
			m.id,
			m.title,
			m.year,
			m.description,
			m.image_url,
			STRING_AGG(DISTINCT g.name, ', ') AS genres,
			STRING_AGG(DISTINCT c.name, ', ') AS cast
			FROM 
			(SELECT * FROM Movies WHERE id > $1 LIMIT 10) AS m
			LEFT JOIN 
			movies_genres mg ON m.id = mg.movie_id
			LEFT JOIN 
			genres g ON mg.genre_id = g.id
			LEFT JOIN 
			movies_casting ma ON m.id = ma.movie_id
			LEFT JOIN 
			casting c ON ma.casting_id = c.id
			GROUP BY 
			m.id, m.title, m.year, m.description, m.image_url
		`

	rows, err := database.Db.Query(query, lastId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var movies []Movie
	var movie Movie
	for rows.Next() {
		rows.Scan(&movie.ID, &movie.Title, &movie.Year, &movie.Description, &movie.ImageUrl, &movie.Genres, &movie.Cast)
		movies = append(movies, movie)
	}

	c.JSON(http.StatusOK, gin.H{"movies": movies})
}

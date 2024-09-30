package movies

import (
	"encoding/json"
	"fmt"
	"io"
	"movie-reservation-system/database"
	"movie-reservation-system/users"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/lib/pq"
)

type Movie struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Year        int    `json:"year"`
	Description string `json:"description"`
	ImageUrl    string `json:"image_url"`
	Genres      string `json:"genres"`
	Cast        string `json:"cast"`
}

const DEFAULT_ID = "0"

func GetMovies(c *gin.Context) {
	lastIdParam := c.Query("last_id")
	lastId := DEFAULT_ID

	if lastIdParam != "" {
		lastId = lastIdParam
	}

	query := fmt.Sprintf(`
		SELECT
			m.id,
			m.title,
			m.year,
			m.description,
			m.image_url,
			STRING_AGG(DISTINCT g.name, ', ') AS genres,
			STRING_AGG(DISTINCT c.name, ', ') AS cast
			FROM 
			(SELECT * FROM Movies WHERE id > %s LIMIT 10) AS m
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
		`,
		lastId,
	)

	rows, err := database.Db.Query(query)
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

type UserClaims struct {
	ID int `json:"_id"`
	jwt.MapClaims
}

func generalError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "an error ocurred"})
}

type ReserveBody struct {
	Seats []string `json:"seats"`
	Date  string   `json:"date"`
}

func ReserveMovie(c *gin.Context) {
	var reserveBody ReserveBody
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	json.Unmarshal(jsonData, &reserveBody)

	if len(reserveBody.Seats) == 0 || reserveBody.Date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date and seat are required"})
		return
	}

	if len(reserveBody.Seats) > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "max 5 seats per reservation"})
		return
	}

	date := reserveBody.Date

	userId := users.ExtractUserIdFromClaims(c)
	movieId := c.Param("id")

	tx, err := database.Db.Begin()
	if err != nil {
		fmt.Println(err)
		generalError(c)
		return
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	exists := false
	err = tx.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM Reservation
			WHERE movie_id = $1
				AND seat = ANY($2)
			FOR UPDATE
		)
	`, movieId, pq.Array(reserveBody.Seats)).Scan(&exists)
	if err != nil {
		fmt.Println(err)
		generalError(c)
		return
	}

	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "seat already reserved"})
		return
	}

	for _, seat := range reserveBody.Seats {
		_, err = tx.Exec(`
			INSERT INTO Reservation (movie_id, user_id, date, seat)
			VALUES ($1, $2, $3, $4)
		`, movieId, userId, date, seat)
		if err != nil {
			generalError(c)
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "an error ocurred"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"date": date, "seats": reserveBody.Seats})
}

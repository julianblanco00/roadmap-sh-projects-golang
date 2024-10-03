package users

import (
	"movie-reservation-system/database"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type reservation struct {
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImageUrl    string    `json:"image_url"`
	Date        time.Time `json:"date"`
	Seat        string    `json:"seat"`
}

func GetAllMovieReservations(c *gin.Context) {
	movieId := c.Param("id")
	query := `
    SELECT 
      u.name,
      u.email,
      m.title,
      m.description,
      m.image_url,
      r.date,
      r.seat
    FROM reservation r
    JOIN users u ON r.user_id = u.id
    JOIN movies m ON r.movie_id = m.id
    WHERE movie_id = $1
  `

	rows, err := database.Db.Query(query, movieId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	defer rows.Close()

	var reservations []reservation
	var reservation reservation
	for rows.Next() {
		rows.Scan(
			&reservation.Name,
			&reservation.Email,
			&reservation.Title,
			&reservation.Description,
			&reservation.ImageUrl,
			&reservation.Date,
			&reservation.Seat,
		)
		reservations = append(reservations, reservation)
	}

	c.JSON(http.StatusOK, reservations)
}

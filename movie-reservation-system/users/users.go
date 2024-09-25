package users

import "movie-reservation-system/database"

type User struct {
	ID        int
	Name      string
	Birthdate string
	Email     string
	Password  string
}

func FindUserByEmail(email string) *User {
	row := database.Db.QueryRow("SELECT * FROM users WHERE email = $1", email)
	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Birthdate, &user.Email, &user.Password)
	if err != nil {
		return nil
	}

	return &user
}

func FindUserById(id int) *User {
	row := database.Db.QueryRow("SELECT * FROM users WHERE id = $1", id)
	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Birthdate, &user.Email, &user.Password)
	if err != nil {
		return nil
	}

	return &user
}

func GetUsers() ([]User, error) {
	rows, err := database.Db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Birthdate, &user.Email, &user.Password)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

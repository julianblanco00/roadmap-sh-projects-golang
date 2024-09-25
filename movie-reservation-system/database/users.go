package database

type User struct {
	ID        int
	Name      string
	Birthdate string
	Email     string
	Password  string
}

func FindUserByEmail(email string) *User {
	db, err := GetDB()
	if err != nil {
		return nil
	}

	row := db.QueryRow("SELECT * FROM users WHERE email = $1", email)
	var user User
	err = row.Scan(&user.ID, &user.Name, &user.Birthdate, &user.Email, &user.Password)
	if err != nil {
		return nil
	}

	return &user
}

func GetUsers() ([]User, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT * FROM users")
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

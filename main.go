package main
import (
	"github.com/gin-gonic/gin";
	"net/http";
	// "fmt";
	"log";
	"database/sql";
	_ "github.com/lib/pq"
)

func check(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "server is runing")
}

func getUserByName(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")

		var id int
		var dbName string

		// Query to retrieve a user by name
		err := db.QueryRow("SELECT id, name FROM task_user WHERE name = $1", name).Scan(&id, &dbName)
		if err != nil {
			if err == sql.ErrNoRows {
				c.IndentedJSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
			}
			log.Fatal(err)
		}

		user := gin.H{"id": id, "name": dbName}
		c.IndentedJSON(http.StatusOK, user)
	}
}

func createUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newUser struct {
			UserName string `json:"username"`
		}

		if err := c.ShouldBindJSON(&newUser); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Body must contain 'username'"})
			return
		}

		if newUser.UserName == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Username field cannot be empty"})
			return
		}

		// Check if the user already exists
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM task_user WHERE name = $1", newUser.UserName).Scan(&count)
		if err != nil {
			log.Fatal(err)
		}

		if count > 0 {
			c.IndentedJSON(http.StatusConflict, gin.H{"error": "User already exists"})
			return
		}

		// User does not exist, create a new user
		_, err = db.Exec("INSERT INTO task_user (name) VALUES ($1)", newUser.UserName)
		if err != nil {
			log.Fatal(err)
		}

		c.IndentedJSON(http.StatusCreated, gin.H{"message": "User created successfully"})
	}
	// return func(c *gin.Context) {
	// 	var newUser struct {
	// 		UserName string `json:"username"`
	// 	}

	// 	if err := c.ShouldBindJSON(&newUser); err != nil {
	// 		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Body must contain 'username'"})
	// 		return
	// 	}

	// 	if newUser.UserName == "" {
	// 		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Username field cannot be empty"})
	// 		return
	// 	}

	// 	// Insert the new user into the database
	// 	_, err := db.Exec("INSERT INTO task_user (name) VALUES ($1)", newUser.UserName)
	// 	if err != nil {
	// 		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
	// 		log.Fatal(err)
	// 		return
	// 	}

	// 	c.IndentedJSON(http.StatusCreated, gin.H{"message": "User created successfully"})
	// }
}

func deleteUser(db *sql.DB) gin.HandlerFunc{
	return func(c *gin.Context) {
		var deleteUser struct {
			UserName string `json:"username"`
		}

		if err := c.ShouldBindJSON(&deleteUser); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Body must contain 'username'"})
			return
		}

		if deleteUser.UserName == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Username field cannot be empty"})
			return
		}

		// Check if the user exists before deletion
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM task_user WHERE name = $1", deleteUser.UserName).Scan(&count)
		if err != nil {
			log.Fatal(err)
		}

		if count == 0 {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "User does not exist"})
			return
		}

		// User exists, proceed to delete
		_, err = db.Exec("DELETE FROM task_user WHERE name = $1", deleteUser.UserName)
		if err != nil {
			log.Fatal(err)
		}

		c.IndentedJSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	}
}

func updateUser(db *sql.DB) gin.HandlerFunc{
	return func(c *gin.Context) {
		userName := c.Param("name")

		var updateUser struct {
			NewUserName string `json:"username"`
		}

		if err := c.ShouldBindJSON(&updateUser); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Body must contain 'new_username'"})
			return
		}

		if updateUser.NewUserName == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "New username field cannot be empty"})
			return
		}

		// Check if the user exists before updating
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM task_user WHERE name = $1", userName).Scan(&count)
		if err != nil {
			log.Fatal(err)
		}

		if count == 0 {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "User does not exist"})
			return
		}

		// Update the user's name
		_, err = db.Exec("UPDATE task_user SET name = $1 WHERE name = $2", updateUser.NewUserName, userName)
		if err != nil {
			log.Fatal(err)
		}

		c.IndentedJSON(http.StatusOK, gin.H{"message": "User updated successfully"})
	}
}

func main() {
	connStr := "postgres://postgres:pass123@localhost:5433/postgres?sslmode=disable"

	// Open a connection to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Check if the connection is successful
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS task_user(
		id serial PRIMARY KEY,
		name varchar(255) UNIQUE
	)`

	_, err = db.Exec(createTableSQL)
	if err != nil {
    log.Fatal(err)
	}

	// Query to retrieve all data from the user table
	rows, err := db.Query("SELECT * FROM task_user")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	router := gin.Default()
	router.GET("/", check)
	
	// CRUD

	getUserHandler := getUserByName(db)
	router.GET("/user/:name", getUserHandler)

	createUserHandler := createUser(db)
	router.POST("/create", createUserHandler)

	deleteUserhandler := deleteUser(db)
	router.DELETE("/delete", deleteUserhandler)

	updateUserHandler := updateUser(db)
	router.PUT("/user/:name", updateUserHandler)



	router.Run("localhost:8080")
}
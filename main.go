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

	// 	// Check if the user already exists
	// 	var count int
	// 	err := db.QueryRow("SELECT COUNT(*) FROM task_user WHERE name = $1", newUser.UserName).Scan(&count)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	if count > 0 {
	// 		c.IndentedJSON(http.StatusConflict, gin.H{"error": "User already exists"})
	// 		return
	// 	}

	// 	// User does not exist, create a new user
	// 	_, err = db.Exec("INSERT INTO task_user (name) VALUES ($1)", newUser.UserName)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	c.IndentedJSON(http.StatusCreated, gin.H{"message": "User created successfully"})
	// }
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

		c.IndentedJSON(http.StatusCreated, gin.H{"message": "User created successfully"})
	}
}

func main() {
	// PostgreSQL connection string
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

	// Query to retrieve all data from the user table
	rows, err := db.Query("SELECT * FROM task_user")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	router := gin.Default()
	router.GET("/", check)
	
	// CRUD

	// Accepts GET requests with a name parameter in the URL to retrieve the corresponding user.
	getUserHandler := getUserByName(db)
	router.GET("/user/:name", getUserHandler)

	// Accepts GET requests with a name parameter in the URL to retrieve the corresponding user.
	createUserHandler := createUser(db)
	router.POST("/create", createUserHandler)

	router.Run("localhost:8080")
}
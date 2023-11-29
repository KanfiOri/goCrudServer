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

	// Iterate through the result set and print data
	// for rows.Next() {
	// 	var id int
	// 	var name string
	// 	err := rows.Scan(&id, &name)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Printf("ID: %d, Name: %s\n", id, name)
	// }

	// // Check for errors during row iteration
	// if err = rows.Err(); err != nil {
	// 	log.Fatal(err)
	// }

	router := gin.Default()
	router.GET("/", check)
	
	// CRUD
	// /user/{name} Accepts GET requests with a name parameter in the URL to retrieve the corresponding user.

	getUserHandler := getUserByName(db)
	router.GET("/user/:name", getUserHandler)

	router.Run("localhost:8080")
}
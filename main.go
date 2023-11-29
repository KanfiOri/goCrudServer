package main
import (
	"github.com/gin-gonic/gin";
	"net/http"
)

func check(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "server is runing")
}

func main() {
	router := gin.Default()
	router.GET("/", check)

	router.Run("localhost:8080")
}
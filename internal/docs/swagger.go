// Package docs contains the Swagger/OpenAPI documentation for the API
package docs

import (
	"fmt"
	"os"
)

// @title Linkeun Go API
// @version 1.0
// @description API for managing various resources including animals
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://linkeunid.com/support
// @contact.email support@linkeunid.com

// @license.name GNU General Public License v2.0
// @license.url https://www.gnu.org/licenses/old-licenses/gpl-2.0.en.html

// @host localhost:4445
// @BasePath /api/v1
// @schemes http https

// GetSwaggerHost returns the host for Swagger UI based on environment variables
func GetSwaggerHost() string {
	host := os.Getenv("API_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return fmt.Sprintf("%s:%s", host, port)
}

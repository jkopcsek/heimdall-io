package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/hubidu/e2e-backend/report-lib/db"
	"github.com/hubidu/e2e-backend/report-lib/middlewares"
	projects "github.com/hubidu/e2e-backend/report-service/routes"
	reports "github.com/hubidu/e2e-backend/report-service/routes"
	screenshots "github.com/hubidu/e2e-backend/report-service/routes"
)

func init() {
	db.Connect()

	db.InitializeIndexes()
}

func main() {
	r := gin.Default()

	r.RedirectTrailingSlash = true
	r.RedirectFixedPath = true

	// Middlewares
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(middlewares.Connect)
	r.Use(middlewares.ErrorHandler)
	// TODO Use production config
	r.Use(cors.Default())

	// Routes
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/reports")
	})

	r.GET("/reports", reports.List)
	r.GET("/reports/:_id", reports.Get)
	r.GET("/report-categories", reports.GetReportCategories)
	r.GET("/report-categories/:hashcategory", reports.GetReportsByCategory)

	r.GET("/screenshots/:path/:file", screenshots.GetScreenshotImg)
	r.GET("/screenshot-categories/:hashids", screenshots.GetScreenshotsByCategory)

	r.GET("/projects/:ownerkey", projects.GetByOwnerkey)

	r.Run("0.0.0.0:8000")
}

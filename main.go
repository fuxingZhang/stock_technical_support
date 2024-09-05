package main

import (
	"app/pkg"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
)

func registerHTMLRoutes(router *gin.Engine, baseDir string, prefix string) error {
	files, err := os.ReadDir(baseDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		path := filepath.Join(baseDir, file.Name())
		if file.IsDir() {
			err := registerHTMLRoutes(router, path, prefix+file.Name()+"/")
			if err != nil {
				return err
			}
		} else if strings.HasSuffix(strings.ToLower(file.Name()), ".html") {
			urlPath := prefix + strings.TrimSuffix(file.Name(), ".html")
			router.GET(urlPath, func(c *gin.Context) {
				c.File(path)
			})
		}
	}
	return nil
}

func main() {
	r := gin.Default()
	//router.StaticFS("/", http.Dir(filepath.Join(dir, "./html")))

	r.GET("/cookie", func(c *gin.Context) {
		var domain = c.Query("domain")
		var key = c.Query("key")

		if domain == "" {
			c.AbortWithStatus(400)
			return
		}

		result, err := pkg.GetCookie(domain)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}

		if key != "" {
			c.String(http.StatusOK, result[key])
			return
		}

		c.JSON(http.StatusOK, result)
	})

	r.GET("/", func(c *gin.Context) {
		c.File("./html/index.html")
	})

	r.Static("/static", "./static")

	err := registerHTMLRoutes(r, "./html", "/")
	if err != nil {
		log.Fatalf("Failed to register HTML routes: %v", err)
	}

	go func() {
		r.Run(":80")
	}()

	log.Fatal(autotls.Run(r, "zhangfuxing.icu", "zhangfuxing.asia"))
}

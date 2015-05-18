package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

// Initialize the database
var db = initDB()

func main() {

	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./assets")
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	r.GET("/", func(c *gin.Context) {
		shames := []Shame{}
		check(db.Table("shames").Order("created_at desc").Find(&shames).Error)

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"shames": shames,
		})
	})

	r.POST("/shame", func(c *gin.Context) {
		var payload struct {
			URL string `json:"url" binding:"required"`
		}
		c.Bind(&payload)

		owner, repo, ref, path, start, end, err := parseURL(payload.URL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		if start < 1 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errors.New("Start line must be 1 or greater")})
			return
		}

		client := github.NewClient(nil)
		result, _, _, err := client.Repositories.GetContents(owner, repo, path, &github.RepositoryContentGetOptions{Ref: ref})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		content, err := result.Decode()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		lines := strings.Split(string(content), "\n")
		if end >= len(lines) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errors.New("End line must be " + strconv.Itoa(len(lines)-1) + " or lower")})
			return
		}

		if end <= start {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errors.New("End line must be after start line")})
			return
		}

		snippet := lines[start-1 : end]
		indent := len(snippet[0]) - len(strings.TrimSpace(snippet[0]))
		for i, l := range snippet {
			if len(l) >= indent {
				snippet[i] = l[indent:]
			}
		}

		name := *result.Name

		shame := Shame{
			URL:       payload.URL,
			Reponame:  owner + "/" + repo,
			Path:      path,
			Content:   []byte(strings.Join(snippet, "\n")),
			Beginline: start,
			Endline:   end,
			Filename:  name,
		}

		if err := db.Table("shames").FirstOrCreate(&shame, Shame{URL: payload.URL}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		c.JSON(http.StatusOK, shame)
	})

	log.Fatalln(r.Run(":" + os.Getenv("PORT")))
}

func initDB() gorm.DB {
	// Initialize
	db, err := gorm.Open("postgres", os.Getenv("DATABASE_URL"))
	fatal(err)
	err = db.DB().Ping()
	fatal(err)

	// Migrate
	db.AutoMigrate(&Shame{})

	// Setup
	db.Model(&Shame{}).AddUniqueIndex("idx_shames_url", "url")

	return db
}

var urlMatcher = regexp.MustCompile(`^https://github.com/(?P<owner>[^/]+)/(?P<repo>[^/]+)/blob/(?P<ref>[^/]+)/(?P<path>[^#]+)#L([0-9]+)-L([0-9]+)`)

func parseURL(url string) (string, string, string, string, int, int, error) {
	match := urlMatcher.FindStringSubmatch(url)
	if match == nil || len(match) < 7 {
		return "", "", "", "", -1, -1, errors.New("Invalid URL")
	}

	start, _ := strconv.Atoi(match[5])
	end, _ := strconv.Atoi(match[6])

	return match[1], match[2], match[3], match[4], start, end, nil
}

package main

import (
	"net/http"

	"github.com/curiouscat2018/helloworld-api/common"
	"github.com/curiouscat2018/helloworld-api/config"
	"github.com/curiouscat2018/helloworld-api/database"
	"github.com/curiouscat2018/helloworld-api/middleware"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var myDB database.Database

func main() {
	r := gin.New()
	r.Use(middleware.ContextMiddleware())
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.HeaderMiddleware())
	r.Use(middleware.ErrorMiddleware())
	r.Use(gin.Recovery())

	r.GET("/", handleIndex)
	r.GET("/ping", handlePing)
	r.NoRoute(handleNotFound)

	if config.Config.IsMockEnv {
		prepareMockEnv(&myDB)
	} else {
		prepareProdEnv(&myDB)
	}

	common.TraceInfo(nil).Msgf("start listening helloworld-api: is mock env: %v", config.Config.IsMockEnv)
	common.TraceError(nil).Err(r.Run(":http")).Msg("helloworld-api stopped")
}

func prepareProdEnv(db *database.Database) {
	common.TraceInfo(nil).Msg("preparing Azure Database")
	tempDB, err := database.NewAzureDatabase(config.Config.DB_URL)
	if err != nil {
		common.TraceFatal(nil).Err(err).Msg("")
	}
	*db = tempDB
}

func prepareMockEnv(db *database.Database) {
	common.TraceInfo(nil).Msg("preparing mock Database")
	tempDB := database.NewMockDatabase()
	*db = tempDB
}

func handleNotFound(c *gin.Context) {
	setHttpError(c, http.StatusNotFound, errors.New("page not found"))
}

func handlePing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func handleIndex(c *gin.Context) {
	entry, err := myDB.GetEntry(c)
	if err != nil {
		setHttpError(c, http.StatusInternalServerError, err)
		return
	}

	response := struct {
		database.Entry
		Hostname string `json:"hostname"`
	}{
		Entry:    *entry,
		Hostname: config.Config.HostName(),
	}

	c.JSON(http.StatusOK, &response)
}

func setHttpError(c *gin.Context, code int, err error) {
	c.Status(code)
	c.Error(err)
}

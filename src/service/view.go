package service

import (
	"embed"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"me.daily/src/db"
	"me.daily/src/log"

	"github.com/patrickmn/go-cache"
)

// swagger embed files

// 日期格式
const dateFormat = "2006-01-02"

// 過期時間
const expiredTime = 8 * 60 * 60

type Service struct {
	a   gin.Accounts
	c   *cache.Cache
	d   *db.Db
	fsh http.Handler
	s   *gin.Engine
}

func NewService(host, user, password, dbname, authUser, authPw string, fs embed.FS) *Service {
	a := make(gin.Accounts)
	a[authUser] = authPw

	return &Service{
		a:   a,
		c:   cache.New(expiredTime*time.Second, 60*time.Minute),
		d:   db.NewDb(host, user, password, dbname),
		fsh: http.FileServer(http.FS(fs)),
		s:   gin.New(),
	}
}

// swag init
// http://localhost:8080/swagger/index.html
// go get -u github.com/swaggo/gin-swagger
// go get -u github.com/swaggo/files
func (s *Service) Start() {
	s.s.RedirectFixedPath = true

	s.s.Use(gin.Recovery(), log.LogHistory.Func)

	// Access-Control-Allow-Origin
	s.s.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", c.Request.Host)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// docs.SwaggerInfo.BasePath = "/"
	// s.s.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// static files
	s.s.Any("/public/*any", func(c *gin.Context) {
		s.fsh.ServeHTTP(c.Writer, c.Request)
	})

	// resource
	{
		gInfo := s.s.Group("info")

		gInfo.Use(gin.BasicAuth(s.a))

		gInfo.GET("/log", s.getLog)
	}

	// Api
	{
		gApi := s.s.Group("/api")

		gApi.Use(s.checkAuth)

		gApi.GET("/main", s.getMainType)
		gApi.GET("/sub/:main_id", s.getSubType)
		gApi.GET("/all", s.getAll)
		gApi.GET("/item/:item_id", s.getItem)
		gApi.GET("/items", s.getItems)
		gApi.GET("/spend/month/:count", s.getSpendByLastMonthly)
		gApi.GET("/sum/main", s.getSumByMainType)

		gApi.GET("logout", s.logout)
		gApi.POST("/login", s.login)
		gApi.POST("/user", s.createUser)
		gApi.POST("/main", s.createMainType)
		gApi.POST("/sub", s.createSubType)
		gApi.POST("/item", s.createItem)

		gApi.PUT("/main", s.updateMainType)
		gApi.PUT("/sub", s.updateSubType)
		gApi.PUT("/item", s.updateItem)

		gApi.DELETE("/main/:main_id", s.deleteMainType)
		gApi.DELETE("/sub/:sub_id", s.deleteSubType)
		gApi.DELETE("/item/:item_id", s.deleteItem)
	}

	s.s.Run(":80")
}

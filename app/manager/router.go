package manager

import (
	"net/http"

	"gin-demo/app/docs"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	routerManager *RouterManager
)

func GetRouter() *RouterManager {
	return routerManager
}

type RouterManager struct {
	router *gin.Engine
	config RouterConfig
}

type RouterConfig struct {
	Version string
}

func (manager *RouterManager) Setup(config RouterConfig) (err error) {

	router := gin.New()

	routerManager = &RouterManager{router: router, config: config}
	routerManager.setupRouter()

	return nil
}

func (manager *RouterManager) Close() {

}

func (manager *RouterManager) GetHandler() *gin.Engine {
	return manager.router
}

func (manager *RouterManager) setupRouter() {

	docs.SwaggerInfo.Version = manager.config.Version

	gin.SetMode(gin.ReleaseMode)

	pprof.Register(manager.router)

	//manager.router.LoadHTMLGlob("static/*.html")

	manager.router.Use(gin.LoggerWithWriter(gin.DefaultWriter, "/health", "/version"))
	manager.router.Use(gin.Recovery(), gin.Logger(), CORSMiddleware(), ErrorMiddleware())

	manager.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	manager.router.GET("/health", HealthHandler)
	manager.router.GET("/version", VersionHandler)
}

func CORSMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers",
			"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-HTC-Account-Id, AuthKey")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func ErrorMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Next()
		err := c.Errors.Last()
		if err == nil {
			return
		}
	}
}

// HealthHandler is health checker API
//
//	@Tags		Default
//	@Success	200	{string}	string	"ok"
//	@Router		/health [get]
func HealthHandler(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}

// VersionHandler is version checker API
//
//	@Tags		Default
//	@Success	200	{string}	string	"1.0.0"
//	@Router		/version [get]
func VersionHandler(c *gin.Context) {
	version := routerManager.config.Version
	c.String(http.StatusOK, version)
}

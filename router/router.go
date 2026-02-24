package router

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gitlab.com/cinemae/cine_stream/app/controller"
	"gitlab.com/cinemae/cine_stream/filter"
)

// router 路由相关

var (
	// GinEngine Gin 框架实例
	GinEngine = gin.New()
	// routerHandleTables 路由处理表，新增一个接口配置一条
	routerHandleTables = []routerHandle{
		{group: "/demo", relativePath: "/index", method: http.MethodGet, controllerHandle: controller.DemoIndex},
		// TS切片相关接口
		{group: "/video_ts", relativePath: "/save", method: http.MethodPost, controllerHandle: controller.VideoTsSave},
		{group: "/video_ts", relativePath: "/list", method: http.MethodGet, controllerHandle: controller.VideoTsList},

		// 播放相关
		{group: "/play", relativePath: "/:video_id", method: http.MethodGet, controllerHandle: controller.Play},
		{group: "/play", relativePath: "/:video_id/index.m3u8", method: http.MethodGet, controllerHandle: controller.PlayHlsIndexM3u8},
		{group: "/play", relativePath: "/key/:video_id", method: http.MethodGet, controllerHandle: controller.PlayHlsIndexEncKey},

		// cine 播放器私有协议
		{group: "/play", relativePath: "/:video_id/index.c3u8", method: http.MethodGet, controllerHandle: controller.PlayCineHlsIndexC3u8},

		// 资源站点接口
		{group: "/provide", relativePath: "/json", method: http.MethodGet, controllerHandle: controller.ProvideIndex},
		{group: "/provide", relativePath: "/xml", method: http.MethodGet, controllerHandle: controller.ProvideIndex},
		{group: "/provide", relativePath: "/save", method: http.MethodPost, controllerHandle: controller.ProvideSave},
	}
)

// routerHandle 路由处理配置
type routerHandle struct {
	group            string
	relativePath     string
	method           string
	controllerHandle controller.HandleFunc
}

// Init 初始化路由
func Init() {
	// 使用跨域中间件，允许 credentials（cookies）
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowCredentials = true
	corsConfig.AllowAllOrigins = true // 如果需要限制域名，可以改为 AllowOrigins
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Cookie"}
	corsConfig.ExposeHeaders = []string{"Content-Length", "Content-Range"}
	GinEngine.Use(cors.New(corsConfig))
	// 使用 logger 中间件
	GinEngine.Use(gin.Logger())
	// 使用 recover 中间件
	GinEngine.Use(gin.Recovery())
	// 自定义 公共参数解析 中间件
	GinEngine.Use(filter.RequestParse())
	// 自定义 auth 登录认证 中间件
	GinEngine.Use(filter.AuthLoginJWT())
	// 自定义 打印耗时 中间件
	GinEngine.Use(filter.DebugCosTime())

	// Swagger 静态文件服务
	GinEngine.StaticFS("/swagger", http.Dir("./swagger"))

	// 初始化路由表
	initRouter()
}

// RegisterHandle 注册一个路由处理
func RegisterHandle(groupPath string, relativePath string, method string, controllerHandle controller.HandleFunc) {
	routerHandleTables = append(routerHandleTables, routerHandle{
		group:            groupPath,
		relativePath:     relativePath,
		method:           method,
		controllerHandle: controllerHandle,
	})
}

func initRouter() {
	// 转化称 group => handles
	routerGroupTables := make(map[string][]routerHandle)
	for _, routerHandleItem := range routerHandleTables {
		if _, ok := routerGroupTables[routerHandleItem.group]; !ok {
			routerGroupTables[routerHandleItem.group] = []routerHandle{}
		}
		routerGroupTables[routerHandleItem.group] = append(routerGroupTables[routerHandleItem.group], routerHandleItem)
	}
	for routerGroup, routerHandles := range routerGroupTables {
		groupRouter := GinEngine.Group(routerGroup)
		for _, routerHandleItem := range routerHandles {
			groupRouter.Handle(
				routerHandleItem.method,
				routerHandleItem.relativePath,
				HandleWrapper(routerHandleItem.controllerHandle),
			)
		}
	}
}

// HandleWrapper 控制器 Wrapper
func HandleWrapper(controllerHandle controller.HandleFunc) func(*gin.Context) {
	return func(context *gin.Context) {
		_ = controllerHandle(context)
	}
}

package core

import (
	"fmt"
	"gin-vue-admin/global"
	"gin-vue-admin/initialize"
	"go.uber.org/zap"
	"time"
)

type server interface {
	ListenAndServe() error
}

// @zgz
func InitRouter() {
	if global.GVA_CONFIG.System.UseMultipoint {
		// 初始化redis服务
		initialize.Redis()
	}
	initialize.InitWkMode()
	Router := initialize.Routers()
	Router.Static("/form-generator", "./resource/page")

	global.GVA_ROUTER = Router
}

func RunWindowsServer() {

	address := fmt.Sprintf(":%d", global.GVA_CONFIG.System.Addr)
	s := initServer(address, global.GVA_ROUTER)
	// 保证文本顺序输出
	// In order to ensure that the text order output can be deleted
	time.Sleep(10 * time.Microsecond)
	global.GVA_LOG.Info("server run success on ", zap.String("address", address))

	fmt.Printf(`
	欢迎使用 GVA
	当前版本:V2.3.8
	默认自动化文档地址:http://127.0.0.1%s/swagger/index.html
	默认前端文件运行地址:http://127.0.0.1:8080
`, address)
	global.GVA_LOG.Error(s.ListenAndServe().Error())
}

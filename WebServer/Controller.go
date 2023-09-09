package WebServer

func (w *WebServer) ApiRouterRules(pre string) {
	r := w.e.Group(pre)
	r = w.e.Group("/api")

	r.Use(w.keyAuthMiddleware)
	//r.POST("/devInfo/all", WebApi.GetAllDevInfo)
	//r.POST("/roastInfo/devId", WebApi.GetRoastInfoById)
	////r.POST("/roastInfo/devId", WebApi.GetAllDevInfo)
	//r.POST("/get/tempData/devId", WebApi.GetTempInfo)

}

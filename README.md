# BlockChainBackEnd

物流管理系统
一个前后端分离的轻量化物流管理系统，后端使用Gin+Mysql+Redis+JWT 等。[前端传送门](https://gitee.com/dcolor/traceability-system)

go 1.18

运行程序
```
 go run main.go
```
生产环境需要修改
在命令行设置
> export GIN_MODE=release

同时在程序中设置ReleaseMode即可即可
```
gin.SetMode(gin.ReleaseMode)
```

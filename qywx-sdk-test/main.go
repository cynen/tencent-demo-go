package main

import (
	"github.com/gin-gonic/gin"
	"github.com/xen0n/go-workwx"
	"os"
)

// 使用 qywx-sdk 简单示范
func main() {
	// 创建企业微信的实例
	Qywx := workwx.New("ww***")
	// 创建app实例.
	WxApp := Qywx.WithApp("-FbzHCvy***********************Itc", 1000002)

	//创建web服务器
	engine := gin.Default()

	rec := workwx.Recipient{}
	rec.UserIDs = append(rec.UserIDs, "xiaohao")
	engine.GET("/send", func(context *gin.Context) {
		msg := context.DefaultQuery("msg", "未收到")
		err := WxApp.SendTextMessage(&rec, msg, false)
		if err != nil {
			context.String(501, "", err.Error())
		}
		context.String(200, "OK")
	})

	// 启动端口
	engine.Run(":8888")

}

// 上传临时文件.获取mediaid
// 的确方便很多.
func uploadfile(WxApp *workwx.WorkwxApp, file *os.File) string {
	uploadMedia, _ := workwx.NewMediaFromFile(file)
	uploadR, _ := WxApp.UploadTempFileMedia(uploadMedia)
	return uploadR.MediaID
}

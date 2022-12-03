package service

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func PostUpload(c *gin.Context) {
	// 默认能上传文件的大小为32MB，可以通过 iris.WithPostMaxMemory(maxSize) 设置，
	// 比如10MB = 10 * 1024 * 1024 =maxSize，iris.WithPostMaxMemory(maxSize)
	// 3个参数，1是文件句柄，2是文件头，3是错误信息
	uploadFile, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"msg":     "获取文件信息失败!" + err.Error(),
		})
	}
	if uploadFile != nil { // 记得及时关闭文件，避免内存泄漏
		defer uploadFile.Close()
	}

	// 获取原文件名字
	fname := fileHeader.Filename
	// 创建1个相同名字的文件，存放在upload目录里面
	// 假定本地已经有名字为 upload 的目录，没有的话会报错
	out, err := os.OpenFile("./upload/"+fname, os.O_WRONLY|os.O_CREATE, 0666)
	defer out.Close()

	io.Copy(out, uploadFile)

	// 读取excel表格的数据
	re, err := excelize.OpenFile("./upload/" + fname)
	if err != nil {
		log.Println("读取表格数据错误：", err)
		return
	}
	//
	//  获取行数，[][]string 类型返回值
	rows, _ := re.GetRows("Sheet1")
	for _, row := range rows {
		// 用 " | " 间隔拼接每一行数据
		str := "| "
		for _, val := range row {
			str += val + " | "
		}
		fmt.Println(str)
	}
}

func GetDownload(c *gin.Context) {
	file := excelize.NewFile()
	// 1个单元格1个单元格添加
	file.SetCellValue("Sheet1", "A1", "IP")
	file.SetCellValue("Sheet1", "B1", "ASN")
	file.SetCellValue("Sheet1", "C1", "国家")
	file.SetCellValue("Sheet1", "D1", "省份")
	file.SetCellValue("Sheet1", "E1", "城市")
	file.SetCellValue("Sheet1", "F1", "域名")

	// 1行添加
	row := []interface{}{"118.122.233.42", "4134", "中国", "四川省", "成都市", "chinatelecom.com.cn"}
	file.SetSheetRow("Sheet1", "A2", &row) // 传递切片指针

	// 暂存再 tmp 目录，之后再删掉
	file.SaveAs("./tmp/IP 批量上传表.xlsx")
	defer os.Remove("./tmp/IP 批量上传表.xlsx")

	f, _ := os.Open("./tmp/IP 批量上传表.xlsx")
	defer f.Close()

	// 将文件读取出来
	data, err := ioutil.ReadAll(f)
	if err != nil {
		// log.Fatal(err)
		c.JSON(http.StatusOK, gin.H{
            "success": false,
            "msg":     "获取文件信息失败!" + err.Error(),
        })
	}
	// 设置头信息：Content-Disposition ，消息头指示回复的内容该以何种形式展示，
	// 是以内联的形式（即网页或者页面的一部分），还是以附件的形式下载并保存到本地
	// Content-Disposition: inline
	// Content-Disposition: attachment
	// Content-Disposition: attachment; filename="filename.后缀"
	// 第一个参数或者是inline（默认值，表示回复中的消息体会以页面的一部分或者
	// 整个页面的形式展示），或者是attachment（意味着消息体应该被下载到本地；
	// 大多数浏览器会呈现一个“保存为”的对话框，将filename的值预填为下载后的文件名，
	// 假如它存在的话）。

	fileContentDisposition := `attachment; filename="IP 批量上传表.xlsx"`
	c.Header("Content-Type", "application/octet-stream") // 这里是压缩文件类型 .zip
	c.Header("Content-Disposition", fileContentDisposition)
	c.Data(http.StatusOK, ".xlsx", data)
}

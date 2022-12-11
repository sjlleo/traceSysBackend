package service

import (
	"io"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/sjlleo/traceSysBackend/models"
	"github.com/xuri/excelize/v2"
)

func GetReviewList(c *gin.Context) {
	u := GetRole(c)
	p := models.PaginationQ{}
	if err := c.ShouldBind(&p); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	if err := u.SearchReview(&p); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(200, p)
	}
}

func PassReview(c *gin.Context) {
	u := GetRole(c)
	reviewID, err := strconv.Atoi(c.Param("review_id"))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	u.PassReview(uint(reviewID))
	c.JSON(200, gin.H{"code": 200, "msg": "提交成功"})
}

func DeclineReview(c *gin.Context) {
	u := GetRole(c)
	reviewID, err := strconv.Atoi(c.Param("review_id"))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	u.DeclineReview(uint(reviewID))
	c.JSON(200, gin.H{"code": 200, "msg": "提交成功"})
}

func DeleteReview(c *gin.Context) {
	reviewID, err := strconv.Atoi(c.Param("review_id"))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	models.DeleteReview(uint(reviewID))
	c.JSON(200, gin.H{"code": 200, "msg": "删除成功"})
}

func AddReview(c *gin.Context) {
	var review models.IPReviews

	if err := c.ShouldBind(&review); err!= nil {
		c.JSON(500, gin.H{"error": err.Error()})
        return
	}

	if err := models.AddReview(review); err!= nil {
		c.JSON(500, gin.H{"error": err.Error()})
        return
	}

	c.JSON(200, gin.H{"code": 200, "msg": "添加成功，请耐心等待管理员审核"})
}

func PostUpload(c *gin.Context) {
	uploadFile, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"msg":     "获取文件信息失败!",
		})
	}
	if uploadFile != nil {
		defer uploadFile.Close()
	}

	nano_id, _ := gonanoid.New()
	out, _ := os.OpenFile("./upload/"+nano_id, os.O_WRONLY|os.O_CREATE, 0666)
	defer out.Close()
	defer os.Remove("./upload/" + nano_id)

	io.Copy(out, uploadFile)

	// 读取excel表格的数据
	re, err := excelize.OpenFile("./upload/" + nano_id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"msg":     "读取表格数据失败!",
		})
		return
	}
	//
	//  获取行数，[][]string 类型返回值
	rows, _ := re.GetRows("Sheet1")
	for _, row := range rows {

		// 校验模块
		if len(row) != 7 {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"msg":     "格式与模板不匹配，请检查是否删除或增加一些字段",
			})
			return
		}
		if row[0] == "IP" {
			continue
		}
		// IP 校验
		if net.ParseIP(row[0]) == nil {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"msg":     "IP 格式不合法",
			})
			return
		}
		// Prefix 校验
		prefix_int, err := strconv.Atoi(row[1])
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"msg":     "Prefix 不为数字",
			})
			return
		}
		// ASN 校验
		asn_int, err := strconv.Atoi(row[2])
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"msg":     "ASN 不为数字",
			})
			return
		}
		r_data := models.IPReviews{
			IP:         row[0],
			Prefix:     uint(prefix_int),
			ASN:        uint(asn_int),
			Country:    row[3],
			Province:   row[4],
			City:       row[5],
			Domain:     row[6],
			Authorized: 0,
		}
		if err := models.AddReview(r_data); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"msg":     "数据插入失败",
			})
			return
		}
	}
}

func GetDownload(c *gin.Context) {
	file := excelize.NewFile()
	// 1个单元格1个单元格添加
	file.SetCellValue("Sheet1", "A1", "IP")
	file.SetCellValue("Sheet1", "B1", "Prefix")
	file.SetCellValue("Sheet1", "C1", "ASN")
	file.SetCellValue("Sheet1", "D1", "国家")
	file.SetCellValue("Sheet1", "E1", "省份")
	file.SetCellValue("Sheet1", "F1", "城市")
	file.SetCellValue("Sheet1", "G1", "域名")
	file.SetColWidth("Sheet1", "A", "A", 20)
	file.SetColWidth("Sheet1", "G", "G", 20)

	// 1行添加
	row := []interface{}{"118.122.233.0", "24", "4134", "中国", "四川省", "成都市", "chinatelecom.com.cn"}
	file.SetSheetRow("Sheet1", "A2", &row) // 传递切片指针
	styleID, _ := file.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
		},
	})
	file.SetColStyle("Sheet1", "A:G", styleID)

	// 暂存再 tmp 目录，之后再删掉
	file.SaveAs("./tmp/IP 批量上传表.xlsx")
	defer os.Remove("./tmp/IP 批量上传表.xlsx")

	f, _ := os.Open("./tmp/IP 批量上传表.xlsx")
	defer f.Close()

	// 将文件读取出来
	data, err := io.ReadAll(f)
	if err != nil {
		// log.Fatal(err)
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"msg":     "生成表格错误!",
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

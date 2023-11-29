package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Gin struct {
	C *gin.Context
}

func appGin(c *gin.Context) *Gin {
	return &Gin{C: c}
}

func (g *Gin) RespData(data interface{}) {
	//前端需要返回data
	g.C.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "ok",
		"data": data,
	})
}
func (g *Gin) RespOK() {
	g.C.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "ok",
	})
	return
}

func (g *Gin) RespOKData(msg string) {
	g.C.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  msg,
	})
}

func (g *Gin) RespErr(err error) {
	g.C.JSON(http.StatusOK, gin.H{
		"code": "500",
		"msg":  err.Error(),
	})

}

func (g *Gin) RespErrData(msg string) {
	g.C.JSON(http.StatusOK, gin.H{
		"code": 500,
		"msg":  msg,
	})
}

// 处理请求体为json的格式
func ParseJSON(c *gin.Context, target interface{}) error {
	if err := c.ShouldBindJSON(target); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}
	return nil
}

package registry

import "github.com/gin-gonic/gin"

type RegisterOperation interface {
	HttpBindOp(engine *gin.Engine)
	Up(context gin.Context)
	Down(context gin.Context)
}

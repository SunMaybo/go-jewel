package registry

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type EtcRegisterOperation struct {
	Registry Registry
}

func (op EtcRegisterOperation) HttpBindOp(engine *gin.Engine) {
	engine.GET("/start", op.Up)
	engine.GET("/stop", op.Down)
	engine.GET("/services", op.Services)
}

func (op EtcRegisterOperation) Up(context *gin.Context) {
	_, err := op.Registry.Up()
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"status": gin.H{
				"error": err.Error(),
				"code":  500,
			},
		})
	} else {
		context.JSON(http.StatusOK, gin.H{"message": "api invoke success",
		})
	}
}
func (op EtcRegisterOperation) Services(context *gin.Context) {
	services, err := op.Registry.Services()
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"status": gin.H{
				"error": err.Error(),
				"code":  500,
			},
		})
	} else {
		context.JSON(http.StatusOK, gin.H{"message": "api invoke success",
			"result": services,
		})
	}
}

func (op EtcRegisterOperation) Down(context *gin.Context) {
	err := op.Registry.Down()
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"status": gin.H{
				"error": err.Error(),
				"code":  500,
			},
		})
	} else {
		context.JSON(http.StatusOK, gin.H{"message": "api invoke success",
		})
	}
}

package registry

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/SunMaybo/go-jewel/context"
	"github.com/SunMaybo/jewel-inject/inject"
)

type EtcRegisterOperation struct {
	Registry Registry
}

func (op EtcRegisterOperation) HttpBindOp(router *gin.RouterGroup, injector *inject.Injector) {
	jewel := injector.Service(&context.JewelProperties{}).(context.JewelProperties)
	manager := jewel.Jewel.Server.Manager
	var r *gin.RouterGroup
	if manager.Enabled != nil && *manager.Enabled {
		accounts := make(gin.Accounts)
		accounts[*manager.User] = *manager.Password
		r = router.Group("admin", gin.BasicAuth(accounts))
	} else {
		r = router.Group("admin")
	}
	r.GET("/start", op.Up)
	r.GET("/stop", op.Down)
	r.GET("/services", op.Services)
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

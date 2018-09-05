package jewel

import (
	"testing"
	"github.com/gin-gonic/gin"
)

func TestApplicationBoot(t *testing.T) {
	jewel := NewHttp()
	jewel.HttpStart(func(engine *gin.Engine) {

	})
}

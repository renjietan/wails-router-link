package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"wails-router-link/service/types"
)

func SUCCESS(c *gin.Context, values ...interface{}) {
	if values != nil {
		c.JSON(http.StatusOK, types.Response{
			Code:    types.Success,
			Message: types.SuccessMsg,
			Data:    values[0],
			Detail:  "",
		})
	} else {
		c.JSON(http.StatusOK, types.Response{
			Code:    types.Success,
			Message: types.SuccessMsg,
			Data:    nil,
			Detail:  "",
		})
	}

}

func ERROR(c *gin.Context, messages ...string) {
	if messages != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    types.Failed,
			Message: types.FailedMsg,
			Data:    nil,
			Detail:  messages[0],
		})
	} else {
		c.JSON(http.StatusBadRequest, types.Response{
			Code:    types.Failed,
			Message: types.FailedMsg,
			Data:    nil,
			Detail:  "",
		})
	}
}

func NotAuth(c *gin.Context, messages ...string) {
	if messages != nil {
		c.JSON(http.StatusUnauthorized, types.Response{
			Code:    types.NotAuthorized,
			Message: types.FailedMsg,
			Data:    nil,
			Detail:  messages[0],
		})
	} else {
		c.JSON(http.StatusUnauthorized, types.Response{
			Code:    types.NotAuthorized,
			Message: types.NoAuth,
			Data:    nil,
			Detail:  "",
		})
	}
}

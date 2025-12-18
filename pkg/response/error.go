package response

import (
	"net/http"

	"github.com/ductran999/letobserv/pkg/xcontext"
	"github.com/gin-gonic/gin"
)

type ErrorDetails struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResp struct {
	Success bool         `json:"success"`
	Error   ErrorDetails `json:"error"`
	TraceID string       `json:"trace_id"`
}

func BadRequest(c *gin.Context, code string, err error) {
	resp := ErrorResp{
		Success: false,
		Error: ErrorDetails{
			Code:    code,
			Message: err.Error(),
		},
		TraceID: c.GetString(xcontext.TraceIDKey),
	}
	c.JSON(http.StatusBadRequest, resp)
}

func InternalServerError(c *gin.Context, code string, err error) {
	resp := ErrorResp{
		Success: false,
		Error: ErrorDetails{
			Code:    code,
			Message: err.Error(),
		},
		TraceID: c.GetString(xcontext.TraceIDKey),
	}
	c.JSON(http.StatusInternalServerError, resp)
}

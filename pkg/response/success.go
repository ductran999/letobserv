package response

import (
	"net/http"
	"strings"

	"github.com/ductran999/letobserv/pkg/xcontext"
	"github.com/gin-gonic/gin"
)

type SuccessResp struct {
	Success bool    `json:"success"`
	Data    any     `json:"data"`
	Message *string `json:"message,omitempty"`
	TraceID string  `json:"trace_id"`
}

func OK(c *gin.Context, data any, message ...string) {
	msg := strings.Join(message, ", ")
	resp := SuccessResp{
		Success: true,
		Data:    data,
		Message: &msg,
		TraceID: c.GetString(xcontext.TraceIDKey),
	}
	c.JSON(http.StatusOK, resp)
}

func Created(c *gin.Context, data any, message ...string) {
	msg := strings.Join(message, ", ")
	resp := SuccessResp{
		Success: true,
		Data:    data,
		Message: &msg,
		TraceID: c.GetString(xcontext.TraceIDKey),
	}
	c.JSON(http.StatusCreated, resp)
}

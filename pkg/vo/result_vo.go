package vo

import (
	"errors"

	"github.com/gin-gonic/gin"
	cerr "github.com/narcissus1949/narcissus-blog/internal/error"
)

type Result struct {
	*cerr.Error
	Data      any         `json:"data"`
	RequestId interface{} `json:"requestId"`
}

func Success(c *gin.Context, data any) Result {
	return Result{
		Error: cerr.NewSuccess(),
		Data:  data,
		// RequestId: c.Request.Response.Header.Get("X-Request-Id"),
	}
}

func Fail(c *gin.Context, data any, err error) Result {
	var cerrErr *cerr.Error
	if ok := errors.As(err, &cerrErr); ok {
		return Result{
			Error: cerrErr,
			Data:  data,
			// RequestId: c.Request.Response.Header.Get("X-Request-Id"),
		}
	}

	return Result{
		Error: cerr.New(cerr.SYSTEM_ERROR),
		Data:  data,
		// RequestId: c.Request.Response.Header.Get("X-Request-Id"),
	}
}

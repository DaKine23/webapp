package middleware

import (
	"net/http"

	"github.bus.zalan.do/ale/gocore/uuid"
	"github.com/gin-gonic/gin"
)

func extractFlowId(r *http.Request) string {
	flowId := r.Header.Get(FlowIDHeaderKey)
	if len(flowId) == 0 {
		flowId, _ = uuid.New()
	}
	return flowId
}

func FlowID() gin.HandlerFunc {

	return func(c *gin.Context) {

		flowid := extractFlowId(c.Request)
		if c.Keys == nil {
			c.Keys = make(map[string]interface{})
		}
		c.Keys[FlowIDKey] = flowid

		c.Next()

	}
}

package google

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/apex/log"
)

func LogRequestFieldProvider(ctx *gin.Context) log.Fielder {
	traceHeader := ctx.Request.Header.Get("X-Cloud-Trace-Context")
	traceParts := strings.Split(traceHeader, "/")
	if len(traceParts) > 0 && len(traceParts[0]) > 0 {
		trace := fmt.Sprintf("projects/%s/traces/%s", os.Getenv("GC_PROJECT_ID"), traceParts[0])
		return log.Fields{"logging.googleapis.com/trace": trace}
	}
	return nil
}

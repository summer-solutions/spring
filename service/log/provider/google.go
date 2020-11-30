package provider

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/summer-solutions/spring/gin"

	"github.com/apex/log"
)

func Google(ctx context.Context) log.Interface {
	var trace string
	g := gin.FromContext(ctx)
	traceHeader := g.Request.Header.Get("X-Cloud-Trace-Context")
	traceParts := strings.Split(traceHeader, "/")
	if len(traceParts) > 0 && len(traceParts[0]) > 0 {
		trace = fmt.Sprintf("projects/%s/traces/%s", os.Getenv("GC_PROJECT_ID"), traceParts[0])
	}
	return log.WithField("logging.googleapis.com/trace", trace)
}

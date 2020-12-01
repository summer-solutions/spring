package google

import (
	j "encoding/json"
	"io"
	"os"
	"sync"

	"github.com/apex/log"
)

var Default = New(os.Stderr)

type Handler struct {
	*j.Encoder
	mu sync.Mutex
}

func New(w io.Writer) *Handler {
	return &Handler{
		Encoder: j.NewEncoder(w),
	}
}

func (h *Handler) HandleLog(e *log.Entry) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	row := e.Fields
	row["severity"] = "ERROR"
	row["message"] = e.Message
	row["timestamp"] = e.Timestamp
	return h.Encoder.Encode(row)
}

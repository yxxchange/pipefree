package internal

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type WatchServer struct {
	param                WatchParam
	ExpectConnectionType string
}

func (s *WatchServer) Serve(c *gin.Context) {
	s.ServeHTTP(c.Writer, c.Request)
}

func (s *WatchServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f, ok := w.(http.Flusher)
	if !ok {
		// TODO: log
		return
	}
	s.beginChunkedStream(w, f)
	// todo: get event for watch kind node
}

type WatchParam struct {
	Namespace string `uri:"namespace" binding:"required"`
	Name      string `uri:"name" binding:"required"`
	Kind      string `query:"kind" binding:"required"`
}

func (s *WatchServer) beginChunkedStream(w http.ResponseWriter, f http.Flusher) {
	w.Header().Set("Content-Type", s.ExpectConnectionType)
	w.Header().Set("Transfer-Encoding", "chunked")
	w.WriteHeader(http.StatusOK)
	f.Flush()

	return
}

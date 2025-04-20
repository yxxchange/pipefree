package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/yxxchange/pipefree/helper/log"
	"net/http"
)

type WatchServer struct {
	param WatchParam
}

func LaunchServer(c *gin.Context, param WatchParam) {
	server := &WatchServer{
		param: param,
	}
	server.Serve(c)
}

func (s *WatchServer) Serve(c *gin.Context) {
	s.ServeHTTP(c.Writer, c.Request)
}

func (s *WatchServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f, ok := w.(http.Flusher)
	if !ok {
		log.Errorf("type asset failed: %T not implement http.Flusher", w)
		return
	}
	s.beginChunkedStream(w, f)
	// todo: get event for watch kind node
	log.Info("start to watch")
}

type WatchParam struct {
	Namespace string `uri:"namespace" binding:"required"`
	Name      string `uri:"name" binding:"required"`
	Kind      string `query:"kind" binding:"required"`
}

func (s *WatchServer) beginChunkedStream(w http.ResponseWriter, f http.Flusher) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.WriteHeader(http.StatusOK)
	f.Flush()

	return
}

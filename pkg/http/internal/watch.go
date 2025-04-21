package internal

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/yxxchange/pipefree/helper/log"
	"github.com/yxxchange/pipefree/pkg/pipe/model"
	"github.com/yxxchange/pipefree/pkg/pipe/orca"
	"net/http"
)

type WatchServer struct {
	ctx   *gin.Context
	param WatchParam
}

func LaunchServer(c *gin.Context, param WatchParam) {
	server := &WatchServer{
		param: param,
		ctx:   c,
	}
	server.Serve()
}

func (s *WatchServer) Serve() {
	s.ServeHTTP(s.ctx.Writer, s.ctx.Request)
}

func (s *WatchServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f, ok := w.(http.Flusher)
	if !ok {
		log.Errorf("type asset failed: %T not implement http.Flusher", w)
		return
	}
	s.beginChunkedStream(w, f)
	log.Info("start to watch")
	eg := model.EngineGroup{
		Engine:    s.param.Engine,
		Namespace: s.param.Namespace,
		Kind:      s.param.Kind,
	}
	done := r.Context().Done()
	ch := orca.GetOrchestrator(context.Background()).Register(eg).Channel()
	for {
		select {
		case b, ok := <-ch:
			if !ok {
				return
			}
			_, err := w.Write(b)
			if err != nil {
				log.Errorf("write data error in watching, disconnected")
				return
			}
			if len(ch) == 0 {
				f.Flush()
			}
		case <-done:
			return
		}
	}
}

type WatchParam struct {
	Namespace string     `uri:"namespace" binding:"required"`
	Name      string     `uri:"name" binding:"required"`
	Kind      model.Kind `query:"kind" binding:"required"`
	Engine    string     `query:"engine" binding:"required"`
}

func (s *WatchServer) beginChunkedStream(w http.ResponseWriter, f http.Flusher) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.WriteHeader(http.StatusOK)
	f.Flush()

	return
}

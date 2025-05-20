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
	idf := model.Schema{
		ApiVersion: s.param.ApiVersion,
		Namespace:  s.param.Namespace,
		Kind:       s.param.Kind,
		Operation:  s.param.Operation,
	}
	done := r.Context().Done()
	ch := orca.GetOrchestrator(context.Background()).Register(idf).Channel()
	for {
		select {
		case b, ok := <-ch.Chan():
			if !ok {
				return
			}
			_, err := w.Write(b)
			if err != nil {
				log.Errorf("write data error in watching, disconnected")
				return
			}
			if len(ch.Chan()) == 0 {
				f.Flush()
			}
		case <-done:
			return
		}
	}
}

type WatchParam struct {
	ApiVersion string     `uri:"apiVersion"`
	Namespace  string     `uri:"namespace"`
	Name       string     `uri:"name"`
	Kind       model.Kind `query:"kind"`
	Operation  string     `query:"operation"`
}

func (s *WatchServer) beginChunkedStream(w http.ResponseWriter, f http.Flusher) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.WriteHeader(http.StatusOK)
	f.Flush()

	return
}

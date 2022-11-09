package gin

//import (
//	"fmt"
//	"net/http"
//	"os"

//	"github.com/gin-gonic/gin"

//	"github.com/realnighthawk/bucky/server"
//)

//const (
//	ApplicationRouteGroup = "application"
//	MonitoringRouteGroup  = "monitoring"
//)

//type Server struct {
//	Options
//	handlers  map[string]*gin.Engine
//	addresses server.Addresses
//}

//type Options struct {
//	Name        string           `json:"name,omitempty" yaml:"name,omitempty"`
//	Environment string           `json:"environment,omitempty" yaml:"environment,omitempty"`
//	Version     string           `json:"version,omitempty" yaml:"version,omitempty"`
//	Addresses   server.Addresses `json:"addresses,omitempty" yaml:"addresses,omitempty"`
//}

//func New(opts Options) (*Server, error) {
//	gin.SetMode(gin.ReleaseMode)
//	if os.Getenv("DEBUG") == "true" {
//		gin.SetMode(gin.DebugMode)
//	}

//	svc := &Server{
//		Options: opts,
//	}

//	svc.addresses = opts.Addresses
//	svc.handlers = make(map[string]*gin.Engine, len(opts.Addresses))
//	for key, _ := range svc.addresses {
//		svc.handlers[key] = gin.New()
//		if key == "application" {
//			svc.handlers[key].GET("/health", healthHandler)
//			svc.handlers[key].GET("/ping", pingHandler)
//		}
//	}

//	return svc, nil
//}

//func (h *Server) RegisterGeneric(method string, route string, handler http.Handler) {
//	for key, _ := range h.handlers {
//		h.handlers[key].Handle(method, route, gin.WrapH(handler))
//		break
//	}

//}

//func (h *Server) Register(method string, route string, handler gin.HandlerFunc) {
//	for key, _ := range h.handlers {
//		h.handlers[key].Handle(method, route, handler)
//		break
//	}

//}

//func (h *Server) RegisterGenericWithGroup(name string, method string, route string, handler http.Handler) {
//	h.handlers[name].Group(name).Handle(method, route, gin.WrapH(handler))
//}

//func (h *Server) RegisterWithGroup(name string, method string, route string, handler gin.HandlerFunc) {
//	h.handlers[name].Group(name).Handle(method, route, handler)
//}

//func (h *Server) Run(runCh chan error) {
//	for name, handler := range h.handlers {
//		go func(hd *gin.Engine, host, port string) {
//			err := hd.Run(fmt.Sprintf("%s:%s", host, port))
//			if err != nil {
//				runCh <- err
//			}
//		}(handler, h.addresses[name].Host, h.addresses[name].Port)
//	}
//}

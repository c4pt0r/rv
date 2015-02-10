package rv

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	log "github.com/ngaut/logging"
)

var verbose bool = false

func SetVerbose(b bool) {
	verbose = b
}

type Matcher interface {
	match(r *http.Request) http.Handler
	reload(param interface{})
}

type Server struct {
	matcher    Matcher
	configFile string
	addr       string
}

func NewServer(configFile string) *Server {
	cfg, err := loadConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}
	if verbose {
		log.Info("config file:", cfg)
	}

	ret := &Server{
		matcher:    newDefaultMatcher(cfg),
		configFile: configFile,
		addr:       cfg.Addr,
	}
	ret.regReloadSignalHandler()
	return ret
}

func (s *Server) regReloadSignalHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGUSR1)
	go func() {
		for {
			<-c
			cfg, err := loadConfig(s.configFile)
			if err != nil {
				log.Error(err)
				continue
			}
			if verbose {
				log.Info("reload config", cfg)
			}
			s.matcher.reload(cfg)
		}
	}()
}

func (s *Server) onRequest(w http.ResponseWriter, r *http.Request) {
	hdlr := s.matcher.match(r)
	if hdlr != nil {
		hdlr.ServeHTTP(w, r)
		return
	}
	http.NotFound(w, r)
}

func (s *Server) Serve() {
	http.HandleFunc("/", s.onRequest)
	err := http.ListenAndServe(s.addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

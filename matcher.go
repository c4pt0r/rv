package rv

import (
	"net/http"
	"strings"
	"sync"

	log "github.com/ngaut/logging"
)

type defaultMatcher struct {
	vhosts []*vhostInfo
	lck    sync.RWMutex
}

func loadHostInfo(cfg Config) []*vhostInfo {
	// create vhost detail from config
	var vhosts []*vhostInfo
	for _, h := range cfg.VHost {
		v := newVhostInfo(h)
		if v != nil {
			vhosts = append(vhosts, v)
		}
	}
	return vhosts
}

func newDefaultMatcher(cfg Config) *defaultMatcher {
	return &defaultMatcher{
		vhosts: loadHostInfo(cfg),
	}
}

func (m *defaultMatcher) reload(cfg interface{}) {
	m.lck.Lock()
	defer m.lck.Unlock()
	m.vhosts = loadHostInfo(cfg.(Config))
}

func (m *defaultMatcher) match(r *http.Request) http.Handler {
	hostName := strings.Split(r.Host, ":")[0]
	log.Info(r.Host, r.URL.RawQuery)

	m.lck.RLock()
	vhosts := m.vhosts
	m.lck.RUnlock()
	for _, vhostInfo := range vhosts {
		if vhostInfo.regexp.MatchString(hostName) {
			return vhostInfo.nextHandler()
		}
	}
	return nil
}

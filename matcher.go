package rv

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
	"sync"

	log "github.com/ngaut/logging"
)

type vhostInfo struct {
	VHost
	regexp *regexp.Regexp
	hdlrs  []http.Handler
}

func newVhostInfo(v VHost) *vhostInfo {
	regexp, err := regexp.Compile(v.hostNamePattern)
	if err != nil {
		log.Error(err)
		return nil
	}

	var hdlrs []http.Handler
	if len(v.Static) > 0 {
		hdlrs = append(hdlrs, http.FileServer(http.Dir(v.Static)))
	} else if len(v.Upstreams) > 0 {
		for _, addr := range v.Upstreams {
			url, err := url.Parse(addr)
			if err != nil {
				log.Error(err)
				continue
			}
			hdlr := httputil.NewSingleHostReverseProxy(url)
			hdlrs = append(hdlrs, hdlr)
		}
	}

	return &vhostInfo{
		VHost:  v,
		regexp: regexp,
		hdlrs:  hdlrs,
	}
}

type defaultMatcher struct {
	vhosts []*vhostInfo
	lck    sync.RWMutex
}

func loadHostInfo(cfg Config) []*vhostInfo {
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

	m.lck.RLock()
	vhosts := m.vhosts
	m.lck.RUnlock()
	for _, v := range vhosts {
		if v.regexp.MatchString(hostName) {
			return v.hdlrs[0]
		}
	}
	return nil
}

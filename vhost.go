package rv

import (
	"container/ring"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"sync"

	log "github.com/ngaut/logging"
)

type vhostInfo struct {
	VHost
	regexp *regexp.Regexp
	hdlrs  *ring.Ring
	lck    sync.RWMutex
}

func slice2ring(values []http.Handler) *ring.Ring {
	if values == nil {
		return nil
	}
	r := ring.New(len(values))
	i, n := 0, r.Len()
	for p := r; i < n; p = p.Next() {
		p.Value = values[i]
		i++
	}
	return r
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

	slice2ring(hdlrs)
	return &vhostInfo{
		VHost:  v,
		regexp: regexp,
		hdlrs:  slice2ring(hdlrs),
	}
}

func (v *vhostInfo) nextHandler() http.Handler {
	if v.hdlrs.Len() == 0 {
		return nil
	}

	ret, _ := v.hdlrs.Value.(http.Handler)
	v.hdlrs = v.hdlrs.Next()
	return ret
}

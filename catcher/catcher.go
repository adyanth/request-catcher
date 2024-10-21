package catcher

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/op/go-logging"
)

type Catcher struct {
	config   *Configuration
	router   *mux.Router
	upgrader websocket.Upgrader
	logger   *logging.Logger

	hostsMu sync.Mutex
	hosts   map[string]*Host

	stats struct {
		processStart       time.Time
		requestsIndex      atomic.Uint64
		requestsCaught     atomic.Uint64
		requestsIgnored    atomic.Uint64
		requestsClientInit atomic.Uint64
	}
}

func NewCatcher(config *Configuration) *Catcher {
	c := &Catcher{
		config: config,
		router: mux.NewRouter(),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		logger: logging.MustGetLogger("request-catcher"),

		hosts: make(map[string]*Host),
	}
	c.stats.processStart = time.Now()
	c.router.HandleFunc("/", c.rootHandler).Host(c.config.RootHost)
	c.router.HandleFunc("/", c.indexHandler)
	c.router.HandleFunc("/init-client", c.initClient)
	c.router.PathPrefix("/assets").Handler(http.StripPrefix("/assets",
		withCacheHeaders(http.FileServer(http.Dir(config.FrontendDir)))))
	c.router.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, config.Favicon)
	})
	c.router.HandleFunc("/statusz", c.status)
	c.router.NotFoundHandler = http.HandlerFunc(c.catchRequests)
	return c
}

func withCacheHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		oneYear := time.Now().Add(time.Hour * 24 * 365).Format(time.RFC3339)
		w.Header().Set("Cache-Control", "public, max-age=31536000")
		w.Header().Set("Expires", oneYear)
		h.ServeHTTP(w, r)
	})
}

func (c *Catcher) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if strings.HasPrefix(req.Host, "www.") {
		rw.Header().Set("Connection", "close")
		url := strings.TrimPrefix(req.Host, "www.") + req.URL.String()
		http.Redirect(rw, req, url, http.StatusTemporaryRedirect)
		return
	}

	c.router.ServeHTTP(rw, req)
}

func (c *Catcher) host(hostString string) *Host {
	hostString = hostWithoutPort(hostString)

	c.hostsMu.Lock()
	defer c.hostsMu.Unlock()
	if host, ok := c.hosts[hostString]; ok {
		return host
	}
	host := newHost(hostString)
	c.hosts[hostString] = host
	return host
}

func (c *Catcher) rootHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, c.config.FrontendDir+"/root.html")
}

func (c *Catcher) indexHandler(w http.ResponseWriter, r *http.Request) {
	// Some people mistakenly expect requests to the index of the subdomain
	// to be caught. For now, just catch those as well. Later I should move
	// the index to be hosted at requestcatcher.com.
	c.Catch(r)
	c.stats.requestsIndex.Add(1)
	http.ServeFile(w, r, c.config.FrontendDir+"/index.html")
}

func (c *Catcher) catchRequests(w http.ResponseWriter, r *http.Request) {
	if c.Catch(r) {
		c.stats.requestsCaught.Add(1)
		fmt.Fprintf(w, "request caught")
		return
	}
	c.stats.requestsIgnored.Add(1)
	// No one is listening to requests to this subdomain.
	if c.config.RedirectDest != "" {
		http.Redirect(w, r, c.config.RedirectDest, http.StatusSeeOther)
		return
	}
	fmt.Fprintf(w, "request ignored")
}

func (c *Catcher) initClient(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	c.stats.requestsClientInit.Add(1)

	ws, err := c.upgrader.Upgrade(w, r, nil)
	if err != nil {
		c.logger.Error(err.Error())
		return
	}

	clientHost := c.host(r.Host)
	c.logger.Infof("Initializing a new client on host %v", clientHost.Host)
	clientHost.clients.Store(c, newClient(c, clientHost, ws))
}

func (c *Catcher) Catch(r *http.Request) (caught bool) {
	hostString := hostWithoutPort(r.Host)
	c.hostsMu.Lock()
	host, ok := c.hosts[hostString]
	c.hostsMu.Unlock()

	if !ok {
		// No one is listening, so no reason to catch it.
		return false
	}

	// Broadcast it to everyone listening for requests on this host
	caughtRequest := convertRequest(r)
	host.broadcast <- caughtRequest
	return true
}

func (c *Catcher) status(w http.ResponseWriter, r *http.Request) {
	c.hostsMu.Lock()
	countHosts := len(c.hosts)
	c.hostsMu.Unlock()
	fmt.Fprintf(w, "uptime: %d\n", int(time.Since(c.stats.processStart).Seconds()))
	fmt.Fprintf(w, "hosts: %d\n", countHosts)
	fmt.Fprintf(w, "index: %d\n", c.stats.requestsIndex.Load())
	fmt.Fprintf(w, "caught: %d\n", c.stats.requestsCaught.Load())
	fmt.Fprintf(w, "ignored: %d\n", c.stats.requestsIgnored.Load())
	fmt.Fprintf(w, "client-init: %d\n", c.stats.requestsClientInit.Load())
}

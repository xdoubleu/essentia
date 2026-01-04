package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/xdoubleu/essentia/v2/internal/helpers"
	"github.com/xdoubleu/essentia/v2/pkg/communication/httptools"
	"golang.org/x/time/rate"
)

var cleanerActive bool                 //nolint: gochecknoglobals //need this
var mu sync.RWMutex                    //nolint: gochecknoglobals //need this
var clients = make(map[string]*client) //nolint: gochecknoglobals //need this

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimit is middleware used to rate limit requests by clients identified by IP.
func RateLimit(
	rps rate.Limit,
	bucketSize int,
	cleanupTimer time.Duration,
	removeAfter time.Duration,
) helpers.Middleware {
	if !cleanerActive {
		cleanerActive = true
		go func() {
			for {
				time.Sleep(cleanupTimer)

				mu.RLock()

				for ip, client := range clients {
					if time.Since(client.lastSeen) > removeAfter {
						mu.RUnlock()
						mu.Lock()

						delete(clients, ip)

						mu.Unlock()
						mu.RLock()
					}
				}

				mu.RUnlock()
			}
		}()
	}

	return func(next http.Handler) http.Handler {
		return rateLimit(&mu, clients, rps, bucketSize, next)
	}
}

func rateLimit(
	mu *sync.RWMutex,
	clients map[string]*client,
	rps rate.Limit,
	bucketSize int,
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			httptools.ServerErrorResponse(w, r, err)
			return
		}

		mu.Lock()

		if _, found := clients[ip]; !found {
			//nolint:exhaustruct //lastSeen is set later
			clients[ip] = &client{
				limiter: rate.NewLimiter(rps, bucketSize),
			}
		}

		clients[ip].lastSeen = time.Now()

		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			httptools.RateLimitExceededResponse(w, r)
			return
		}

		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}

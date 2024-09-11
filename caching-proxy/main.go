package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

type CachedResponse []byte

type Caching struct {
	n         int
	origin    string
	responses map[string]CachedResponse
}

func NewCaching() *Caching {
	return &Caching{
		responses: make(map[string]CachedResponse),
	}
}

func (c *Caching) newPort(n int) {
	c.n = n
}

func (c *Caching) newOrigin(u string) {
	c.origin = u
}

func (c *Caching) cacheResponse(k string, v CachedResponse) {
	c.responses[k] = v
}

func (c *Caching) getResponse(k string) CachedResponse {
	return c.responses[k]
}

func (c *Caching) clearCache() {
	c.responses = make(map[string]CachedResponse)
}

type ContextKey struct {
	name string
}

func (c *Caching) NewContextKey(r *http.Request) ContextKey {
	return ContextKey{name: c.buildFullOriginUrl(r)}
}

func (c *Caching) buildFullOriginUrl(r *http.Request) string {
	return c.origin + r.URL.Path + "?" + r.URL.RawQuery
}

func checkCacheMiddleware(caching *Caching, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remoteUrl := caching.buildFullOriginUrl(r)

		cachedResponse := caching.getResponse(remoteUrl)
		if cachedResponse != nil {
			w.Header().Add("X-Cache", "HIT")
			w.Write(cachedResponse)
			return
		}

		resp, err := http.Get(remoteUrl)
		if err != nil {
			http.Error(w, "Error fetching data from origin", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		parsedResp, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Error reading data from origin", http.StatusInternalServerError)
			return
		}

		caching.cacheResponse(remoteUrl, parsedResp)

		contextKey := caching.NewContextKey(r)
		ctx := context.WithValue(r.Context(), contextKey, parsedResp)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func main() {
	caching := NewCaching()

	reader := bufio.NewReader(os.Stdin)

	go func() {
		fmt.Print("To clear the cache type: --clear-cache\n")
		for {
			text, _ := reader.ReadString('\n')
			if text == "--clear-cache\n" {
				caching.clearCache()
				fmt.Println("Cache cleared")
			} else {
				fmt.Println("Invalid command")
			}
		}
	}()

	for i, arg := range os.Args {
		if i == len(os.Args)-1 {
			break
		}

		argVal := os.Args[i+1]
		if arg == "--port" {
			n, err := strconv.Atoi(argVal)
			if err != nil {
				fmt.Println("Port must be a number")
				os.Exit(1)
			}

			if n < 1 || n > 65535 {
				fmt.Println("Port must be between 1 and 65535")
				os.Exit(1)
			}

			caching.newPort(n)
			continue
		}

		if arg == "--origin" {
			_, err := url.ParseRequestURI(argVal)
			if err != nil {
				fmt.Println("Origin must be a valid URL")
				os.Exit(1)
			}
			caching.newOrigin(argVal)
			continue
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		contextKey := caching.NewContextKey(r)
		resp := r.Context().Value(contextKey).([]byte)
		w.Header().Add("X-Cache", "MISS")
		w.Write(resp)
	})

	println("Server running on port: ", caching.n)

	err := http.ListenAndServe(":"+strconv.Itoa(caching.n), checkCacheMiddleware(caching, mux))
	if err != nil {
		fmt.Println("Error starting server: ", err)
		os.Exit(1)
	}

	// select {}
}

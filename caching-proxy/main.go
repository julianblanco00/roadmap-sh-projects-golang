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
	"strings"
	"time"
)

type CachedResponse struct {
	Data     []byte
	LastRead time.Time
}

type Caching struct {
	Port      int
	Origin    string
	Responses map[string]*CachedResponse
}

func NewCaching() *Caching {
	return &Caching{}
}

func (c *Caching) NewPort(port int) {
	c.Port = port
}

func (c *Caching) NewOrigin(o string) {
	c.Origin = o
}

func (c *Caching) CacheResponse(k string, v []byte) {
	if c.Responses == nil {
		c.Responses = make(map[string]*CachedResponse)
	}

	c.Responses[k] = &CachedResponse{
		Data:     v,
		LastRead: time.Now(),
	}
}

func (c *Caching) UpdateLastRead(k string) {
	c.Responses[k].LastRead = time.Now()
}

func (c *Caching) GetResponse(k string) *CachedResponse {
	return c.Responses[k]
}

func (c *Caching) ClearCache() {
	c.Responses = make(map[string]*CachedResponse)
}

type ContextKey struct {
	Name string
}

func (c *Caching) NewContextKey(r *http.Request) ContextKey {
	return ContextKey{Name: c.BuildFullOriginUrl(r)}
}

func (c *Caching) BuildFullOriginUrl(r *http.Request) string {
	return c.Origin + r.URL.Path + "?" + r.URL.RawQuery
}

func checkCacheMiddleware(caching *Caching, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remoteUrl := caching.BuildFullOriginUrl(r)

		cachedResponse := caching.GetResponse(remoteUrl)

		if cachedResponse != nil && cachedResponse.Data != nil {
			caching.UpdateLastRead(remoteUrl)

			w.Header().Add("X-Cache", "HIT")
			w.Write(cachedResponse.Data)
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

		caching.CacheResponse(remoteUrl, parsedResp)

		contextKey := caching.NewContextKey(r)
		ctx := context.WithValue(r.Context(), contextKey, parsedResp)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func listenCommands(caching *Caching) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("To clear the cache type: --clear-cache\n")
	for {
		text, _ := reader.ReadString('\n')
		if text == "--clear-cache\n" {
			caching.ClearCache()
			fmt.Println("Cache cleared")
		} else if len(strings.TrimSpace(text)) == 0 {
			continue
		} else if text == "exit\n" {
			os.Exit(0)
		}
	}
}

func clearOldCache(caching *Caching) {
	for {
		time.Sleep(30 * time.Second)

		if len(caching.Responses) == 0 {
			continue
		}

		for k, v := range caching.Responses {
			if time.Since(v.LastRead) > 10*time.Minute {
				delete(caching.Responses, k)
				fmt.Println("Cache cleared for: ", k)
			}
		}
	}
}

func readArguments(caching *Caching) {
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

			caching.NewPort(n)
			continue
		}

		if arg == "--origin" {
			_, err := url.ParseRequestURI(argVal)
			if err != nil {
				fmt.Println("Origin must be a valid URL")
				os.Exit(1)
			}
			caching.NewOrigin(argVal)
			continue
		}
	}

}

func startCachingServer(caching *Caching) {
	mux := http.NewServeMux()
	mux.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		contextKey := caching.NewContextKey(r)
		resp := r.Context().Value(contextKey).([]byte)
		w.Header().Add("X-Cache", "MISS")
		w.Write(resp)
	})

	println("Server running on port: ", caching.Port)

	err := http.ListenAndServe(":"+strconv.Itoa(caching.Port), checkCacheMiddleware(caching, mux))
	if err != nil {
		fmt.Println("Error starting server: ", err)
		os.Exit(1)
	}
}

func main() {
	caching := NewCaching()

	go listenCommands(caching)
	go clearOldCache(caching)
	readArguments(caching)
	startCachingServer(caching)
}

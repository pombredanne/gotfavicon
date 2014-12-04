package main

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/subosito/shorticon"
	"github.com/tsileo/disklru"
)

func ipAddrFromRemoteAddr(s string) string {
	idx := strings.LastIndex(s, ":")
	if idx == -1 {
		return s
	}
	return s[:idx]
}

func getIpAddress(r *http.Request) string {
	hdr := r.Header
	hdrRealIp := hdr.Get("X-Real-Ip")
	hdrForwardedFor := hdr.Get("X-Forwarded-For")
	if hdrRealIp == "" && hdrForwardedFor == "" {
		return ipAddrFromRemoteAddr(r.RemoteAddr)
	}
	if hdrForwardedFor != "" {
		log.Println(hdrForwardedFor)
		// X-Forwarded-For is potentially a list of addresses separated with ","
		parts := strings.Split(hdrForwardedFor, ",")
		for i, p := range parts {
			parts[i] = strings.TrimSpace(p)
		}
		return parts[len(parts)-1]
	}
	return hdrRealIp
}

var client = &http.Client{}

var lru *disklru.DiskLRU

func LRUFunc(key string, url interface{}) ([]byte, error) {
	log.Printf("url: %v", url.(string))
	ux, err := shorticon.NewScraper(url.(string)).Favicon()
	if err != nil {
		return nil, err
	}
	log.Printf("f: %v", ux.String())
	resp, err := client.Get(ux.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func init() {
	var err error
	lru, err = disklru.New("./favicons_cache", LRUFunc, 100000000) // 100MB
	if err != nil {
		panic(err)
	}
}

func respHandler(res http.ResponseWriter, req *http.Request) {
	url := req.URL.Query().Get("url")
	log.Printf("req [url=%v] from %v", url, getIpAddress(req))
	if url == "" {
		return
	}
	hurl := fmt.Sprintf("%x", sha1.Sum([]byte(url)))
	fav, fetched, err := lru.Get(hurl, url)
	if err != nil {
		panic(err)
	}
	log.Printf("[favicon_url=%v, fetched=%v]", hurl, fetched)
	res.Header().Set("Content-Type", "image/x-icon")
	res.Write(fav)
}

func main() {
	http.HandleFunc("/", respHandler)
	log.Fatal(http.ListenAndServe(":8008", nil))
}

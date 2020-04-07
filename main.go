package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	sourceKey = "source"
	targetKey = "target"
)

type httpHandler struct {
	db db
}

func main() {
	addr := env("HTTP_ADDR", ":5545")
	dbPath := env("DB_PATH", "db")

	log.Printf("start HTTP_ADDR:%q DB_PATH:%q", addr, dbPath)
	db, err := newDB(dbPath)
	if err != nil {
	    log.Fatalf("failed to open DB: %v", err)
	}
	defer db.close()

	handle := httpHandler{db: db}

	http.HandleFunc("/", handle.handleResolve)
	http.HandleFunc("/add", handle.handleAdd)
	http.HandleFunc("/list", handle.handleList)

	log.Printf("start listening on %q", addr)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("listen error: %v", err)
	}
}

func (h httpHandler) handleResolve(w http.ResponseWriter, r *http.Request) {
	scheme := r.URL.Scheme
	if scheme == "" {
		scheme = "http"
	}

	path := strings.TrimPrefix(r.URL.Path, "/")
	source := fmt.Sprintf("%s://%s/%s", scheme, r.Host, path)
	target, err := h.db.get(source)
	if err != nil  {
		if err == ErrNotFound {
			log.Printf("%q not found", source)
			http.NotFound(w, r)
			return
		}
		status := http.StatusInternalServerError
		http.Error(w, http.StatusText(status), status)
	}

	if r.Method == http.MethodDelete {
		err := h.db.delete(source)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("deleted source:%q", source)
		_, _ = w.Write([]byte("deleted ok"))
		return
	}

	log.Printf("redirecting %q -> %q", source, target)
	http.Redirect(w, r, target, http.StatusTemporaryRedirect)
}

func (h httpHandler) handleAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.handleResolve(w, r,)
		return
	}

	source := r.FormValue(sourceKey)
	if source == "" || !validURL(source) {
		http.Error(w, "invalid source", http.StatusBadRequest)
		return
	}

	target := r.FormValue(targetKey)
	if target == "" || !validURL(target) {
		http.Error(w, "invalid target", http.StatusBadRequest)
		return
	}

	err := h.db.add(source, target)
	if err != nil {
	    http.Error(w, err.Error(), http.StatusInternalServerError)
	    return
	}

	log.Printf("added source:%q target:%q", source, target)
	_, _ = w.Write([]byte("ok"))
}

func (h httpHandler) handleList(w http.ResponseWriter, r *http.Request) {
	redirects, err := h.db.list()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = fmt.Fprintf(w, "%d entr(ies)\n", len(redirects))
	i := 0
	for source, target := range redirects {
		i++
		_, _ = fmt.Fprintf(w, "[%3d] %s -> %s\n", i, source, target)
	}
}

func env(key string, def string) string {
	val := os.Getenv(key)
	if val == "" {
		val = def
	}
	return val
}

func validURL(s string) bool {
	u, err := url.ParseRequestURI(s)
	return err == nil && u.Scheme != "" && u.Host != ""
}

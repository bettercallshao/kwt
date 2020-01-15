package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	for _, url := range []string{
		"https://cdnjs.cloudflare.com/ajax/libs/highlight.js/9.15.8/styles/default.min.css",
		"https://cdn.jsdelivr.net/npm/vue@2.6.11/dist/vue.js",
		"https://code.jquery.com/jquery-3.3.1.slim.min.js",
		"https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css",
		"https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css.map",
		"https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/js/bootstrap.min.js",
		"https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/js/bootstrap.min.js.map",
		"https://unpkg.com/axios/dist/axios.min.js",
		"https://unpkg.com/axios/dist/axios.min.map",
		"https://unpkg.com/vue-router@3.1.3/dist/vue-router.js",
	} {
		dl(url)
	}
}

func check(err error) {
	if err != nil {
		os.Exit(1)
	}
}

func dl(url string) {
	fn := strings.Replace(url, "https://", "cmd/kutd/assets/third/", -1)

	check(os.MkdirAll(filepath.Dir(fn), os.ModePerm))
	out, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE, 0644)
	check(err)
	defer out.Close()

	resp, err := http.Get(url)
	check(err)
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	check(err)
}

package downloader

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDownloader(t *testing.T) {
	ts := startTestServer()
	defer ts.Close()

	urls := []string{
		ts.URL + "/file1",
		ts.URL + "/file2",
		ts.URL + "/file3",
	}
	numWorkers := 2
	timeout := 10 * time.Second
	downloadDir := "downloads"

	Downloader(urls, numWorkers, timeout, downloadDir)

	// Проверяем, что файлы были созданы
	for _, url := range urls {
		filename := filepath.Join(downloadDir, filepath.Base(url))
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			t.Errorf("Файл %s не существует", filename)
		}
	}
}

// временный HTTP сервер для тестирования
func startTestServer() *httptest.Server {
	handler := http.NewServeMux()

	handler.HandleFunc("/file1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Файл 1")
	})

	handler.HandleFunc("/file2", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Файл 2")
	})

	handler.HandleFunc("/file3", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Файл 3")
	})

	return httptest.NewServer(handler)
}

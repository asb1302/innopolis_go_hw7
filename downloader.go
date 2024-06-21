package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func Downloader(urls []string, numWorkers int, timeout time.Duration, downloadDir string) {
	if err := os.MkdirAll(downloadDir, os.ModePerm); err != nil {
		fmt.Printf("Ошибка при создании папки для хранения скачанных файлов: %v\n", err)

		return
	}

	urlChan := make(chan string, len(urls))
	results := make(chan string, len(urls))

	var wg sync.WaitGroup

	// Запуск воркеров (Worker Pool)
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)

		go Worker(i, urlChan, results, &wg, timeout, downloadDir)
	}

	// Отправка URL в канал urlChan (Fan-out)
	for _, url := range urls {
		urlChan <- url
	}
	close(urlChan)

	// Ожидание завершения всех воркеров
	wg.Wait()
	close(results)

	// Печать результатов (Fan-in)
	for result := range results {
		fmt.Println(result)
	}
}

func Worker(id int, urls <-chan string, results chan<- string, wg *sync.WaitGroup, timeout time.Duration, downloadDir string) {
	defer wg.Done()
	for url := range urls {
		resultChan := make(chan string)
		go func(url string) {
			filename, err := DownloadFile(url, downloadDir, timeout)

			if err != nil {
				resultChan <- fmt.Sprintf("Воркер %d: ошибка при скачивани из URL %s: %v", id, url, err)
			} else {
				resultChan <- fmt.Sprintf("Воркер %d: скачивание успешно выполнено из URL %s в файл %s", id, url, filename)
			}
		}(url)

		// Используем select для обработки (timeout)
		select {
		case res := <-resultChan:
			results <- res
		case <-time.After(timeout):
			results <- fmt.Sprintf("Воркер %d: timeout при скачивании из URL %s", id, url)
		}
	}
}

func DownloadFile(url, downloadDir string, timeout time.Duration) (string, error) {
	client := http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	filename := filepath.Join(downloadDir, filepath.Base(url))
	out, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	return filename, nil
}

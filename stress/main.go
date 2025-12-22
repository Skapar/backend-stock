package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

var (
	baseURL     = flag.String("url", "http://localhost:8080/api", "API base url")
	email       = flag.String("email", "test@mail.com", "login email")
	password    = flag.String("password", "123456", "login password")
	endpoint    = flag.String("endpoint", "/stocks/", "endpoint to test")
	method      = flag.String("method", "GET", "HTTP method")
	concurrency = flag.Int("c", 1000, "concurrency level")
	requests    = flag.Int("n", 10000, "total requests")
)

type LoginResponse struct {
	Token string `json:"token"`
}

func login() (string, error) {
	body := map[string]string{
		"email":    *email,
		"password": *password,
	}

	b, _ := json.Marshal(body)

	req, _ := http.NewRequest(
		"POST",
		*baseURL+"/login",
		bytes.NewReader(b),
	)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		data, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("login failed: %s", data)
	}

	var lr LoginResponse
	json.NewDecoder(resp.Body).Decode(&lr)

	return lr.Token, nil
}

func worker(
	wg *sync.WaitGroup,
	token string,
	jobs <-chan struct{},
	success *int64,
	fail *int64,
	totalLatency *int64,
) {
	defer wg.Done()

	client := &http.Client{}

	for range jobs {
		start := time.Now()

		req, _ := http.NewRequest(
			*method,
			*baseURL+*endpoint,
			nil,
		)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := client.Do(req)
		latency := time.Since(start).Microseconds()

		atomic.AddInt64(totalLatency, latency)

		if err != nil {
			atomic.AddInt64(fail, 1)
			continue
		}

		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			atomic.AddInt64(success, 1)
		} else {
			atomic.AddInt64(fail, 1)
		}
	}
}

func main() {
	flag.Parse()

	fmt.Println("Logging in...")
	token, err := login()
	if err != nil {
		panic(err)
	}

	jobs := make(chan struct{}, *requests)

	var success int64
	var fail int64
	var totalLatency int64

	start := time.Now()

	var wg sync.WaitGroup
	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go worker(&wg, token, jobs, &success, &fail, &totalLatency)
	}

	for i := 0; i < *requests; i++ {
		jobs <- struct{}{}
	}
	close(jobs)

	wg.Wait()
	elapsed := time.Since(start)

	fmt.Println("------ RESULT ------")
	fmt.Println("Endpoint:", *endpoint)
	fmt.Println("Requests:", *requests)
	fmt.Println("Concurrency:", *concurrency)
	fmt.Println("Success:", success)
	fmt.Println("Fail:", fail)
	fmt.Println("Elapsed:", elapsed)
	fmt.Printf("RPS: %.2f\n", float64(*requests)/elapsed.Seconds())

	if success > 0 {
		fmt.Printf(
			"Avg latency: %.2f ms\n",
			float64(totalLatency)/float64(success)/1000,
		)
	}
}

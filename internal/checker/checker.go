// Package checker faz a verificação de saúde (health check) de serviços HTTP.
package checker

import (
	"net/http"
	"sync"
	"time"
)

// Service representa um serviço a ser monitorado.
type Service struct {
	Name           string `json:"name"`
	URL            string `json:"url"`
	ExpectedStatus int    `json:"expected_status,omitempty"`
}

// Result é o resultado de uma verificação.
type Result struct {
	Service    Service       `json:"service"`
	Up         bool          `json:"up"`
	StatusCode int           `json:"status_code"`
	Latency    time.Duration `json:"latency_ms"`
	Err        string        `json:"error,omitempty"`
}

// Check verifica um único serviço e devolve o resultado.
func Check(client *http.Client, svc Service) Result {
	expected := svc.ExpectedStatus
	if expected == 0 {
		expected = http.StatusOK // 200 por padrão
	}

	start := time.Now()
	resp, err := client.Get(svc.URL)
	latency := time.Since(start)

	if err != nil {
		return Result{Service: svc, Up: false, Latency: latency, Err: err.Error()}
	}
	defer resp.Body.Close()

	return Result{
		Service:    svc,
		Up:         resp.StatusCode == expected,
		StatusCode: resp.StatusCode,
		Latency:    latency,
	}
}

// CheckAll verifica todos os serviços EM PARALELO usando goroutines.
// A ordem dos resultados é preservada em relação à lista de entrada.
func CheckAll(services []Service, timeout time.Duration) []Result {
	results := make([]Result, len(services))
	client := &http.Client{Timeout: timeout}

	var wg sync.WaitGroup
	for i, svc := range services {
		wg.Add(1)
		// Cada serviço é checado em sua própria goroutine.
		go func(idx int, s Service) {
			defer wg.Done()
			results[idx] = Check(client, s)
		}(i, svc)
	}
	wg.Wait() // espera todas as goroutines terminarem

	return results
}

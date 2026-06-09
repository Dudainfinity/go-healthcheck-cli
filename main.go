// Comando healthcheck: CLI que verifica a saúde de vários serviços HTTP
// em paralelo (goroutines) e reporta o status.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/Dudainfinity/go-healthcheck-cli/internal/checker"
)

func main() {
	configPath := flag.String("config", "", "Arquivo JSON com a lista de serviços")
	timeout := flag.Duration("timeout", 5*time.Second, "Tempo máximo de espera por serviço")
	interval := flag.Duration("interval", 0, "Se > 0, roda em loop nesse intervalo (modo watch). Ex: 10s")
	asJSON := flag.Bool("json", false, "Imprime o resultado em JSON")
	flag.Usage = usage
	flag.Parse()

	services, err := loadServices(*configPath, flag.Args())
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(2)
	}
	if len(services) == 0 {
		usage()
		os.Exit(2)
	}

	// Modo watch: roda indefinidamente no intervalo informado.
	if *interval > 0 {
		for {
			run(services, *timeout, *asJSON)
			time.Sleep(*interval)
		}
	}

	// Modo único: roda uma vez e sai com código != 0 se algo estiver fora.
	allUp := run(services, *timeout, *asJSON)
	if !allUp {
		os.Exit(1) // permite falhar um pipeline de CI quando há serviço fora
	}
}

// loadServices carrega os serviços do arquivo de config OU dos argumentos (URLs).
func loadServices(configPath string, args []string) ([]checker.Service, error) {
	if configPath != "" {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("lendo config: %w", err)
		}
		var cfg struct {
			Services []checker.Service `json:"services"`
		}
		if err := json.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("parseando config JSON: %w", err)
		}
		return cfg.Services, nil
	}

	// Sem config: cada argumento é tratado como uma URL.
	services := make([]checker.Service, 0, len(args))
	for _, url := range args {
		services = append(services, checker.Service{Name: url, URL: url})
	}
	return services, nil
}

// run executa uma rodada de checagem e imprime o resultado. Retorna true se
// todos os serviços estiverem no ar.
func run(services []checker.Service, timeout time.Duration, asJSON bool) bool {
	results := checker.CheckAll(services, timeout)

	if asJSON {
		out, _ := json.MarshalIndent(results, "", "  ")
		fmt.Println(string(out))
	} else {
		printTable(results)
	}

	for _, r := range results {
		if !r.Up {
			return false
		}
	}
	return true
}

func printTable(results []checker.Result) {
	w := tabwriter.NewWriter(os.Stdout, 0, 2, 3, ' ', 0)
	fmt.Fprintln(w, "SERVIÇO\tSTATUS\tHTTP\tLATÊNCIA\tDETALHE")
	upCount := 0
	for _, r := range results {
		status := "DOWN"
		detail := r.Err
		if r.Up {
			status = "UP"
			upCount++
		}
		httpCode := "-"
		if r.StatusCode > 0 {
			httpCode = fmt.Sprintf("%d", r.StatusCode)
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%dms\t%s\n",
			r.Service.Name, status, httpCode, r.Latency.Milliseconds(), detail)
	}
	w.Flush()
	fmt.Printf("\n%d/%d no ar\n", upCount, len(results))
}

func usage() {
	fmt.Fprintln(os.Stderr, "healthcheck — verifica a saúde de serviços HTTP em paralelo")
	fmt.Fprintln(os.Stderr, "\nUso:")
	fmt.Fprintln(os.Stderr, "  healthcheck -config services.json")
	fmt.Fprintln(os.Stderr, "  healthcheck https://google.com https://github.com")
	fmt.Fprintln(os.Stderr, "  healthcheck -config services.json -interval 10s   # modo watch")
	fmt.Fprintln(os.Stderr, "\nOpções:")
	flag.PrintDefaults()
}

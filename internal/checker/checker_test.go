package checker

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCheck_ServicoNoAr(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client := &http.Client{Timeout: 2 * time.Second}
	res := Check(client, Service{Name: "teste", URL: srv.URL})

	if !res.Up {
		t.Fatalf("esperava Up=true, recebi false (erro=%s)", res.Err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("esperava status 200, recebi %d", res.StatusCode)
	}
}

func TestCheck_ServicoForaDoAr(t *testing.T) {
	client := &http.Client{Timeout: 1 * time.Second}
	// Porta 0 em localhost não aceita conexão: simula serviço fora do ar.
	res := Check(client, Service{Name: "offline", URL: "http://127.0.0.1:0"})

	if res.Up {
		t.Fatal("esperava Up=false para um serviço fora do ar")
	}
	if res.Err == "" {
		t.Fatal("esperava uma mensagem de erro preenchida")
	}
}

func TestCheck_StatusInesperado(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError) // 500
	}))
	defer srv.Close()

	client := &http.Client{Timeout: 2 * time.Second}
	// Esperamos 200, mas o serviço devolve 500 -> deve ficar Up=false.
	res := Check(client, Service{Name: "erro", URL: srv.URL, ExpectedStatus: 200})

	if res.Up {
		t.Fatal("esperava Up=false quando o status difere do esperado")
	}
}

func TestCheckAll_Paralelo(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	services := []Service{
		{Name: "a", URL: srv.URL},
		{Name: "b", URL: srv.URL},
		{Name: "c", URL: srv.URL},
	}

	results := CheckAll(services, 2*time.Second)
	if len(results) != len(services) {
		t.Fatalf("esperava %d resultados, recebi %d", len(services), len(results))
	}
	// A ordem deve ser preservada.
	for i, r := range results {
		if r.Service.Name != services[i].Name {
			t.Fatalf("ordem trocada: posicao %d esperava %q, recebi %q", i, services[i].Name, r.Service.Name)
		}
		if !r.Up {
			t.Fatalf("servico %q deveria estar no ar", r.Service.Name)
		}
	}
}

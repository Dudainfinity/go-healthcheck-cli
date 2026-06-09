# healthcheck — CLI de monitoramento em Go

Ferramenta de linha de comando, escrita em **Go**, que verifica a saúde de
**vários serviços HTTP em paralelo** (usando goroutines) e reporta o status de
cada um: no ar ou fora, código HTTP e latência. Pensada para rodar no terminal,
em um cron ou dentro de um **pipeline de CI** (sai com código de erro se algum
serviço estiver fora).

> Projeto de portfólio unindo **Go** e **DevOps/Automação** — construir a
> ferramenta, não só usar.

---

## ✨ Recursos

- ⚡ **Checagem em paralelo** com goroutines + `sync.WaitGroup`.
- 📄 **Configuração por arquivo JSON** ou URLs passadas direto no terminal.
- 🔁 **Modo watch** (`-interval`) para monitorar continuamente.
- 🧪 **Testes automatizados** com `httptest`.
- 🐳 **Imagem Docker** mínima (build multi-stage).
- 🤖 **CI no GitHub Actions** (vet + testes + build).
- 🚦 **Exit code != 0** quando há serviço fora — integra com pipelines.

---

## 🚀 Uso

### Compilando

```bash
go build -o healthcheck .
```

### Passando URLs direto

```bash
./healthcheck https://www.google.com https://github.com
```

### Usando um arquivo de configuração

```bash
cp config.example.json services.json
./healthcheck -config services.json
```

Exemplo de saída:

```
SERVIÇO   STATUS   HTTP   LATÊNCIA   DETALHE
Google    UP       200    142ms
GitHub    UP       200    210ms
Offline   DOWN     -      5001ms     Get "https://...": dial tcp: lookup ...

2/3 no ar
```

### Modo watch (monitoramento contínuo)

```bash
./healthcheck -config services.json -interval 10s
```

### Saída em JSON (para integrar com outras ferramentas)

```bash
./healthcheck -config services.json -json
```

### Opções

| Flag | Padrão | Descrição |
|---|---|---|
| `-config` | — | Arquivo JSON com a lista de serviços |
| `-timeout` | `5s` | Tempo máximo de espera por serviço |
| `-interval` | `0` | Se > 0, roda em loop (modo watch). Ex: `10s` |
| `-json` | `false` | Imprime o resultado em JSON |

---

## 🐳 Docker

```bash
docker build -t healthcheck .
docker run --rm healthcheck            # usa a config de exemplo embutida
docker run --rm -v $(pwd)/services.json:/config.json healthcheck
```

---

## 🧪 Testes

```bash
go test -v ./...
```

---

## 📄 Formato da configuração

```json
{
  "services": [
    { "name": "Google", "url": "https://www.google.com", "expected_status": 200 },
    { "name": "GitHub", "url": "https://github.com" }
  ]
}
```

`expected_status` é opcional (padrão `200`). Se o serviço responder um status
diferente do esperado, ele é considerado **fora do ar**.

---

## 📁 Estrutura

```
go-healthcheck-cli/
├── main.go                       # CLI: flags, saída em tabela/JSON, modo watch
├── internal/checker/
│   ├── checker.go                # Checagem em paralelo (goroutines)
│   └── checker_test.go           # Testes
├── config.example.json
├── Dockerfile                    # Build multi-stage
└── .github/workflows/ci.yml      # CI: vet + testes + build
```

---

## 🧰 Stack

`Go` · `Goroutines` · `net/http` · `Docker` · `GitHub Actions`

---

Feito por **Maria Eduarda** — foco em DevOps & Cloud (AWS).
[GitHub](https://github.com/Dudainfinity)

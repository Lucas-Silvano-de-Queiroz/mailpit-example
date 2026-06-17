# Mailpit Examples

Este projeto demonstra como integrar o **Mailpit** para captura de e-mails em desenvolvimento com exemplos em Node.js e Go.

## 📁 Estrutura do Projeto
- `docker-compose.yml`: Sobe o servidor Mailpit localmente.
- `node/`: Exemplo usando Node.js e Nodemailer.
- `go/`: Exemplo usando Go e o pacote nativo `net/smtp`.

---

## 🚀 Como Iniciar o Mailpit

Certifique-se de ter o Docker instalado e execute na raiz do projeto:
```bash
docker-compose up -d
```
- **SMTP**: `localhost:1025`
- **Dashboard Web**: [http://localhost:8025](http://localhost:8025)

---

## 🟢 Exemplo Node.js

1. Entre na pasta: `cd node`
2. Instale: `npm install`
3. Execute: `node index.js`

## 🔵 Exemplo Go

1. Entre na pasta: `cd go`
2. Execute: `go run main.go`


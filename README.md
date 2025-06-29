# 🔐 access-control

Uma API de **autenticação e autorização** desenvolvida em **Golang**, utilizando boas práticas, JWT e conexão com banco de dados relacional.

---

## 🚀 Tecnologias Utilizadas

- 🐹 **Go** – Linguagem principal da API
- 🌶️ **Gin** – Framework web leve e performático
- 🐘 **PostgreSQL 15** – Banco de dados relacional
- 🔑 **JWT (JSON Web Token)** – Gerenciamento de autenticação

---
## 🚀 Rodar local 
# Configurações do servidor
export SRV_PORT=8080
export SRV_MODE=DEVELOPER
export SRV_JWT_SECRET_KEY=LinuxRust162!
export SRV_JWT_TOKEN_EXP=60       
export SRV_JWT_REFRESH_EXP=90
export SRV_DB_HOST=aws-0-sa-east-1.pooler.supabase.com
export SRV_DB_NAME=postgres
export SRV_DB_USER=postgres.uldkaiigwtybxrxrvpxd
export SRV_DB_PASS=LinuxJava!162
export SRV_DB_PORT=5432
export SRV_DB_SSL_MODE=require

export SRV_RDB_HOST=localhost
export SRV_RDB_PORT=6379
export SRV_RDB_USER=
export SRV_RDB_PASS=lalal
export SRV_RDB_DB=0

  
## 📁 Estrutura Geral

```bash
.
├── cmd/                  # Ponto de entrada da aplicação
├── internal/             # Pacotes internos (handlers, config, etc.)
├── pkg/                  # Lógicas reutilizáveis (serviços, modelos, etc.)
├── go.mod                # Gerenciador de dependências
└── README.md             # Documentação do projeto
https://login-1-vfr6.onrender.com/swagger//index.html

backendkatana@gmail.com

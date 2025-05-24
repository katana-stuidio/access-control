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
export SRV_JWT_TOKEN_EXP=5          # 5 minutos
export SRV_JWT_REFRESH_EXP=30       # 30 minutos
export SRV_DB_HOST=0.0.0.0  
export SRV_DB_USER=postgres
export SRV_DB_PASS=supersenha
export SRV_DB_NAME=katana_db

  
## 📁 Estrutura Geral

```bash
.
├── cmd/                  # Ponto de entrada da aplicação
├── internal/             # Pacotes internos (handlers, config, etc.)
├── pkg/                  # Lógicas reutilizáveis (serviços, modelos, etc.)
├── go.mod                # Gerenciador de dependências
└── README.md             # Documentação do projeto

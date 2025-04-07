# 🔐 access-control

Uma API de **autenticação e autorização** desenvolvida em **Golang**, utilizando boas práticas, JWT e conexão com banco de dados relacional.

---

## 🚀 Tecnologias Utilizadas

- 🐹 **Go** – Linguagem principal da API
- 🌶️ **Gin** – Framework web leve e performático
- 🐘 **PostgreSQL 15** – Banco de dados relacional
- 🔑 **JWT (JSON Web Token)** – Gerenciamento de autenticação

---

## 📁 Estrutura Geral

```bash
.
├── cmd/                  # Ponto de entrada da aplicação
├── internal/             # Pacotes internos (handlers, config, etc.)
├── pkg/                  # Lógicas reutilizáveis (serviços, modelos, etc.)
├── go.mod                # Gerenciador de dependências
└── README.md             # Documentação do projeto

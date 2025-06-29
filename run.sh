#!/bin/bash

# Para o script imediatamente se um comando falhar
set -e

# O caminho para o seu arquivo main.go
MAIN_PATH="./cmd/api"

# Verifica se o go.mod existe, se não, inicializa o módulo.
if [ ! -f "go.mod" ]; then
    echo "▶️ Inicializando o módulo Go..."
    go mod init github.com/katana-stuidio/access-control
fi

# Sincroniza as dependências com base no seu go.mod e código fonte.
# Este comando é mais eficiente que 'go get' para projetos com módulos.
echo "▶️ Sincronizando dependências..."
go mod tidy

# Compila a aplicação, especificando o caminho correto para o main.go.
# O output (-o) será um executável chamado 'app' na raiz do projeto.
echo "⚙️ Compilando a aplicação..."
go build -o app $MAIN_PATH

# Executa a aplicação compilada.
echo "🚀 Executando a aplicação..."
./app
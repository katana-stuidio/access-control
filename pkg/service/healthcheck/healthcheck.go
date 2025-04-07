package healthcheck

import (
	"database/sql"
	"fmt"
)

// Interface para o serviço de healthcheck
type HealthcheckServiceInterface interface {
	CheckDB() (bool, error)
}

// Estrutura para o serviço de healthcheck
type HealthcheckService struct {
	db *sql.DB
}

// Função para criar uma nova instância do serviço de healthcheck
func NewHealthcheckService(db *sql.DB) *HealthcheckService {
	return &HealthcheckService{
		db: db,
	}
}

// Implementação do método CheckDB para verificar se o banco está online
func (h *HealthcheckService) CheckDB() (bool, error) {
	// Verifica se o banco está acessível com uma consulta simples
	err := h.db.Ping()
	if err != nil {
		return false, fmt.Errorf("failed to ping database: %w", err)
	}

	// Executa uma consulta simples para garantir a conexão
	var result int
	err = h.db.QueryRow("SELECT 1").Scan(&result)
	if err != nil || result != 1 {
		return false, fmt.Errorf("failed to execute healthcheck query: %w", err)
	}

	return true, nil
}

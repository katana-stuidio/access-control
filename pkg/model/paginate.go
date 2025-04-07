package model

import (
	"math"
)

type Paginate struct {
	Total    int64       `json:"total"`
	Currente int64       `json:"current_page"`
	Last     int64       `json:"last_page"`
	Data     interface{} `json:"data"`
	Limit    int64       `json:"-"`
	Page     int64       `json:"-"`
}

func (p *Paginate) GetPaginatedOpts() (cursor uint64, match string, count int64) {
	// O cursor inicial é 0.
	cursor = 0
	// Match é o padrão para correspondência (ex.: "prefixo:*")
	match = "*"
	// Limit para o comando SCAN
	count = p.Limit

	return cursor, match, count
}

// Função para paginar os dados
func (p *Paginate) Paginate(data interface{}) {
	p.Data = data
	p.Currente = p.Page
	d := float64(p.Total) / float64(p.Limit)
	p.Last = int64(math.Ceil(d))
}

// Função para criar um novo paginador
func NewPaginate(limit, page, total int64) *Paginate {
	var limitL, pageL int64 = 10, 1

	if limit > 0 {
		limitL = limit
	}
	if page > 0 {
		pageL = page
	}

	return &Paginate{
		Limit: limitL,
		Page:  pageL,
		Total: total,
	}
}

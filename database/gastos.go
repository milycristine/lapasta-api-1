package sql

import (
	"fmt"
	"lapasta/internal/models"
	"strings"
	"time"
)

func (r *SQLStr) TotalGastos(tipo string) (float64, error) {
	switch strings.ToLower(tipo) {
	case "avista":
		return r.TotalRecebimentosAvistaMesAtual()
	case "boleto":
		return r.TotalBoletosPagosMesAtual()
	default:
		totalAvista, _ := r.TotalRecebimentosAvistaMesAtual()
		totalBoletos, _ := r.TotalBoletosPagosMesAtual()
		return totalAvista + totalBoletos, nil
	}
}
func (r *SQLStr) ListarGastos(tipo, inicioData, fimData string) ([]models.Gasto, error) {
	var gastos []models.Gasto
	layout := "2006-01-02"

	if inicioData == "" || fimData == "" {
		hoje := time.Now().Format(layout)
		inicioData = hoje
		fimData = hoje
	}

	start, err := time.Parse(layout, inicioData)
	if err != nil {
		return nil, fmt.Errorf("data de início inválida: %w", err)
	}

	end, err := time.Parse(layout, fimData)
	if err != nil {
		return nil, fmt.Errorf("data de fim inválida: %w", err)
	}

	end = end.AddDate(0, 0, 1).Add(-time.Nanosecond)

	isAvista := func(fp string) bool {
		forma := strings.ToLower(fp)
		forma = strings.ReplaceAll(forma, "à", "a") 
		forma = strings.ReplaceAll(forma, " ", "")  
		return forma == "avista" || forma == "pix" || forma == "dinheiro" || forma == "debito"
	}

	if tipo == "avista" || tipo == "" {
		rcvs, err := r.FiltrarDataRecebimentos(start, end)
		if err != nil {
			return nil, fmt.Errorf("erro ao filtrar recebimentos: %w", err)
		}

		for _, rcv := range rcvs {
			if !isAvista(rcv.FormaPagamento) {
				continue
			}

			produto := ""
			if rcv.Produto != nil {
				produto = *rcv.Produto
			}

			fornecedor := ""
			if rcv.NomeFornecedor != nil {
				fornecedor = *rcv.NomeFornecedor
			}

			gastos = append(gastos, models.Gasto{
				Id:         rcv.Id,
				Data:       rcv.Dia,
				Descricao:  produto,
				Valor:      rcv.Valor,
				Tipo:       "À vista",
				Fornecedor: fornecedor,
			})
		}
	}

	if tipo == "boleto" || tipo == "" {
		boletos, err := r.FiltrarBoletosPagosPorData(inicioData, fimData)
		if err != nil {
			return nil, fmt.Errorf("erro ao filtrar boletos pagos: %w", err)
		}

		for _, b := range boletos {
			gastos = append(gastos, models.Gasto{
				Id:         b.Id,
				Data:       b.DataPagamentoOrDefault(),
				Descricao:  "Boleto",
				Valor:      b.Valor,
				Tipo:       "Boleto",
				Fornecedor: b.FornecedorNome,
			})
		}
	}


	return gastos, nil
}

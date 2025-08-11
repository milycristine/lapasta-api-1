package tiny

import (
	"encoding/json"
	sql "lapasta/database"
	"lapasta/internal/models"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GerarNotasMotoristasHandler(sqlConn *sql.SQLStr) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idAgrupamento := r.URL.Query().Get("idAgrupamento")
		if idAgrupamento == "" {
			http.Error(w, "Parâmetro idAgrupamento obrigatório", http.StatusBadRequest)
			return
		}

		tinyClient, err := NewTinyClient(sqlConn)
		if err != nil {
			http.Error(w, "Erro ao inicializar TinyClient: "+err.Error(), http.StatusInternalServerError)
			return
		}

		resp, err := tinyClient.BuscarExpedicoesPorAgrupamento(idAgrupamento)
		if err != nil {
			http.Error(w, "Erro ao buscar expedições Tiny: "+err.Error(), http.StatusInternalServerError)
			return
		}

		motoristas, err := sqlConn.BuscarMotoristasAtivos()
		if err != nil {
			http.Error(w, "Erro ao buscar motoristas: "+err.Error(), http.StatusInternalServerError)
			return
		}

		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(motoristas), func(i, j int) { motoristas[i], motoristas[j] = motoristas[j], motoristas[i] })

		notasMap := make(map[int]*models.ExpedicaoNotaMotorista)
		motoristaIndex := 0

		processadas := make(map[string]bool)

		for _, agrup := range resp.Retorno.Agrupamentos {
			for _, expWrapper := range agrup.Expedicoes {
				exp := expWrapper.Expedicao

				if processadas[exp.Identificacao] {
					continue
				}
				processadas[exp.Identificacao] = true

				motorista := motoristas[motoristaIndex]
				motoristaIndex = (motoristaIndex + 1) % len(motoristas)

				nota, exists := notasMap[motorista.Id]
				if !exists {
					nota = &models.ExpedicaoNotaMotorista{
						MotoristaCPF:  motorista.CPF,
						MotoristaNome: motorista.Nome,
						NumExpedicoes: 0,
						ValorTotal:    0,
						PesoTotal:     0,
						Expedicoes:    []models.ExpedicaoTiny{},
					}
					notasMap[motorista.Id] = nota
				}

				if nota.NumExpedicoes >= 2 {
					continue
				}

				nota.NumExpedicoes++
				nota.ValorTotal += parseDecimalString(exp.ValorDeclarado)
				nota.PesoTotal += parseDecimalString(exp.PesoBruto)
				nota.Expedicoes = append(nota.Expedicoes, exp)

				emissao := models.EmissaoNota{
					NumeroNota:         exp.Identificacao,
					Valor:              parseDecimalString(exp.ValorDeclarado),
					DataEmissao:        formatarDataSQL(exp.DataEmissao),
					Descricao:          "Expedição Tiny - Agrupamento " + agrup.IdAgrupamento,
					MotoristaId:        motorista.Id,
					MotoristaNome:      motorista.Nome,
					IdStatusLancamento: 1,
				}
				if err := sqlConn.SalvarEmissaoNota(emissao); err != nil {
					http.Error(w, "Erro ao salvar emissão: "+err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}

		var notas []models.ExpedicaoNotaMotorista
		for _, v := range notasMap {
			notas = append(notas, *v)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(notas)
	}
}

func parseDecimalString(value string) float64 {
	clean := strings.ReplaceAll(value, ",", ".")
	f, err := strconv.ParseFloat(clean, 64)
	if err != nil {
		return 0
	}
	return f
}

func formatarDataSQL(data string) string {
	layouts := []string{"02/01/2006", "02-01-2006"}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, data); err == nil {
			return t.Format("2006-01-02") 
		}
	}
	return data 
}

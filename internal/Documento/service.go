package documento

import (
	"time"
)

type DocumentoService interface {
	CriarDocumento(documento *Documento) error
	ListarDocumentos(page int) ([]Documento, error)
	FiltrarDataDocumento(inicioData, fimData time.Time) ([]Documento, error)
}

type documentoService struct {
	repo DocumentoRepository
}

func NovoDocumentoService(repo DocumentoRepository) DocumentoService {
	return &documentoService{
		repo: repo,
	}
}

func (s *documentoService) CriarDocumento(documento *Documento) error {
	return s.repo.CriarDocumento(documento)
}

func (s *documentoService) ListarDocumentos(page int) ([]Documento, error) {
	return s.repo.ListarDocumentos(page)
}
func (s *documentoService) FiltrarDataDocumento(inicioData, fimData time.Time) ([]Documento, error) {
	return s.repo.FiltrarDataDocumento(inicioData, fimData)
}

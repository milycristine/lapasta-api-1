package server

import (
	"embed"
	"lapasta/config"
	auth "lapasta/internal/AUTH"
	boleto "lapasta/internal/Boletos"
	documento "lapasta/internal/Documento"
	fornecedor "lapasta/internal/Fornecedores"
	funcionario "lapasta/internal/Funcionario"
	motorista "lapasta/internal/Motorista"
	nota "lapasta/internal/Notas"
	pagamento "lapasta/internal/Pagamento"
	ponto "lapasta/internal/Ponto"
	recebimento "lapasta/internal/Recebimento"
	utils "lapasta/internal/Utils"
	valetransporte "lapasta/internal/ValeTransporte"
	tiny "lapasta/internal/tiny"

	"log"
	"net/http"
)

var fs embed.FS

func Controllers() {
	repo := auth.NewAuthRepository(utils.ConnectionDb)
	authService := auth.NewAuthService(repo)

	recebimentoRepo := recebimento.NovoRecebimentoRepository(utils.ConnectionDb)
	recebimentoService := recebimento.NovoRecebimentoService(recebimentoRepo)
	recebimentoHandler := recebimento.NovoRecebimentoHandler(recebimentoService)

	pontoRepo := ponto.NovoPontoRepository(utils.ConnectionDb)
	pontoService := ponto.NovoPontoService(pontoRepo)
	pontoHandler := ponto.NovoPontoHandler(pontoService)

	notaRepo := nota.NovoNotaRepository(utils.ConnectionDb)
	notaService := nota.NovaNotaService(notaRepo)
	notaHandler := nota.NovaNotaHandler(notaService)

	documentoRepo := documento.NovoDocumentoRepository(utils.ConnectionDb)
	documentoService := documento.NovoDocumentoService(documentoRepo)
	documentoHandler := documento.NovoDocumentoHandler(documentoService)

	pagamentoRepo := pagamento.NovoPagamentoRepository(utils.ConnectionDb)
	pagamentoService := pagamento.NovoPagamentoService(pagamentoRepo)
	pagamentoHandler := pagamento.NovoPagamentoHandler(pagamentoService)

	funcionarioRepo := funcionario.NovoFuncionarioRepository(utils.ConnectionDb)
	funcionarioService := funcionario.NovoFuncionarioService(funcionarioRepo)
	funcionarioHandler := funcionario.NovoFuncionarioHandler(funcionarioService)

	valeRepo := valetransporte.NovoValeRepository(utils.ConnectionDb)
	valeService := valetransporte.NovoValeService(valeRepo)
	valeHandler := valetransporte.NovoValeHandler(valeService)

	fornecedorRepo := fornecedor.NovoFornecedorRepository(utils.ConnectionDb)
	fornecedorService := fornecedor.NovoFornecedorService(fornecedorRepo)
	fornecedorHandler := fornecedor.NovoFornecedorHandler(fornecedorService)

	boletoRepo := boleto.NovoBoletoRepository(utils.ConnectionDb)
	boletoService := boleto.NovoBoletoService(boletoRepo)
	boletoHandler := boleto.NovoBoletoHandler(boletoService)

	motoristaRepo := motorista.NovoMotoristaRepository(utils.ConnectionDb)
	motoristaService := motorista.NovoMotoristaService(motoristaRepo)
	motoristaHandler := motorista.NovoMotoristaHandler(motoristaService)

	log.Printf("Iniciando servidor na porta: %s", config.Yml.API.Port)

	http.HandleFunc("/login", auth.LoginHandler(authService))

	http.HandleFunc("/recebimento", recebimentoHandler.CriarRecebimento)
	http.HandleFunc("/listarRecebimento", recebimentoHandler.ListarRecebimentos)
	http.HandleFunc("/filtrarDataRecebimento", recebimentoHandler.FiltrarDataRecebimentos)
	http.HandleFunc("/buscarPorNota", recebimentoHandler.BuscarDadosRecebimentoPorNumeroNota)

	http.HandleFunc("/listarPonto", pontoHandler.ListarPontos)
	http.HandleFunc("/listarPontoId", pontoHandler.ListarPontosPorId)
	http.HandleFunc("/listarPontoIdEDia", pontoHandler.ListarPontosPorIdEDia)
	http.HandleFunc("/listarPontoPorData", pontoHandler.ListarPontosPorData)
	http.HandleFunc("/listarPontoPorDataId", pontoHandler.ListarPontosPorDataId)
	http.HandleFunc("/horaEntrada", pontoHandler.RegistrarEntrada)
	http.HandleFunc("/horaSaidaAlmoco", pontoHandler.RegistrarSaidaAlmoco)
	http.HandleFunc("/horaRetornoAlmoco", pontoHandler.RegistrarRetornoAlmoco)
	http.HandleFunc("/horaSaida", pontoHandler.RegistrarSaida)
	http.HandleFunc("/relatorioExcel", pontoHandler.GerarRelatorioMensal)

	http.HandleFunc("/nota", notaHandler.CriarNota)
	http.HandleFunc("/listarNotas", notaHandler.ListarNotas)
	http.HandleFunc("/filtrarDataNota", notaHandler.FiltrarDataNota)
	http.HandleFunc("/buscarNota", notaHandler.BuscarNotasPorNumero)

	http.HandleFunc("/documento", documentoHandler.CriarDocumento)
	http.HandleFunc("/listarDocumento", documentoHandler.ListarDocumentos)
	http.HandleFunc("/filtrarDataDocumento", documentoHandler.FiltrarDataDocumento)

	http.HandleFunc("/pagamento", pagamentoHandler.CriarPagamento)
	http.HandleFunc("/listarPagamento", pagamentoHandler.ListarPagamentos)
	http.HandleFunc("/listarPagamentoPorDia", pagamentoHandler.ListarPagamentosPorDia)
	http.HandleFunc("/atualizarStatusPagamento", pagamentoHandler.AtualizarPagamento)
	http.HandleFunc("/listarPagamentoPorMes", pagamentoHandler.ListarPagamentosPorMes)

	http.HandleFunc("/listarVale", valeHandler.ListarVales)
	http.HandleFunc("/atualizarStatusVale", valeHandler.AtualizarVale)
	http.HandleFunc("/listarValePorSemana", valeHandler.ListarValesDaSemana)

	http.HandleFunc("/funcionario", funcionarioHandler.CriarFuncionario)
	http.HandleFunc("/editarFuncionario", funcionarioHandler.EditarFuncionario)
	http.HandleFunc("/listarFuncionario", funcionarioHandler.ListarFuncionarios)
	http.HandleFunc("/buscarFuncionario", funcionarioHandler.BuscarFuncionarioPorCPF)
	http.HandleFunc("/buscarFuncionarioPorId", funcionarioHandler.BuscarFuncionarioPorID)
	http.HandleFunc("/statusFuncionario", funcionarioHandler.AtualizarStatusFuncionario)

	http.HandleFunc("/fornecedores", fornecedorHandler.CriarFornecedor)
	http.HandleFunc("/editarFornecedores", fornecedorHandler.EditarFornecedor)
	http.HandleFunc("/listarFornecedores", fornecedorHandler.ListarFornecedores)
	http.HandleFunc("/fornecedores/buscar", fornecedorHandler.BuscarFornecedorPorCNPJouNome)
	http.HandleFunc("/pedidoFornecedor", fornecedorHandler.CriarPedido)
	http.HandleFunc("/listarPedido", fornecedorHandler.ListarPedidosPorFornecedor)
	http.HandleFunc("/buscarPedidoId", fornecedorHandler.BuscarPedidosFornecedorPorDescricaoOuId)

	http.HandleFunc("/boleto", boletoHandler.CriarBoleto)
	http.HandleFunc("/listarPorPedido", boletoHandler.ListarBoletosPorRecebimento)
	http.HandleFunc("/listarPorFornecedor", boletoHandler.ListarBoletosPorFornecedor)
	http.HandleFunc("/boletodoDia", boletoHandler.ListarBoletosDoDia)
	http.HandleFunc("/boletoAPagar", boletoHandler.PagarBoleto)
	http.HandleFunc("/boletoPagos", boletoHandler.ListarBoletosPagos)
	http.HandleFunc("/boletoVencidos", boletoHandler.ListarBoletosVencidos)
	http.HandleFunc("/boletoPendentes", boletoHandler.ListarBoletosPendentes)
	http.HandleFunc("/atualizarStatusBoleto", boletoHandler.AtualizarBoleto)
	http.HandleFunc("/gerarRelatorio", boletoHandler.GerarEEnviarRelatorioBoletos)
	http.HandleFunc("/boletostotais", boletoHandler.TotaisBoletos)

	http.HandleFunc("/motorista", motoristaHandler.CriarMotorista)
	http.HandleFunc("/editarMotorista", motoristaHandler.EditarMotorista)
	http.HandleFunc("/statusmotorista", motoristaHandler.AtualizarStatusMotorista)
	http.HandleFunc("/listarMotoristas", motoristaHandler.ListarMotoristas)
	http.HandleFunc("/buscasMotoristaPorId", motoristaHandler.BuscarMotoristaPorID)
	http.HandleFunc("/emissaoNota", motoristaHandler.CriarEmissaoNota)
	http.HandleFunc("/listarEmissaoNotas", motoristaHandler.ListarEmissaoNotas)
	http.HandleFunc("/listarEmissaoNotasPorMotorista", motoristaHandler.ListarEmissaoNotasPorMotorista)
	http.HandleFunc("/buscarEmissao", motoristaHandler.BuscarEmissaoNotas)
	http.HandleFunc("/filtrarDataEmissao", motoristaHandler.FiltrarDataEmissaoNota)

	http.HandleFunc("/notaMotorista", motoristaHandler.CriarNotaMotorista)
	http.HandleFunc("/listarNotasPorMotorista", motoristaHandler.ListarNotasPorMotorista)
	http.HandleFunc("/atualizarStatusLancamentoNota", motoristaHandler.AtualizarStatusLancamentoNotaMotorista)
	http.HandleFunc("/filtrarDataNotaMotorista", motoristaHandler.FiltrarNotasMotoristaPorData)

	http.HandleFunc("/verificarSeLancouTodas", motoristaHandler.MotoristaLancouTodasAsNotas)
	http.HandleFunc("/pagamentoMotorista", motoristaHandler.CriarPagamentoMotorista)
	http.HandleFunc("/listarPagamentoMotorista", motoristaHandler.ListarPagamentosMotorista)
	http.HandleFunc("/calculoPagamento", motoristaHandler.CalcularPagamentoMotorista)
	http.HandleFunc("/atualizarStatusPagamentoMotorista", motoristaHandler.AtualizarStatusPagamentoMotorista)
	http.HandleFunc("/buscarMotoristaPorCPFouNome", motoristaHandler.BuscarMotoristaPorCPFouNome)

	http.HandleFunc("/gerarNotasMotoristas", tiny.GerarNotasMotoristasHandler(utils.ConnectionDb))

	http.Handle("/html/", http.StripPrefix("/html/", http.FileServer(http.FS(fs))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./public/images"))))

	log.Fatal(http.ListenAndServe(":"+config.Yml.API.Port, nil))
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"sghss-backend/handlers"
	"sghss-backend/middleware"
	"sghss-backend/repository"
	"sghss-backend/services"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Inicializar repositório em memória
	repo := repository.NewMemoryRepository()

	// Inicializar serviços
	authService := services.NewAuthService(repo)
	pacienteService := services.NewPacienteService(repo)
	profissionalService := services.NewProfissionalService(repo)
	consultaService := services.NewConsultaService(repo)
	prontuarioService := services.NewProntuarioService(repo)

	// Inicializar handlers
	authHandler := handlers.NewAuthHandler(authService)
	pacienteHandler := handlers.NewPacienteHandler(pacienteService)
	profissionalHandler := handlers.NewProfissionalHandler(profissionalService)
	consultaHandler := handlers.NewConsultaHandler(consultaService)
	prontuarioHandler := handlers.NewProntuarioHandler(prontuarioService)

	// Inicializar router
	router := mux.NewRouter()

	// Configurar CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	// === ROTAS PÚBLICAS ===
	
	// Rota de teste
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "SGHSS Backend API - Sistema de Gestão Hospitalar e de Serviços de Saúde", "version": "1.0.0", "status": "online"}`)
	}).Methods("GET")

	// Rotas de autenticação (públicas)
	router.HandleFunc("/auth/signup", authHandler.Signup).Methods("POST")
	router.HandleFunc("/auth/login", authHandler.Login).Methods("POST")

	// === ROTAS PROTEGIDAS ===
	
	// Subrouter para rotas protegidas
	protected := router.PathPrefix("/api").Subrouter()
	protected.Use(middleware.AuthMiddleware)

	// Rotas de pacientes (protegidas)
	protected.HandleFunc("/pacientes", pacienteHandler.CreatePaciente).Methods("POST")
	protected.HandleFunc("/pacientes", pacienteHandler.GetPacientes).Methods("GET")
	protected.HandleFunc("/pacientes/{id}", pacienteHandler.GetPacienteByID).Methods("GET")
	protected.HandleFunc("/pacientes/{id}", pacienteHandler.UpdatePaciente).Methods("PUT")
	protected.HandleFunc("/pacientes/{id}", pacienteHandler.DeletePaciente).Methods("DELETE")

	// Rotas de profissionais (protegidas)
	protected.HandleFunc("/profissionais", profissionalHandler.CreateProfissional).Methods("POST")
	protected.HandleFunc("/profissionais", profissionalHandler.GetProfissionais).Methods("GET")
	protected.HandleFunc("/profissionais/{id}", profissionalHandler.GetProfissionalByID).Methods("GET")
	protected.HandleFunc("/profissionais/{id}", profissionalHandler.UpdateProfissional).Methods("PUT")
	protected.HandleFunc("/profissionais/{id}", profissionalHandler.DeleteProfissional).Methods("DELETE")

	// Rotas de consultas (protegidas)
	protected.HandleFunc("/consultas", consultaHandler.CreateConsulta).Methods("POST")
	protected.HandleFunc("/consultas", consultaHandler.GetConsultas).Methods("GET")
	protected.HandleFunc("/consultas/{id}", consultaHandler.GetConsultaByID).Methods("GET")
	protected.HandleFunc("/consultas/{id}", consultaHandler.UpdateConsulta).Methods("PUT")
	protected.HandleFunc("/consultas/{id}", consultaHandler.DeleteConsulta).Methods("DELETE")

	// Rotas de prontuários (protegidas)
	protected.HandleFunc("/prontuarios", prontuarioHandler.CreateProntuario).Methods("POST")
	protected.HandleFunc("/prontuarios", prontuarioHandler.GetProntuarios).Methods("GET")
	protected.HandleFunc("/prontuarios/{id}", prontuarioHandler.GetProntuarioByID).Methods("GET")
	protected.HandleFunc("/prontuarios/{id}", prontuarioHandler.UpdateProntuario).Methods("PUT")
	protected.HandleFunc("/prontuarios/{id}", prontuarioHandler.DeleteProntuario).Methods("DELETE")

	// Aplicar CORS ao router
	handler := c.Handler(router)

	// Configurar servidor
	fmt.Println("=== SGHSS Backend API ===")
	fmt.Println("Servidor iniciado na porta 8080")
	fmt.Println("Acesse: http://localhost:8080")
	fmt.Println("")
	fmt.Println("Endpoints disponíveis:")
	fmt.Println("GET    /                     - Status da API")
	fmt.Println("POST   /auth/signup          - Cadastro de usuário")
	fmt.Println("POST   /auth/login           - Login")
	fmt.Println("GET    /api/pacientes        - Listar pacientes")
	fmt.Println("POST   /api/pacientes        - Criar paciente")
	fmt.Println("GET    /api/pacientes/{id}   - Buscar paciente")
	fmt.Println("PUT    /api/pacientes/{id}   - Atualizar paciente")
	fmt.Println("DELETE /api/pacientes/{id}   - Remover paciente")
	fmt.Println("GET    /api/profissionais    - Listar profissionais")
	fmt.Println("POST   /api/profissionais    - Criar profissional")
	fmt.Println("GET    /api/profissionais/{id} - Buscar profissional")
	fmt.Println("PUT    /api/profissionais/{id} - Atualizar profissional")
	fmt.Println("DELETE /api/profissionais/{id} - Remover profissional")
	fmt.Println("GET    /api/consultas        - Listar consultas")
	fmt.Println("POST   /api/consultas        - Criar consulta")
	fmt.Println("GET    /api/consultas/{id}   - Buscar consulta")
	fmt.Println("PUT    /api/consultas/{id}   - Atualizar consulta")
	fmt.Println("DELETE /api/consultas/{id}   - Remover consulta")
	fmt.Println("GET    /api/prontuarios      - Listar prontuários")
	fmt.Println("POST   /api/prontuarios      - Criar prontuário")
	fmt.Println("GET    /api/prontuarios/{id} - Buscar prontuário")
	fmt.Println("PUT    /api/prontuarios/{id} - Atualizar prontuário")
	fmt.Println("DELETE /api/prontuarios/{id} - Remover prontuário")
	fmt.Println("")
	
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", handler))
}



package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sghss-backend/handlers"
	"sghss-backend/middleware"
	"sghss-backend/models"
	"sghss-backend/repository"
	"sghss-backend/services"
	"testing"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// setupTestServer configura um servidor de teste
func setupTestServer() *httptest.Server {
	// Inicializar repositório em memória
	repo := repository.NewMemoryRepository()

	// Inicializar serviços
	authService := services.NewAuthService(repo)
	pacienteService := services.NewPacienteService(repo)
	profissionalService := services.NewProfissionalService(repo)

	// Inicializar handlers
	authHandler := handlers.NewAuthHandler(authService)
	pacienteHandler := handlers.NewPacienteHandler(pacienteService)
	profissionalHandler := handlers.NewProfissionalHandler(profissionalService)

	// Inicializar router
	router := mux.NewRouter()

	// Configurar CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	// Rotas públicas
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "SGHSS Backend API - Sistema de Gestão Hospitalar e de Serviços de Saúde", "version": "1.0.0", "status": "online"}`)
	}).Methods("GET")

	router.HandleFunc("/auth/signup", authHandler.Signup).Methods("POST")
	router.HandleFunc("/auth/login", authHandler.Login).Methods("POST")

	// Rotas protegidas
	protected := router.PathPrefix("/api").Subrouter()
	protected.Use(middleware.AuthMiddleware)

	protected.HandleFunc("/pacientes", pacienteHandler.CreatePaciente).Methods("POST")
	protected.HandleFunc("/pacientes", pacienteHandler.GetPacientes).Methods("GET")
	protected.HandleFunc("/pacientes/{id}", pacienteHandler.GetPacienteByID).Methods("GET")
	protected.HandleFunc("/pacientes/{id}", pacienteHandler.UpdatePaciente).Methods("PUT")
	protected.HandleFunc("/pacientes/{id}", pacienteHandler.DeletePaciente).Methods("DELETE")

	protected.HandleFunc("/profissionais", profissionalHandler.CreateProfissional).Methods("POST")
	protected.HandleFunc("/profissionais", profissionalHandler.GetProfissionais).Methods("GET")
	protected.HandleFunc("/profissionais/{id}", profissionalHandler.GetProfissionalByID).Methods("GET")
	protected.HandleFunc("/profissionais/{id}", profissionalHandler.UpdateProfissional).Methods("PUT")
	protected.HandleFunc("/profissionais/{id}", profissionalHandler.DeleteProfissional).Methods("DELETE")

	handler := c.Handler(router)
	return httptest.NewServer(handler)
}

// TestStatusEndpoint testa o endpoint de status
func TestStatusEndpoint(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/")
	if err != nil {
		t.Fatalf("Erro ao fazer requisição: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status esperado: %d, recebido: %d", http.StatusOK, resp.StatusCode)
	}

	fmt.Println("✅ Teste do endpoint de status: PASSOU")
}

// TestUserSignup testa o cadastro de usuário
func TestUserSignup(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	signupData := models.SignupRequest{
		Email:  "test@example.com",
		Senha:  "password123",
		Perfil: "admin",
	}

	jsonData, _ := json.Marshal(signupData)
	resp, err := http.Post(server.URL+"/auth/signup", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Erro ao fazer requisição: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Status esperado: %d, recebido: %d", http.StatusCreated, resp.StatusCode)
	}

	fmt.Println("✅ Teste de cadastro de usuário: PASSOU")
}

// TestUserLogin testa o login de usuário
func TestUserLogin(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Primeiro, cadastrar um usuário
	signupData := models.SignupRequest{
		Email:  "test@example.com",
		Senha:  "password123",
		Perfil: "admin",
	}
	jsonData, _ := json.Marshal(signupData)
	http.Post(server.URL+"/auth/signup", "application/json", bytes.NewBuffer(jsonData))

	// Agora, fazer login
	loginData := models.LoginRequest{
		Email: "test@example.com",
		Senha: "password123",
	}
	jsonData, _ = json.Marshal(loginData)
	resp, err := http.Post(server.URL+"/auth/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Erro ao fazer requisição: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status esperado: %d, recebido: %d", http.StatusOK, resp.StatusCode)
	}

	var loginResponse models.LoginResponse
	json.NewDecoder(resp.Body).Decode(&loginResponse)

	if loginResponse.Token == "" {
		t.Error("Token não foi retornado no login")
	}

	fmt.Println("✅ Teste de login de usuário: PASSOU")
}

// TestUnauthorizedAccess testa acesso sem autenticação
func TestUnauthorizedAccess(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/api/pacientes")
	if err != nil {
		t.Fatalf("Erro ao fazer requisição: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Status esperado: %d, recebido: %d", http.StatusUnauthorized, resp.StatusCode)
	}

	fmt.Println("✅ Teste de acesso não autorizado: PASSOU")
}

// getAuthToken obtém um token de autenticação para testes
func getAuthToken(serverURL string) (string, error) {
	// Cadastrar usuário
	signupData := models.SignupRequest{
		Email:  "test@example.com",
		Senha:  "password123",
		Perfil: "admin",
	}
	jsonData, _ := json.Marshal(signupData)
	http.Post(serverURL+"/auth/signup", "application/json", bytes.NewBuffer(jsonData))

	// Fazer login
	loginData := models.LoginRequest{
		Email: "test@example.com",
		Senha: "password123",
	}
	jsonData, _ = json.Marshal(loginData)
	resp, err := http.Post(serverURL+"/auth/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var loginResponse models.LoginResponse
	json.NewDecoder(resp.Body).Decode(&loginResponse)
	return loginResponse.Token, nil
}

// TestCreatePaciente testa a criação de paciente
func TestCreatePaciente(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	token, err := getAuthToken(server.URL)
	if err != nil {
		t.Fatalf("Erro ao obter token: %v", err)
	}

	pacienteData := models.Paciente{
		Nome:             "João Silva",
		CPF:              "12345678900",
		HistoricoClinico: "Teste de histórico clínico",
	}

	jsonData, _ := json.Marshal(pacienteData)
	req, _ := http.NewRequest("POST", server.URL+"/api/pacientes", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Erro ao fazer requisição: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Status esperado: %d, recebido: %d", http.StatusCreated, resp.StatusCode)
	}

	fmt.Println("✅ Teste de criação de paciente: PASSOU")
}

// TestGetPacientes testa a listagem de pacientes
func TestGetPacientes(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	token, err := getAuthToken(server.URL)
	if err != nil {
		t.Fatalf("Erro ao obter token: %v", err)
	}

	req, _ := http.NewRequest("GET", server.URL+"/api/pacientes", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Erro ao fazer requisição: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status esperado: %d, recebido: %d", http.StatusOK, resp.StatusCode)
	}

	fmt.Println("✅ Teste de listagem de pacientes: PASSOU")
}

// TestCreateProfissional testa a criação de profissional
func TestCreateProfissional(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	token, err := getAuthToken(server.URL)
	if err != nil {
		t.Fatalf("Erro ao obter token: %v", err)
	}

	profissionalData := models.ProfissionalSaude{
		Nome:            "Dr. Maria Oliveira",
		CRMCOREN:        "CRM-SP 123456",
		Especialidade:   "Cardiologia",
		PerfilUsuarioID: 1,
	}

	jsonData, _ := json.Marshal(profissionalData)
	req, _ := http.NewRequest("POST", server.URL+"/api/profissionais", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Erro ao fazer requisição: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Status esperado: %d, recebido: %d", http.StatusCreated, resp.StatusCode)
	}

	fmt.Println("✅ Teste de criação de profissional: PASSOU")
}

func main() {
	fmt.Println("=== EXECUTANDO TESTES DA API SGHSS ===")
	fmt.Println()

	// Executar testes
	testing.Main(func(pat, str string) (bool, error) { return true, nil },
		[]testing.InternalTest{
			{"TestStatusEndpoint", TestStatusEndpoint},
			{"TestUserSignup", TestUserSignup},
			{"TestUserLogin", TestUserLogin},
			{"TestUnauthorizedAccess", TestUnauthorizedAccess},
			{"TestCreatePaciente", TestCreatePaciente},
			{"TestGetPacientes", TestGetPacientes},
			{"TestCreateProfissional", TestCreateProfissional},
		},
		nil, nil)

	fmt.Println()
	fmt.Println("=== TODOS OS TESTES CONCLUÍDOS ===")
}


package handlers

import (
	"encoding/json"
	"net/http"
	"sghss-backend/models"
	"sghss-backend/services"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// AuthHandler gerencia as requisições de autenticação
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler cria uma nova instância do handler de autenticação
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Signup registra um novo usuário
func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var req models.SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Dados inválidos"}`, http.StatusBadRequest)
		return
	}

	usuario, err := h.authService.Signup(req)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Usuário criado com sucesso",
		"usuario": usuario,
	})
}

// Login autentica um usuário
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Dados inválidos"}`, http.StatusBadRequest)
		return
	}

	response, err := h.authService.Login(req)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// PacienteHandler gerencia as requisições relacionadas a pacientes
type PacienteHandler struct {
	pacienteService *services.PacienteService
}

// NewPacienteHandler cria uma nova instância do handler de pacientes
func NewPacienteHandler(pacienteService *services.PacienteService) *PacienteHandler {
	return &PacienteHandler{pacienteService: pacienteService}
}

// CreatePaciente cria um novo paciente
func (h *PacienteHandler) CreatePaciente(w http.ResponseWriter, r *http.Request) {
	var paciente models.Paciente
	if err := json.NewDecoder(r.Body).Decode(&paciente); err != nil {
		http.Error(w, `{"error": "Dados inválidos"}`, http.StatusBadRequest)
		return
	}

	err := h.pacienteService.CreatePaciente(&paciente)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Paciente criado com sucesso",
		"paciente": paciente,
	})
}

// GetPacientes retorna todos os pacientes
func (h *PacienteHandler) GetPacientes(w http.ResponseWriter, r *http.Request) {
	pacientes := h.pacienteService.GetAllPacientes()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"pacientes": pacientes,
		"total":     len(pacientes),
	})
}

// GetPacienteByID retorna um paciente específico
func (h *PacienteHandler) GetPacienteByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, `{"error": "ID inválido"}`, http.StatusBadRequest)
		return
	}

	paciente, err := h.pacienteService.GetPacienteByID(id)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paciente)
}

// UpdatePaciente atualiza um paciente
func (h *PacienteHandler) UpdatePaciente(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, `{"error": "ID inválido"}`, http.StatusBadRequest)
		return
	}

	var paciente models.Paciente
	if err := json.NewDecoder(r.Body).Decode(&paciente); err != nil {
		http.Error(w, `{"error": "Dados inválidos"}`, http.StatusBadRequest)
		return
	}

	paciente.ID = id
	paciente.AtualizadoEm = time.Now()

	err = h.pacienteService.UpdatePaciente(&paciente)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Paciente atualizado com sucesso",
		"paciente": paciente,
	})
}

// DeletePaciente remove um paciente
func (h *PacienteHandler) DeletePaciente(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, `{"error": "ID inválido"}`, http.StatusBadRequest)
		return
	}

	err = h.pacienteService.DeletePaciente(id)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Paciente removido com sucesso",
	})
}

// ProfissionalHandler gerencia as requisições relacionadas a profissionais de saúde
type ProfissionalHandler struct {
	profissionalService *services.ProfissionalService
}

// NewProfissionalHandler cria uma nova instância do handler de profissionais
func NewProfissionalHandler(profissionalService *services.ProfissionalService) *ProfissionalHandler {
	return &ProfissionalHandler{profissionalService: profissionalService}
}

// CreateProfissional cria um novo profissional
func (h *ProfissionalHandler) CreateProfissional(w http.ResponseWriter, r *http.Request) {
	var profissional models.ProfissionalSaude
	if err := json.NewDecoder(r.Body).Decode(&profissional); err != nil {
		http.Error(w, `{"error": "Dados inválidos"}`, http.StatusBadRequest)
		return
	}

	err := h.profissionalService.CreateProfissional(&profissional)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":     "Profissional criado com sucesso",
		"profissional": profissional,
	})
}

// GetProfissionais retorna todos os profissionais
func (h *ProfissionalHandler) GetProfissionais(w http.ResponseWriter, r *http.Request) {
	profissionais := h.profissionalService.GetAllProfissionais()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"profissionais": profissionais,
		"total":         len(profissionais),
	})
}

// GetProfissionalByID retorna um profissional específico
func (h *ProfissionalHandler) GetProfissionalByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, `{"error": "ID inválido"}`, http.StatusBadRequest)
		return
	}

	profissional, err := h.profissionalService.GetProfissionalByID(id)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profissional)
}

// UpdateProfissional atualiza um profissional
func (h *ProfissionalHandler) UpdateProfissional(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, `{"error": "ID inválido"}`, http.StatusBadRequest)
		return
	}

	var profissional models.ProfissionalSaude
	if err := json.NewDecoder(r.Body).Decode(&profissional); err != nil {
		http.Error(w, `{"error": "Dados inválidos"}`, http.StatusBadRequest)
		return
	}

	profissional.ID = id
	profissional.AtualizadoEm = time.Now()

	err = h.profissionalService.UpdateProfissional(&profissional)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":      "Profissional atualizado com sucesso",
		"profissional": profissional,
	})
}

// DeleteProfissional remove um profissional
func (h *ProfissionalHandler) DeleteProfissional(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, `{"error": "ID inválido"}`, http.StatusBadRequest)
		return
	}

	err = h.profissionalService.DeleteProfissional(id)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Profissional removido com sucesso",
	})
}

// ConsultaHandler gerencia as requisições relacionadas a consultas
type ConsultaHandler struct {
	consultaService *services.ConsultaService
}

// NewConsultaHandler cria uma nova instância do handler de consultas
func NewConsultaHandler(consultaService *services.ConsultaService) *ConsultaHandler {
	return &ConsultaHandler{consultaService: consultaService}
}

// CreateConsulta cria uma nova consulta
func (h *ConsultaHandler) CreateConsulta(w http.ResponseWriter, r *http.Request) {
	var consulta models.Consulta
	if err := json.NewDecoder(r.Body).Decode(&consulta); err != nil {
		http.Error(w, `{"error": "Dados inválidos"}`, http.StatusBadRequest)
		return
	}

	err := h.consultaService.CreateConsulta(&consulta)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Consulta criada com sucesso",
		"consulta": consulta,
	})
}

// GetConsultas retorna todas as consultas
func (h *ConsultaHandler) GetConsultas(w http.ResponseWriter, r *http.Request) {
	consultas := h.consultaService.GetAllConsultas()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"consultas": consultas,
		"total":     len(consultas),
	})
}

// GetConsultaByID retorna uma consulta específica
func (h *ConsultaHandler) GetConsultaByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, `{"error": "ID inválido"}`, http.StatusBadRequest)
		return
	}

	consulta, err := h.consultaService.GetConsultaByID(id)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(consulta)
}

// UpdateConsulta atualiza uma consulta
func (h *ConsultaHandler) UpdateConsulta(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, `{"error": "ID inválido"}`, http.StatusBadRequest)
		return
	}

	var consulta models.Consulta
	if err := json.NewDecoder(r.Body).Decode(&consulta); err != nil {
		http.Error(w, `{"error": "Dados inválidos"}`, http.StatusBadRequest)
		return
	}

	consulta.ID = id
	consulta.AtualizadoEm = time.Now()

	err = h.consultaService.UpdateConsulta(&consulta)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Consulta atualizada com sucesso",
		"consulta": consulta,
	})
}

// DeleteConsulta remove uma consulta
func (h *ConsultaHandler) DeleteConsulta(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, `{"error": "ID inválido"}`, http.StatusBadRequest)
		return
	}

	err = h.consultaService.DeleteConsulta(id)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Consulta removida com sucesso",
	})
}

// ProntuarioHandler gerencia as requisições relacionadas a prontuários
type ProntuarioHandler struct {
	prontuarioService *services.ProntuarioService
}

// NewProntuarioHandler cria uma nova instância do handler de prontuários
func NewProntuarioHandler(prontuarioService *services.ProntuarioService) *ProntuarioHandler {
	return &ProntuarioHandler{prontuarioService: prontuarioService}
}

// CreateProntuario cria um novo prontuário
func (h *ProntuarioHandler) CreateProntuario(w http.ResponseWriter, r *http.Request) {
	var prontuario models.Prontuario
	if err := json.NewDecoder(r.Body).Decode(&prontuario); err != nil {
		http.Error(w, `{"error": "Dados inválidos"}`, http.StatusBadRequest)
		return
	}

	err := h.prontuarioService.CreateProntuario(&prontuario)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "Prontuário criado com sucesso",
		"prontuario": prontuario,
	})
}

// GetProntuarios retorna todos os prontuários
func (h *ProntuarioHandler) GetProntuarios(w http.ResponseWriter, r *http.Request) {
	prontuarios := h.prontuarioService.GetAllProntuarios()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"prontuarios": prontuarios,
		"total":       len(prontuarios),
	})
}

// GetProntuarioByID retorna um prontuário específico
func (h *ProntuarioHandler) GetProntuarioByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, `{"error": "ID inválido"}`, http.StatusBadRequest)
		return
	}

	prontuario, err := h.prontuarioService.GetProntuarioByID(id)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prontuario)
}

// UpdateProntuario atualiza um prontuário
func (h *ProntuarioHandler) UpdateProntuario(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, `{"error": "ID inválido"}`, http.StatusBadRequest)
		return
	}

	var prontuario models.Prontuario
	if err := json.NewDecoder(r.Body).Decode(&prontuario); err != nil {
		http.Error(w, `{"error": "Dados inválidos"}`, http.StatusBadRequest)
		return
	}

	prontuario.ID = id
	prontuario.AtualizadoEm = time.Now()

	err = h.prontuarioService.UpdateProntuario(&prontuario)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "Prontuário atualizado com sucesso",
		"prontuario": prontuario,
	})
}

// DeleteProntuario remove um prontuário
func (h *ProntuarioHandler) DeleteProntuario(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, `{"error": "ID inválido"}`, http.StatusBadRequest)
		return
	}

	err = h.prontuarioService.DeleteProntuario(id)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Prontuário removido com sucesso",
	})
}



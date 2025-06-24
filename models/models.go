package models

import (
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("recurso não encontrado")
	ErrConflict = errors.New("conflito de dados (já existe)")
)

// Usuario representa um usuário do sistema (para autenticação)
type Usuario struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	SenhaHash string `json:"-"` // Não incluir na serialização JSON
	Perfil   string `json:"perfil"` // "admin", "medico", "enfermeiro", "tecnico"
	CriadoEm time.Time `json:"criado_em"`
}

// Paciente representa um paciente no sistema
type Paciente struct {
	ID               int       `json:"id"`
	Nome             string    `json:"nome"`
	CPF              string    `json:"cpf"`
	DataNascimento   time.Time `json:"data_nascimento"`
	HistoricoClinico string    `json:"historico_clinico"`
	CriadoEm         time.Time `json:"criado_em"`
	AtualizadoEm     time.Time `json:"atualizado_em"`
}

// ProfissionalSaude representa um profissional de saúde
type ProfissionalSaude struct {
	ID              int       `json:"id"`
	Nome            string    `json:"nome"`
	CRMCOREN        string    `json:"crm_coren"` // CRM para médicos, COREN para enfermeiros
	Especialidade   string    `json:"especialidade"`
	PerfilUsuarioID int       `json:"perfil_usuario_id"` // Referência ao usuário
	CriadoEm        time.Time `json:"criado_em"`
	AtualizadoEm    time.Time `json:"atualizado_em"`
}

// Consulta representa uma consulta agendada
type Consulta struct {
	ID                int       `json:"id"`
	PacienteID        int       `json:"paciente_id"`
	ProfissionalID    int       `json:"profissional_id"`
	DataHora          time.Time `json:"data_hora"`
	Especialidade     string    `json:"especialidade"`
	Status            string    `json:"status"` // "agendada", "realizada", "cancelada"
	Observacoes       string    `json:"observacoes"`
	CriadoEm          time.Time `json:"criado_em"`
	AtualizadoEm      time.Time `json:"atualizado_em"`
}

// Prontuario representa o prontuário eletrônico de um paciente
type Prontuario struct {
	ID                int       `json:"id"`
	PacienteID        int       `json:"paciente_id"`
	ProfissionalID    int       `json:"profissional_id"`
	DataAtendimento   time.Time `json:"data_atendimento"`
	Diagnostico       string    `json:"diagnostico"`
	Tratamento        string    `json:"tratamento"`
	Medicamentos      string    `json:"medicamentos"`
	Observacoes       string    `json:"observacoes"`
	CriadoEm          time.Time `json:"criado_em"`
	AtualizadoEm      time.Time `json:"atualizado_em"`
}

// LoginRequest representa a estrutura de requisição de login
type LoginRequest struct {
	Email string `json:"email"`
	Senha string `json:"senha"`
}

// LoginResponse representa a resposta de login
type LoginResponse struct {
	Token   string  `json:"token"`
	Usuario Usuario `json:"usuario"`
}

// SignupRequest representa a estrutura de requisição de cadastro
type SignupRequest struct {
	Email  string `json:"email"`
	Senha  string `json:"senha"`
	Perfil string `json:"perfil"`
}



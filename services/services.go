package services

import (
	"errors"
	"sghss-backend/models"
	"sghss-backend/repository"
	"sghss-backend/utils"
	"strings"
)

// AuthService gerencia a autenticação de usuários
type AuthService struct {
	repo *repository.MemoryRepository
}

// NewAuthService cria uma nova instância do serviço de autenticação
func NewAuthService(repo *repository.MemoryRepository) *AuthService {
	return &AuthService{repo: repo}
}

// Signup registra um novo usuário
func (s *AuthService) Signup(req models.SignupRequest) (*models.Usuario, error) {
	// Validações básicas
	if req.Email == "" || req.Senha == "" {
		return nil, errors.New("email e senha são obrigatórios")
	}

	if len(req.Senha) < 6 {
		return nil, errors.New("senha deve ter pelo menos 6 caracteres")
	}

	// Validar perfil
	validPerfis := []string{"admin", "medico", "enfermeiro", "tecnico"}
	perfilValido := false
	for _, perfil := range validPerfis {
		if req.Perfil == perfil {
			perfilValido = true
			break
		}
	}
	if !perfilValido {
		return nil, errors.New("perfil inválido")
	}

	// Hash da senha
	senhaHash, err := utils.HashPassword(req.Senha)
	if err != nil {
		return nil, errors.New("erro ao processar senha")
	}

	// Criar usuário
	usuario := models.Usuario{
		Email:     strings.ToLower(req.Email),
		SenhaHash: senhaHash,
		Perfil:    req.Perfil,
	}

	createdUser, err := s.repo.CreateUser(usuario)
	if err != nil {
		return nil, err
	}

	return &createdUser, nil
}

// Login autentica um usuário
func (s *AuthService) Login(req models.LoginRequest) (*models.LoginResponse, error) {
	// Validações básicas
	if req.Email == "" || req.Senha == "" {
		return nil, errors.New("email e senha são obrigatórios")
	}

	// Buscar usuário por email
	usuario, err := s.repo.GetUserByEmail(strings.ToLower(req.Email))
	if err != nil {
		return nil, errors.New("credenciais inválidas")
	}

	// Verificar senha
	if !utils.CheckPasswordHash(req.Senha, usuario.SenhaHash) {
		return nil, errors.New("credenciais inválidas")
	}

	// Gerar token JWT
	token, err := utils.GenerateJWT(usuario.ID, usuario.Email, usuario.Perfil)
	if err != nil {
		return nil, errors.New("erro ao gerar token")
	}

	return &models.LoginResponse{
		Token:   token,
		Usuario: usuario,
	}, nil
}

// PacienteService gerencia operações relacionadas a pacientes
type PacienteService struct {
	repo *repository.MemoryRepository
}

// NewPacienteService cria uma nova instância do serviço de pacientes
func NewPacienteService(repo *repository.MemoryRepository) *PacienteService {
	return &PacienteService{repo: repo}
}

// CreatePaciente cria um novo paciente
func (s *PacienteService) CreatePaciente(paciente *models.Paciente) error {
	// Validações básicas
	if paciente.Nome == "" || paciente.CPF == "" {
		return errors.New("nome e CPF são obrigatórios")
	}

	_, err := s.repo.CreatePaciente(*paciente)
	return err
}

// GetPacienteByID busca um paciente por ID
func (s *PacienteService) GetPacienteByID(id int) (*models.Paciente, error) {
	paciente, err := s.repo.GetPacienteByID(id)
	return &paciente, err
}

// GetAllPacientes retorna todos os pacientes
func (s *PacienteService) GetAllPacientes() []*models.Paciente {
	pacientes, _ := s.repo.GetAllPacientes()
	var result []*models.Paciente
	for i := range pacientes {
		result = append(result, &pacientes[i])
	}
	return result
}

// UpdatePaciente atualiza um paciente
func (s *PacienteService) UpdatePaciente(paciente *models.Paciente) error {
	if paciente.Nome == "" || paciente.CPF == "" {
		return errors.New("nome e CPF são obrigatórios")
	}
	_, err := s.repo.UpdatePaciente(*paciente)
	return err
}

// DeletePaciente remove um paciente
func (s *PacienteService) DeletePaciente(id int) error {
	return s.repo.DeletePaciente(id)
}

// ProfissionalService gerencia operações relacionadas a profissionais de saúde
type ProfissionalService struct {
	repo *repository.MemoryRepository
}

// NewProfissionalService cria uma nova instância do serviço de profissionais
func NewProfissionalService(repo *repository.MemoryRepository) *ProfissionalService {
	return &ProfissionalService{repo: repo}
}

// CreateProfissional cria um novo profissional
func (s *ProfissionalService) CreateProfissional(profissional *models.ProfissionalSaude) error {
	// Validações básicas
	if profissional.Nome == "" || profissional.CRMCOREN == "" {
		return errors.New("nome e CRM/COREN são obrigatórios")
	}

	_, err := s.repo.CreateProfissional(*profissional)
	return err
}

// GetProfissionalByID busca um profissional por ID
func (s *ProfissionalService) GetProfissionalByID(id int) (*models.ProfissionalSaude, error) {
	profissional, err := s.repo.GetProfissionalByID(id)
	return &profissional, err
}

// GetAllProfissionais retorna todos os profissionais
func (s *ProfissionalService) GetAllProfissionais() []*models.ProfissionalSaude {
	profissionais, _ := s.repo.GetAllProfissionais()
	var result []*models.ProfissionalSaude
	for i := range profissionais {
		result = append(result, &profissionais[i])
	}
	return result
}

// UpdateProfissional atualiza um profissional
func (s *ProfissionalService) UpdateProfissional(profissional *models.ProfissionalSaude) error {
	if profissional.Nome == "" || profissional.CRMCOREN == "" {
		return errors.New("nome e CRM/COREN são obrigatórios")
	}
	_, err := s.repo.UpdateProfissional(*profissional)
	return err
}

// DeleteProfissional remove um profissional
func (s *ProfissionalService) DeleteProfissional(id int) error {
	return s.repo.DeleteProfissional(id)
}

// ConsultaService gerencia operações relacionadas a consultas
type ConsultaService struct {
	repo *repository.MemoryRepository
}

// NewConsultaService cria uma nova instância do serviço de consultas
func NewConsultaService(repo *repository.MemoryRepository) *ConsultaService {
	return &ConsultaService{repo: repo}
}

// CreateConsulta cria uma nova consulta
func (s *ConsultaService) CreateConsulta(consulta *models.Consulta) error {
	// Validações adicionais podem ser feitas aqui, por exemplo, verificar se paciente e profissional existem
	if consulta.PacienteID == 0 || consulta.ProfissionalID == 0 || consulta.DataHora.IsZero() {
		return errors.New("paciente, profissional e data/hora são obrigatórios para a consulta")
	}
	_, err := s.repo.CreateConsulta(*consulta)
	return err
}

// GetConsultaByID busca uma consulta por ID
func (s *ConsultaService) GetConsultaByID(id int) (*models.Consulta, error) {
	consulta, err := s.repo.GetConsultaByID(id)
	return &consulta, err
}

// GetAllConsultas retorna todas as consultas
func (s *ConsultaService) GetAllConsultas() []*models.Consulta {
	consultas, _ := s.repo.GetAllConsultas()
	var result []*models.Consulta
	for i := range consultas {
		result = append(result, &consultas[i])
	}
	return result
}

// UpdateConsulta atualiza uma consulta
func (s *ConsultaService) UpdateConsulta(consulta *models.Consulta) error {
	if consulta.PacienteID == 0 || consulta.ProfissionalID == 0 || consulta.DataHora.IsZero() {
		return errors.New("paciente, profissional e data/hora são obrigatórios para a consulta")
	}
	_, err := s.repo.UpdateConsulta(*consulta)
	return err
}

// DeleteConsulta remove uma consulta
func (s *ConsultaService) DeleteConsulta(id int) error {
	return s.repo.DeleteConsulta(id)
}

// ProntuarioService gerencia operações relacionadas a prontuários
type ProntuarioService struct {
	repo *repository.MemoryRepository
}

// NewProntuarioService cria uma nova instância do serviço de prontuários
func NewProntuarioService(repo *repository.MemoryRepository) *ProntuarioService {
	return &ProntuarioService{repo: repo}
}

// CreateProntuario cria um novo prontuário
func (s *ProntuarioService) CreateProntuario(prontuario *models.Prontuario) error {
	// Validações adicionais podem ser feitas aqui
	if prontuario.PacienteID == 0 || prontuario.ProfissionalID == 0 || prontuario.DataAtendimento.IsZero() || prontuario.Diagnostico == "" {
		return errors.New("paciente, profissional, data de atendimento e diagnóstico são obrigatórios para o prontuário")
	}
	_, err := s.repo.CreateProntuario(*prontuario)
	return err
}

// GetProntuarioByID busca um prontuário por ID
func (s *ProntuarioService) GetProntuarioByID(id int) (*models.Prontuario, error) {
	prontuario, err := s.repo.GetProntuarioByID(id)
	return &prontuario, err
}

// GetAllProntuarios retorna todos os prontuários
func (s *ProntuarioService) GetAllProntuarios() []*models.Prontuario {
	prontuarios, _ := s.repo.GetAllProntuarios()
	var result []*models.Prontuario
	for i := range prontuarios {
		result = append(result, &prontuarios[i])
	}
	return result
}

// UpdateProntuario atualiza um prontuário
func (s *ProntuarioService) UpdateProntuario(prontuario *models.Prontuario) error {
	if prontuario.PacienteID == 0 || prontuario.ProfissionalID == 0 || prontuario.DataAtendimento.IsZero() || prontuario.Diagnostico == "" {
		return errors.New("paciente, profissional, data de atendimento e diagnóstico são obrigatórios para o prontuário")
	}
	_, err := s.repo.UpdateProntuario(*prontuario)
	return err
}

// DeleteProntuario remove um prontuário
func (s *ProntuarioService) DeleteProntuario(id int) error {
	return s.repo.DeleteProntuario(id)
}



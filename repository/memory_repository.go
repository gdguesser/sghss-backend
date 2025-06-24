package repository

import (
	"sghss-backend/models"
	"sync"
	"time"
)

type MemoryRepository struct {
	mu                  sync.RWMutex
	users               map[int]models.Usuario
	pacientes           map[int]models.Paciente
	profissionais       map[int]models.ProfissionalSaude
	consultas           map[int]models.Consulta
	prontuarios         map[int]models.Prontuario

	nextUserID          int
	nextPacienteID      int
	nextProfissionalID  int
	nextConsultaID      int
	nextProntuarioID    int

	pacientesByCPF      map[string]int // CPF -> ID do paciente
	profissionaisByCRMCOREN map[string]int // CRM/COREN -> ID do profissional
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		users:               make(map[int]models.Usuario),
		pacientes:           make(map[int]models.Paciente),
		profissionais:       make(map[int]models.ProfissionalSaude),
		consultas:           make(map[int]models.Consulta),
		prontuarios:         make(map[int]models.Prontuario),
		nextUserID:          1,
		nextPacienteID:      1,
		nextProfissionalID:  1,
		nextConsultaID:      1,
		nextProntuarioID:    1,
		pacientesByCPF:      make(map[string]int),
		profissionaisByCRMCOREN: make(map[string]int),
	}
}

// --- Métodos para Usuário ---

func (r *MemoryRepository) CreateUser(user models.Usuario) (models.Usuario, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user.ID = r.nextUserID
	user.CriadoEm = time.Now()
	r.users[user.ID] = user
	r.nextUserID++
	return user, nil
}

func (r *MemoryRepository) GetUserByID(id int) (models.Usuario, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[id]
	if !ok {
		return models.Usuario{}, models.ErrNotFound
	}
	return user, nil
}

func (r *MemoryRepository) GetUserByEmail(email string) (models.Usuario, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return models.Usuario{}, models.ErrNotFound
}

// --- Métodos para Paciente ---

func (r *MemoryRepository) CreatePaciente(paciente models.Paciente) (models.Paciente, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.pacientesByCPF[paciente.CPF]; exists {
		return models.Paciente{}, models.ErrConflict
	}

	paciente.ID = r.nextPacienteID
	paciente.CriadoEm = time.Now()
	paciente.AtualizadoEm = time.Now()
	r.pacientes[paciente.ID] = paciente
	r.pacientesByCPF[paciente.CPF] = paciente.ID
	r.nextPacienteID++
	return paciente, nil
}

func (r *MemoryRepository) GetPacienteByID(id int) (models.Paciente, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	paciente, ok := r.pacientes[id]
	if !ok {
		return models.Paciente{}, models.ErrNotFound
	}
	return paciente, nil
}

func (r *MemoryRepository) GetAllPacientes() ([]models.Paciente, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	pacientes := make([]models.Paciente, 0, len(r.pacientes))
	for _, p := range r.pacientes {
		pacientes = append(pacientes, p)
	}
	return pacientes, nil
}

func (r *MemoryRepository) UpdatePaciente(paciente models.Paciente) (models.Paciente, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	oldPaciente, ok := r.pacientes[paciente.ID]
	if !ok {
		return models.Paciente{}, models.ErrNotFound
	}

	if oldPaciente.CPF != paciente.CPF {
		if _, exists := r.pacientesByCPF[paciente.CPF]; exists {
			return models.Paciente{}, models.ErrConflict
		}
		delete(r.pacientesByCPF, oldPaciente.CPF)
		r.pacientesByCPF[paciente.CPF] = paciente.ID
	}

	paciente.CriadoEm = oldPaciente.CriadoEm
	paciente.AtualizadoEm = time.Now()
	r.pacientes[paciente.ID] = paciente
	return paciente, nil
}

func (r *MemoryRepository) DeletePaciente(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	paciente, ok := r.pacientes[id]
	if !ok {
		return models.ErrNotFound
	}

	delete(r.pacientesByCPF, paciente.CPF)
	delete(r.pacientes, id)
	return nil
}

// --- Métodos para ProfissionalSaude ---

func (r *MemoryRepository) CreateProfissional(profissional models.ProfissionalSaude) (models.ProfissionalSaude, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.profissionaisByCRMCOREN[profissional.CRMCOREN]; exists {
		return models.ProfissionalSaude{}, models.ErrConflict
	}

	profissional.ID = r.nextProfissionalID
	profissional.CriadoEm = time.Now()
	profissional.AtualizadoEm = time.Now()
	r.profissionais[profissional.ID] = profissional
	r.profissionaisByCRMCOREN[profissional.CRMCOREN] = profissional.ID
	r.nextProfissionalID++
	return profissional, nil
}

func (r *MemoryRepository) GetProfissionalByID(id int) (models.ProfissionalSaude, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	profissional, ok := r.profissionais[id]
	if !ok {
		return models.ProfissionalSaude{}, models.ErrNotFound
	}
	return profissional, nil
}

func (r *MemoryRepository) GetAllProfissionais() ([]models.ProfissionalSaude, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	profissionais := make([]models.ProfissionalSaude, 0, len(r.profissionais))
	for _, p := range r.profissionais {
		profissionais = append(profissionais, p)
	}
	return profissionais, nil
}

func (r *MemoryRepository) UpdateProfissional(profissional models.ProfissionalSaude) (models.ProfissionalSaude, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	oldProfissional, ok := r.profissionais[profissional.ID]
	if !ok {
		return models.ProfissionalSaude{}, models.ErrNotFound
	}

	if oldProfissional.CRMCOREN != profissional.CRMCOREN {
		if _, exists := r.profissionaisByCRMCOREN[profissional.CRMCOREN]; exists {
			return models.ProfissionalSaude{}, models.ErrConflict
		}
		delete(r.profissionaisByCRMCOREN, oldProfissional.CRMCOREN)
		r.profissionaisByCRMCOREN[profissional.CRMCOREN] = profissional.ID
	}

	profissional.CriadoEm = oldProfissional.CriadoEm
	profissional.AtualizadoEm = time.Now()
	r.profissionais[profissional.ID] = profissional
	return profissional, nil
}

func (r *MemoryRepository) DeleteProfissional(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	profissional, ok := r.profissionais[id]
	if !ok {
		return models.ErrNotFound
	}

	delete(r.profissionaisByCRMCOREN, profissional.CRMCOREN)
	delete(r.profissionais, id)
	return nil
}

// --- Métodos para Consulta ---

func (r *MemoryRepository) CreateConsulta(consulta models.Consulta) (models.Consulta, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	consulta.ID = r.nextConsultaID
	consulta.CriadoEm = time.Now()
	consulta.AtualizadoEm = time.Now()
	r.consultas[consulta.ID] = consulta
	r.nextConsultaID++
	return consulta, nil
}

func (r *MemoryRepository) GetConsultaByID(id int) (models.Consulta, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	consulta, ok := r.consultas[id]
	if !ok {
		return models.Consulta{}, models.ErrNotFound
	}
	return consulta, nil
}

func (r *MemoryRepository) GetAllConsultas() ([]models.Consulta, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	consultas := make([]models.Consulta, 0, len(r.consultas))
	for _, c := range r.consultas {
		consultas = append(consultas, c)
	}
	return consultas, nil
}

func (r *MemoryRepository) UpdateConsulta(consulta models.Consulta) (models.Consulta, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	oldConsulta, ok := r.consultas[consulta.ID]
	if !ok {
		return models.Consulta{}, models.ErrNotFound
	}

	consulta.CriadoEm = oldConsulta.CriadoEm
	consulta.AtualizadoEm = time.Now()
	r.consultas[consulta.ID] = consulta
	return consulta, nil
}

func (r *MemoryRepository) DeleteConsulta(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.consultas[id]; !ok {
		return models.ErrNotFound
	}
	delete(r.consultas, id)
	return nil
}

// --- Métodos para Prontuario ---

func (r *MemoryRepository) CreateProntuario(prontuario models.Prontuario) (models.Prontuario, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	prontuario.ID = r.nextProntuarioID
	prontuario.CriadoEm = time.Now()
	prontuario.AtualizadoEm = time.Now()
	r.prontuarios[prontuario.ID] = prontuario
	r.nextProntuarioID++
	return prontuario, nil
}

func (r *MemoryRepository) GetProntuarioByID(id int) (models.Prontuario, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	prontuario, ok := r.prontuarios[id]
	if !ok {
		return models.Prontuario{}, models.ErrNotFound
	}
	return prontuario, nil
}

func (r *MemoryRepository) GetAllProntuarios() ([]models.Prontuario, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	prontuarios := make([]models.Prontuario, 0, len(r.prontuarios))
	for _, p := range r.prontuarios {
		prontuarios = append(prontuarios, p)
	}
	return prontuarios, nil
}

func (r *MemoryRepository) UpdateProntuario(prontuario models.Prontuario) (models.Prontuario, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	oldProntuario, ok := r.prontuarios[prontuario.ID]
	if !ok {
		return models.Prontuario{}, models.ErrNotFound
	}

	prontuario.CriadoEm = oldProntuario.CriadoEm
	prontuario.AtualizadoEm = time.Now()
	r.prontuarios[prontuario.ID] = prontuario
	return prontuario, nil
}

func (r *MemoryRepository) DeleteProntuario(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.prontuarios[id]; !ok {
		return models.ErrNotFound
	}
	delete(r.prontuarios, id)
	return nil
}



# SGHSS Backend API

Sistema de Gestão Hospitalar e de Serviços de Saúde - Backend desenvolvido em Golang

## Descrição

Este projeto implementa a API backend para o Sistema de Gestão Hospitalar e de Serviços de Saúde (SGHSS) da instituição VidaPlus. O sistema gerencia pacientes, profissionais de saúde, consultas, prontuários eletrônicos e inclui funcionalidades completas de autenticação e autorização.

## Tecnologias Utilizadas

- **Linguagem:** Go (Golang) 1.23.0
- **Framework Web:** Gorilla Mux
- **Autenticação:** JWT (JSON Web Tokens)
- **Criptografia:** bcrypt para hash de senhas
- **Banco de Dados:** In-memory (mapas e slices)
- **CORS:** rs/cors

## Arquitetura

O projeto segue uma arquitetura em camadas:

```
├── main.go                 # Ponto de entrada da aplicação
├── models/                 # Estruturas de dados (structs)
├── repository/             # Camada de persistência (in-memory)
├── services/               # Lógica de negócio
├── handlers/               # Controladores HTTP
├── middleware/             # Middlewares (autenticação, etc.)
├── utils/                  # Utilitários (JWT, bcrypt)
├── tests/                  # Testes automatizados
└── diagrams/               # Diagramas UML do projeto
```

## Instalação e Execução

### Pré-requisitos

- Go 1.21+ instalado
- Git (opcional, para clonar o repositório)

### Passos para executar

1. **Clone ou baixe o projeto:**
   ```bash
   git clone <url-do-repositorio>
   cd sghss-backend
   ```

2. **Instale as dependências:**
   ```bash
   go mod tidy
   ```

3. **Compile o projeto:**
   ```bash
   go build -o sghss-backend
   ```

4. **Execute o servidor:**
   ```bash
   ./sghss-backend
   ```

5. **Acesse a API:**
   - URL base: `http://localhost:8080`
   - Status da API: `GET http://localhost:8080/`

## Endpoints da API

### Autenticação (Rotas Públicas)

#### POST /auth/signup
Cadastra um novo usuário no sistema.

**Requisição:**
```json
{
  "email": "admin@vidaplus.com",
  "senha": "senha123",
  "perfil": "admin"
}
```

**Perfis válidos:** `admin`, `medico`, `enfermeiro`, `tecnico`

**Resposta (201 Created):**
```json
{
  "message": "Usuário criado com sucesso",
  "usuario": {
    "id": 1,
    "email": "admin@vidaplus.com",
    "perfil": "admin",
    "criado_em": "2025-06-14T03:30:00Z"
  }
}
```

#### POST /auth/login
Autentica um usuário e retorna um token JWT.

**Requisição:**
```json
{
  "email": "admin@vidaplus.com",
  "senha": "senha123"
}
```

**Resposta (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "usuario": {
    "id": 1,
    "email": "admin@vidaplus.com",
    "perfil": "admin",
    "criado_em": "2025-06-14T03:30:00Z"
  }
}
```

### Pacientes (Rotas Protegidas)

**Autenticação necessária:** Todas as rotas de pacientes requerem o header `Authorization: Bearer <token>`

#### GET /api/pacientes
Lista todos os pacientes cadastrados.

**Resposta (200 OK):**
```json
{
  "pacientes": [
    {
      "id": 1,
      "nome": "João Silva",
      "cpf": "12345678900",
      "data_nascimento": "1990-01-15T00:00:00Z",
      "historico_clinico": "Paciente com histórico de hipertensão",
      "criado_em": "2025-06-14T03:30:00Z",
      "atualizado_em": "2025-06-14T03:30:00Z"
    }
  ],
  "total": 1
}
```

#### POST /api/pacientes
Cria um novo paciente.

**Requisição:**
```json
{
  "nome": "João Silva",
  "cpf": "12345678900",
  "data_nascimento": "1990-01-15T00:00:00Z",
  "historico_clinico": "Paciente com histórico de hipertensão"
}
```

**Validações:**
- CPF deve ser único no sistema
- Nome é obrigatório
- CPF é obrigatório

**Resposta (201 Created):**
```json
{
  "message": "Paciente criado com sucesso",
  "paciente": {
    "id": 1,
    "nome": "João Silva",
    "cpf": "12345678900",
    "data_nascimento": "1990-01-15T00:00:00Z",
    "historico_clinico": "Paciente com histórico de hipertensão",
    "criado_em": "2025-06-14T03:30:00Z",
    "atualizado_em": "2025-06-14T03:30:00Z"
  }
}
```

#### GET /api/pacientes/{id}
Busca um paciente específico por ID.

#### PUT /api/pacientes/{id}
Atualiza um paciente existente.

#### DELETE /api/pacientes/{id}
Remove um paciente do sistema.

### Profissionais de Saúde (Rotas Protegidas)

**Autenticação necessária:** Todas as rotas de profissionais requerem o header `Authorization: Bearer <token>`

#### GET /api/profissionais
Lista todos os profissionais cadastrados.

#### POST /api/profissionais
Cria um novo profissional.

**Requisição:**
```json
{
  "nome": "Dr. Maria Oliveira",
  "crm_coren": "CRM-SP 123456",
  "especialidade": "Cardiologia",
  "perfil_usuario_id": 2
}
```

**Validações:**
- CRM/COREN deve ser único no sistema
- Nome é obrigatório
- CRM/COREN é obrigatório

#### GET /api/profissionais/{id}
Busca um profissional específico por ID.

#### PUT /api/profissionais/{id}
Atualiza um profissional existente.

#### DELETE /api/profissionais/{id}
Remove um profissional do sistema.

### Consultas (Rotas Protegidas)

**Autenticação necessária:** Todas as rotas de consultas requerem o header `Authorization: Bearer <token>`

#### GET /api/consultas
Lista todas as consultas agendadas.

**Resposta (200 OK):**
```json
{
  "consultas": [
    {
      "id": 1,
      "paciente_id": 1,
      "profissional_id": 1,
      "data_hora": "2025-09-01T10:00:00Z",
      "especialidade": "Cardiologia",
      "status": "agendada",
      "observacoes": "Primeira consulta",
      "criado_em": "2025-08-26T10:00:00Z",
      "atualizado_em": "2025-08-26T10:00:00Z"
    }
  ],
  "total": 1
}
```

#### POST /api/consultas
Agenda uma nova consulta.

**Requisição:**
```json
{
  "paciente_id": 1,
  "profissional_id": 1,
  "data_hora": "2025-09-01T10:00:00Z",
  "especialidade": "Cardiologia",
  "status": "agendada",
  "observacoes": "Primeira consulta"
}
```

**Status válidos:** `agendada`, `realizada`, `cancelada`

#### GET /api/consultas/{id}
Busca uma consulta específica por ID.

#### PUT /api/consultas/{id}
Atualiza uma consulta existente.

#### DELETE /api/consultas/{id}
Cancela/remove uma consulta do sistema.

### Prontuários Eletrônicos (Rotas Protegidas)

**Autenticação necessária:** Todas as rotas de prontuários requerem o header `Authorization: Bearer <token>`

#### GET /api/prontuarios
Lista todos os prontuários eletrônicos.

**Resposta (200 OK):**
```json
{
  "prontuarios": [
    {
      "id": 1,
      "paciente_id": 1,
      "profissional_id": 1,
      "data_atendimento": "2025-09-01T11:00:00Z",
      "diagnostico": "Hipertensão arterial",
      "tratamento": "Medicação anti-hipertensiva",
      "medicamentos": "Losartana 50mg",
      "observacoes": "Paciente respondendo bem ao tratamento",
      "criado_em": "2025-08-26T11:00:00Z",
      "atualizado_em": "2025-08-26T11:00:00Z"
    }
  ],
  "total": 1
}
```

#### POST /api/prontuarios
Cria um novo prontuário eletrônico.

**Requisição:**
```json
{
  "paciente_id": 1,
  "profissional_id": 1,
  "data_atendimento": "2025-09-01T11:00:00Z",
  "diagnostico": "Hipertensão arterial",
  "tratamento": "Medicação anti-hipertensiva",
  "medicamentos": "Losartana 50mg",
  "observacoes": "Paciente respondendo bem ao tratamento"
}
```

#### GET /api/prontuarios/{id}
Busca um prontuário específico por ID.

#### PUT /api/prontuarios/{id}
Atualiza um prontuário existente.

#### DELETE /api/prontuarios/{id}
Remove um prontuário do sistema.

## Códigos de Status HTTP

- **200 OK:** Requisição bem-sucedida
- **201 Created:** Recurso criado com sucesso
- **400 Bad Request:** Dados inválidos na requisição
- **401 Unauthorized:** Token de autenticação inválido ou ausente
- **403 Forbidden:** Acesso negado (permissões insuficientes)
- **404 Not Found:** Recurso não encontrado
- **500 Internal Server Error:** Erro interno do servidor

## Segurança e Validações

### Autenticação JWT
- Tokens JWT são válidos por 24 horas
- Chave secreta configurada para assinatura dos tokens
- Middleware de autenticação protege rotas sensíveis

### Hash de Senhas
- Senhas são criptografadas usando bcrypt
- Senhas nunca são armazenadas em texto plano
- Validação de força de senha (mínimo 6 caracteres)

### Validações de Integridade
- **CPF único:** Cada paciente deve ter um CPF único no sistema
- **CRM/COREN único:** Cada profissional deve ter um registro único
- **Validação na criação e atualização:** Impede duplicação de documentos
- **Mensagens de erro específicas:** Retorna erros claros sobre violações de unicidade

### CORS
- Configurado para aceitar requisições de qualquer origem
- Headers e métodos HTTP permitidos configurados

## Estrutura do Banco de Dados (In-Memory)

### Usuários
- ID (int, auto-incremento)
- Email (string, único)
- SenhaHash (string, bcrypt)
- Perfil (string: admin, medico, enfermeiro, tecnico)
- CriadoEm (timestamp)

### Pacientes
- ID (int, auto-incremento)
- Nome (string)
- CPF (string, único)
- DataNascimento (timestamp)
- HistoricoClinico (string)
- CriadoEm (timestamp)
- AtualizadoEm (timestamp)

### Profissionais de Saúde
- ID (int, auto-incremento)
- Nome (string)
- CRMCOREN (string, único)
- Especialidade (string)
- PerfilUsuarioID (int, referência ao usuário)
- CriadoEm (timestamp)
- AtualizadoEm (timestamp)

### Consultas
- ID (int, auto-incremento)
- PacienteID (int, referência ao paciente)
- ProfissionalID (int, referência ao profissional)
- DataHora (timestamp)
- Especialidade (string)
- Status (string: agendada, realizada, cancelada)
- Observacoes (string)
- CriadoEm (timestamp)
- AtualizadoEm (timestamp)

### Prontuários Eletrônicos
- ID (int, auto-incremento)
- PacienteID (int, referência ao paciente)
- ProfissionalID (int, referência ao profissional)
- DataAtendimento (timestamp)
- Diagnostico (string)
- Tratamento (string)
- Medicamentos (string)
- Observacoes (string)
- CriadoEm (timestamp)
- AtualizadoEm (timestamp)

## Testes

O projeto inclui scripts de teste automatizados para validar todas as funcionalidades:

### Scripts Disponíveis
- `test_api.sh` - Testes básicos de CRUD para pacientes e profissionais
- `test_new_features.sh` - Testes das funcionalidades de consultas e prontuários (requer jq)
- `test_simple_no_jq.sh` - Testes completos sem dependência do jq
- `test_uniqueness.sh` - Testes de validação de unicidade

### Executar Testes
```bash
# Testes completos sem dependência externa
./test_simple_no_jq.sh

# Testes de unicidade
./test_uniqueness.sh

# Testes básicos
./test_api.sh
```

### Exemplo de teste com curl:

1. **Cadastrar usuário:**
```bash
curl -X POST http://localhost:8080/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@test.com","senha":"123456","perfil":"admin"}'
```

2. **Fazer login:**
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@test.com","senha":"123456"}'
```

3. **Criar paciente (com token):**
```bash
curl -X POST http://localhost:8080/api/pacientes \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer SEU_TOKEN_AQUI" \
  -d '{"nome":"João Silva","cpf":"12345678900","data_nascimento":"1990-01-15T00:00:00Z","historico_clinico":"Teste"}'
```

4. **Agendar consulta:**
```bash
curl -X POST http://localhost:8080/api/consultas \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer SEU_TOKEN_AQUI" \
  -d '{"paciente_id":1,"profissional_id":1,"data_hora":"2025-09-01T10:00:00Z","especialidade":"Cardiologia","status":"agendada","observacoes":"Primeira consulta"}'
```

## Considerações de Desenvolvimento

### Banco de Dados In-Memory
- Os dados são perdidos quando o servidor é reiniciado
- Adequado para desenvolvimento e prototipagem
- Para produção, seria necessário integrar com um banco de dados persistente (PostgreSQL, MySQL, etc.)

### Funcionalidades Implementadas
- ✅ Sistema completo de autenticação e autorização
- ✅ CRUD completo para pacientes e profissionais
- ✅ Sistema de agendamento de consultas
- ✅ Prontuários eletrônicos
- ✅ Validações de integridade de dados
- ✅ Testes automatizados abrangentes
- ✅ Documentação completa da API

### Melhorias Futuras
- Integração com banco de dados persistente
- Implementação de logs estruturados
- Documentação automática com Swagger
- Rate limiting e throttling
- Implementação de soft delete
- Paginação para listagens
- Filtros e busca avançada
- Notificações para pacientes
- Integração com sistemas de telemedicina

## Diagramas UML

O projeto inclui diagramas UML completos na pasta `diagrams/`:
- **Diagrama de Casos de Uso:** Interações entre usuários e funcionalidades
- **Diagrama de Classes:** Estrutura de dados e relacionamentos
- **Diagrama de Fluxo de Testes:** Estratégia de validação do sistema

## Autor

Desenvolvido como projeto acadêmico para a disciplina de Projeto Multidisciplinar - Trilha Backend.

## Licença

Este projeto é desenvolvido para fins educacionais.


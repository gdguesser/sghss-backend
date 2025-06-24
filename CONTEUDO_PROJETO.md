# Conteúdo do Projeto SGHSS Backend

## Arquivos Principais

### Documentação
- `README.md` - Documentação técnica completa da API
- `PLANO_TESTES.md` - Plano de testes com 27 casos de teste
- `PROJETO_FINAL_COMPLETE.md` - Documento acadêmico final
- `PROJETO_FINAL_COMPLETE.pdf` - Versão PDF do documento acadêmico

### Código Fonte
- `main.go` - Ponto de entrada da aplicação
- `go.mod` / `go.sum` - Dependências do projeto
- `sghss-backend` - Executável compilado

### Estrutura do Código
- `models/` - Estruturas de dados (structs)
- `repository/` - Camada de persistência (in-memory)
- `services/` - Lógica de negócio
- `handlers/` - Controladores HTTP
- `middleware/` - Middlewares (autenticação, etc.)
- `utils/` - Utilitários (JWT, bcrypt)

### Testes
- `test_simple_no_jq.sh` - **Script principal de testes** (sem dependências)
- `test_api.sh` - Testes básicos de CRUD
- `test_new_features.sh` - Testes das funcionalidades avançadas (requer jq)
- `test_uniqueness.sh` - Testes de validação de unicidade
- `tests/` - Estrutura para testes unitários em Go

### Diagramas UML
- `diagrams/use_case_diagram.png/.puml` - Diagrama de Casos de Uso
- `diagrams/class_diagram.png/.puml` - Diagrama de Classes
- `diagrams/test_flow_diagram.png/.puml` - Diagrama de Fluxo de Testes

## Como Usar

1. **Compilar e executar:**
   ```bash
   go build -o sghss-backend
   ./sghss-backend
   ```

2. **Executar testes:**
   ```bash
   ./test_simple_no_jq.sh
   ```

3. **Ver documentação:**
   - Abrir `README.md` para documentação técnica
   - Abrir `PROJETO_FINAL_COMPLETE.pdf` para documento acadêmico

## Funcionalidades Implementadas

- ✅ Sistema completo de autenticação JWT
- ✅ CRUD de pacientes com validação de CPF único
- ✅ CRUD de profissionais com validação de CRM/COREN único
- ✅ Sistema de agendamento de consultas
- ✅ Prontuários eletrônicos completos
- ✅ 15 endpoints funcionais
- ✅ 27 casos de teste automatizados
- ✅ Documentação completa e profissional

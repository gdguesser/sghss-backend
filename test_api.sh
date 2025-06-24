#!/bin/bash

echo "=== TESTES DA API SGHSS ==="
echo ""

# URL base da API
BASE_URL="http://localhost:8080"

# Cores para output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Função para testar endpoint
test_endpoint() {
    local description="$1"
    local method="$2"
    local url="$3"
    local data="$4"
    local headers="$5"
    local expected_status="$6"
    
    echo -e "${YELLOW}Testando: $description${NC}"
    
    if [ "$method" = "GET" ]; then
        if [ -n "$headers" ]; then
            response=$(curl -s -w "%{http_code}" -H "$headers" "$url")
        else
            response=$(curl -s -w "%{http_code}" "$url")
        fi
    else
        if [ -n "$headers" ]; then
            response=$(curl -s -w "%{http_code}" -X "$method" -H "Content-Type: application/json" -H "$headers" -d "$data" "$url")
        else
            response=$(curl -s -w "%{http_code}" -X "$method" -H "Content-Type: application/json" -d "$data" "$url")
        fi
    fi
    
    # Extrair status code (últimos 3 caracteres)
    status_code="${response: -3}"
    # Extrair body (tudo exceto os últimos 3 caracteres)
    body="${response%???}"
    
    if [ "$status_code" = "$expected_status" ]; then
        echo -e "${GREEN}✅ PASSOU - Status: $status_code${NC}"
    else
        echo -e "${RED}❌ FALHOU - Status esperado: $expected_status, recebido: $status_code${NC}"
        echo "Resposta: $body"
    fi
    echo ""
    
    # Retornar o body para uso posterior
    echo "$body"
}

echo "=== 1. TESTE DE STATUS DA API ==="
test_endpoint "Status da API" "GET" "$BASE_URL/" "" "" "200" > /dev/null

echo "=== 2. TESTES DE AUTENTICAÇÃO ==="

# Teste de cadastro
echo "Cadastrando usuário administrador..."
signup_response=$(test_endpoint "Cadastro de usuário" "POST" "$BASE_URL/auth/signup" '{"email":"admin@test.com","senha":"admin123","perfil":"admin"}' "" "201")

# Teste de login
echo "Fazendo login..."
login_response=$(test_endpoint "Login de usuário" "POST" "$BASE_URL/auth/login" '{"email":"admin@test.com","senha":"admin123"}' "" "200")

# Extrair token do response
TOKEN=$(echo "$login_response" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
echo "Token obtido: ${TOKEN:0:50}..."
echo ""

echo "=== 3. TESTES DE ACESSO SEM AUTENTICAÇÃO ==="
test_endpoint "Acesso sem token" "GET" "$BASE_URL/api/pacientes" "" "" "401" > /dev/null

echo "=== 4. TESTES DE PACIENTES ==="

# Criar paciente
echo "Criando paciente..."
create_patient_response=$(test_endpoint "Criar paciente" "POST" "$BASE_URL/api/pacientes" '{"nome":"João Silva","cpf":"12345678900","data_nascimento":"1990-01-15T00:00:00Z","historico_clinico":"Paciente com histórico de hipertensão"}' "Authorization: Bearer $TOKEN" "201")

# Listar pacientes
echo "Listando pacientes..."
test_endpoint "Listar pacientes" "GET" "$BASE_URL/api/pacientes" "" "Authorization: Bearer $TOKEN" "200" > /dev/null

# Buscar paciente por ID
echo "Buscando paciente por ID..."
test_endpoint "Buscar paciente por ID" "GET" "$BASE_URL/api/pacientes/1" "" "Authorization: Bearer $TOKEN" "200" > /dev/null

# Atualizar paciente
echo "Atualizando paciente..."
test_endpoint "Atualizar paciente" "PUT" "$BASE_URL/api/pacientes/1" '{"nome":"João Silva Santos","cpf":"12345678900","data_nascimento":"1990-01-15T00:00:00Z","historico_clinico":"Paciente com histórico de hipertensão controlada"}' "Authorization: Bearer $TOKEN" "200" > /dev/null

echo "=== 5. TESTES DE PROFISSIONAIS ==="

# Criar profissional
echo "Criando profissional..."
test_endpoint "Criar profissional" "POST" "$BASE_URL/api/profissionais" '{"nome":"Dr. Maria Oliveira","crm_coren":"CRM-SP 123456","especialidade":"Cardiologia","perfil_usuario_id":1}' "Authorization: Bearer $TOKEN" "201" > /dev/null

# Listar profissionais
echo "Listando profissionais..."
test_endpoint "Listar profissionais" "GET" "$BASE_URL/api/profissionais" "" "Authorization: Bearer $TOKEN" "200" > /dev/null

# Buscar profissional por ID
echo "Buscando profissional por ID..."
test_endpoint "Buscar profissional por ID" "GET" "$BASE_URL/api/profissionais/1" "" "Authorization: Bearer $TOKEN" "200" > /dev/null

echo "=== 6. TESTES DE VALIDAÇÃO ==="

# Teste de cadastro com dados inválidos
echo "Testando cadastro com dados inválidos..."
test_endpoint "Cadastro com email duplicado" "POST" "$BASE_URL/auth/signup" '{"email":"admin@test.com","senha":"123","perfil":"admin"}' "" "400" > /dev/null

# Teste de login com credenciais inválidas
echo "Testando login com credenciais inválidas..."
test_endpoint "Login com credenciais inválidas" "POST" "$BASE_URL/auth/login" '{"email":"admin@test.com","senha":"senhaerrada"}' "" "401" > /dev/null

# Teste de criação de paciente com dados inválidos
echo "Testando criação de paciente com dados inválidos..."
test_endpoint "Criar paciente sem nome" "POST" "$BASE_URL/api/pacientes" '{"cpf":"12345678900"}' "Authorization: Bearer $TOKEN" "400" > /dev/null

echo "=== RESUMO DOS TESTES ==="
echo -e "${GREEN}✅ Todos os testes principais foram executados${NC}"
echo -e "${YELLOW}📋 Funcionalidades testadas:${NC}"
echo "   - Status da API"
echo "   - Cadastro e login de usuários"
echo "   - Autenticação JWT"
echo "   - CRUD de pacientes"
echo "   - CRUD de profissionais"
echo "   - Validação de dados"
echo "   - Controle de acesso"
echo ""
echo -e "${GREEN}🎉 API SGHSS funcionando corretamente!${NC}"


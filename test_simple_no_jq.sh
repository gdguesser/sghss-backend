#!/bin/bash

# URL base da API
BASE_URL="http://localhost:8080/api"
AUTH_URL="http://localhost:8080/auth"

# Cores para saída do terminal
GREEN="\033[0;32m"
RED="\033[0;31m"
NC="\033[0m" # No Color

# Função para imprimir status
print_status() {
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✔ $1${NC}"
    else
        echo -e "${RED}✖ $1${NC}"
        exit 1
    fi
}

echo "Iniciando testes para as novas funcionalidades de Consulta e Prontuário..."

# 1. Registrar um novo usuário (se não existir)
echo -e "\n--- Teste de Registro de Usuário ---"
curl -s -X POST "$AUTH_URL/signup" \
     -H "Content-Type: application/json" \
     -d '{"email": "testuser@example.com", "senha": "password123", "perfil": "medico"}' > /dev/null
print_status "Registro de usuário (testuser@example.com)"

# 2. Fazer login e obter token
echo -e "\n--- Teste de Login ---"
LOGIN_RESPONSE=$(curl -s -X POST "$AUTH_URL/login" \
     -H "Content-Type: application/json" \
     -d '{"email": "testuser@example.com", "senha": "password123"}'
)

# Extrair token usando sed e grep (mais simples)
TOKEN=$(echo "$LOGIN_RESPONSE" | sed -n 's/.*"token":"\([^"]*\)".*/\1/p')

if [ -z "$TOKEN" ]; then
    echo -e "${RED}✖ Falha ao obter token de autenticação.${NC}"
    echo "Resposta do login: $LOGIN_RESPONSE"
    exit 1
fi
print_status "Login e obtenção de token"

# 3. Criar um paciente para a consulta/prontuário
echo -e "\n--- Teste de Criação de Paciente ---"
PACIENTE_RESPONSE=$(curl -s -X POST "$BASE_URL/pacientes" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer $TOKEN" \
     -d '{"nome": "Paciente Teste Consulta", "cpf": "11122233344", "data_nascimento": "1990-01-01T00:00:00Z", "historico_clinico": "Nenhum"}'
)

# Extrair ID do paciente usando sed (busca por "id": seguido de número)
PACIENTE_ID=$(echo "$PACIENTE_RESPONSE" | sed -n 's/.*"id":\([0-9]*\).*/\1/p' | head -1)
print_status "Criação de paciente (ID: $PACIENTE_ID)"

# 4. Criar um profissional para a consulta/prontuário
echo -e "\n--- Teste de Criação de Profissional ---"
PROFISSIONAL_RESPONSE=$(curl -s -X POST "$BASE_URL/profissionais" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer $TOKEN" \
     -d '{"nome": "Dr. Consulta", "crm_coren": "CRM/SP 123456", "especialidade": "Clínico Geral", "perfil_usuario_id": 1}'
)

# Extrair ID do profissional usando sed
PROFISSIONAL_ID=$(echo "$PROFISSIONAL_RESPONSE" | sed -n 's/.*"id":\([0-9]*\).*/\1/p' | head -1)
print_status "Criação de profissional (ID: $PROFISSIONAL_ID)"

# 5. Testar criação de Consulta
echo -e "\n--- Teste: Criar Consulta ---"
CONSULTA_RESPONSE=$(curl -s -X POST "$BASE_URL/consultas" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer $TOKEN" \
     -d '{"paciente_id": '$PACIENTE_ID', "profissional_id": '$PROFISSIONAL_ID', "data_hora": "2025-09-01T10:00:00Z", "especialidade": "Cardiologia", "status": "agendada", "observacoes": "Primeira consulta"}'
)

# Extrair ID da consulta usando sed
CONSULTA_ID=$(echo "$CONSULTA_RESPONSE" | sed -n 's/.*"id":\([0-9]*\).*/\1/p' | head -1)
print_status "Criação de Consulta (ID: $CONSULTA_ID)"

# 6. Testar obtenção de Consulta por ID
echo -e "\n--- Teste: Obter Consulta por ID ---"
curl -s -X GET "$BASE_URL/consultas/$CONSULTA_ID" \
     -H "Authorization: Bearer $TOKEN" > /dev/null
print_status "Obtenção de Consulta por ID"

# 7. Testar obtenção de todas as Consultas
echo -e "\n--- Teste: Obter Todas as Consultas ---"
curl -s -X GET "$BASE_URL/consultas" \
     -H "Authorization: Bearer $TOKEN" > /dev/null
print_status "Obtenção de Todas as Consultas"

# 8. Testar atualização de Consulta
echo -e "\n--- Teste: Atualizar Consulta ---"
curl -s -X PUT "$BASE_URL/consultas/$CONSULTA_ID" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer $TOKEN" \
     -d '{"id": '$CONSULTA_ID', "paciente_id": '$PACIENTE_ID', "profissional_id": '$PROFISSIONAL_ID', "data_hora": "2025-09-01T11:00:00Z", "especialidade": "Cardiologia", "status": "realizada", "observacoes": "Consulta de retorno"}' > /dev/null
print_status "Atualização de Consulta"

# 9. Testar criação de Prontuário
echo -e "\n--- Teste: Criar Prontuário ---"
PRONTUARIO_RESPONSE=$(curl -s -X POST "$BASE_URL/prontuarios" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer $TOKEN" \
     -d '{"paciente_id": '$PACIENTE_ID', "profissional_id": '$PROFISSIONAL_ID', "data_atendimento": "2025-09-01T11:00:00Z", "diagnostico": "Gripe", "tratamento": "Repouso e hidratação", "medicamentos": "Paracetamol", "observacoes": "Paciente com sintomas leves"}'
)

# Extrair ID do prontuário usando sed
PRONTUARIO_ID=$(echo "$PRONTUARIO_RESPONSE" | sed -n 's/.*"id":\([0-9]*\).*/\1/p' | head -1)
print_status "Criação de Prontuário (ID: $PRONTUARIO_ID)"

# 10. Testar obtenção de Prontuário por ID
echo -e "\n--- Teste: Obter Prontuário por ID ---"
curl -s -X GET "$BASE_URL/prontuarios/$PRONTUARIO_ID" \
     -H "Authorization: Bearer $TOKEN" > /dev/null
print_status "Obtenção de Prontuário por ID"

# 11. Testar obtenção de todos os Prontuários
echo -e "\n--- Teste: Obter Todos os Prontuários ---"
curl -s -X GET "$BASE_URL/prontuarios" \
     -H "Authorization: Bearer $TOKEN" > /dev/null
print_status "Obtenção de Todos os Prontuários"

# 12. Testar atualização de Prontuário
echo -e "\n--- Teste: Atualizar Prontuário ---"
curl -s -X PUT "$BASE_URL/prontuarios/$PRONTUARIO_ID" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer $TOKEN" \
     -d '{"id": '$PRONTUARIO_ID', "paciente_id": '$PACIENTE_ID', "profissional_id": '$PROFISSIONAL_ID', "data_atendimento": "2025-09-01T11:00:00Z", "diagnostico": "Gripe (melhora)", "tratamento": "Repouso e hidratação", "medicamentos": "Paracetamol", "observacoes": "Paciente em recuperação"}' > /dev/null
print_status "Atualização de Prontuário"

# 13. Testar exclusão de Consulta
echo -e "\n--- Teste: Excluir Consulta ---"
curl -s -X DELETE "$BASE_URL/consultas/$CONSULTA_ID" \
     -H "Authorization: Bearer $TOKEN" > /dev/null
print_status "Exclusão de Consulta"

# 14. Testar exclusão de Prontuário
echo -e "\n--- Teste: Excluir Prontuário ---"
curl -s -X DELETE "$BASE_URL/prontuarios/$PRONTUARIO_ID" \
     -H "Authorization: Bearer $TOKEN" > /dev/null
print_status "Exclusão de Prontuário"

# 15. Limpar dados de teste (opcional, mas boa prática)
echo -e "\n--- Limpeza de Dados de Teste ---"
curl -s -X DELETE "$BASE_URL/pacientes/$PACIENTE_ID" \
     -H "Authorization: Bearer $TOKEN" > /dev/null
print_status "Exclusão de Paciente de Teste"

curl -s -X DELETE "$BASE_URL/profissionais/$PROFISSIONAL_ID" \
     -H "Authorization: Bearer $TOKEN" > /dev/null
print_status "Exclusão de Profissional de Teste"

echo -e "\nTodos os testes de Consulta e Prontuário concluídos."


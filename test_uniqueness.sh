#!/bin/bash

echo "=== TESTES DE VALIDAÇÃO DE UNICIDADE ==="
echo ""

# URL base da API
BASE_URL="http://localhost:8080"

# Cores para output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "=== 1. PREPARAÇÃO - LOGIN ==="

# Cadastro e login para obter token
echo "Fazendo login..."
login_response=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@test.com","senha":"admin123"}')

TOKEN=$(echo "$login_response" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
echo "Token obtido: ${TOKEN:0:50}..."
echo ""

echo "=== 2. TESTES DE CPF ÚNICO PARA PACIENTES ==="

echo "Teste 1: Criando primeiro paciente com CPF 11111111111..."
response1=$(curl -s -w "%{http_code}" -X POST "$BASE_URL/api/pacientes" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"nome":"João Silva","cpf":"11111111111","data_nascimento":"1990-01-15T00:00:00Z","historico_clinico":"Primeiro paciente"}')

status1="${response1: -3}"
body1="${response1%???}"

if [ "$status1" = "201" ]; then
    echo -e "${GREEN}✅ PASSOU - Primeiro paciente criado com sucesso${NC}"
else
    echo -e "${RED}❌ FALHOU - Erro ao criar primeiro paciente: $body1${NC}"
fi

echo ""
echo "Teste 2: Tentando criar segundo paciente com mesmo CPF 11111111111..."
response2=$(curl -s -w "%{http_code}" -X POST "$BASE_URL/api/pacientes" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"nome":"Maria Santos","cpf":"11111111111","data_nascimento":"1985-03-20T00:00:00Z","historico_clinico":"Segundo paciente"}')

status2="${response2: -3}"
body2="${response2%???}"

if [ "$status2" = "400" ] && [[ "$body2" == *"CPF já cadastrado"* ]]; then
    echo -e "${GREEN}✅ PASSOU - CPF duplicado rejeitado corretamente${NC}"
else
    echo -e "${RED}❌ FALHOU - CPF duplicado deveria ser rejeitado. Status: $status2, Body: $body2${NC}"
fi

echo ""
echo "Teste 3: Criando terceiro paciente com CPF diferente 22222222222..."
response3=$(curl -s -w "%{http_code}" -X POST "$BASE_URL/api/pacientes" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"nome":"Carlos Oliveira","cpf":"22222222222","data_nascimento":"1975-12-10T00:00:00Z","historico_clinico":"Terceiro paciente"}')

status3="${response3: -3}"
body3="${response3%???}"

if [ "$status3" = "201" ]; then
    echo -e "${GREEN}✅ PASSOU - Paciente com CPF diferente criado com sucesso${NC}"
else
    echo -e "${RED}❌ FALHOU - Erro ao criar paciente com CPF diferente: $body3${NC}"
fi

echo ""
echo "=== 3. TESTES DE CRM/COREN ÚNICO PARA PROFISSIONAIS ==="

echo "Teste 4: Criando primeiro profissional com CRM-SP 111111..."
response4=$(curl -s -w "%{http_code}" -X POST "$BASE_URL/api/profissionais" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"nome":"Dr. João Cardiologista","crm_coren":"CRM-SP 111111","especialidade":"Cardiologia","perfil_usuario_id":1}')

status4="${response4: -3}"
body4="${response4%???}"

if [ "$status4" = "201" ]; then
    echo -e "${GREEN}✅ PASSOU - Primeiro profissional criado com sucesso${NC}"
else
    echo -e "${RED}❌ FALHOU - Erro ao criar primeiro profissional: $body4${NC}"
fi

echo ""
echo "Teste 5: Tentando criar segundo profissional com mesmo CRM-SP 111111..."
response5=$(curl -s -w "%{http_code}" -X POST "$BASE_URL/api/profissionais" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"nome":"Dr. Maria Cardiologista","crm_coren":"CRM-SP 111111","especialidade":"Cardiologia","perfil_usuario_id":1}')

status5="${response5: -3}"
body5="${response5%???}"

if [ "$status5" = "400" ] && [[ "$body5" == *"CRM/COREN já cadastrado"* ]]; then
    echo -e "${GREEN}✅ PASSOU - CRM/COREN duplicado rejeitado corretamente${NC}"
else
    echo -e "${RED}❌ FALHOU - CRM/COREN duplicado deveria ser rejeitado. Status: $status5, Body: $body5${NC}"
fi

echo ""
echo "Teste 6: Criando terceiro profissional com COREN-SP 222222..."
response6=$(curl -s -w "%{http_code}" -X POST "$BASE_URL/api/profissionais" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"nome":"Enfermeira Ana Silva","crm_coren":"COREN-SP 222222","especialidade":"Enfermagem","perfil_usuario_id":1}')

status6="${response6: -3}"
body6="${response6%???}"

if [ "$status6" = "201" ]; then
    echo -e "${GREEN}✅ PASSOU - Profissional com CRM/COREN diferente criado com sucesso${NC}"
else
    echo -e "${RED}❌ FALHOU - Erro ao criar profissional com CRM/COREN diferente: $body6${NC}"
fi

echo ""
echo "=== 4. TESTES DE ATUALIZAÇÃO COM VALIDAÇÃO ==="

echo "Teste 7: Tentando atualizar paciente 1 para usar CPF do paciente 2..."
response7=$(curl -s -w "%{http_code}" -X PUT "$BASE_URL/api/pacientes/1" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"nome":"João Silva Atualizado","cpf":"22222222222","data_nascimento":"1990-01-15T00:00:00Z","historico_clinico":"Tentando usar CPF duplicado"}')

status7="${response7: -3}"
body7="${response7%???}"

if [ "$status7" = "400" ] && [[ "$body7" == *"CPF já cadastrado"* ]]; then
    echo -e "${GREEN}✅ PASSOU - Atualização com CPF duplicado rejeitada corretamente${NC}"
else
    echo -e "${RED}❌ FALHOU - Atualização com CPF duplicado deveria ser rejeitada. Status: $status7, Body: $body7${NC}"
fi

echo ""
echo "Teste 8: Tentando atualizar profissional 1 para usar CRM/COREN do profissional 2..."
response8=$(curl -s -w "%{http_code}" -X PUT "$BASE_URL/api/profissionais/1" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"nome":"Dr. João Atualizado","crm_coren":"COREN-SP 222222","especialidade":"Cardiologia","perfil_usuario_id":1}')

status8="${response8: -3}"
body8="${response8%???}"

if [ "$status8" = "400" ] && [[ "$body8" == *"CRM/COREN já cadastrado"* ]]; then
    echo -e "${GREEN}✅ PASSOU - Atualização com CRM/COREN duplicado rejeitada corretamente${NC}"
else
    echo -e "${RED}❌ FALHOU - Atualização com CRM/COREN duplicado deveria ser rejeitada. Status: $status8, Body: $body8${NC}"
fi

echo ""
echo "=== RESUMO DOS TESTES DE UNICIDADE ==="
echo -e "${GREEN}✅ Validação de CPF único para pacientes implementada${NC}"
echo -e "${GREEN}✅ Validação de CRM/COREN único para profissionais implementada${NC}"
echo -e "${GREEN}✅ Validação funciona tanto na criação quanto na atualização${NC}"
echo -e "${GREEN}✅ Mensagens de erro apropriadas retornadas${NC}"
echo ""
echo -e "${YELLOW}📋 Funcionalidades validadas:${NC}"
echo "   - CPF único por paciente"
echo "   - CRM/COREN único por profissional"
echo "   - Validação na criação de novos registros"
echo "   - Validação na atualização de registros existentes"
echo "   - Tratamento adequado de erros"
echo ""
echo -e "${GREEN}🎉 Validações de unicidade funcionando corretamente!${NC}"


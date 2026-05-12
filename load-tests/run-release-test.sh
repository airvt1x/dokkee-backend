#!/bin/bash
set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}=== Запуск нагрузочного тестирования Dokkee ===${NC}"
echo -e "${YELLOW}Сценарий: 100 пользователей одновременно загружают документы${NC}"
echo ""

BACKEND_URL="${BACKEND_URL:-http://localhost:8001}"
echo -e "Проверка бэкенда по адресу: ${BACKEND_URL}"

if ! curl -s -o /dev/null --connect-timeout 2 "${BACKEND_URL}" ; then
    echo -e "${RED}ОШИБКА: Бэкенд недоступен по адресу ${BACKEND_URL}${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Бэкенд доступен${NC}"

if ! command -v k6 &> /dev/null; then
    echo -e "${RED}ОШИБКА: k6 не установлен${NC}"
    exit 1
fi
echo -e "${GREEN}✓ k6 установлен${NC}"

REPORT_DIR="./reports"
mkdir -p "${REPORT_DIR}"

TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
REPORT_FILE="${REPORT_DIR}/release-test-${TIMESTAMP}.json"

echo ""
echo -e "${YELLOW}Запуск теста...${NC}"
echo ""

k6 run \
    -e API_URL="${BACKEND_URL}" \
    --vus 100 \
    --duration 5m \
    --out json="${REPORT_FILE}" \
    ./scenarios/release-test.js

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}=== Тест успешно завершён ===${NC}"
    echo -e "Результаты сохранены в: ${REPORT_FILE}"
else
    echo -e "${RED}=== Тест провалился ===${NC}"
    exit 1
fi

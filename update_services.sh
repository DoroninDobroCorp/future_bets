#!/bin/bash

# --- Настройки ---
BACKEND_DIR="livebet_backend-main" # Название директории с бэкендом

# --- Цвета для красивого вывода ---
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# --- Проверка, что мы находимся в правильной директории ---
if [ ! -d "$BACKEND_DIR" ]; then
    echo -e "${RED}ОШИБКА: Скрипт должен быть запущен из директории 'future_bets', содержащей папку '$BACKEND_DIR'.${NC}"
    exit 1
fi

# --- Функция для вывода заголовков ---
print_header() {
    echo -e "\n${BLUE}--- $1 ---${NC}"
}

# --- Функция для выполнения и логирования команды ---
# Принимает директорию и саму команду как аргументы
run_command() {
    local dir=$1
    shift # Сдвигаем аргументы, чтобы $@ содержал только команду
    
    echo -e "${YELLOW}В директории '$dir' ВЫПОЛНЯЕТСЯ:${NC} $@"
    (cd "$dir" && "$@")
    
    if [ $? -ne 0 ]; then
        echo -e "${RED}ОШИБКА: Команда провалилась: '$@' в директории '$dir'${NC}"
        exit 1
    fi
}

# --- Основная логика ---

# 1. Пересобираем и перезапускаем Analyzer
print_header "Обновление сервиса Analyzer"
ANALYZER_PATH="$BACKEND_DIR/analyzer"
run_command "$ANALYZER_PATH" docker compose down
run_command "$ANALYZER_PATH" docker compose up -d --build

# 2. Перезапускаем зависимые сервисы
print_header "Перезапуск зависимых сервисов: Calculator и Auto-Matcher"
CALCULATOR_PATH="$BACKEND_DIR/calculator"
run_command "$CALCULATOR_PATH" docker compose restart

AUTO_MATCHER_PATH="$BACKEND_DIR/auto_matcher"
run_command "$AUTO_MATCHER_PATH" docker compose restart

# 3. Уведомление о парсерах
print_header "Информация по парсерам"
echo -e "${GREEN}Для обновления кода в парсерах используйте панель управления в Runner (вкл/выкл).${NC}"
echo -e "${GREEN}Этот скрипт не трогает парсеры, чтобы не сбивать их работу.${NC}"

# 4. Показываем логи обновленного Analyzer
print_header "Отображение логов Analyzer (последние 50 строк)"
docker logs --tail 50 analyzer

echo -e "\n${GREEN}✅ Все сервисы успешно обновлены и перезапущены!${NC}"
echo -e "${YELLOW}Для просмотра логов Analyzer в реальном времени, выполните: 'docker logs -f analyzer'${NC}\n"
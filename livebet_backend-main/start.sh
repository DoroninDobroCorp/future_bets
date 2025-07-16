#!/bin/bash
echo "--- Запускаем ВСЕ бэкенд-сервисы и инфраструктуру ---"

# Список всех директорий с docker-compose.yaml
SERVICES=(
    "dev" "analyzer" "auto_matcher" "calculator" "monitoring" "results"
    "statistic" "tg_livebot" "tg_manager" "tg_testbot" "runner"
)

for service_dir in "${SERVICES[@]}"; do
    echo "--> Запускаем сервис в папке: $service_dir"
    (cd "$service_dir" && docker compose up -d)
done

echo ""
echo "--- Бэкенд запущен! ---"
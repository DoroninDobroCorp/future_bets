#!/bin/bash
echo "--- Останавливаем все бэкенд-сервисы ---"

SERVICES=(
    "dev" "analyzer" "auto_matcher" "calculator" "monitoring" "results"
    "statistic" "tg_livebot" "tg_manager" "tg_testbot" "runner"
)

for service_dir in "${SERVICES[@]}"; do
    echo "--> Останавливаем $service_dir"
    (cd "$service_dir" && docker compose down)
done

echo "--- Бэкенд остановлен ---"
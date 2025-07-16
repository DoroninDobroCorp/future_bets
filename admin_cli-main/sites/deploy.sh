#!/bin/bash

# 1. Проверяем синтаксис Apache
echo "Проверка конфигурации Apache..."
if ! sudo apache2ctl -t; then
    echo
    echo "Ошибка в конфигурации Apache!"
    exit 1
fi

# 2. Перезагружаем Apache
echo "Перезагрузка Apache..."
sudo systemctl reload apache2

# 3. Проверяем результат
if [ $? -eq 0 ]; then
    systemctl status apache2
    echo
    echo "Successfully Apache reload!"
 else
    echo "Apache reload error!"
    exit 1
fi
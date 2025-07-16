# Используем Python 3.11 на базе Alpine
FROM python:3.11-alpine

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файл requirements.txt в рабочую директорию
COPY statistic/requirements.txt .

# Устанавливаем зависимости из файла requirements.txt
RUN pip install --no-cache-dir -r requirements.txt

# Копируем остальные файлы приложения в контейнер
COPY statistic/ .

CMD ["python", "main.py"]

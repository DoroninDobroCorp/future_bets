README: Проект "Future of Bets"
1. Общее описание проекта
Future of Bets — это комплексная платформа для value-беттинга (ставок с перевесом). Система автоматизирует поиск и анализ букмекерских линий, находит выгодные ставки ("вилки" или "коридоры") и предоставляет операторам интерфейс для их обработки и размещения.
Проект состоит из двух основных частей:
livebet_backend-main: Бэкенд, реализованный на микросервисной архитектуре с использованием Go и Python.
admin_cli-main: Фронтенд (административная панель) для операторов, разработанный на Next.js.
2. Архитектура и компоненты
2.1. Фронтенд: admin_cli-main
Это клиентское приложение, с которым взаимодействуют операторы.
Технологии: Next.js, TypeScript, Effector.js (для управления состоянием), Axios, CSS Modules.
Назначение: Предоставление пользовательского интерфейса для мониторинга, анализа и размещения ставок.
Ключевые директории:
pages/: Точки входа и маршрутизация. Основная страница — /client.
components/: UI-компоненты (Calculator, Match, ClientMatchFilter и т.д.).
stores/: Модели состояния (stores) на Effector.js. Хранят данные о букмекерах, лигах, парах матчей и т.д.
services/: Слой для взаимодействия с API бэкенда.
http/: Инстансы Axios для каждого микросервиса бэкенда.
public/html/: Статические HTML-страницы, включая investments.html (презентация проекта).
2.2. Бэкенд: livebet_backend-main (Микросервисы)
Бэкенд состоит из набора независимых сервисов, каждый из которых решает свою задачу. Они общаются между собой через HTTP API и WebSocket.
Сервис	Язык	Назначение	API/Порты (известные)	Взаимодействует с
analyzer	Go	Ядро системы. Принимает данные от парсеров. Сопоставляет матчи от разных букмекеров, находит "вилки", рассчитывает ROI, маржу. Отправляет найденные пары (forks) на фронтенд.	WS: 7300, 7301<br>HTTP: 7005, 7006	Парсеры (входящие данные), admin_cli (отдает по WS), Postgres
auto_matcher	Go	Ручное и автоматическое (AI) сопоставление лиг и команд между букмекерами.	HTTP: 7001	analyzer (для данных), Postgres
calculator	Go	Рассчитывает рекомендуемый размер ставки для оператора на основе ROI и других параметров.	HTTP: 7010	analyzer (для данных), tg_manager
runner	Go	Диспетчер парсеров. Запускает, останавливает и мониторит состояние парсеров.	HTTP: 9200	Парсеры (управляет)
parse_*	Go	Набор сервисов-парсеров для каждого букмекера (Pinnacle, Lobbet, Fonbet и т.д.).	Разные порты (9000-9112)	analyzer (отправляет данные)
tg_manager	Python	API-шлюз и бот-менеджер. Управляет сотрудниками, их доступами, логами, принимает подтвержденные ставки от admin_cli.	HTTP: 7020	Postgres, admin_cli, calculator
tg_livebot	Go	Telegram-бот, который сканирует директорию /logs/bets и отправляет новые CSV-файлы с отчетами в чаты.	-	Файловая система, Postgres
tg_testbot	Go	Аналогичен tg_livebot, но для тестовых ставок (/logs/testbets).	-	Файловая система, Postgres
results	Go	Сервис постобработки. Определяет исход сделанных ставок (выигрыш/проигрыш), сверяясь с Pinnacle, и обновляет БД.	-	Postgres, Pinnacle API
statistic	Python	Генерирует CSV-отчеты со статистикой на основе данных из БД.	-	Postgres, calculator
monitoring	-	Система мониторинга на базе Prometheus и Grafana.	Grafana: 24600	Собирает метрики со всех сервисов
3. Потоки данных (Data Flow)
3.1. Обнаружение и отображение "вилки" (Live/Prematch)
Запуск: runner запускает сервисы parse_*.
Парсинг: parse_* собирают данные с сайтов букмекеров.
Передача в анализатор: Данные по HTTP/WS отправляются в сервис analyzer.
Анализ: analyzer кэширует данные, сопоставляет матчи по лигам и командам (используя данные из auto_matcher), находит общие исходы и рассчитывает ROI.
Трансляция: analyzer транслирует найденные пары (объекты IPair) по WebSocket на ws://ibet.team/api/analyzer/live или .../prematch.
Отображение: admin_cli на странице /client слушает WebSocket. При получении новых данных, PairStore (Effector) обновляется, и UI перерисовывается, показывая новый Match компонент.
3.2. Процесс размещения и логирования ставки
Инициация: Оператор на admin_cli видит выгодную пару и кликает на исход.
Калькулятор: Открывается компонент Calculator, который делает запрос в сервис calculator для получения рекомендуемой суммы ставки.
Подтверждение: Оператор вводит фактическую сумму, коэффициент, время и нажимает "Отправить ставку".
Логирование: admin_cli отправляет POST-запрос со всеми данными ставки (объект AcceptBet) в tg_manager на эндпоинт /log_bet.
Сохранение: tg_manager записывает информацию о ставке в таблицу calculator.log_bet_accept в PostgreSQL.
Уведомление: tg_manager также может отправлять уведомление в Telegram-группу о новой ставке.
3.2.1 Контракт данных: Парсер -> Анализатор
Все сервисы-парсеры (parse_*) обязаны отправлять данные в сервис analyzer в строго определенном формате JSON. Этот формат соответствует структуре GameData из shared/game-data.go.
Структура Go:
Generated go
type GameData struct {
    Pid        int64  `json:"Pid"`
    LeagueName string `json:"LeagueName"`
    HomeName   string `json:"homeName"`
    AwayName   string `json:"awayName"`
    MatchId    string `json:"MatchId"`
    IsLive     bool   `json:"isLive"`
    HomeScore  float64      `json:"HomeScore"`
    AwayScore  float64      `json:"AwayScore"`
    Periods    []PeriodData `json:"Periods"`
    Source     Parser    `json:"Source"`    // Имя букмекера (e.g., "Lobbet")
    SportName  SportName `json:"SportName"` // Имя спорта (e.g., "Soccer")
    CreatedAt  time.Time `json:"CreatedAt"`
    Raw        interface{} `json:"Raw"`
}

type PeriodData struct {
    Win1x2           Win1x2Struct             `json:"Win1x2"`
    Games            map[string]*Win1x2Struct `json:"Games"`
    Totals           map[string]*WinLessMore  `json:"Totals"`
    Handicap         map[string]*WinHandicap  `json:"Handicap"`
    FirstTeamTotals  map[string]*WinLessMore  `json:"FirstTeamTotals"`
    SecondTeamTotals map[string]*WinLessMore  `json:"SecondTeamTotals"`
}

// ... и другие вложенные структуры (Odd, Win1x2Struct, etc.)
Use code with caution.
Go
Пример объекта JSON, отправляемого парсером:
Generated json
{
  "Pid": 22487824,
  "LeagueName": "SPAIN, LaLiga2",
  "homeName": "Huesca",
  "awayName": "Tenerife",
  "MatchId": "22487824",
  "isLive": true,
  "HomeScore": 1,
  "AwayScore": 0,
  "Periods": [
    {
      "Win1x2": { "Win1": { "value": 1.3 }, "WinNone": { "value": 4.75 }, "Win2": { "value": 13.5 } },
      "Totals": {
        "2.5": { "WinMore": { "value": 2.45 }, "WinLess": { "value": 1.5 } }
      },
      "Handicap": {
         "-1.5": { "Win1": { "value": 2.75 }, "Win2": { "value": 0 } }
      },
      "FirstTeamTotals": {},
      "SecondTeamTotals": {}
    }
  ],
  "Source": "Lobbet",
  "SportName": "Soccer",
  "CreatedAt": "2023-12-19T21:15:00Z"
}
Use code with caution.
Json
Эта структура является единственной точкой соприкосновения между парсерами и анализатором. Любые изменения в ней потребуют обновления всех парсеров.
3.3. Сопоставление сущностей (лиги и команды)
Ручное сопоставление: Оператор на admin_cli (стр. /leagues или /matches) выбирает вид спорта и двух букмекеров.
Запрос данных: Фронтенд запрашивает у auto_matcher список несопоставленных лиг/команд.
Создание пары: Оператор выбирает две сущности и нажимает "Создать пару".
Сохранение: Фронтенд отправляет POST-запрос в auto_matcher, который создает запись о сопоставлении в таблицах analyzer.leagues_merge или analyzer.teams_merge.
Автоматическое сопоставление: Сервис auto_matcher периодически запускает AIMatcherService, который запрашивает данные у onlineMatcherService, формирует промпт для Claude API, получает пары и автоматически их сохраняет.
4. CI/CD и развертывание
Технологии: Docker, Docker Compose, GitHub Actions, Apache (как reverse proxy).
Процесс:
Разработчик пушит изменения в одну из веток (analyzer, tg_manager, admin_cli и т.д.).
GitHub Actions (.github/workflows/*.yaml) запускает соответствующий workflow.
Ворфлоу подключается к серверу по SSH.
На сервере выполняется git pull для нужной ветки.
Запускается docker-compose up -d --build для конкретного измененного сервиса. Это пересобирает образ и перезапускает контейнер.
Сетевая конфигурация:
Все сервисы работают в единой Docker-сети livebets.
Наружу "смотрит" Apache, который настроен как reverse proxy (sites/ibet.team.conf). Он проксирует запросы на нужные порты контейнеров, включая HTTP API и WebSocket.
5. Как работать над этим проектом (инструкция для AI)
Чтобы эффективно работать над задачами, не требуя полного контекста каждый раз, используй следующий алгоритм:
Ознакомься с этой README. Это твой "мозг" и "память" по данному проекту.
Получи от пользователя четкую задачу. Например: "Нужно добавить кнопку 'Отмена' в компонент CheckList".
Проанализируй задачу, используя эту README.
Определи, какие компоненты/сервисы будут затронуты. Для задачи с CheckList это только фронтенд, компонент admin_cli-main/components/CheckList/CheckList.tsx и его .props.ts и .module.css.
Оцени возможные побочные эффекты. Изменение компонента CheckList затронет страницы LeagueCandidatesPageComponent и MatchCandidatesPageComponent, где он используется. Это нужно учесть.
Сформулируй запрос на необходимый код. Запроси у пользователя содержимое только тех файлов, которые необходимы для выполнения задачи.
Пример твоего ответа пользователю: "Для добавления кнопки 'Отмена' мне понадобятся следующие файлы: CheckList.tsx, CheckList.props.ts, CheckList.module.css. Изменения затронут страницы кандидатов лиг и матчей, но не сломают их. Пожалуйста, предоставьте их содержимое."
Выполни задачу на основе предоставленных фрагментов кода и верни результат.
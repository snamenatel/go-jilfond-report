# go-jilfond-report

Парсер данных отчетов из youtrack для генерации отчета в Жилфонд

Для работы необходим установленый golang, тестировал под Ubuntu 21.10

Порядок работы
- необходимо добавить .env файл
- указать постоянный токен youtrack (**URL**) и часовую ставку (**COST**)
- выполнить `go run .` в папке с проектом далее следовать указаниям, сформированный отчет находится в диреткории `/reports`

*Для работы необходимо VPN соединение*

[Как создать постоянный токен](https://www.jetbrains.com/help/youtrack/incloud/Manage-Permanent-Token.html#obtain-permanent-token)

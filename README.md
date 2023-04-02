

# Домашнее задание 1

- Создать скелеты трёх сервисов по описанию АПИ из файла contracts.md
- Структуру проекта сделать с учетом разбиения на слои, бизнес-логику писать отвязанной от реализаций клиентов и хендлеров
- Все хендлеры отвечают просто заглушками
- Сделать удобный враппер для сервера по тому принципу, по которому делали на воркшопе
- Придумать самостоятельно удобный враппер для клиента
- Все межсервисные вызовы выполняются. Если хендлер по описанию из contracts.md должен ходить в другой сервис, он должен у вас это успешно делать в коде.
- Общение сервисов по http-json-rpc
- должны успешно проходить make precommit и make run-all в корневой папке
- Наладить общение с product-service (в хендлере Checkout.listCart). Токен для общения с product-service получить, написав в личку @pav5000

# Домашнее задание 2

Ваша задача перевести всё взаимодействие между вашими сервисами на протокол gRPC. 
То есть взаимодействие по http мы полностью выпиливаем и оставляем только gRPC.  
Для вашего удобства и удобства тьютора в каждом проекте заведите Makefile (если ещё нет) и там укажите полезные команды: 
  генерация кода из proto файла и скачивание нужных зависимостей.

Теперь кратко:

1. Переводим всё взаимодействие на gRPC.
2. В Makefile реализуем команды generate (если есть, что. генерить), vendor-proto (если используете вендоринг)

P. S. Gateway и proto-валидацию прикручивать НЕ нужно.

Ссылка на код из workshop (ветка master):

[https://gitlab.ozon.dev/go/classroom-5/Week-2/workshop](https://gitlab.ozon.dev/go/classroom-5/Week-2/workshop)

# Домашнее задание 3

Задание №3 (https://gitlab.ozon.dev/go/classroom-5/Week-3/Homework):

1) Для каждого сервиса(где необходимо что-то сохранять/брать) поднять отдельную БД в docker-compose.

2) Сделать миграции в каждом сервисе (достаточно папки миграций и скрипта). 
   Создать необходимые таблицы.

3) Реализовать логику репозитория в каждом сервисе. 
   В качестве query builder-а можно использовать любую либу(согласовать индивидуально с тьютором). 
   Рекомендуется https://github.com/Masterminds/squirrel.

4) Драйвер для работы с postgresql: только pgx(версия v4) (pool). 
   В одном из сервисов сделать транзакционность запросов (как на воркшопе).

5) Задание с *: Для каждой БД поднять свой балансировщик (pgbouncer или odyssey, можно и то и то).
   Сервисы ходят не на прямую в БД, а через балансировщик

# Домашнее задание 4

https://gitlab.ozon.dev/go/classroom-5/Week-4/Homework

1) Ускорить Checkout.listCart (т.е. уменьшить время ответа этой ручки). При использовании worker pool запрашивать не более 5 sku одновременно. Worker pool нужно написать самостоятельно. Обязательное требование - читаемость и покрытие кода комментариями

2) Во всем сервисе при общении с Product Service необходимо использовать рейт лимит на клиентской стороне (10 RPS). Допускается использование библиотечных рейт лимитеров

3) Во всех слоях сервиса необходимо прокинуть контекст в интерфейсах, если этого не было сделано ранее

4) Аннулирование заказов старше 10 минут в фоне (в отдельной горутине, рекомендуется применять воркер пул для общения с базой)

Доп. задание на алмазик:

5) Написать собственный рейт-лимитер (читаемый код + комментарии обязательны).

# Домашнее задание 5

Необходимо обеспечить 100% покрытие бизнес-логики ручек ListCart и Purchase. Если вдруг ваши слои до сих пор не изолированны друг от друга через интерфейсы, то стоит этим озаботиться, так как возникнут проблемы с генерацией моков.
В качестве генератора моков можете использовать, что душе угодно, хоть minimock, хоть gomock, хоть что-то другое.

Задание на алмазик: обеспечить покрытие всей бизнес-логики на 80+%. Только бизнес-логики, sql тестировать не нужно, интеграционные тесты не нужны, условный пакет utils также тестить не требуется.

# Домашнее задание 6

Материалы: https://gitlab.ozon.dev/go/classroom-5/Week-6/kafka
ДЗ:
- развернуть кафка кластер в docker-compose
- LOMS пишет в Кафку изменения статусов заказов
- Сервис нотификаций должен их вычитывать и отправлять нотификации об изменениях статуса заказа. (писать в лог)
- Нотификация должна быть доставлена гарантированно и ровно один раз
- Нотификации должны доставляться в правильном порядке(по ключу заказа)

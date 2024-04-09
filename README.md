Запуск программы:

migration-create

Установка переменных окружения:

DBNAME=test;HOST=localhost;PASSWORD=test;PORT=5432

Запуск docker-compose.yaml

Функционал программы:

Программа представляет из себя контроллер пунктов выдачи заказов (прием, выдача и возварт заказов) и мененджер самих пунктов (добавление ПВЗ и список ПВЗ)

Программа работает с флагами и аргументами командной строки, доступные команды и необходимые флаги выведутся при запуске программы без каких либо флагов

Есть 3 режима работы программы:

1) Режим работы с заказами, вся необходимая информация передается через флаги, в качестве базы данных выступает текстовый файл
    
    (В данный момент не поддерживается из за отсутствия обработки упаковки заказов)

3) Режим работы с ПВЗ в интерактивном режиме при указании флага -с=pvz, в качестве базы данных выступает текстовый файл

4) Режим RESTful веб-сервера - работа с ПВЗ,  в качестве базы данных выступает PostgreSQL. Ниже приведены примеры примеры курлов



curl-ы

curl -u user1:password1 -X POST -d '{"name":"pvz1", "address":"pvz1", "contact":"pvz1"}' http://localhost:9000/pvz

curl -u user1:password1 -X POST -d '{"name":"pvz2", "address":"pvz2", "contact":"pvz2"}' http://localhost:9000/pvz

curl -u user1:password1 -X POST -d '{"name":"pvz3", "address":"pvz3", "contact":"pvz3"}' http://localhost:9000/pvz

curl -u user1:password1 -X POST -d '{"name":"pvz4", "address":"pvz4", "contact":"pvz4"}' http://localhost:9000/pvz

curl -u user1:password1 http://localhost:9000/pvz/1

curl -u user1:password1 http://localhost:9000/pvz/2

curl -u user1:password1 -X PUT -d '{"name":"pvz2_change", "address":"pvz2_change", "contact":"pvz2_change"}' http://localhost:9000/pvz/2

curl -u user1:password1 http://localhost:9000/pvz

curl -u user1:password1 -X DELETE http://localhost:9000/pvz/1



**Шаблон Стратегия ([Strategy Pattern](https://ru.wikipedia.org/wiki/Стратегия_(шаблон_проектирования)))**:
    - Используется для обработки различных типов упаковки (`PackingStrategy`), что позволяет легко добавлять новые стратегии для других типов упаковки без изменения основного кода приложения.
    - Каждая стратегия отвечает за проверку определенных условий и изменение стоимости заказа в соответствии с типом упаковки.

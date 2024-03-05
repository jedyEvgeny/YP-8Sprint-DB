Упрощённая реализация сервиса отслеживания посылок.

Функционал:
- регистрация посылки,
- получение списка посылок клиента,
- изменение статуса посылки,
- изменение адреса доставки,
- удаление посылки.
- 
Информация о посылке хранится в БД (SQLite). Посылка имеет три состояния: зарегистрирована, отправлена, доставлена. При регистрации посылки создаётся новая запись в БД.
У только что зарегистрированной должен быть статус «зарегистрирована». Трек-номер посылки равен её идентификатору в таблице.
Если посылка в статусе «зарегистрирована», можно изменить адрес доставки или удалить посылку.
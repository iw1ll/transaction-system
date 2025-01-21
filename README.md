# Transaction System
## Система обработки транзакций платёжной системы

### Запуск:

1. Создать .env по примеру .env.example
2. Собрать проект

    ```
    docker-compose build
    ```
3. Запуск:

    ```
    docker-compose up
    ```

### Описание методов:
-   POST http://localhost:8080/api/send \
Отправка средств с одного из кошельков на указанный кошелек

-   GET http://localhost:8080/api/transactions?count=N \
Возвращает информацию о N последних по времени переводах средств

-   GET http://localhost:8080/api/wallet/{address}/balance \
Возвращает информацию о балансе кошелька 

### Ссылка на проект в GitHub:
https://github.com/iw1ll/transaction-system
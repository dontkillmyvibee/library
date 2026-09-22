# Library

## Все возможности работы с API описаны коллекцией в папке postman

### Books

| Method | Endpoint | Description |
|---|---|---|
| GET | `/books` | Получить список книг с фильтрацией |
| GET | `/books/{id}` | Получить книгу по ID |
| POST | `/books` | Создать книгу |
| PUT | `/books/{id}` | Обновить книгу |
| DELETE | `/books/{id}` | Удалить книгу |

### Authors

| Method | Endpoint | Description |
|---|---|---|
| GET | `/authors` | Получить список авторов |
| GET | `/authors/{id}` | Получить автора по ID |
| POST | `/authors` | Создать автора |
| PUT | `/authors/{id}` | Обновить автора |
| DELETE | `/authors/{id}` | Удалить автора |

### Readers

| Method | Endpoint | Description |
|---|---|---|
| GET | `/readers` | Получить список читателей |
| GET | `/readers/{id}` | Получить читателя по ID |
| POST | `/readers` | Создать читателя |
| PUT | `/readers/{id}` | Обновить читателя |
| DELETE | `/readers/{id}` | Удалить читателя |

### Book lending

| Method | Endpoint | Description |
|---|---|---|
| POST | `/readers/{readerID}/books/{bookID}` | Выдать книгу читателю |
| DELETE | `/readers/{readerID}/books/{bookID}` | Вернуть книгу |
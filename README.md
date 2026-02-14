# url-shortener

создать .env используя пример из env.example

по дефолту в .evn устанавливается переменная STORAGE_TYPE=persistent

для запуска с использованием in-memory
в docker-compose.yml добавить в окружение
service/app/environment:
STORAGE_TYPE=inmemory

запуск: docker-compose up --build



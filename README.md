# BankApp Простое банковское приложение на Go, позволяющее управлять пользователями, счетами и транзакциями. 

## 🛠 Технологии и стек 
- Язык программирования: Go (Golang) 
- API: RESTful 
- Логирование: logrus 

## 🚀 Функциональность 
1. Пользователи (Users) 
- Регистрация нового пользователя 
- Аутентификация (логин) 
2. Получение данных о пользователе 
- Счета (Accounts) 
- Создание нового счета 
- Получение списка счетов пользователя 
- Получение данных о конкретном счете 
- Удаление счета 
3. Транзакции (Transfers) 
- Перевод средств между счетами 
- Просмотр истории транзакций 

## 📌 Описание API 
Openapi спецификация представлена в директории api

## Запуск проекта
Запуск проекта существляется из корневой директории проекта командой:
Доступ к приложению будет по адресу localhost:8080
docker-compose -f .\deploy\compose.yml -p "banking_app" up --build  app

# Shopify clone Backend Platform

Pet-project backend для платформы создания интернет-магазинов с поддержкой multi-tenancy и интеграцией платежной системы.

Это backend-часть проекта. Frontend разрабатывался отдельно и другим человеком.

Авто API документация https://documenter.getpostman.com/view/50487053/2sB3dPQpCY

## Технологии

- **Язык:** Go 1.24.5
- **Базы данных:** PostgreSQL, Redis
- **Основные зависимости:** 
  - Gorilla Mux (роутинг)
  - pgx (драйвер PostgreSQL)
  - go-redis (клиент Redis)
  - JWT для аутентификации

## Архитектура  

Проект построен по принципам Clean Architecture с четким разделением на слои:
- **internal/**
  - handlers/ # HTTP обработчики запросов
  - middleware/ # Цепочка middleware
  - repositories/ # Работа с источниками данных
  - services/ # Бизнес-логика
  - models/ # Доменные модели


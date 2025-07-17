# task-manager

## Создание файла с кредами:

Нужно создать файл .env и скопировать туда содержимое .env.example

## Запуск микросервиса:

```
make up
```

## Восстановление структуры базы данных:

```
make migrations-apply
```

## Сброс базы данных

```
migrations-rollback
make migrations-apply
```

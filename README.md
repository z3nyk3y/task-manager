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
make migrations-rollback
make migrations-apply
```

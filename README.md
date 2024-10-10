## Запуск:
в корневой директории проекта
```shell
docker compose up --build
```

## Документация:
``` 
http://127.0.0.1:<порт указанный в .env>/docs/swagger/
```

## Логи:
С запущенным приложением в контейнере -
```shell 
docker exec -it music_library-app-1 /bin/sh

cd /var/log/app 

cat muslib.log
```   

## Тесты:
```shell
go test ./tests -v
```
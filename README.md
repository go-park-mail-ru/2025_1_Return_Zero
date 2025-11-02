# 2025_1_Return_Zero
Бэкэнд проекта "Spotify" команды "Return Zero"

## Авторы

[Артём Абрамов](https://github.com/zeritonik) - _Frontend-разработчик_

[Иван Сысоев](https://github.com/OlegWhiteRose) - _Frontend-разработчик_

[Роман Инякин](https://github.com/Mockird31) - _Backend-разработчик_

[Александр Гароев](https://github.com/derletzte256) - _Backend-разработчик_

## Менторы

[Алик Нигматуллин](https://github.com/BigBullas) - _Frontend_

[Илья Жиленков](https://github.com/ilyushkaaa) - _Backend_

- UX

## Ссылки

[Фронтенд проекта](https://github.com/frontend-park-mail-ru/2025_1_Return_Zero)

[Сайт](http://returnzero.live/)

## Документация

[Документация API на Swagger](https://returnzero.live/api/v1/docs/)

Для регенерации документации используйте следующие команды:

```bash
make swag
```

## Docker Compose

Для запуска Docker Compose с базами данных используйте следующую команду:
```bash
cd deploy
make deploy
```

Для запуска на сервере в первый раз используйте следующую команду:
```bash
cd deploy
make deploy-prod
```

Далее watchtower будет подхватывать новые версии с docker hub

Как развернуть демо:
---

Используемый образ:  
```bash
p1hub/qilin-crm-api:demo-sdk
```

URL для открытия страницы с кнопкой играть: `/integration/demo/parent/game`

Порты: 
```
[8080, 1443]
```

Режимы запуска
---

Есть 3 режима: `parent`, `qilin`, `dev`. 

Нужно поднять демон для каждого режима. 

Режим `parent`
---

Команда для запуска в docker-compose:
```
"sdk","-c","configs/sdk-parent.yaml","-b",":8080","-d"
```

Прокидываемые URLs:
```
https://tst.qilin.super.com/integration/demo/parent/{pattern}
```
нужно прокинуть в демон в контейнере (на порту 8080):
```
<host>:8080/{pattern}
```
```
https://tst.qilin.super.com/integration/game/iframe 
https://tst.qilin.super.com/integration/game/billing
```
нужно прокинуть соответственно в 
```
<host>:8080/integration/game/iframe
<host>:8080/integration/game/billing
```

Режим `dev`
---

Команда для запуска в docker-compose:
```
"sdk","-c","configs/sdk-dev.yaml","-b",":8080","-d"
```

Прокидываемые URLs:
```
https://tst.qilin.super.com/integration/demo/dev/{pattern}
```
нужно прокинуть в демон в контейнере (на порту 8080):
```
/integration/demo/dev/{pattern} -> :8080/{pattern}
```

Нужно прокинуть запросы к поддомену `demo-game.tst.qilin.super.com` на порт `:1443` контейнера.


Режим `qilin`
---

Команда для запуска в docker-compose:
```
"sdk","-c","configs/sdk.yaml","-b",":8080","-d"
```

Прокидываемые URLs:
```
https://tst.qilin.super.com/integration/demo/qilin/{pattern}
```
нужно прокинуть в демон в контейнере (на порту 8080):
```
:8080/{pattern}
```


Пример файла docker-compose
---

Названия контейнеров и режимы:
- qilin-crm-sdk-parent: режим parent
- qilin-crm-sdk-dev: режим dev
- qilin-crm-sdk: режим qilin

```
version: '3.7'
services:
  crm-sdk:
    container_name: qilin-crm-sdk
    image: p1hub/qilin-crm-api:demo-sdk
    restart: always
    ports:
      - 8086:8080
    command: ["sdk","-c","configs/sdk.yaml","-b",":8080","-d"]
    networks:
      - default
  crm-sdk-dev:
    container_name: qilin-crm-sdk-dev
    image: p1hub/qilin-crm-api:demo-sdk
    restart: always
    ports:
      - 8085:8080
      - 1443:1443
    command: ["sdk","-c","configs/sdk-dev.yaml","-b",":8080","-d"]
    networks:
      - default
  crm-sdk-parent:
    container_name: qilin-crm-sdk-parent
    image: p1hub/qilin-crm-api:demo-sdk
    restart: always
    ports:
      - 8084:8080
    command: ["sdk","-c","configs/sdk-parent.yaml","-b",":8080","-d"]
    networks:
      - default
```
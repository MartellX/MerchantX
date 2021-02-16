### Что использовал
- Golang
- Веб-фреймворк: [echo](https://github.com/labstack/echo)
- ORM: [gorm](https://github.com/go-gorm/gorm)
- Тестирование с использованием: [sqlmock](https://github.com/DATA-DOG/go-sqlmock), [gomock](https://github.com/golang/mock), [gomega](https://github.com/onsi/gomega)

### Запуск
1. Клонировать репозиторий: `git clone https://github.com/MartellX/MerchantX`
2. Прописать необходимые параметры в `.env`
3. Запустить контейнер `docker compose up`

P.S. При первом запуски `gorm` проведет необходимые автомиграции

### Описание запросов
1. **POST** /tasks - создание нового задания
    - Тело:
        - `seller_id` - id продавца
        - `url` - ссылка на xlsx таблицу
    - Возвращает созданное задание
    - Пример запроса:
      ```shell 
      curl -L -X POST 'http://localhost:1323/tasks' \
      -H 'Content-Type: application/x-www-form-urlencoded' \
      --data-urlencode 'seller_id=2' \
      --data-urlencode 'url=https://docs.google.com/spreadsheets/d/1IqTYDGuPnFc40sMaKF4KEbnGWholL2Fp4ISQhMcsPD4/export?format=xlsx'
    - Пример ответа:
         ```json
        {
        "task_id": "ac71f2d9-49d3-4ba2-8069-078b945be570",
        "status": "Created",
        "status_code": 201,
        "info": {}
        }
    
2. **GET** /tasks - получение задания
    - Параметры:
        - `task_id` - id задания
    - Возвращает найденное задание
    - Пример запроса:
      ```shell
      curl -L -X GET 'http://localhost:1323/tasks?task_id=1cc82fee-2658-4a6b-97d0-7fffeabdf988'
    - Пример ответа:
      ```json 
      {
      "task_id": "1cc82fee-2658-4a6b-97d0-7fffeabdf988",
      "status": "Completed",
      "status_code": 200,
      "info": {
      "created": 10989,
      "updated": 10,
      "errors": 9
         }
      }
    
3. **GET** /offers - получение товаров по заданным параметрам
    - Параметры (Необязательные):
        - `offer_id`- id товара
        - `seller_id` - id продавца
        - `name` - построка названия товара
      
    - Возвращает товары по указанным параметрам (если не указаны, то возвращает все)
    - Пример запроса:
      ```shell
      curl -L -X GET 'http://localhost:1323/offers?name=phone'
    - Пример ответа:
      ```json
      {
       "count": 18,
       "items": [
           {
               "offer_id": 2312,
               "seller_id": 121231,
               "name": "iPhone",
               "price": 123,
               "quantity": 12,
               "available": true
           },
           {
               "offer_id": 54353,
               "seller_id": 121231,
               "name": "iphone",
               "price": 4333,
               "quantity": 2,
               "available": true
           },
           {
               "offer_id": 12312,
               "seller_id": 121231,
               "name": "telephone",
               "price": 4242,
               "quantity": 1,
               "available": true
           },
            ...
         ]
      }
      ```
   
### Примечания
- При разработке в качестве тестового файла использовался [этот](https://docs.google.com/spreadsheets/d/1IqTYDGuPnFc40sMaKF4KEbnGWholL2Fp4ISQhMcsPD4/export?format=xlsx)
- При тестировании пакета `serivces` поднимается небольшой сервер для проверки обработки ссылок. Используется порт `1234`, хотя, думаю, его лучше задавать через переменные окружения
- Думал для каждой распарсенной строки создавать отдельную горутину для формирования запросов к БД, но на больших файлах работа быстро становилась нестабильной. Пробовал ограничить одновременную загрузку нескольких строк, но оказалось, что это очень замедляет работу (я так и не очень разобрался почему). Поэтому вернулся к идее последовательной загрузки
- Нагрузочное тестирование не проводил, но ~10 тысяч строк обрабатываются и добавляются в БД примерно за 20 секунд
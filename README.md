## Шаблон: python38-graphql

```
Python 3.8.6 
pip 19.3.1
```
Для корректной работы в функции должен быть скрипт `handler/handler.py` с функцией `handle` принимающей на вход `FunctionRequest` и возвращающий `FunctionResponse`.  
В `FunctionsRequest` содержится `payload`(пришедший запрос) и `headers`(заголовки).  
Дополнительная информация (path, queryParams, Content-Type и т.д.) содержатся в заголовках.  
В `FunctionsResponse` можно передать тело ответа, статус и заголовки.  

Для мониторинга метрик используется библиотека [prometheus/client_python](https://github.com/prometheus/client_python).  
Пример регистрации метрики можно увидеть в [`handler.py`](./handler/handler.py).

### Подключение внешних зависимостей
Внешние зависимости могут быть определены в [`requirements.txt`](./requirements.txt) проекта функции.

### Взаимодействие с DATASPACE
⚠️ **При запуске должны быть определены переменные среды DATASPACE_URL, APP_KEY, APP_SECRET** ⚠️

Dataspace поддерживает взаимодействие по протоколу GraphQL. Поэтому для вызова в данном шаблоне используется библиотека [gql](https://github.com/graphql-python/gql)

Вызовы DataSpace, идущие через API-Gateway, должны быть подписаны ключом (APP_KEY) и секретом (APP_SECRET) из переменных среды.  
Для этого используется [signer.py](./handler/apig_sdk/signer.py).  

В [`apolloClient`](./handler/handler.py) создается клиент `client` настроенный на DATASPACE_URL.  
Через него можно вызывать сервис Dataspace, выполнять различные операции - query, mutate и т.д.  
Запросы проходящие через этого клиента подписываются классом `DataspaceAuth`.

Пример подписи и вызова Dataspace можно увидеть в [handler.py](./handler/handler.py)

### Конфигурирование через ConfigMap:
Для возможности конфигурирования функции через ConfigMap пользователю необходимо:
1) Определить свойства в файле, например: `properties.ini`
2) В ```functions.yaml``` указать конфигурационные файлы для монтирования.
```yaml
    configs:
      - name: python38-graphql-example              # Имя конфигурации. В UI OSE ConfigMap будет называться <имя-функции>-cm-<имя конфигурации>
        files:                                # Список файлов для монтирования
          - properties.ini    # Полный путь до файла относительно директории с функцией
```
Данный файл будет смонтирован в папку `/app/config`.  
Далее в функции можете использовать этот файл, например через библиотеку [configparser](https://docs.python.org/3/library/configparser.html)

### Тестирование
⚠️ **При сборке функции будет запущен скрипт [`./handler_tests.py`](./test_handler.py). Если тесты не пройдут, то функция не развернется.** ⚠️
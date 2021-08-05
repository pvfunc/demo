## Шаблон: java11-springboot

Шаблон Java 11 Spring Boot использует maven в качестве системы сборки проекта.

Maven: 3.6.3

Java: OpenJDK 11

Spring Boot: 2.3.1.RELEASE


### Структура проекта

    root
     ├ src                              maven source route
     │ └ main
     │   ├ java                         implement your function here
     │   │ ├ ...
     │   │ ...
     │   └ resources
     │     └ config.yaml               write configuration properties used in scripts in this file
     │                                  config.yaml must not be renamed
     │
     └ pom.xml                          place dependencies of your function here


### Подключение внешних зависимостей

Внешние зависимости могут быть определены в ```./pom.xml``` проекта функции.

    
### Конфигурирование через ConfigMap:
Для возможности конфигурирования функции через ConfigMap пользователю необходимо:
1) Определить свойства в файле ```src/main/resources/config.yaml```
2) В ```functions.yaml``` указать конфигурационные файлы для монтирования.
```yaml
    configs:
      - name: springboot-example              # Имя конфигурации. В UI OSE ConfigMap будет называться <имя-функции>-cm-<имя конфигурации>
        files:                                # Список файлов для монтирования
          - src/main/resources/config.yaml    # Полный путь до файла относительно директории с функцией
```

## Локальная отладка функции
### Prerequisites
На рабочей станции разработчика должны быть установлены:
 - `JDK 11`
 - `Maven 3.5.0+` 

### Локальный запуск приложения 
Если вы используете переменные в файле [config.yaml](./src/main/resources/config.yaml) то для запуска Springboot приложения необходимо определить переменную среды:  
`--spring.config.location=<путь до config.yaml в папке resources>`

Далее можно запустить Main класс [App.java](./src/main/java/sbp/ts/faas/templates/springboot/pure/App.java) и вызывать по HTTP (с помощью CURL, Postman и т.д.) необходимые контроллеры.

### Локальное тестирование
Если вы используете переменные в файле [config.yaml](./src/main/resources/config.yaml) то для запуска тестов необходимо указать спрингу на него. Например, как это сделано в [ControllerTest](./src/test/java/sbp/ts/faas/templates/springboot/pure/ControllerTest.java).

Во всем остальном тестирование фукнции никак не отличается от тестирования любого Springboot приложения.  
Документация по тестированию в Springboot: [https://docs.spring.io](https://docs.spring.io/spring-boot/docs/2.3.1.RELEASE/reference/html/spring-boot-features.html#boot-features-testing)



Во время деплоя функции будут запущены тесты с помощью мавена, если тесты не пройдут - функция не задеплоится.
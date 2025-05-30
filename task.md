#### `Сервер` должен реализовывать следующую бизнес-логику:
- регистрация, аутентификация и авторизация пользователей;
- хранение приватных данных;
- синхронизация данных между несколькими авторизованными клиентами одного владельца;
- передача приватных данных владельцу по запросу.

#### `Клиент` должен реализовывать следующую бизнес-логику:
- аутентификация и авторизация пользователей на удалённом сервере;
- доступ к приватным данным по запросу.
- Функции, реализация которых остаётся на усмотрение исполнителя:
- создание, редактирование и удаление данных на стороне сервера или клиента;
- формат регистрации нового пользователя;
- выбор хранилища и формат хранения данных;
- обеспечение безопасности передачи и хранения данных;
- протокол взаимодействия клиента и сервера;
- механизмы аутентификации пользователя и авторизации доступа к информации.
 
####  * Дополнительные требования:
- клиент должен распространяться в виде CLI-приложения с возможностью запуска на платформах Windows, Linux и Mac OS;
- клиент должен давать пользователю возможность получить информацию о версии и дате сборки бинарного файла клиента.


#### `Типы` хранимой информации
- пары логин/пароль;
- произвольные текстовые данные;
- произвольные бинарные данные;
- данные банковских карт.

Для любых данных должна быть возможность хранения произвольной текстовой метаинформации (принадлежность данных к веб-сайту, личности или банку, списки одноразовых кодов активации и прочее).

### Абстрактная схема взаимодействия с системой
Ниже описаны базовые сценарии взаимодействия пользователя с системой. Они не являются исчерпывающими — решение отдельных сценариев (например, разрешение конфликтов данных на сервере) остаётся на усмотрение исполнителя.

#### Для нового пользователя:
- Пользователь получает клиент под необходимую ему платформу.
- Пользователь проходит процедуру первичной регистрации.
- Пользователь добавляет в клиент новые данные.
- Клиент синхронизирует данные с сервером.
 
#### Для существующего пользователя:
- Пользователь получает клиент под необходимую ему платформу.
- Пользователь проходит процедуру аутентификации.
- Клиент синхронизирует данные с сервером.
- Пользователь запрашивает данные.
- Клиент отображает данные для пользователя.

#### Тестирование и документация
- Код всей системы должен быть покрыт юнит-тестами не менее чем на 80%. 
- Каждая экспортированная функция, тип, переменная, а также пакет системы должны содержать исчерпывающую документацию.
- Необязательные функции

#### Перечисленные ниже функции необязательны к имплементации, однако позволяют лучше оценить степень экспертизы исполнителя. 
Исполнитель может реализовать любое количество из представленных ниже функций на свой выбор:
- поддержка данных типа OTP (one time password );
- поддержка терминального интерфейса (TUI — terminal user interface);
- использование бинарного протокола;
- наличие функциональных и/или интеграционных тестов;
- описание протокола взаимодействия клиента и сервера в формате Swagger.
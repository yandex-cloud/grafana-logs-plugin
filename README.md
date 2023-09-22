# Плагин Yandex Cloud Logging для Grafana

Плагин Yandex Cloud Logging для Grafana — расширение для Grafana, с помощью которого можно добавить [Cloud Logging](https://cloud.yandex.ru/services/logging) в качестве источника данных.

Плагин позволяет сервисам и приложениям читать записи из лог-групп, используя [язык фильтрующих выражений](https://cloud.yandex.ru/docs/logging/concepts/filter).

## Требования

1. [Создайте сервисный аккаунт](https://cloud.yandex.ru/docs/iam/operations/sa/create#create-sa) и назначьте ему роль `logging.reader` на каталог.
1. [Создайте авторизованный ключ](https://cloud.yandex.ru/docs/iam/operations/authorized-key/create) для сервисного аккаунта, чтобы аутентифицироваться в Cloud Logging API.
1. [Создайте лог-группу](https://cloud.yandex.ru/docs/logging/operations/create-group) и [добавьте](https://cloud.yandex.ru/docs/logging/operations/write-logs) в нее записи.

## Поддержка

Чтобы задать вопросы о плагине или рассказать о проблемах, которые возникли при работе с ним, заведите issue в репозитории [yandex-cloud/grafana-logs-plugin](https://github.com/yandex-cloud/grafana-logs-plugin).

## Контрибьюты

Мы рады развивать плагин совместно с сообществом. Перед тем как сделать PR, ознакомьтесь с [руководством](https://github.com/yandex-cloud/grafana-logs-plugin/blob/master/CONTRIBUTING.md).

## Лицензия

[Лицензия Apache, версия 2.0](https://github.com/yandex-cloud/grafana-logs-plugin/blob/master/CONTRIBUTING.md)

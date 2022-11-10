# Агрегатор статей по Go
Собирает статьи с заданных ресурсов и предоставляет их через единый интерфейс

Требования к разработке нового парсера:
* Исходники парсера должны располагаться в отдельной папке по [пути](./internal/parser) `./internal/parser`, пример
  * `./internal/parser/myparser`
* Рекомендуемое наименование пакета парсера (имени файла) `parser` (`parser.go`)
* Парсер должен реализовывать интерфейс ArticlesParser - [пакет](./internal/parser/parser.go) `./internal/parser/parser.go`
* Парсер должен реализовать функцию создания парсера с таким контрактом
  * ```go
    func NewParser(cfg parser.Config) parser.ArticlesParser
    ```
* Для авторегистрации парсера необходимо
  * добавить функцию init в своем пакете, в которой зарегистрировать свой парсер, указав его наименование и функцию создания, пример
    * ```go
      func init() {
          parser.RegisterParser("myparser", NewParser)
      }
      ```
  * импортировать пакет парсера черeз *underscore* в [файле](./internal/parser/autoregister/register.go) `./internal/parser/autoregister/register.go`, пример
    * ```go
      _ "github.com/indikator/aggregator_lets_go/internal/parser/myparser"
      ```
  * указать в [конфиге](./configs/config.yaml) `./configs/config.yaml` секцию с настройками своего парсера, пример
    ```yaml
    - myparser:
        url: https://myparser.org/news
    ```

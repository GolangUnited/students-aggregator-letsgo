# Github CI
* Описание переменных Secrets:
  * Доступ к DockerHub
    * DOCKERHUB_USERNAME - пользователь DockerHub
    * DOCKERHUB_TOKEN - токен для доступа к аккаунту DockerHub
  * Доступ к удаленного серверу
    * SSH_REMOTE_HOST - адрес удаленного сервера
    * SSH_REMOTE_HOST_PORT - порт удаленного сервера
    * SSH_PRIVATE_KEY - приватный ключ для доступа к удаленному серверу по SSH
    * SSH_KNOWN_HOST - запись из файла known_host для добавления удаленного сервера

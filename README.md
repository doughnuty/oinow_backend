# REST API BACKEND
#### for BTSD HackNU oinow web app
https://github.com/Zhalkhas/hacknu_react_aitu

## UserInit
Инициализировать пользователя. Данная команда нужна при первой авторизации юзера и дальнейшей ассоциации его aituID с номером телефона. Эта самая ассоциация в свою очередь нужна в связи с тем что метод getContacts() не передает aituID пользователей из контактной книги.

### Request

`POST /rest/oinow/profile/`
#####
    {
      "aituID": "ABC",
      "first_name": "Kafka",
      "last_name": "Kafka",
      "phone":  "+78005553535"
    }

### Response
Возвращает очки юзера если таковые есть в бд. Если юзер новый, возвращает 0
    HTTP/1.1 201 Created
    Status: 201 Created
    Content-Type: application/json
    Content-Length: 2
    
    {0}


## UserGetScore
Получить количество очков пользователя по его уникальному ID
### Request

`GET /rest/oinow/profile/{aituID}`

### Response

    HTTP/1.1 200 OK
    Status: 200 OK
    Content-Type: application/json
    Content-Length: 2
    
    {0}

## CreateGame
Создать игру по ее названию

### Request

`POST /rest/oinow/games/`

### Response
Возвращает Success при благополучном создании. Добавляет игру в дб

    HTTP/1.1 201 Created
    Status: 201 Created
    Content-Type: application/json
    Content-Length: 2
    
    {"Success"}
    
## GetLeaderboard
Возвращает список пользователей сортированный по их баллам
###### ! На данный момент высылается все что есть в структуре (поле телефона пустое т.к. не генерируется без необходимости), но достаточно будет ограничиться ФИ и скором !
### Request

`GET /rest/oinow/leaderboard/
  
### Response

    HTTP/1.1 200 OK
    Status: 200 OK
    Content-Type: application/json
    Content-Length: 194

    [
    {
        "id": 3,
        "first_name": "Kafka",
        "last_name": "Kafka",
        "aituID": "CBA",
        "score": 14,
        "style": 0,
        "phone": ""
    },
    {
        "id": 1,
        "first_name": "Kafka",
        "last_name": "Josh",
        "aituID": "ABC",
        "score": 13,
        "style": 0,
        "phone": ""
    }
    ]

##### Остальные команды на https://www.notion.so/Backend-f3b84a754e784f0bb86750ad36217350



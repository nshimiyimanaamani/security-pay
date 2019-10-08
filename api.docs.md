# paypack http api documentation
This document is organised such that each service's endpoints are listed in their own section


## version

Get version information
*`GET /version`

## users service

Registering a new user.
* `POST /users`:
    - contentType:`application/json`
    
    - request body: 
    ```
    {
        "email":"user@example.com", "password":"12345", "cell:admin"
    }
    ```
    - response body:
    ```
    {
        "id": "555955b4-a2dc-48c2-83d8-ce736e7bb24a"
    }
    ```

Login a user.
* `POST /users`:
    - contentType:`application/json`
    
    - request body: 
    ```
    {
        "email":"user@example.com", "password":"12345"
    }
    ```
    - response body:
    ```
    {
        "token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NjE2NzcwMjksImlhdCI6MTU2MTY0MTAyOSwiaXNzIjoicGF5cGFjayIsInN1YiI6ImphbWVzQGdtYWlsLmNvbSJ9.fkLjioB4hLugMTIc7FbO_sLBhvsdkasJq4wuoxWw198"
    }
    ```

List users(TODO)
* `GET /users`:
    - contentType:`application/json`

## transactions(historic) service

Record a new transaction.
* `POST /transactions`:
    - contentType:`application/json`

    - HEADERS:`["Authorization"]`

    - request body: 
    ```
    {
        "amount":"1000", 
        "method":"MTN", 
        "madefor":"83232d60-c527-4b92-a45a-c451ca217a4e"
        "madeby=f21265be-3d4b-4961-aef2-1a3514f07684"
    }
    ```
    - response body: `
    ```
    {
        "id": "c48e8607-1834-4b81-a935-7cb30d4e7416"
    }
    ```

Note that the `property` must be a valid uuid and you can generate them with at https://www.uuidgenerator.net/

Get a transaction given it's id
* `GET /transactions/:id`
    - example: `transactions/c48e8607-1834-4b81-a935-7cb30d4e7416`

    -  HEADERS:`["Authorization"]`

    - response body: 
    ```
    {   
        "id": "c48e8607-1834-4b81-a935-7cb30d4e7416",
        "amount": "1000",
        "method": "MTN",
        "property": "83232d60-c527-4b92-a45a-c451ca217a4e"
        "owner":"Johnny Evans"
        "date":""
    }
    ```
Get a list a subset of transactions given an offset and the limit
* `GET /transactions/?offset=0&limit=5`
    - example using httpie: `http  "localhost:8081/api/transactions/?offset=0&limit=5"`

     HEADERS:`["Authorization"]`

    - response body:
``` 
{
    "limit": 5,
    "offset": 0,
    "total": 4,
    "transactions": [
        {
            "amount": "1000",
            "id": "32f8ebc7-67a2-41dc-a8f1-3f06f3b58b84",
            "method": "MTN",
            "property": "c49ca697-de03-4798-b2cb-845c3c3f2e7f"
            "owner":"Johnny Evans"
            "date":""
        },
        {
            "amount": "1000",
            "id": "c48e8607-1834-4b81-a935-7cb30d4e7416",
            "method": "MTN",
            "property": "83232d60-c527-4b92-a45a-c451ca217a4e"
            "owner":"Johnny Evans"
            "date":""
        },
        {
            "amount": "1000",
            "id": "d1756a50-010d-4f57-b41e-34b1acf6dcf9",
            "method": "MTN",
            "property": "83232d60-c527-4b92-a45a-c451ca217a4e"
            "owner":"Johnny Evans"
            "date":""
        },
        {
            "amount": "1000",
            "id": "fbc7a2bd-6a78-448e-9bd6-3dfcc6436f55",
            "method": "MTN",
            "property": "83232d60-c527-4b92-a45a-c451ca217a4e"
            "owner":"Johnny Evans"
            "date":""
        }
        ]
    }
```
## properties(houses) endpoints.
**properties**: are the general properties endpoints. 
Add a property object to a owner portofolio given their uid
* `/properties/`
    - method: `POST`

    -  HEADERS:`["Authorization"]`
    
    - request body: 
    ```
    {
        "cell": "gishushu",
        "owner": "54f0a0a5-373a-439f-8831-2e5151535679",
        "due":"1000",
        "sector": "remera",
        "village": "ingabo"
    }
    ```
    - response body:
    ```
    {
        "id":"9f27518b-9023-4f4d-b949-c097511b66e7"
    }
    ```
Retrieve a property given it's uid
* `/properties/:id`

**owners**: endpoints return help manage property owners.
Add a new owner to the reperitory(require before adding properties)
* `/properties/owners/`
    - method: `POST`

    - HEADERS:`["Authorization"]`

    - request body: 
    ```
    {
        "fname":"James", "lname":"Tamal", "phone":"0784555222"
    }
    ```
    - response body:
    ```
    {
        "id": "c48e8607-1834-4b81-a935-7cb30d4e7416"
    }
    ```
Retrieve a property given it's uid
* `/properties/owners/:id`
    - method: `GET`

    -  HEADERS:`["Authorization"]`

    - request body: `empty`

    - response body: 
    ```
    {
        "fname": "Tucky",
        "id": "54f0a0a5-373a-439f-8831-2e5151535679",
        "lname": "Bucky",
        "phone": "0784577882"
    }
    ```
Search and retrieve an owner given their first name and phone number
* `/properties/owners/search/?fname=n&lname=m&phone=o`
    - method: `GET`

    -  HEADERS:`["Authorization"]`

    - request body: `empty`
    - response body: 
    ```
    {
        "fname": "Tucky",
        "id": "54f0a0a5-373a-439f-8831-2e5151535679",
        "lname": "Bucky",
        "phone": "0784577882"
    }
    ```
Retrieve a subset of owners as a list given an offset and a limit
* `/properties/owners/?offset=n&limit=m`
    - method: `GET`
    - HEADERS:`["Authorization"]`
    - request body: `empty`
    - response body:
    ```
    {
        "limit": 4,
        "offset": 0,
        "total": 3
        "owners": [
            {
                "fname": "Tucky",
                "id": "54f0a0a5-373a-439f-8831-2e5151535679",
                "lname": "Bucky",
                "phone": "0784577882"
            },
            {
                "fname": "Tucky",
                "id": "652c5f01-3259-4297-a359-99203997c532",
                "lname": "Tucky",
                "phone": "0784577882"
            },
            {
                "fname": "jason",
                "id": "e56b3456-2466-460c-b987-3ec3dab2f4a6",
                "lname": "Born",
                "phone": "0734577882"
            }
        ],
    }

    ```
Retrieve properties given the owner
* `/properties/owners/properties/{owner}?offset=0&limit=5`
    - method: `GET`
    - HEADERS:`["Authorization"]`
    - request body: `empty`
    - response body: 
    ```
    {
        "limit": 5,
        "offset": 0,
        "total": 1
        "properties": [
            {
                "cell": "gishushu",
                "id": "9f27518b-9023-4f4d-b949-c097511b66e7",
                "owner": "54f0a0a5-373a-439f-8831-2e5151535679",
                "due":"1000",
                "sector": "remera",
                "village": "ingabo"
            }
        ],
    }

    ```

**admin blocks**: endpoints return properties within certain administration blocks

Retrieve properties given the sector of their location
* `"/properties/sectors/:sector/?offset=n&limit=m`

Retrieve properties given the sector and cell of their location
* `"/properties/sectors/:sector/cells/:cell/?offset=n&limit=m`

Retrieve properties given the sector, cell and village of their location
* `"/properties/sectors/:sector/cells/:cell/villages/:village/?offset=n&limit=m`

***payment**: endpoints to make and validate payment
* `/payment/initialize`
    - request_body:
        {
            "code":"3124jifr",
            "amount": 5,
            "phone":"+250789000111"
        }

**Notice**: 
* all the endpoints except the users endpoints now require an`Authorization` header which contains the token 
acquired after a successful login.
    
* To try it out check this cloud run endpoint https://paypack-backend-qahoqfdr3q-uc.a.run.app/api/

# paypack http api documentation
This document is organised such that each service's endpoints are listed in their own section

## users service

Registering a new user.
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

    - request body: 
    ```
    {
        "amount":"1000", "method":"MTN", "property":"83232d60-c527-4b92-a45a-c451ca217a4e"
    }
    ```
    - response body: `
    ```
    {
        "id": "c48e8607-1834-4b81-a935-7cb30d4e7416"
    }
    ```

Note that the `property` must be a valid uuid and you can generate them with at https://www.uuidgenerator.net/

* `GET /transactions/:id`
    - example: `transactions/c48e8607-1834-4b81-a935-7cb30d4e7416`
    - response body: 
    ```
    {
        "amount": "1000",
        "id": "c48e8607-1834-4b81-a935-7cb30d4e7416",
        "method": "MTN",
        "property": "83232d60-c527-4b92-a45a-c451ca217a4e"
    }
    ```
* `GET /transactions/?offset=0&limit=5`
    - example using httpie: `http  "localhost:8081/api/transactions/?offset=0&limit=5"`
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
        },
        {
            "amount": "1000",
            "id": "c48e8607-1834-4b81-a935-7cb30d4e7416",
            "method": "MTN",
            "property": "83232d60-c527-4b92-a45a-c451ca217a4e"
        },
        {
            "amount": "1000",
            "id": "d1756a50-010d-4f57-b41e-34b1acf6dcf9",
            "method": "MTN",
            "property": "83232d60-c527-4b92-a45a-c451ca217a4e"
        },
        {
            "amount": "1000",
            "id": "fbc7a2bd-6a78-448e-9bd6-3dfcc6436f55",
            "method": "MTN",
            "property": "83232d60-c527-4b92-a45a-c451ca217a4e"
        }
        ]
    }
```
## properties(houses) endpoints.
**properties**: are the general properties endpoints. 

* `/properties/`
    - method: `POST`
    - request body: 
    ```
    {
        "cell": "gishushu",
        "owner": "54f0a0a5-373a-439f-8831-2e5151535679",
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

* `/properties/:id`

**owners**: endpoints return help manage property owners.

* `/properties/owners/`
    - method: `POST`
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

* `/properties/owners/:id`
    - method: `GET`
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
* `/properties/owners/search/?fname=n&lname=m&phone=o`
    - method: `GET`
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

* `/properties/owners/:owner/?offset=n&limit=m`
    - method: `GET`
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
* `/properties/owners/properties/{owner}?offset=0&limit=5`
    - method: `GET`
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
                "sector": "remera",
                "village": "ingabo"
            }
        ],
    }

    ```

**admin blocks**: endpoints return properties within certain administration blocks

* `"/properties/sectors/:sector/?offset=n&limit=m`

* `"/properties/sectors/:sector/cells/:cell/?offset=n&limit=m`

* `"/properties/sectors/:sector/cells/:cell/villages/:village/?offset=n&limit=m`


**Note**: To try it out check this cloud run endpoint https://paypack-backend-qahoqfdr3q-uc.a.run.app/api/

# paypack http api documentation
This document is organised such that each service's endpoints are listed in their own section

## users service

Registering a new user.
* `POST /users`:
    - contentType:`application/json`
    
    - request body: `{"email":"user@example.com", "password":"12345"}`
    - response body:`{"id": "555955b4-a2dc-48c2-83d8-ce736e7bb24a"}`

Login a user.
* `POST /users`:
    - contentType:`application/json`
    
    - request body: `{"email":"user@example.com", "password":"12345"}`
    - response body:
        `{
            "token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NjE2NzcwMjksImlhdCI6MTU2MTY0MTAyOSwiaXNzIjoicGF5cGFjayIsInN1YiI6ImphbWVzQGdtYWlsLmNvbSJ9.fkLjioB4hLugMTIc7FbO_sLBhvsdkasJq4wuoxWw198"
        }`

List users(TODO)
* `GET /users`:
    - contentType:`application/json`

## transactions(historic) service

Record a new transaction.
* `POST /transactions`:
    - contentType:`application/json`

    - request body: `{"amount":"1000", "method":"MTN", "property":"83232d60-c527-4b92-a45a-c451ca217a4e"}`
    - response body: `{"id": "c48e8607-1834-4b81-a935-7cb30d4e7416"}`

Note that the `property` must be a valid uuid and you can generate them with at https://www.uuidgenerator.net/

* `GET /transactions/:id`
    - example: `transactions/c48e8607-1834-4b81-a935-7cb30d4e7416`
    - response body: 
    
    `{
        "amount": "1000",
        "id": "c48e8607-1834-4b81-a935-7cb30d4e7416",
        "method": "MTN",
        "property": "83232d60-c527-4b92-a45a-c451ca217a4e"
    }`
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

**Note**: To try it out check this cloud run endpoint https://paypack-backend-qahoqfdr3q-uc.a.run.app/api/

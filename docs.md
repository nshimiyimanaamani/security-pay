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

**Note**: To try it out check this cloud run endpoint https://paypack-backend-qahoqfdr3q-uc.a.run.app/api/

# Contributing 

## Issues

Open an issue describing what you wanted to do, what you expected and what you got.
For new feature you can also open an issue describing why it's needed.

## Commit messages

## Run paypack api on your machine

1. make sure there a .env file in the root directory of the repo:

    ```$ touch .env```

    Add the following variables:

    ```
    DATABASE_HOST=database
    DATABASE_NAME=test
    DATABASE_PORT=5432
    DATABASE_USER=postgres
    DATABASE_PASSWORD=password
    LOG_LEVEL=info
    PAYPACK_SECRET=e9Z0e23e23r23r23tdqdqwfqwfqwfq
    PAYPACK_PAYMENT_ENDPOINT=https://novapay.rw/api/v1/novapay/initialize-payment
    PAYPACK_PAYMENT_TOKEN=sDclmEn3b0oh3M1QfQN5hLYf1ATWpKtQ

    ```
2. run the following command to start the stack`[database | api]`:
    
   `$ make dev`

    paypack should be listening on at:
    
    `localhost:8081`

3. ounce you are finished you can stop the backend stack with:

    `$ make dev-teardown`

4. For more commands run:
    `$ make help`

***Note**: make sure you have both docker and docker-compose installed and make(optinonal if you want to run commands manualy)
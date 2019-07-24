[![CircleCI](https://circleci.com/gh/rugwirobaker/paypack-backend.svg?style=shield&circle-token=6f5aa06254f06fbeccf22a77d54ee272b197fbd8)](https://circleci.com/gh/rugwirobaker/paypack-backend)

# paypack-backend
The repo for the paypack backend monolithic service.

## About
Paypack is a custom payment service built for the local government needs for collection of contributions from the community.
 * Language: Go
 * Organization: [Quarks Group]()
 * Lead Developer: Rugwiro Mbanda Valentin: [github.com](https://github.com/rugwirobaker/), [email](rugwiorbaker@gmail.com), [Twitter](https://twitter.com/acodechef)

 This repository is part of a bigger projects composed of the following repositories:
  * **docs**: main documentation at [paypack-docs](https://github.com/rugwirobaker/paypack-docs).
  * **backend**: backend code which is this one.
  * **frontend**: frontend code at [paypack-frontend](https://github.com/rugwirobaker/paypack-frontend).
  * **compose-manifests**: docker compose deployments manifests at [paypack-compose-manifests](https://github.com/rugwirobaker/paypack-compose-manifests).
  * **k8s-manifests**:  kubernetes deployments manifests at [paypack-kubernetes-manifests](https://github.com/rugwirobaker/paypack-kubernetes-manifests)

 ## Contributions
 We require all new contributions to go through a pull request as per the following instructions:
  * Create a new branch. For details on branch naming check [contribution guidelines](CONTRIBUTORS.md)
  * push after adding yout changes push the changes to origin/`branch`
  * create a pull request with descriptions of the changes made.
  * assign a reviewer.

## Infrastructure.
This project tries to follow the clean architecture guidelines. Where code is subdivided into `infrastucure`, `business logic`, and `transport` going from bottom to top. For this repo those layers have been  subdived into different packeges as follow:
 * **infastructure**: store
 * **business logic**: app(application)
 * **transport**: api`[http, grpc, websocket,...]`
 * **shared**: the models package contains type definitions that are shared among the 3 layers.
 * **bootstrap**: cmd package contains server bootstrapp logic.

Other directories:
 * **.circleci**: circleci configuration
 * **bin**: the bin directory is created when you run `make build` to create the applciation binary.
 * **secrets**: the secrets directory is created when you run `make secrets` to create secrets and certficates.

 ## To Get started.
 For a locally runnable version of paypack checkout the compose manifest rep at  https://github.com/rugwirobaker/paypaack-compose-manifests
 ## API docs
 Checkout the REST API docs at [api](api.docs.md)


 check the endpoint documentantion detials at https://github.com/rugwirobaker/paypack-backend/blob/master/docs.md

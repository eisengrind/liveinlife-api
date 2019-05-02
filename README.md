![Banner](https://github.com/51st-state/api/blob/master/assets/banner.png?raw=true)

[![GitHub license](https://img.shields.io/github/license/51st-state/api.svg)](https://github.com/51st-state/api/blob/master/LICENSE)
[![Documentation](https://godoc.org/github.com/51st-state/api?status.svg)](https://godoc.org/github.com/51st-state/api)
[![Go Report Card](https://goreportcard.com/badge/github.com/51st-state/api)](https://goreportcard.com/report/github.com/51st-state/api)

api is an API service for the Eisengrind RageMP RPG server. It provides scalability, stability and availability of important API functions to be called by the server and its clients and thus, reducing the load on the main game server.

## Concepts

One of the most basic concept of this api is to move the read requests of a CRUD-model to the client. Other operations such as create, update and delete can be partially moved to the client. The main difficulty is that some actions (e.g. picking up an item) are limited to the game itself and thus there has to be a trusted authority (in this case the game server itself) to manage those requests (picking up item -> move to player inventory for example).

Another concept is the interoperability between services. This means that we want other players to be able to interact with our API to create services extending the virtual roleplaying world. In addition to that, we try our best to provide OAuth2-specific authentication for external service to enable contributors to be aple to interact with our API.

### Scalability

Scalability has a big impact on the games performance - whether it is the server or the client. So we try the best we can to create non-monolithic applications. The main goal is to create redundant services which are able to keep up with growing player amounts. In addition to that, we run our services in a Kubernetes Cluster. Therefore, some of the microservices provided by this API use message queues to provide another layer of scalability.

## Usage

Just deploy given docker containers in your k8s cluster. (list of Docker containers will be created in the future!)

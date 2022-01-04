# Ecommerce User Service

This is the GRPC user service that handles all activities in the ecommerce application that has to do with users (retrieve, update etc).

## Applications

* ## [Jaeger](https://www.jaegertracing.io/)

  * Jaeger is an **open source software for tracing transactions between distributed services**.
* ## [Nats](https://nats.io)

  * NATS is **an open-source messaging system** (sometimes called message-oriented middleware).
  * NATS is used in the user service for communicating with the notification service when a new user is added.
* ## [MongoDB](https://www.mongodb.com/)

  * MongoDB is an open source NoSQL database management program. NoSQL is used as an alternative to traditional relational databases.
  * The user service stores users in MongoDB.

### Usage

To install / run the user microservice run the command below:

```bash
docker-compose up
```

## Requirements

The application requires the following:

* Go (v 1.5+)
* Docker (v3+)
* Docker Compose

## Other Micro-Services / Resources

* #### [Product Service](https://github.com/wisdommatt/ecommerce-microservice-product-service)
* #### [Notification Service](https://github.com/wisdommatt/ecommerce-microservice-notification-service)
* #### [Cart Service](https://github.com/wisdommatt/ecommerce-microservice-cart-service)
* #### [Shared](https://github.com/wisdommatt/ecommerce-microservice-shared)

## Public API

The public graphql API that interacts with the microservices internally can be found in [https://github.com/wisdommatt/ecommerce-microservice-public-api](https://github.com/wisdommatt/ecommerce-microservice-public-api).

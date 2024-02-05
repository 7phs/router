# Destination Measurement

A service is responsible for collecting information from 3rd party map-services and decorating it into internal domain model representation.

**Functional Requirements:**

1. Handles REST API requests with a list of geographic coordinates and orders these points by time to travel to them.
2. The service should support easy switsching to new services.

**Non-functional Requirements:**

1. Throttles requests to 3rd party map-services to align with their limits.
2. Response time should be less than 20 ms on average for a working day in a specific location.
3. Needs to respond with already known destinations even if 3rd party services are not available at the current moment.

# Architecture

![Architecture of the service](./docs/architecture.wsd.svg)

## REST API server

**Functional Requirements:**

1. Fetch data from bridge of routing services;
2. Represents data in JSON response.

**Non-functional Requirements:**

1. Responses as much data as fetched to response even some of the routes processed with errors. 

## Bridge of routing services

**Functional Requirements:**

1. Requests an external services if no data in cache
2. Stores fetched data of routing in cache

**Non-functional Requirements:**

1. There are no requirements to limit the number of requests to external services. It is okay to do several requests simultaneously and then store or replace them in the cache.

## Cache of routing data

In-memory cache used in the current implementation to reduce the number of services included in deployments. The second reason is that the number of possible routes is limited for a certain city or area of the city. A solution can be scaled horizontally to reduce the latency of response.

**Functional Requirements:**

1. Stores routing data with a key as a pair of latitude and longitude;
2. Expires cache data based on configurable TTL.

**Non-functional Requirements:**

1. Needs to make it simple to replace the current in-memory cache with distributed cache based on a NoSQL DB (such as Redis);
2. In-memory cache should be cleaned based on two rules - TTL and number of stored data as a preventive measure of OOM.

## Decorator of external routing services

**Functional Requirements:**

1. Requests routing data from external services;
2. Transforms a response of external service into internal representation.

**Non-functional Requirements:**

1. Needs to make it simple to replace http://project-osrm.org with another solution;
2. Applies throttling and other limitations aligning with the limits of external services.

# Building/running/deploying

## Requirements

- Go 1.21+

## Configuration

A service supports the following parameters of configuration with environment variables:

* `ROUTER_HTTP_PORT` - a port that a service listens to handle HTTP requests. Default: `8080`;
* `ROUTER_HTTP_TIMEOUT` - a timeout of handling HTTP requests. Default: `30s` (30 seconds);
* `ROUTER_OSRM_HOST` - a host of OSRM server. Default: `https://router.project-osrm.org`;
* `ROUTER_OSRM_REQUEST_TIMEOUT` - a client timeout to request OSRM server. Default: `20s` (20 seconds).

## Building/running

There it no requirement for external tools needs to be installed.

Uses the following command to run a service locally:

```bash
go run ./cmd/router/
```

## Running/deploying

Uses Docker image to run or to deploy this service.

You can use the following command to build and run this service with Docker:

```bash
docker build . -t demo-router:latest
docker run -p 8080:8080 demo-router:latest
 ```
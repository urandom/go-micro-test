# go-micro-test
A playground for exploring writing services using go-micro

## Services

Three main rudimentary services interact with one another and provide data to external users. There's an additional small http server that exposes endpoints via a simple REST api, though the `micro api` command may also be used to query the service data via its RPC handler.

#### The DB service

A simple 'db' service is provided under ./db/dummy-db, with a small set of hardcoded data. It implements two services defined using protobufs in ./db. The Todo service retrieves todo items for a given user. The User service has a method for fetching a user's profile data.

#### The Auth service

A JWT authentication service is provided under ./auth/jwt. It implements two services defined in ./auth - a Token service for generating and checking authentication tokens, and a User service for retrieving a user profile given an authentication token. The profile data is obtained from a running DB service.

#### The Todo service

The Todo service is a simple service that produces a list of todo items given an authentication token. It uses the auth service to check if the token is valid and obtain the user for the token. It then queries the DB service for that user's todo items.

#### The http router

The ./router starts an http server that proxies requests to the various services. It exposes the following endpoints:

   * GET /user - it returns a user profile as json. It expects a JWT token in the Authorization header.
   * POST /user/login - it expects a user and password parameters and returns a JWT token if the credentials are valid.
   * GET /user/todo - for a valid JWT token in the Authorization header, it returns a list of todo items formatted as json.
   
The `start-services.sh` script is a bash script that starts consul and all 4 services.

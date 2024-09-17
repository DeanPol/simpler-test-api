# simpler-test-api

This is my implementation of a RESTful JSON API CRUD microservice in Go.
Within the code you will find many comments that I have placed whilst learning about Go, as I'm not too familiar with the language.
I have created some simple routes that allow the user modify a basic Product model. The output can be viewed either via the barebones React application I have created, or through an HTTP client (Postman).
The frontend should be available on http://localhost:80. Pagination variables are set as environment variables.
I have also included some test cases, within the '/tests/' directory, that can be tested with the command `go test ./... -v`.

How to run:

- Create a new .env file following the provided .env.example and insert your database connection details.
- Have Docker installed and ready. Run `docker compose up --build` in the project root directory.
- Use the provided sql file to populate your table with some entries OR make a POST request with the following template of a body :

  {
  "name": "Sample Product",
  "description": "A description of the sample product.",
  "price": 19.99,
  "stock": 50
  }

- Navigate to your localhost to view the table entries, or use an HTTP client for more options.

# Changelog

## 0.0.1

Setting up basic directory layout, .env variables, database connectivity, schema + migration system.

## 0.0.2

Defining our model, creating our CRUD methods and registering our handlers.

## 0.0.3

Pagination and more .env variables.

## 0.0.4

Basic React template to display list of products. Routes/Handlers directory layout refactor.

## 0.0.5

Remove the total_products handler and attach the product total count and limit to the get products payload.
Added test cases.

## 0.0.6

Added docker-compose. Removed migrations as I couldn't figure out how to dockerize.

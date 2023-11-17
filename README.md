<div align="center">
  <!-- Your Project Title -->
  <h1 align="center">HTTP Check API (Golang)</h1>

  <!-- Optional: Add a brief description of your API here -->
  <p align="center">A REST API for managing HTTP checks, written in Golang.</p>

</div>
  
## Table of Contents

- [Table of Contents](#table-of-contents)
- [About](#about)
- [API Endpoints](#api-endpoints)
  - [Create HTTP Check](#create-http-check)
  - [Update HTTP Check](#update-http-check)
  - [Delete HTTP Check](#delete-http-check)
  - [Get All HTTP Checks](#get-all-http-checks)
  - [Get HTTP Check by ID](#get-http-check-by-id)
- [Authentication and Authorization](#authentication-and-authorization)
- [Getting Started](#getting-started)

## About

RESTful service to chek the health status

## API Endpoints

### Create HTTP Check

- **Endpoint:** `/http-check`
- **Method:** POST
- **Description:** Create a new HTTP check.
- **Request Body:**
  - `name` (string): Name of the HTTP check.
  - `url` (string): URL to check.
  - `interval` (integer): Time interval (in minutes) for performing the check.
  - `timeout` (integer): Timeout duration (in seconds) for the check.
- **Response:** 201 Created

### Update HTTP Check

- **Endpoint:** `/http-check/{checkId}`
- **Method:** PUT
- **Description:** Update an existing HTTP check.
- **Request Body:**
  - `name` (string): Name of the HTTP check (optional).
  - `url` (string): URL to check (optional).
  - `interval` (integer): Time interval (in minutes) for performing the check (optional).
  - `timeout` (integer): Timeout duration (in seconds) for the check (optional).
- **Response:** 200 OK

### Delete HTTP Check

- **Endpoint:** `/http-check/{checkId}`
- **Method:** DELETE
- **Description:** Delete an existing HTTP check.
- **Response:** 204 No Content

### Get All HTTP Checks

- **Endpoint:** `/http-check`
- **Method:** GET
- **Description:** Retrieve a list of all HTTP checks.
- **Response:** 200 OK

### Get HTTP Check by ID

- **Endpoint:** `/http-check/{checkId}`
- **Method:** GET
- **Description:** Retrieve details of a specific HTTP check by its ID.
- **Response:** 200 OK

## Authentication and Authorization

- Implement user authentication to control access to the Golang API.
- Ensure that only the user who created an HTTP check can update or delete it.

[Optional: Add details about your Golang authentication and authorization mechanism.]

## Getting Started

```bash


# Access the project directory
$ cd http-check

# Build and run the Golang API
$ go build
$ ./http-check
```

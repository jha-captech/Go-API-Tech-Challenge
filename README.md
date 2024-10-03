![CapTech Banner](resources/images/CaptechLogo.png)

# Go API Tech Challenge: Assignment

## Table of Contents

- [Go API Tech Challenge: Assignment](#go-api-tech-challenge-assignment)
  - [Table of Contents](#table-of-contents)
    - [Setting up a database instance](#setting-up-a-database-instance)
  - [Tech Challenge Assignment](#tech-challenge-assignment)
    - [Summary](#summary)
    - [Web Framework](#web-framework)
    - [Data](#data)
    - [Project Structure](#project-structure)
    - [`api/course`](#apicourse)
    - [`api/person`](#apiperson)
  - [Project Requirements Checklist](#project-requirements-checklist)

### Setting up a database instance

For this Tech Challenge, we will be using Postgres running in a docker container as our database.
This has already been configured for you using a docker-compose file. To start the data base, run
the following:

```bassh
make db_up
```

## Tech Challenge Assignment

### Summary

For this Tech Challenge, you will create a web API that exposes endpoints that read from and write
to a database that represents a fictional college and contains courses and people.

### Web Framework

You have the freedom to build this API with whatever tool you would like. With that being said, you
are strongly encouraged to use the standard library and/or chi where possible.

### Data

The data for this project lives inside of the database we created in the
previous section. As part of this Tech Challenge, you will need to access this data.

Please note that you will need to establish relationships between each table in order to complete
this challenge.

### Project Structure

Your project will need to define three routes, one each for `courses` and `person`. Each route will
need to include handlers for `get all`, `get by id`, `update by id`, `add`, and `delete` actions.

You should follow idiomatic go principles for your project. This includes project structure, code
organization, and naming conventions.

Your final product should include the web server, a dockerfile, and integration for that docker file
in the provided docker compose. You should also update the makefile with any needed commands to get
your app running. In addition, any documentation that an outside developer would need to get your
app up and running should be included.

Bellow are further details for each endpoint.

---

### `api/course`

| Request Type | Endpoint                              | Query Parameters | Request Body                                         | Response Type                                                        | Instructions                                                                                                                  |
|--------------|---------------------------------------|------------------|------------------------------------------------------|----------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------|
| GET          | http://localhost:8000/api/course      | *none*           | *none*                                               | JSON-formatted string representing  a list of  `Course` objects      | Return all `Course` objects from the database.                                                                                |
| GET          | http://localhost:8000/api/course/{id} | *none*           | *none*                                               | JSON-formatted string representing  a `Course` object                | Return a given `Course` object based on `id`.                                                                                 |
| PUT          | http://localhost:8000/api/course/{id} | *none*           | JSON-formatted string representing a `Course` object | JSON-formatted string representing  an updated `Course` object       | Update a given `Course` object in the database based on `id`. The `Course` object passed to the endpoint should be validated. |
| POST         | http://localhost:8000/api/course      | *none*           | JSON-formatted string representing a `Course` object | JSON-formatted string representing  a the new `Course` object's `id` | Add a new `Course` object to the database. `id` does not need to be provided as the database will generate it.                |
| DELETE       | http://localhost:8000/api/course/{id} | *none*           | *none*                                               | JSON-formatted string representing  a deletion confirmation message  | Delete a given `Course` object from the database based on `id`.                                                               |

Here is the schema for a `Course` object
| Column Name | Column Type |
| ----------- | ----------- |
| `id`        | integer |
| `name`      | string |

---

### `api/person`

| Request Type | Endpoint                                | Query Parameters                 | Request Body                                         | Response Type                                                        | Instructions                                                                                                                                                                                             |
|--------------|-----------------------------------------|----------------------------------|------------------------------------------------------|----------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| GET          | http://localhost:8000/api/person        | `name`: string<br>`age`: integer | *none*                                               | JSON-formatted string representing  a list of  `Person` objects      | Return all `People` objects from the database. If query parameters are passed to the endpoint, filter off of them.                                                                                       |
| GET          | http://localhost:8000/api/person/{name} | *none*                           | *none*                                               | JSON-formatted string representing  a `Person` object                | Return a given `Person` based off of `name`.                                                                                                                                                             |
| PUT          | http://localhost:8000/api/person/{name} | *none*                           | JSON-formatted string representing a `Person` object | JSON-formatted string representing  an updated `Person` object       | Update a given `Person` in the database based on `name`. The `Person` object passed to the endpoint should be validated.                                                                                 |
| POST         | http://localhost:8000/api/person        | *none*                           | JSON-formatted string representing a `Person` object | JSON-formatted string representing  a the new `Person` object's `id` | Add a new `Person` to the database. `id` does not need to be provided as the database will generate it. If any `Course` objects `id`s are passed in, that association should be updated in the database. |
| DELETE       | http://localhost:8000/api/person/{name} | *none*                           | *none*                                               | JSON-formatted string representing  a deletion confirmation message  | Delete a given `Person` object from the database based on `name`.                                                                                                                                        |

Here is the schema for a `Person` object:
| Column Name | Column Type | Notes |
| ------------ | ------------------ | -------------------------------------------- |
| `id`         | integer | primary key |
| `first_name` | string | N/A |
| `last_name`  | string | N/A |
| `type`       | string | possible values are `student` and `professor` |
| `age`        | integer | N/A |
| `courses`    | list of integers | list of course ids |

---

## Project Requirements Checklist

- [ ] Your API should use port `8000`.
- [ ] Your API should have a single entry point.
- [ ] Each endpoint should return the appropriate statues code with each response.
- [ ] If an error is encountered by the application, an informative error message should be returned
  to the client and your application should also log details.
- [ ] Your project should include unit tests with 80% unit test coverage.

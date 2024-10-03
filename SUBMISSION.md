
# John Flynn Tech Challenge Submission

## MakeFile

Due to this being worked on using a captech device, docker cannot be installed due to licensing. `podman` has been used in place. 
see podman-makefile a second makefile updated with podman specific commands. I.E `make -f podman-makefile run`. 

`make run` runs the application locally, but will still require a database instance to be running `make db_up`

> **Note:** When running the database be sure to remove old instances of `courses-db-container` as the schema has been changed!


The app can also be run locally using go directly.
```bash
DATABASE_NAME=coursesDB DATABASE_USER=courses-db-user DATABASE_PASSWORD=courses-db-password DATABASE_HOST=localhost DATABASE_PORT=5432 go run cmd/api/main.go 
```

## Docker

To build the docker image for this project, use the make command 

```bash
make build_image
```
This will build a docker image to your local registry. 

Next, copy the .env.local file to .env ```cp .env.local .env```

Now you can start the database and api using the make command

```bash
make run_app
```

## Code Overview

The external libraries used by this project has been kept to a minimum. 

- Http Routing is done using the standard net/http library. 

- [FX](https://github.com/uber-go/fx) a dependency injection library is used to manage singletons and app dependencies

- [Gorm](https://github.com/go-gorm/gorm ) is used as an ORM solution for database querying

- [Google uuid](https://github.com/google/uuid) for generating `guids`. 

- [testify](https://github.com/stretchr/testify) is also used for unit testing and creating mocks. 

- mockery was used to generate mocks. 

I made a  minor change to the schema and routing in the requirements. Most noteably, the introduction of a `Guid` field in place of the `id` and `name` in the route of `/course/{id}` and `/person/{name}`.

This project is composed of a few different packages with dedicated responsibilities. 

`apperror` - specific error types for handling errors and http status codes. Also includes utilities for handling unexpected error cases, multi errors, and ensure that stack trace data is not leaked to the caller. 

`applog` - Wrapper around default logging library to add colors. 

`config` - utility for converting environment variables into a configuration struct

`database` - function for handling creating and connecting to the database

`handler` - contains http.Handler functions to be used by ServeMux net/http packages. Handles encoding and decoding logic.

`models` - Gorm Data Model structs.

`services` - Contains business logic functions for validating data and calling repository functions.

`repository` - Contains functions for saving, updating and retrieving data from gorm. 


# API Overview and examples

# Models

### Course

|Key|Type|                      |
|---|----|----------------------|
| Guid| string| Identifier of course |
|Name| string| Name of course       |

### Course Input

|Key|Type|  |
|---|----|-------------|
|Name| string| Must Not Be Blank|


### Person

| Key       | Type                     |                                             |
|-----------|--------------------------|---------------------------------------------|
| Guid      | string                   | Person Unique Identifier                    |
| FirstName | string                   | First Name                                  |
| LastName  | string                   | Last Name                                   |
| Email     | string                   | Email Address                               |
| Age       | int                      | Age                                         |
| Type      | string                   | "student" or "professor"                    |
| Courses   | array<[Course](#Course)> | array of of courses the user is enrolled in |


### Person Input


| Key          | Type          |                                        |
|--------------|---------------|----------------------------------------|
| FirstName    | string        | Must Not Be Blank                      |
| LastName     | string        | Must Not Be Blank                      |
| Email        | string        | Must Be a valid email address          |
| Age          | int           | Must be at least 10                    |
| Type         | string        | Must be either "student" or "professor" |
| Course Guids | array[string] | array of strings of course guids for the user to be enrolled in |

# Error Handling

This project will handle providing all errors to the user if they make multiple mistakes in a request.
For example if a user provides a single incorrect query parameter when querying `/api/person`
I.E

`/api/person?FirstNam=Foo`

They will get the following 400 error.
```json
{"Message":"Invalid Request Parameter: FirstNa"}
```
If they provide > 1 incorrect query parameter
I.E
`/api/person?FirstNam=Foo&SomethingElse=Bar&anotherError=FooBar`

The response will be a multi error response
```json
{
    "Message": "Multiple Errors:",
    "Errors": [
        {
            "Message": "Invalid Request Parameter: anotherError"
        },
        {
            "Message": "Invalid Request Parameter: FirstNam"
        },
        {
            "Message": "Invalid Request Parameter: SomethingElse"
        }
    ]
}
```

The same error convention applies to data validation on *PUTS* and *POSTS*

For example, Creating a person with a blank last name and incorrect email address will result in the following

```json
{
    "Message": "Multiple Errors:",
    "Errors": [
        {
            "Message": "Last Name must not be blank"
        },
        {
            "Message": "Email must be a valid email address"
        }
    ]
}
```

## Course

### Get All Courses

*GET* `/api/course` 

Returns a list of courses. 

`Query Params`
- Name: (Optional), wild card search for courses by name

Returns: 
- 200 ok with courses if query was successful
- 400 bad request if any query parameters are invalid
- 500 internal server error if any errors occurred

Responds with: Array<[Course](#course)>

Example

```bash
curl 'http://localhost:8000/api/course?Name=Programming'
```
Response
```json
[
    {
        "Guid": "123a",
        "Name": "Programming"
    }
]
```

### Get Course By Guid 

*GET* `/api/course/{guid}`

Returns a single course by guid.

`Query Params` _None_

Returns:
- 200 ok with course if id is valid
- 404 not found if the course guid does not exist.
- 500 internal server error if any errors occurred.

Responds with: [Course](#course)

Example

```bash
curl 'http://localhost:8000/api/course/456b'
```

Response

```json
{
    "Guid": "456b",
    "Name": "Databases"
}
```

### Add Course

*POST* `/api/course`

Adds a new course.

`Query Params` _None_

`Request Body`

|Key|Type|  |
|---|----|-------------|
|Name| string| Must Not Be Blank|

Returns 
- 201 created if course is created successfully
- 400 bad request if the name is blank
- 500 internal server error if any issues occurred

Accepts: [CourseInput](#course-input)

Responds with: [Course](#course)

Example 

```bash
curl http://localhost:8000/api/course -X POST -H 'Content-Type: application/json' -d '{"Name": "How to turn a 44 billion dollar company into a $9.5 billion dollar company"}'
```

Response 

```json
{
    "Guid": "b91c7fdf-53ed-4114-a3a0-50401607f93e",
    "Name": "How to turn a 44 billion dollar company into a 9.5 billion dollar company"
}
```

### Update Course

*PUT* `/api/course/{guid}`

Updates a course by `guid'

`Query Params` _None_

`Request Body`


Returns
- 200 if updated successfully
- 404 if course by guid is not found
- 400 if name is blank
- 500 if an internal error occurred

Accepts: [CourseInput](#course-input)

Responds with: [Course](#course)

Example

```bash
curl http://localhost:8000/api/course/123a -X PUT -H 'Content-Type: application/json' -d '{"Name":"Programming In Golang"}'
```
Response
```json
{
    "Guid": "123a",
    "Name": "Programming In Golang"
}
```

### Delete Course

*DELETE* `/api/course/{guid}`

Deletes a course. Will also un enrol any people that were linked to that course.

`Query Params` _None_

Returns
- 200 if deleted successfully
- 404 if guid is not found
- 500 if an internal error occurred


Responds with: String

String - "OK" confirmation message.

Example

```bash
curl http://localhost:8000/api/course/789c -X DELETE
```

Response
```
"OK"
```

## Person

### Get All People

*GET* `/api/person`

Returns a list of people. Query Parameters can be combined to search results.

`QueryParams`
- FirstName: (Optional) Wild card search for people with matching first name
- LastName:  (Optional) Wild card search for people with matching last name
- Email:     (Optional) Wild card search for people with matching email

Returns: 
- 200 ok with people if query was successful
- 400 bad request if any query parameters are invalid
- 500 internal server error if any unexpected errors occurred

Responds with: Array<[Person](#person)>

Example

```bash
curl 'http://localhost:8000/api/person?FirstName=J&Email=test'
```

Response

```json
[
    {
        "Guid": "efgh",
        "FirstName": "Jeff",
        "LastName": "Bezos",
        "Email": "jbezos@test.com",
        "Age": 60,
        "Type": "professor",
        "Courses": null
    },
    {
        "Guid": "uvwx",
        "FirstName": "John",
        "LastName": "Flynn",
        "Email": "jflynn@test.com",
        "Age": 52,
        "Type": "student",
        "Courses": null
    }
]
```

### Get Person By guid

*GET* `/api/person/{guid}`

Return a single person with matching `guid` and their courses.

QueryParams: _None_

Returns:
- 200 ok if found
- 404 if person by guid does not exist
- 500 internal server error if any unexpected errors occurred

Example
```bash
curl http://localhost:8000/api/person/abcd
```

Accepts: [PersonInput](#person-input)

Response

```json
{
    "Guid": "abcd",
    "FirstName": "Steve",
    "LastName": "Jobs",
    "Email": "sjobs@test.com",
    "Age": 56,
    "Type": "professor",
    "Courses": [
        {
            "Guid": "123a",
            "Name": "Programming In Golang"
        },
        {
            "Guid": "456b",
            "Name": "Databases"
        }
    ]
}
```

### Add Person

*POST* `/api/person`

Create a new person with courses.

QueryParams: _None_

Returns:
- 201 created if the person is created successfully
- 400 bad request if validation rules fail
- 500 internal server error if any unexpected errors occurred

Accepts: [PersonInput](#person-input)

Responds with: [Person](#person)

Example

```bash
curl http://localhost:8000/api/person -X POST -H 'Content-Type: application/json' -d '{"FirstName": "Foo", "LastName": "Bar", "Email": "foobar@test.com", "Age": 102, "Type": "student", "CourseGuids": ["456b", "123a"]}'
```

Response

```json
{
  "Guid": "681d0819-9d46-473d-bc24-2e64cb3a76dd",
  "FirstName": "Foo",
  "LastName": "Bar",
  "Email": "foobar@test.com",
  "Age": 102,
  "Type": "student",
  "Courses": [
    {
      "Guid": "456b",
      "Name": "Databases"
    },
    {
      "Guid": "123a",
      "Name": "Programming In Golang"
    }
  ]
}
```

### Update Person

*PUT* `/api/person/{guid}`

Updates a person with matching `guid`. Accepts a list of `CourseGuids` to be updated to the person.
Courses will be added or removed from the person based on the values provided.

QueryParams: _None_

Returns:
- 200 ok if the person and courses were updated
- 404 if person with guid does not exist or any course guids do not exist
- 400 if any validation rules fail.
- 500 internal server error if any unexpected errors occurred

Accepts: [PersonInput](#person-input)

Responds with: [Person](#person)

Example
```bash
curl http://localhost:8000/api/person/qrst -X PUT -H 'Content-Type: application/json' -d '{"FirstName":"Elon","LastName":"Musk","Email":"misinformationmaster@test.com","Age":14,"Type":"professor","CourseGuids":["b91c7fdf-53ed-4114-a3a0-50401607f93e"]}'
```

Response
```json
{
  "Guid": "qrst",
  "FirstName": "Elon",
  "LastName": "Musk",
  "Email": "misinformationmaster@test.com",
  "Age": 14,
  "Type": "professor",
  "Courses": [
    {
      "Guid": "b91c7fdf-53ed-4114-a3a0-50401607f93e",
      "Name": "How to turn a 44 billion dollar company into a $9.5 billion dollar company"
    }
  ]
}
```

### Delete Person

*DELETE* `/api/person/{guid}`

Deletes a person by `guid`. Un enrols them from any courses they were in.

`Query Params` _None_

Returns
- 200 if deleted successfully
- 404 if guid is not found
- 500 internal server error if any unexpected errors occurred

Responds with: String

Example

```bash
curl 'http://localhost:8000/api/person/yzab' -X DELETE
```

Response

```json
"OK"
```
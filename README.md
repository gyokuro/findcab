![Image](./screen1.png?raw=true "screenshot")

Available at [http://192.241.249.201]

# Overview

Should you choose to accept it, your challenge is to build an HTTP server that can quickly find
the nearest cabs to a particular client location. Utilize python and you can use a web service
framework of your choice. Your server must be able to handle the following CRUD operations:

# 1 Create / Update request

Insert a new record or update the latitude and longitude of a particular cab

HTTP Method | URL | Returns | Normal Response
:----------:|:---:|:-------:|:--------------:
PUT | /cabs/(cab_id) | No body | 200 OK

## Parameters

Required | Name | Type | Default | Description | Example
:-------:|:----:|:----:|:-------:|:-----------:|:------:
Yes | latitude | Float | | GPS |
Yes | longitude | Float | | GPS |

# 2 Get request

Get the full details of the cab

HTTP Method | URL | Returns | Normal Response
:----------:|:---:|:-------:|:--------------:
GET | /cabs/(cab_id) | Cab | 200 OK

## Parameters

This action does not have any parameters.

# 3 Query request

Search for nearest cabs. The returned data does not have to be sorted and should be a list of
cab records. DIstance should be calculated based on the Haversine formula.

HTTP Method | URL | Returns | Normal Response
:----------:|:---:|:-------:|:--------------:
GET | /cabs | Cab | 200 OK

## Query Parameters

Required | Name | Type | Default | Description | Example
:-------:|:----:|:----:|:-------:|:-----------:|:------:
Yes | latitude | Float | |  Client GPS |
Yes | longitude | Float | | Client GPS |
|| limit | Integer | 8 | The total |
|| radius | Float || The radius |

## Sample Request

    GET /cabs?latitude=37.763658&longitude=-122.427521=&radius=1000&limit=15

## Sample Response

    [
     {
       "id": 23706134,
       "latitude": 37.788654783559,
       "longitude": -122.50747748978
     },
     {
       "id": 61344818,
       "latitude": 37.778952285851,
       "longitude": -122.43865835511
     },
     {
       "id": 19485186,
       "latitude": 37.778665475753,
       "longitude": -122.39094602609
     },
     ...
    ]

# 4 Destroy request

Destroy a cab

HTTP Method | URL | Returns | Normal Response
:----------:|:---:|:-------:|:--------------:
DELETE | /cabs/(cab_id) | No body | 200 OK

## Parameters

This action does not have any parameters.

# 5 Destroy all request

Destroy all cab records

HTTP Method | URL | Returns | Normal Response
:----------:|:---:|:-------:|:--------------:
DELETE | /cabs | No body | 200 OK

## Parameters

This action does not have any parameters.


# Evaluation / Additional Notes

Beyond these basic API specs, the challenge is left intentionally open-ended. The are no
requirements for which coding language, database, or design you must use, but we do expect
that you can defend the decisions you make. The best solutions will be those that go beyond
something that “just-works” and have thought about the implications of real-world production
environment. For example, which endpoints do you think will be used most frequently?

All latitude and longitude coordinate are based on the WGS84 coordinate system (the same one
used by GPS enabled devices).

# Why GO

Go is used as the language of implementation because:

+ It's statically typed language but fast compilation makes it feels like a scripting language,
while allowing the code base to scale gracefully over time.
+ Its performance approaches that of C.
+ Manages concurrency elegantly thanks to language constructs like go routines and channels.
+ Compiles into a single binary executable. No other dependencies!  This makes production deployment trivial:
    + No dependencies on runtime libraries that needs to be installed on each VM host
    + Even for a webapp, the static images, css, js files can be compiled into go source code and embedded in the binary.

# Compiling and Building

Steps to build:

+ Install Go

+ For easier deployment, the files required for a web UI (e.g. .js and .css files, etc.) can
be compiled into the executable binary.  To do so, a `embedfs` utility (written by me) is run
that will generate .go source files for each one of the static content files:

    util/embedfs -match=".*" -generated=true -destDir=resources webapp/

This will generate .go source files for each file matched under the `webapp` directory.  The
go source files will be placed in the `resources/webapp` directory.

Note that util/embedfs is a binary compiled for MacOSX.  For linux, you can build embedfs from
source, from my other github project: [https://github.com/qorio/embedfs]

+ From the top level directory, build the binary:

    go build main/findcab.go

This will produce the final executable, `findcab`, in the current directory.

# Test

There are unit tests in different packages that can be run to verify functionality.  Do

    go test -v

Also, the server can be tested using curl:

    curl -i -H 'Accept:application/json' -X PUT -d '{"id":1234, "latitude":55.00, "longitude":55.00}' "http://localhost:8080/cabs/1234"

    curl -i -X GET "http://localhost:8080/cabs/1234"
    HTTP/1.1 200 OK
    Content-Type: application/json
    Content-Length: 40
    Date: Tue, 17 Dec 2013 20:33:22 GMT

    {"id":1234,"latitude":55,"longitude":55}


    curl -i -X GET "http://localhost:8080/cabs?latitude=55&longitude=55&radius=500"
    HTTP/1.1 200 OK
    Content-Type: application/json
    Content-Length: 42
    Date: Tue, 17 Dec 2013 20:36:02 GMT

    [{"id":1234,"latitude":55,"longitude":55}]

    curl -i -X GET "http://localhost:8080/cabs?latitude=55&longitude=65&radius=500"
    HTTP/1.1 200 OK
    Content-Type: application/json
    Content-Length: 2
    Date: Tue, 17 Dec 2013 20:36:39 GMT

    []

    curl -i -X DELETE "http://localhost:8080/cabs/1234"
    curl -i -X DELETE "http://localhost:8080/cabs"

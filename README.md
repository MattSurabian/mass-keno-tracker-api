# Mass Keno Tracker API

This repo contains a relatively simplistic REST-like API application written in 
[Golang]. The structure of the repo and overall developer workflow/deployment 
pipeline is based on the [go-build-template project].

The production version of this application is hosted on [Joyent's Triton Platform] 
and can be accessed via: [masskenotracker.com]

## Quick Start

### Get the code

Thanks to Go's built in vendoring support you can clone or unzip this repo into 
any directory on your system and effectively interact with this project. If you 
prefer to keep things "by the book" follow the standard Go pattern and put this 
repo in:

```
$GOPATH/src/github.com/mattsurabian/mass-keno-tracker
```

### Build & Run

From inside the root directory of this repository run `make`. 

That will compile the application in an ad-hoc Docker container, write the 
resulting binary to the `bin` directory, use the `Dockerfile` to compose a 
Docker image containing that binary, and run that image in a new container. 
All necessary dependencies like a Redis store will be spun up in companion 
Docker containers, automatically configured and populated with test data.

If successful you should see:

```
Mass Keno Tracker API now running on port 8088
curl localhost:8088/api/v1/health && echo
{
    "status_code": 200,
    "status_text": "OK"
}

```

To stop the api use `make stop`. To re-build and restart the api use `make start`.

**_To cleanup build pipeline artifact sprawl in your environment run:_** `make clean`

### Dependencies

All dependencies are managed and vendored using [gvt].

## The Stack

This application is written in Go and utilizes the [Gin HTTP framework] and 
[Manners HTTP server]. It also leverages [Redis] as a data store. [Docker networking]
ties things together for local development.

## The Network

This application uses [Docker networking] to provide an isolated network environment
for both local development convenience and production mirroring. By default the 
network allows all inter-container communication.





[Golang]: https://golang.org/
[go-build-template project]: https://github.com/thockin/go-build-template
[Joyent's Triton Platform]: https://www.joyent.com/triton/compute
[masskenotracker.com]: http://masskenotracker.com
[gvt]: https://github.com/FiloSottile/gvt
[Gin HTTP framework]: https://github.com/gin-gonic/gin
[Manners HTTP server]: https://github.com/braintree/manners
[Redis]: https://redis.io/
[Docker networking]: https://docs.docker.com/engine/userguide/networking/
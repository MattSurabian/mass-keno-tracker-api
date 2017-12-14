# Mass Keno Tracker API

This repo contains a relatively simplistic REST-like API application written in 
[Golang]. The structure of the repo and overall developer workflow/deployment 
pipeline is based on the [go-build-template project].

The production version of this application ~~is~~ will be hosted on [Joyent's Triton Platform] 
and available at: [masskenotracker.com]

## Why Does This Repo Exist?

This project is a way for me to demonstrate a few things I've learned writing
APIs and ETL pipelines in Go with Docker. Specifically the initial API server 
setup with Gin and Manners; along with the concurrency setup for ingesting
historic data found in the `keno-tracker-etl` package. 

This repo also provides an implementation of a workflow that wrangles dependent 
technologies like Redis or Elasticsearch using Docker.

The idea being that when folks ask me questions about this kind of stuff, I can point to
some code examples and talk in the context of an actual application instead of 
having to hand wave; which I already do enough of.

That said, not everything in this repo is gold. It's just a toy application I like to 
play around with in my free time and use with friends at dive bars I occasionally 
find myself at. So you know, your mileage may vary.

Hopefully something in here is useful to you even if you don't much care about Keno.
Feel free to shoot questions my way via Github issues.

## Quick Start

### Get the code

Thanks to Go's built in vendoring support you can clone or unzip this repo into 
any directory on your system and effectively interact with this project. If you 
prefer to keep things "by the book" follow the standard Go pattern and put this 
repo in:

```
$GOPATH/src/github.com/mattsurabian/mass-keno-tracker-api
```

### Build & Run

**_I developed this in a Linux environment, so it's possible things like Docker volume sharing and user id mapping won't work quite as well if you're on Windows or Mac. Please open an issue if you encounter problems building this!_**

From inside the root directory of this repository run `make`. 

That will compile the application in an ad-hoc Docker container, write the 
resulting binary to the `bin` directory, use the `Dockerfile` to compose a 
Docker image containing that binary, and run that image in a new container. 
All necessary dependencies like a Redis store will be spun up in companion 
Docker containers, automatically configured and populated with test data.

To stop the api use `make stop`. To re-build and restart the api use `make start`.

**_To cleanup build pipeline artifact sprawl in your environment run:_** `make clean`

**_To completely tear down the API server, redis cache, docker network, and all build artifacts; effectively starting fresh, run:_** `make destroy`

### Other Commands
The make file also provides specific control helpers (most using the shell scripts inside `tasks`) to do thinks like start, stop, and remove redis `redis-start, redis-stop, redis-rm`; compile just a binary `make build`; and build a container for the API `make container` to name a few.

### Dependencies

All dependencies are managed and vendored using [gvt].

## The Stack

This application is written in Go and utilizes the [Gin HTTP framework] and 
[Manners HTTP server]. It also leverages [Redis] as a data store. [Docker networking]
ties things together for local development. At boot time all historic data is ingested 
into Redis and backups are cut to a git ignored directory to ensure fast restarts. 
In the near future Elasticsearch will be integrated to demonstrate how to ingest 
from external services while effectively leveraging network request caching.

## The Network

This application uses [Docker networking] to provide an isolated network environment
for both local development convenience and production mirroring. By default the 
network allows all inter-container communication.

## The API

The API is definitely a work in progress. Right now four types of endpoints are 
supported, but that will change as a front-end consumer for this API is written.

In addition to supporting the endpoints documented below, the API server also 
performs a self check and updates it's health status every minute, and kicks off
it's ETL pipeline ingest every 12 hours; making this totally self contained and 
not reliant on external job servers or cron. The benefit of this is the ability 
in to deploy this inside a scratch Docker container that has no underlying
operating system in the future.

### Todays `/v1/todays`

This endpoint returns the latest Keno manifest from today with draws
in the form:

```
{
    "min": "1939153",
    "max": "1939322",
    "date": "2017-08-13",
    "draws": [
        {
            "draw_id": "1939153",
            "draw_date": "2017-08-13",
            "day_normalized_id": "0",
            "winning_num": "03-13-18-19-31-38-41-44-46-50-51-53-59-62-66-68-69-73-75-78",
            "bonus": "3x"
        },
        ...
    ]
}
```

### History `/v1/history/YYYY/MM/DD`

This endpoint returns a full days worth of data in the form:

```
{
    "min": "1935911",
    "max": "1936210",
    "date": "2017-08-02",
    "draws": [
        {
            "draw_id": "1935911",
            "draw_date": "2017-08-02",
            "day_normalized_id": "0",
            "winning_num": "01-09-13-14-18-20-28-32-33-34-37-49-56-60-61-64-69-73-75-77",
            "bonus": "3x"
        },
        {
            "draw_id": "1935912",
            "draw_date": "2017-08-02",
            "day_normalized_id": "1",
            "winning_num": "04-18-24-26-28-32-35-37-41-50-60-61-62-67-68-71-74-75-76-77",
            "bonus": "3x"
        },        
        ...
    ]
}
```

### Occurences `/v1/occurences/FULL-DRAW`

Returns all instances of a given draw in the form:

**Request URL:** ` /v1/occurences/02-04-06-07-12-15-22-23-25-26-30-33-38-43-46-49-57-61-74-77`

**Response:**

```
{
    "winning_num": "02-04-06-07-12-15-22-23-25-26-30-33-38-43-46-49-57-61-74-77",
    "draws": [
        {
            "draw_id": "1938254",
            "draw_date": "2017-08-10",
            "day_normalized_id": "17",
            "winning_num": "02-04-06-07-12-15-22-23-25-26-30-33-38-43-46-49-57-61-74-77",
            "bonus": "No Bonus"
        }
    ]
}
```

Future versions of this endpoint will support wild carding, so you can check how
many games your specific number has come up in. For example in the future you could 
pass `01-02-03-04*` to see every game with draws starting with the numbers one through four.
Currently only a full draw of twenty numbers is supported.

### Health `/v1/health` or `/health`

Currently there are no real version specific health checks in place. So you will
get the same status whether you request the versioned or unversioned health endpoint.
Currently the API server performs a self health assessment every minute and updates
its status accordingly. At present the only factor affecting service health is whether
or not the cache is "volatile". The cache becomes volatile when an ETL process is 
running because the service cannot guarantee it will be able to return complte
results. During periods of volatility the health endpoints will return 503. Individual
endpoints may still respond but the data will be of unknown quality.

[Golang]: https://golang.org/
[go-build-template project]: https://github.com/thockin/go-build-template
[Joyent's Triton Platform]: https://www.joyent.com/triton/compute
[masskenotracker.com]: http://masskenotracker.com
[gvt]: https://github.com/FiloSottile/gvt
[Gin HTTP framework]: https://github.com/gin-gonic/gin
[Manners HTTP server]: https://github.com/braintree/manners
[Redis]: https://redis.io/
[Docker networking]: https://docs.docker.com/engine/userguide/networking/
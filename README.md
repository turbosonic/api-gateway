# Turbosonic: api-gateway
A lightweight, sub-milisecond api-gateway intended for microservices on docker

## What does it do?
Acts as the conduit between the outside world and your internal microservice ecosystem, using a simple yaml configuration file you can set up your routing to take external requests and forward them on to internal endpoints.

## Getting started
### Create a configuration file
You need to create a configuration yaml file to be used when the api gateway first starts, full details are in the wiki, but here's a starter for 10 (it's the same as `config.yaml` on this project's root).

```yaml
name: /web/v1 
endpoints:
  - url: /example
    methods:
      - method: GET
        destination:
          name: service
          host: http://my-first-service
          url: /things
```

This configuration will expose one endpoint, `web/v1/example` (which, when following these instructions, can be accessed via `http://localhost:8080/web/v1/example`) which will allow access to the `/things` endpoint on the internal service hosted at `http://my-first-service`
### Docker networking
If you have some other containers running APIs which you want to expose using the gateway, you'll need to create a network, that's pretty easy.
```bash
$ docker network create turbosonic
```

### Run the docker image
All of Turbosonic's services run on a scratch docker image, which means they're very small and incredibly lightweight (7mb of magic!), you'll need a linux docker machine running, then run the following command:
```bash
$ docker run -d -p 8080:8080 --name api-gateway --net turbosonic -v /d/go/src/github.com/turbosonic/api-gateway/config.yaml:/config.yaml turbosonic/api-gateway
```
This is just standard docker stuff:
* `-d` is to run to in the background (daemon mode)
* `-p {external port}:8080` will expose the gateway to the port you add
* `--name {a nice name}` gives the container a nice name
* `-net {name of docker network}` the name of the docker network (created in the last step)
* `-v {path to yaml file}:/config.yaml` mounts your config file in to the container to be used
* `turbosonic/api-gateway` is the name of the image on docker hub

## Things to come
This is the start of a journey to create a simple, secure, scalable and production ready api-gateway, to allow developers to focus on the core functionality of their systems, below is a list of things to be added in the coming months.

### Coming soon...

* **Authentication** - it would be pretty useless if a gateway didn't keep out the bad guys
* **Authorization**  - limit access to certain endpoints based on the claims of the user accessing it
* **Request bundling** - collate multiple internal requests in to one external request, improving latency for devices with shoddy internet connections
* **Multiple configurations** - so you can have `/web/v2` and `/mobile/v1` etc, running simlutaneously
* **Caching** - hold responses in memory for ultra fast response times
* **Mocking** - return mocked data for fast prototyping and concurrent development of front and back end
* **Rate limiting** - Protect your internal endpoints by adding a forcefield when users start spamming you
* **Retries** - If at first you don't succeed; try, try again
* **Logging** - To monitor usage, find errors, and to show off how fast your API is.

### A bit further away...
* **Websockets** - allow clients to connect to your application via websockets rather than http, this is when things get silly, silly fast.
* **Remote configuration** - create your APIs via a web interface, then just run the docker image, it will go pick up your configuration from the web (securely) to make life so much easier.
* **Self documentation** - Swagger on demand

## Want to help?
This is our first foray in to the world of Golang, we used to do this stuff in node, but Go has blown our socks off with its simplicity, speed and itsy bitsy teeny weeny container images.

That means we could probably do with a bit of help, take a fork and create a pull request, but please, if you do, add and link an issue to what you are doing, we won't accept it otherwise!
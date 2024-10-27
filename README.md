# Reverse Proxy

This is a simple HTTP reverse proxy written in Go as part of my technical challenge for Traefik

To test the code locally, do the following:

### 1. Start the demo server running traefik/whoami

```
make start
```

This will start a container running traefik/whoami on localhost:9000. Try accessing this with curl to see if it's up.

```bash
$ curl -v localhost:9000
*   Trying [::1]:9000...
* Connected to localhost (::1) port 9000
> GET / HTTP/1.1
> Host: localhost:9000
> User-Agent: curl/8.4.0
> Accept: */*
>
< HTTP/1.1 200 OK
< Date: Sat, 26 Oct 2024 16:47:38 GMT
< Content-Length: 160
< Content-Type: text/plain; charset=utf-8
<
Hostname: 1bab8fed0c66
IP: 127.0.0.1
IP: 172.17.0.2
RemoteAddr: 192.168.65.1:31331
GET / HTTP/1.1
Host: localhost:9000
User-Agent: curl/8.4.0
Accept: */*
```

### 2. Start the reverse proxy server

```bash
go run cmd/reverse-proxy/main.go
```

The reverse proxy is configured to run on port 8080. This is currently configured from the `config.yaml` file. You can define routers and services from this file. An example config.yaml file is as shown below:

```yaml
port: 8080

routers:
  - host: iamfoo.localhost:8080
    service: whoami
  - host: testing.localhost:8080
    service: whoami

services:
  - name: whoami
    url: "http://localhost:9000"
```

Under services you can define your services to which you want to proxy requests to. These are like services in the Traefik configuration. Under routers you can define how to map different urls to these services. This is similar to the Traefik router configuration or the IngressRoute CRD in a way. 

If you access `iamfoo.localhost:8080` after the reverse proxy is started with the above configuration, your request will be proxied to traefik/whoami running on `localhost:9000` and you'll get the response from traefik/whoami. You can try this for the different paths served by traefik/whoami such as `/api`, `/bench` and so on across the different routers that you've defined.
 
```bash
$ curl -v localhost:8080
*   Trying [::1]:8080...
* Connected to localhost (::1) port 8080
> GET / HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/8.4.0
> Accept: */*
>
< HTTP/1.1 200 OK
< Content-Length: 183
< Content-Type: text/plain; charset=utf-8
< Date: Sat, 26 Oct 2024 16:49:40 GMT
<
Hostname: 1bab8fed0c66
IP: 127.0.0.1
IP: 172.17.0.2
RemoteAddr: 192.168.65.1:31452
GET / HTTP/1.1
Host: localhost:9000
User-Agent: curl/8.4.0
Accept: */*
Accept-Encoding: gzip
```

### 3. Cleanup

Cleanup the traefik/whoami container.

```bash
make stop
```

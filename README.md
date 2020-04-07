# redir
tiny http redirection using badger as KV store

## Setup

```
docker run -it --rm -v "$(PWD)/db:/app/db" -p "5545:5545" yulrizka/redirect:latest
```

## Operation
## Add
```
$ curl -v -X POST "localhost:5545/add?source=http://localhost:5545/blog&target=https://labs.yulrizka.com"
```

redirecting `http://localhost:5545/blog` -> `https://labs.yulrizka.com`
```
curl -v localhost:5545/blog
*   Trying ::1...
* TCP_NODELAY set
* Connected to localhost (::1) port 5545 (#0)
> GET /blog HTTP/1.1
> Host: localhost:5545
> User-Agent: curl/7.54.0
> Accept: */*
>
< HTTP/1.1 307 Temporary Redirect
< Content-Type: text/html; charset=utf-8
< Location: https://labs.yulrizka.com
< Date: Tue, 07 Apr 2020 17:25:21 GMT
< Content-Length: 61
<
<a href="https://labs.yulrizka.com">Temporary Redirect</a>.
```

## List
```
$ curl -v -X POST localhost:5545/list

...
< HTTP/1.1 200 OK
< Date: Tue, 07 Apr 2020 17:04:09 GMT
< Content-Length: 70
< Content-Type: text/plain; charset=utf-8
<
1 entr(ies)
[  1] http://localhost:5545/ -> https://labs.yulrizka.com
```

## Delete
```
$ curl -v -X DELETE localhost:5545/blog
> DELETE /blog HTTP/1.1
> Host: localhost:5545
> User-Agent: curl/7.54.0
> Accept: */*
>
< HTTP/1.1 200 OK
< Date: Tue, 07 Apr 2020 17:27:22 GMT
< Content-Length: 10
< Content-Type: text/plain; charset=utf-8
<
* Connection #0 to host localhost left intact
deleted ok⏎
```

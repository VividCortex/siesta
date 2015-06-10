Authentication example
===
This program demonstrates the use of Siesta's contexts and middleware chaining to handle authentication. In addition, there are also other features like request identification and logging that are extremely useful in practice.

Suppose we have some state with the following data:

| Token | User |
| ----- | ---- |
| abcde | alice |
| 12345 | bob |

| User | Resource ID | Value |
| ---- | ----------- | ----- |
| alice | 1 | foo |
| alice | 2 | bar |
| bob | 3 | baz |

Users of the API have to supply a valid token to be able to access the secured resources that they are assigned to.

There is a single endpoint: `GET /resources/:resourceID`

The token will be provided by the user for every request as the HTTP basic authentication username. This is similar to [Stripe](https://stripe.com/docs/api#authentication)'s API authentication.

Example requests
---
```
$ curl -i localhost:8080
HTTP/1.1 401 Unauthorized
X-Request-Id: 4d65822107fcfd52
Date: Wed, 10 Jun 2015 13:03:36 GMT
Content-Length: 27
Content-Type: text/plain; charset=utf-8

{"error":"token required"}
```

```
$ curl -i localhost:8080/resources/1 -u abcde:
HTTP/1.1 200 OK
Content-Type: application/json
X-Request-Id: 55104dc76695721d
Date: Wed, 10 Jun 2015 13:04:23 GMT
Content-Length: 15

{"data":"foo"}
```

```
$ curl -i localhost:8080/resources/3 -u 12345:
HTTP/1.1 200 OK
Content-Type: application/json
X-Request-Id: 380704bb7b4d7c03
Date: Wed, 10 Jun 2015 13:05:07 GMT
Content-Length: 15

{"data":"baz"}
```

```
$ curl -i localhost:8080/resources/2 -u 12345:
HTTP/1.1 404 Not Found
X-Request-Id: 365a858149c6e2d1
Date: Wed, 10 Jun 2015 13:05:28 GMT
Content-Length: 22
Content-Type: text/plain; charset=utf-8

{"error":"not found"}
```

Logging
---
You'll notice that the server supplies a `X-Request-Id` header. This ID is generated for every request and is provided in the log output.

```
$ ./authentication 
2015/06/10 09:03:24 Listening on :8080
2015/06/10 09:03:36 [Req 4d65822107fcfd52] GET /
2015/06/10 09:03:36 [Req 4d65822107fcfd52] Did not provide a token
2015/06/10 09:04:19 [Req 78629a0f5f3f164f] GET /resources/1
2015/06/10 09:04:19 [Req 78629a0f5f3f164f] Provided a token for: bob
2015/06/10 09:04:23 [Req 55104dc76695721d] GET /resources/1
2015/06/10 09:04:23 [Req 55104dc76695721d] Provided a token for: alice
2015/06/10 09:05:07 [Req 380704bb7b4d7c03] GET /resources/3
2015/06/10 09:05:07 [Req 380704bb7b4d7c03] Provided a token for: bob
2015/06/10 09:05:28 [Req 365a858149c6e2d1] GET /resources/2
2015/06/10 09:05:28 [Req 365a858149c6e2d1] Provided a token for: bob
```

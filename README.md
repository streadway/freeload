# Load me up some data URIs

Accepts an aggregate list of URLs and returns a JSON structure of a URI keyed
object to either data URI or original URI.

Useful for doing batch downloads of images and JSON based expansion of the
results into CSS backround URLs.

# Request API

Assuming you're hosting this on localhost:7433 format take your list of origin
URLs from the same host, and strip out the shared prefix and suffix and make a
request with the origin paths encoded in the query string.

Given you wish to request these URLs:

```
http://origin/images/avatar-00123.jpg
http://origin/images/avatar-00124.jpg
http://origin/images/avatar-01234.jpg
```

You would make a request like:

```
http://localhost:7433/json?p=http://origin/images/avatar-0&s=.jpg&i=0123&i=0124&i=1234
```

# Results

The Result JSON and HTTP headers are designed for partial failures on origin
requests.  The format for for successful origin requests contains the
cache-control header max-age of the minimum max-age of the origin max-age
headers.  If any of the origin requests do not contain publically cachable
content, the aggregate response will not be cachable.

```
HTTP 200 OK
Cache-Control: public,max-age=XXX

{
"http://origin/images/avatar-00123.jpg":{"uri":"data:image/jpeg;base64,8QAHAAAAgIDAQEAA..."},
"http://origin/images/avatar-00124.jpg":{"uri":"data:image/jpeg;base64,8QAHAAAAgIDAQEAA..."},
"http://origin/images/avatar-01234.jpg":{"uri":"data:image/jpeg;base64,8QAHAAAAgIDAQEAA..."},
}
```

If any of the origins fail to respond within time, the response will be mixed:

```
HTTP 200 OK
Cache-Control: private,no-store,max-age=0

{
"http://origin/images/avatar-00123.jpg":{"uri":"data:image/jpeg;base64,8QAHAAAAgIDAQEAA..."},
"http://origin/images/avatar-00124.jpg":{"err":"timeout 500ms"},
"http://origin/images/avatar-01234.jpg":{"uri":"data:image/jpeg;base64,8QAHAAAAgIDAQEAA..."},
}
```


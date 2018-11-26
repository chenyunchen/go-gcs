# go-gcs

* [1. Quck Start](#header_1)
* [2. Overview](#header_2)
* [3. API](#header_3)
	* [3.1 Sign Url](#header_3_1)
		* [3.1.1 Request](#header_3_1_1)
		* [3.1.2 Response](#header_3_1_2)
	* [3.2 Resize Image](#header_3_2)
		* [3.2.1 Request](#header_3_1_1)
		* [3.2.2 Response](#header_3_1_2)

<a name="header_1"></a>
## Quick Start (Server Side)

1. Generate new JSON key First: [Create JSON](https://console.developers.google.com/project/<your-project-id>/apiui/credential)

2. Put it into Config file config/local.json .

3. Run the testing.

```bash
go test ./src/...
```

4. Run the server.

```bash
go run src/cmd/gcs-server/main.go
```

<a name="header_2"></a>
## Overview

<!---
```
msc {
  wordwraparcs=true;
  
  auth [linecolor=green, arclinecolor=green, label="Auth"],
  file [linecolor=grey, arclinecolor=grey, label="File"],
  app [linecolor=red, arclinecolor=red, label="App"],
  gcs [linecolor=blue, arclinecolor=blue, label="GCS"];
  
  app => auth [label="GetFileServerURLandJWTToken()"];
  auth => app [label="{\"url\":\"fileserver.jello.com.tw\", 
  \"jwt\":<jwt_token>}"];
  
  app => file [label="POST fileserver.jello.com.tw with info"];
  file => app [label="{\"url\":<upload_url>}"];
  app => gcs [label="PUT <url>"];
  gcs => app [label="GET <url>"]; 
}
```

[mscgenjs](https://mscgen.js.org/)
-->

![](https://i.imgur.com/dalZEaf.png)

Upload URL would be

`https://google_storage/jello_bucket/<url>`

File server would return upload url for app side to put.

`https://google_storage/jello_bucket/<url>?Signature=...`

For more detail, please look at the next section.

<a name="header_3"></a>
## API

<a name="header_3_1"></a>
### Sign Url


<a name="header_3_1_1"></a>
#### Request

```
POST /v1/storage/signurl
```

|Request Header|Description|
|:-:|:-:|
|Content-Type| application/json |
| Authorization | Bearer {jwt token} |

**FYI:** [contentType: MIME](https://en.wikipedia.org/wiki/MIME)

Example:

* image/jpeg
* image/png
* video/mp4 

```json
{
    "fileName": <file_name(string)>,
    "contentType": <content_type(string)>,
    "tag": "single" or "group",
    "payload" : <payload(string)>
}
```

**FYI:** Remember to marshal payload to the string!

When `"tag" = "single"`

```json
"payload": {
   "from": <from_uuid(string)>,
   "to": <to_uuid(string)>
}
```

When `"tag" = "group"`

```json
"payload": {
   "from": <from_uuid(string)>,
   "groupId": <group_uuid(string)>
}
```

Example:

```json
{
	"fileName": "cat.jpg",
	"contentType": "image/jpeg",
	"tag": "single",
	"payload": "{\r\n  \"to\": \"singleawesomeId\" \r\n}"
}
```

<a name="header_3_1_2"></a>
#### Response

```json
{
    "url": <url(string)>,
    "uploadHeaders": {
        "Content-Type": <content_type(string)>,
        "x-goog-content-length-range": <content_length_range(string)>
    },
    "uploadQueries": {
        "Expires": <expires(string)>,
        "GoogleAccessId": <google_access_id(string)>,
        "Signature": <signature(string)>,
    }
}
```

Example:

```json
{
    "baseUrl": "https://storage.googleapis.com/jkopay-test/Group/groupawesomeId/myawesomeId/3f2d6655_abc-123.txt",
    "uploadHeaders": {
        "Content-Type": "text/plain",
        "x-goog-content-length-range": "0,200000000"
    },
    "uploadQueries": {
        "Expires": "1542867051",
        "GoogleAccessId": "jkopay@jkopay-5566.iam.gserviceaccount.com",
        "Signature": "sTBJLfg0Failw9RihUpw2xFgEss4zwmqQQ/ob17e9zJ2xMYUgRIupqiGaMJNGN3cfQxO7nNf/L/LyCoEvwy2ioRflAg4LoNULO3GSCQSokhOgrXbhy44Ie2+ZAKMkWCxsTL9UgWaivWfN62b81HTbQtBYzBWLa8+QAMJd/qvDoqDsgzyYWAkBCGliTQ0x4o6DMcVWVIGeYLrx6FP2v2vvgWSwYfOTbkVcyWoLQjzHdWbr2uURCzCNln9th+8ius8hjCys8nGboCwx7Jy2tNgYC2Ee0RlRiCRlYumGY5mVUzDTCZ7VkV2AHmq6fXb83UBWBB9GOuunn7qXLSxMXqjWQ=="
    }
}
```

**Upload URL:**

`PUT URL = url + uploadQueries(Remember to encode queries)`

|Request Header|Description|
|:-:|:-:|
|Content-Type| {Content-Type} |
| x-goog-content-length-range | {x-goog-content-length-range} |

Example:

```
https://storage.googleapis.com/jkopay-test/Group/groupawesomeId/myawesomeId/3f2d6655_abc-123.txt
?Expires=1542189287
&GoogleAccessId=jkopay%40jkopay-5566.iam.gserviceaccount.com%0A
&Signature=Qb9CDIR6M3OiRFzXbIFP7xWuIFh77B5kgy7Q3gpbYZ3jDh9SxFjZLGHVB%2FLeXcKaCzTs9nOyrNfTWc5A0cX%2BaQztPB7ZKvKE0qf89FTERI6g8hWCCG%2BOEktICXPUgqeBZr1Xm5g6oJRKkXn4BmnSiwrd5TGTUtCyC4qsJWtFwXGHsoy%2F%2Bb41Q6HDRcHHDbXeS8BdyeklMGGHDFpHZVnQMmf7UiIYgZWhY4lKQ2JuU7eTZF4YyLjvsZrHvfPVupgF8O0lF6f2h%2FrgwrT3nR72dgCSMYNwxxcAqQIKw1PH1DLpXrA9GX0vPYkeZHJCIScPOXFyNhNSWGnfwBq8DFvu1g%3D%3D
```

**Download URL:**

`GET URL = url`

Example:

```
https://storage.googleapis.com/jkopay-test/Group/groupawesomeId/myawesomeId/3f2d6655_abc-123.txt
```

<a name="header_3_2"></a>
### Resize Image

<a name="header_3_2_1"></a>
#### Request

```
POST /v1/storage/resize/image
```

```json
{
   "url": <url(string)>,
   "contentType": <contentType(string)>
}
```
Example:

```json
{
	"url": "https://storage.googleapis.com/jkopay-test/test/cat.jpg",
	"contentType": "image/jpeg"
}
```

<a name="header_3_2_2"></a>
#### Response

```json
{
   "100": <100(string)>,
   "150": <150(string)>,
   "300": <300(string)>,
   "640": <640(string)>,
   "1080": <1080(string)>,
	"origin": <origin<string>>
}
```
Example:

```json
{
    "100": "https://storage.googleapis.com/jkopay-test_img_resize/test/cat.jpg_100",
    "150": "https://storage.googleapis.com/jkopay-test_img_resize/test/cat.jpg_150",
    "300": "https://storage.googleapis.com/jkopay-test_img_resize/test/cat.jpg_300",
    "640": "https://storage.googleapis.com/jkopay-test_img_resize/test/cat.jpg_640",
    "1080": "https://storage.googleapis.com/jkopay-test_img_resize/test/cat.jpg_1080",
    "origin": "https://storage.googleapis.com/jkopay-test_img_resize/test/cat.jpg"
}
```
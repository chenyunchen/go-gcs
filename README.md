# Go-GCS

1. Generate new P12 key First: [Create P12](https://console.developers.google.com/project/<your-project-id>/apiui/credential)

2. Convert it into a PEM file and put it on go-gcs/.

```bash
$ openssl pkcs12 -in key.p12 -passin pass:notasecret -out key.pem -nodes
```

3. go test ./src/...

4. go run src/cmd/gcs-server/main.go


## Overview

![](https://i.imgur.com/dalZEaf.png)


`https://google_storage/jello_proj/<url>`

Upload URL would be

`https://google_storage/jello_proj/<url>?key=...`

File server would return upload url for app side to put.

## API


### Request

```
POST /v1/upload
```

```json
{
    "fileName": <file_name(string)>,
    "contentType": <content_type(string)>,
    "tag": "single" or "group",
    "payload" : <payload_string(string)>
}
```

#### When `"tag" = "single"`

```json
"payload": {
   "from": <from_uuid(string)>,
   "to": <to_uuid(string)>
}
```

#### When `"tag" = "group"`

```json
"payload": {
   "from": <from_uuid(string)>,
   "groupId": <group_uuid(string)>
}
```

### Request Headers

|Request header|Description|
|:-:|:-:|
|Content-Type| application/json |
| Authorization | Bearer {jwt token} |


### Response

```json
{
    "url": <url(string)>,
    "uploadQuerys": {
        "Expires": <expires(string)>,
        "GoogleAccessId": <google_access_id(string)>,
        "Signature": <signature(string)>,
    }
}
```

## URL Format

### Response Example:

```json
{
    "url": "https://storage.googleapis.com/jkopay-test/Group/groupawesomeId/myawesomeId/3f2d6655_abc-123.txt",
    "uploadQuerys": {
        "Expires": "1542189287",
        "GoogleAccessId": "jkopay@jkopay-5566.iam.gserviceaccount.com",
        "Signature": "s5NocDqomEEtPcjdcbq5jSAqcyFrNaDnsILSRabdTseDrn/gQE9Qm8vbPbOHzWs2oe6bZiJi0vXW8Sh/Wf5KYoTXUloeKAsmao9StcLz2ShYJ6ZvSaLz2bccwu/j1KV/AKDirihnYlBgDue/HS59mKE6swALYgzlojxATCXpIKgAkcRC5VSAIrH+o2DlQ6gCn+xTDZBHiqsB8XM3sjtvy23elKjCfCpK7duuQU/6t24cEhN9gvaK69kBQmEi687+XX618WoH8d85KgebcyuYNsFNSF6BgJZj2qwNkOxxVKBFgxmk1MfP+/qaY7TeiqhhxiTCQGS7NJ/Fr92HcPSblQ=="
    }
}
```

### UploadURL:
```
https://storage.googleapis.com/jkopay-test/Group/groupawesomeId/myawesomeId/3f2d6655_abc-123.txt
?Expires=1542189287
&GoogleAccessId=jkopay@jkopay-5566.iam.gserviceaccount.com
&Signature=s5NocDqomEEtPcjdcbq5jSAqcyFrNaDnsILSRabdTseDrn/gQE9Qm8vbPbOHzWs2oe6bZiJi0vXW8Sh/Wf5KYoTXUloeKAsmao9StcLz2ShYJ6ZvSaLz2bccwu/j1KV/AKDirihnYlBgDue/HS59mKE6swALYgzlojxATCXpIKgAkcRC5VSAIrH+o2DlQ6gCn+xTDZBHiqsB8XM3sjtvy23elKjCfCpK7duuQU/6t24cEhN9gvaK69kBQmEi687+XX618WoH8d85KgebcyuYNsFNSF6BgJZj2qwNkOxxVKBFgxmk1MfP+/qaY7TeiqhhxiTCQGS7NJ/Fr92HcPSblQ==
```
DownloadURL:
```
https://storage.googleapis.com/jkopay-test/Group/groupawesomeId/myawesomeId/3f2d6655_abc-123.txt
```

### Single

`/single/<A>/<B>/<hash_id>_filename`

* A send a file with filename, `filename`.
* Server side would generate a `hash_id` to prepend to the filename to avoid collision of file names.

It would be easy for A to get all the files sent from A.

* `/single/<A>/<B>`
* `/single/<A>/<C>`
* `/single/<A>/<F>`


### Group

`/group/<group_id>/<A>/<hash_id>_filename`

* A upload a file, `filename`, to group with id `group_id`
* Server side would generate a `hash_id` to prepend to the filename to avoid collision of file names.

It would be easy for a group to get all related files

`/group/<group_id>`


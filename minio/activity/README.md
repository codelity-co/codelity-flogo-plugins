<!--
title: MinIO
weight: 4705
-->
# MinIO
This activity allows you to manage MinIO object.

## Installation

### Flogo CLI
```bash
flogo install github.com/codelity-co/codelity-flogo-plugins/minio/activity
```

## Configuration

### Settings:
  | Name                | Type   | Description
  | :---                | :---   | :---
  | endpoint            | string | The MinIO endpoint - ***REQUIRED***
  | accessKey           | string | Access Key
  | secretKey           | string | Secret Key
  | enableSSL           | bool   | Enable SSL connection

### Handler Settings
  | Name                | Type   | Description
  | :---                | :---   | :---
  | method              | string | Minio SDK method - ***REQUIRED***
  | parmas              | object | Method parameters - ***REQUIRED***

#### Method

MinIO provides many methods to maintain buckets. Here are the list of all avialabe methods of this plugins:

* PutObject

#### Method opitons

Method options are documented in [MinIO Golang SDK](https://docs.min.io/docs/golang-client-api-reference#).  Please reference specific method options based on your setting.

## Example

```json
{
  "id": "nats-trigger",
  "name": "NATS Trigger",
  "ref": "github.com/codelity-co/codelity-flogo-plugins/minio/activity",
  "settings": {
    "endpoint" : "minio:9000",
    "accessKey": "minio",
    "secretKey": "minio123",
    "enableSsl": false,


  },
  "input": {
    "method": "PutObject",
    "params": {
      "bucketName": "flogo",
      "objectName": "inbox/testing.json",
      "stringData": "{\"abc\": \"123\"}"
    }
  },
  "output": {
    "status": "SUCCESS",
    "result": {
      "bytes": 14,
    }
  }
}
```
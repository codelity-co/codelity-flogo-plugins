<!--
title: MinIO
weight: 4705
-->
# MinIO

**This plugin is stll IN PROGRESS**

This activity allows you to manage MinIO object.

## Installation

### Flogo CLI
```bash
flogo install github.com/codelity-co/codelity-flogo-plugins/activity/minio
```

## Configuration

### Settings:
  | Name                | Type   | Description
  | :---                | :---   | :---
  | endpoint            | string | The MinIO endpoint - ***REQUIRED***
  | accessKey           | string | Access Key
  | secretKey           | string | Secret Key
  | enableSSL           | bool   | Enable SSL connection
  | bucketName          | string | MinIO bucket name
  | region              | string | MinIO Region/Zone name
  | methodName          | string | MinIO SDK method name
  | methodOptions       | object | MinIO method options

### Handler Settings
  | Name                | Type   | Description
  | :---                | :---   | :---
  | objectName          | string | Minio object name - ***REQUIRED***
  | data                | any    | data - ***REQUIRED***

## Output
  | Name                | Type   | Description
  | :---                | :---   | :---
  | status              | string | status text - ***REQUIRED***
  | result              | object | result object - ***REQUIRED***

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
  "ref": "github.com/codelity-co/codelity-flogo-plugins/activity/minio",
  "settings": {
    "endpoint" : "minio:9000",
    "accessKey": "minioadmin",
    "secretKey": "minioadmin",
    "enableSsl": false,
    "bucketName": "flogo",
    "region": "zone-0",
    "methodName": "PutObject"
  },
  "input": {
    "objectName": "inbox/test.json",
    "data": "{\"abc\": \"123\"}"
  },
  "output": {
    "status": "SUCCESS",
    "result": {
      "bytes": 14,
    }
  }
}
```
<!--
title: NATS
weight: 4705
-->
# NATS
This trigger allows you to listen to messages on NATS.

## Installation

### Flogo CLI
```bash
flogo install github.com/project-flogo/edge-contrib/nats/trigger
```

## Configuration

### Settings:
  | Name                | Type   | Description
  | :---                | :---   | :---
  | clusterUrls         | string | The NATS cluster URL start with nats://, comma seperated - ***REQUIRED***
  | connName            | string | NATS connection name
  | auth                | object | Auth configuration
  | reconnect           | object | Reconnect configuration
  | sslConfig           | object | SSL configuration

 #### *auth* Object:
  | Property            | Type   | Description
  |:---                 | :---   | :---     
  | username            | string | The user's name
  | password            | string | The user's password
  | token               | string | NATS token
  | nkeySeedfile        | string | NKey seed file path
  | credFile            | string | Credential file path

 #### *reconnect* Object:
  | Property            | Type   | Description
  |:---                 | :---   | :---     
  | autoReconnect       | bool   | Enable Auto-Reconnect
  | maxReconnects       | int    | Max reconnect attemtps
  | dontRandomize       | bool   | Don't randomize reconnect
  | reconnectWait       | int    | Reconnect wait seconds
  | reconnectBufferSize | int    | Reconnect buffer size in bytes

 #### *sslConfig* Object:
  | Property            | Type   | Description
  |:---                 | :---   | :---     
  | skipVerify          | bool   | Skip SSL validation, defaults to true
  | caFile              | string | The path to PEM encoded root certificates file
  | certFile            | string | The path to PEM encoded client certificate
  | keyFile             | string | The path to PEM encoded client key

### Handler Settings
  | Name                | Type   | Description
  | :---                | :---   | :---
  | subject             | string | The subject to listen on - ***REQUIRED***
  | queue               | string | Subscriber queue group


#### Subject
NATS provides two wildcards that can take the place of one or more elements in a dot-separated subject. Subscribers can use these wildcards to listen to multiple subjects with a single subscription but Publishers will always use a fully specified subject, without the wildcard.

The first wildcard is * which will match a single token. For example, if an application wanted to listen for eastern time zones, they could subscribe to time.*.east, which would match time.us.east and time.eu.east.

The second wildcard is > which will match one or more tokens, and can only appear at the end of the subject. For example, time.us.> will match time.us.east and time.us.east.atlanta, while time.us.* would only match time.us.east since it can't match more than one token.

### Output:

| Name          | Type   | Description
| :---          | :---   | :---
| message       | string | The message received
| subject       | string | The NATS subject

### Reply:

| Name  | Type   | Description
| :---  | :---   | :---
| data  | object | The data recieved

## Example

```json
{
  "id": "nats-trigger",
  "name": "NATS Trigger",
  "ref": "github.com/project-flogo/edge-contrib/nats/trigger",
  "settings": {
      "clusterUrls" : "nats://localhost:4222",
     	"connName":"NATS connection"
  },
  "handlers": {
    "settings": {
      "async": true,
    	"subject": "flogo",
    }
  }
}
```
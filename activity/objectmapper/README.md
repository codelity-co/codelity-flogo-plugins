<!--
title: ObjectMapper
weight: 4705
-->
# ObjectMapper

**This plugin is in ALPHA stage**

This trigger allows you to transform input based on Flogo mapping.

## Installation

### Flogo CLI
```bash
flogo install github.com/codelity-co/codelity-flogo-plugins/activity/objectmapper
```

## Configuration

### Settings:
No settings is required

### Input
  | Name                | Type   | Description
  | :---                | :---   | :---
  | in                  | object | input variable, please see Mappings and Example - ***REQUIRED***

### Output:
  | Name          | Type   | Description
  | :---          | :---   | :---
  | out           | any    | output object - ***REQUIRED***

#### Mappings
Mapping syntax is based on [Flogo documented syntax](https://tibcosoftware.github.io/flogo/development/flows/mapping/)

## Example

```json
{
  "id": "codelity-objectmapper-activity",
  "name": "Codelity Object Mapper Activity",
  "ref": "github.com/codelity-co/codelity-flogo-plugins/activity/objectmapper",
  "settings": {},
  "input": {
    "in": {
      "mapping": {
        "subject": "=json.path(\"$.subject\", coerce.toObject($flow.payload))",
        "message": "=json.path(\"$.message\", coerce.toObject($flow.payload))",
        "receivedTimestamp": "=json.path(\"$.receivedTimestamp\", coerce.toObject($flow.payload))"
      }
    }
  }
}
```
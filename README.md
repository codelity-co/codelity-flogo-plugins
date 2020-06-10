# Codelity Flogo Plugins

Codelity is now contributing some plugins for system-level integration

## Flogo

Project Flogo is an open source framework to simplify building efficient & modern serverless functions and edge microservices and this repository is the core library used to create and extend those Flogo Applications.

Please go [here](https://github.com/project-flogo/core) for details.

## List of Flogo plugins

### Triggers

* [MQTT Trigger](https://github.com/codelity-co/flogo-mqtt-trigger) (forked from github.com/project-flogo/edge-contrib)
* [NATS Trigger](https://github.com/codelity-co/flogo-nats-trigger) (with STAN support)
* [Zeebe Task Trigger](https://github.com/codelity-co/flogo-zeebetask-trigger)


### Activities

* [NATS Activity](https://github.com/codelity-co/flogo-nats-activity) (with STAN support)
* [ObjectMapper Activity](https://github.com/codelity-co/flogo-objectmapper-activity) (similar with Mapper but it is for activities within the same flow)
* [CockroachDB Activity](https://github.com/codelity-co/flogo-cockroachdb-activity)
* [MinIO Activity](https://github.com/codelity-co/flogo-minio-activity)
* [Zeebe Workflow Activity](https://github.com/codelity-co/flogo-zeebeworkflow-activity)

### Functions
* [datetimex Functions](https://github.com/codelity-co/flogo-datetimex-function) (Extented functions for datetime)
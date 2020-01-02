# Integration test

This folder acts as a small photon go example and a full integration test of the runtime binary fetching.

This integration test folder needs to be run with docker:

```shell script
docker build . -f docker/integration.dockerfile -t integration && docker run integration
```

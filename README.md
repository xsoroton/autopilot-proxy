# Autopilot API caching middleware server

Implement Golang caching middleware server for Autopilot API /contact method (you can create a free trial account and read about Autopilot API on https://autopilot.docs.apiary.io/#introduction/getting-help):

- Create GET / POST / PUT /contact endpoint
- Retrieve a requested contact from redis, if it is not present retrieve it from Autopilot API and store in redis
- Create / Update a contact and invalidate redis cache after POST / PUT requests
- Cover the necessary methods with unit tests
- Write README file with instructions how to run and test it

We understand that this is a test task and it should not take more than a couple of hours to do it, but we would like you to take an approach as it was a usual production task to show your developer capabilities.
# Autopilot API caching middleware server

Implement Golang caching middleware server for Autopilot API /contact method (you can create a free trial account and read about Autopilot API on https://autopilot.docs.apiary.io/#introduction/getting-help):

- Create GET / POST / PUT /contact endpoint
    - _Note: **PUT in not supported by API Docs**, but POST support idempotency, and act like PUT_
        ```
        # PUT response example at API
        {
            "code": "MethodNotAllowedError",
            "message": "PUT is not allowed"
        }
        ```
- Retrieve a requested contact from redis, if it is not present retrieve it from Autopilot API and store in redis
- Create / Update a contact and invalidate redis cache after POST / PUT requests
- Cover the necessary methods with unit tests
- Write README file with instructions how to run and test it

---
## Run Reverse Proxy

Before running service please make that Environment Variables are set correctly and point correct services, by default it set to:
```bash
export AUTOPILOT_API="https://private-1f378a-autopilot.apiary-proxy.com"
export REDIS_HOST="localhost:6379"
```

Please run server.go file and listen Reverse Proxy on port 8080.
```bash
$ go get -v
$ go run ./
```

#### Post Contact
request example
```bash
curl --location --request POST 'http://0.0.0.0:8080/v1/contact' \
--header 'autopilotapikey: xxxxxxxxxxxxxxx' \
--header 'Content-Type: application/json' \
--data-raw '{
  "contact": {
    "FirstName": "Bob",
    "LastName": "Goodman",
    "Email": "test@gtest.com",
    "custom": {
      "string--Test--Field": "This is a test"
    }
  }
}'
```

expected response
```json
{
    "contact_id": "person_BEDEF3B9-8B84-4F5F-AA58-22D025DDA683"
}
```
---
#### Get Contact
request example
```bash
curl --location --request GET 'http://0.0.0.0:8080/v1/contact/person_BEDEF3B9-8B84-4F5F-AA58-22D025DDA683' \
--header 'autopilotapikey: xxxxxxxxxxxxxxxxxxxxxxx'
```

expected response
```json
{
    "Email": "test@gtest.com",
    "created_at": "2021-10-31T04:21:50.000Z",
    "updated_at": "2021-10-31T04:21:50.000Z",
    "api_originated": true,
    "custom_fields": [
        {
            "kind": "Test Field",
            "value": "This is a test",
            "fieldType": "string",
            "deleted": false
        }
    ],
    "Name": "Bob Goodman",
    "LastName": "Goodman",
    "FirstName": "Bob",
    "contact_id": "person_BEDEF3B9-8B84-4F5F-AA58-22D025DDA683"
}
```
---

## Test Cover

to run tests
```bash
 go test -v -cover ./...
```

Expected output
```
=== RUN   TestGetContact

  Test GET Contact ✔✔✔✔✔
    Test cache is set ✔✔


7 total assertions

--- PASS: TestGetContact (0.10s)
=== RUN   TestPostContact

  Test POST Contact ✔✔✔✔✔✔
    Test POST invalidate cache ✔


14 total assertions

--- PASS: TestPostContact (0.11s)
PASS
coverage: 70.4% of statements
ok      github.com/xsoroton/autopilot-proxy/handlers    0.218s  coverage: 70.4% of statements
=== RUN   TestStoreMem

  Test Store Mem ✔✔✔✔✔✔


6 total assertions

--- PASS: TestStoreMem (0.00s)
PASS
coverage: 65.2% of statements
ok      github.com/xsoroton/autopilot-proxy/store       (cached)        coverage: 65.2% of statements

```

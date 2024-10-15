# 1: Prime Time
https://protohackers.com/problem/1

- JSON based protocol
- Each Request is single line contaning a JSON Object
- `{"method":"isPrime","number":123}`
- Request Object: must have `method` field and `number` field
- Response Object
- {"method":"isPrime","prime":false}
-  A response is malformed if it is not a well-formed JSON object, if any required field is missing, if the method name is not "isPrime", or if the prime value is not a boolean.\
- Note that non-integers can not be prime.
- Whenever you receive a malformed request, send back a single malformed response, and disconnect the client. 





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

Takeaways:
- Always consider all the constraints and input type
    - Constraints relating to the number range and considering BigInts (integer value greater int64 or int32)
    - Consider the input type: number -> floats + int
- Go through the RFC to understand the Constraints
- Logging is always -> Important 
- BigFloats/Int in go cannot be unmarshalled (suprising!)
- It might help to serialise JSON data first to string and then converting if need be
- Custom Data Type for JSON deserilisation. `BigInteger`
- `fmt` and `log`: both are different and used for differnt purposes.
- Remember when to close the connection and how to.
- Remember the lifecycle of the connection. for ex. considering taking mulitple request from the connection befre closing the connection
- Make it more readable, more use of chanels and go routine.

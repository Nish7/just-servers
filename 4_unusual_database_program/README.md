# 4: Unusual Database Program 
[Protohackers - Unusual Database Program](https://protohackers.com/problem/4)

## Design

![General-Architecture](architecture.png | width=500) 

### Useful Links
- https://www.ietf.org/rfc/rfc768.txt - UDP RFC768

## Takeaways
- UDP Basics
  - UDP treats each message or piece of data as a separate and complete unit, called a "datagram" or "packet." Each datagram is sent independently of the others.
  - Imagine sending individual postcards through the mail. Each postcard is its own message, and you don't know if they'll arrive in order, at all, or duplicated.
  - So there is no concept of connection in UDP, each packet is sent indepdently, with all the information it needs to reach the destination.
  - It is completly stateless.
- `WriteTo` vs `Write`, as UDP does not maintain a connection, each connection packet is discrete, `Write` doesnt seem to work as expected, `WriteTo` is expeted to be used to send to a given addr, which originally sent the request
    - Even though UDP is connectionless, the Connect method is a convenience provided by the OS. When you call Connect on a UDP socket:
    - The OS associates the socket with a specific remote address.
    - The kernel restricts the socket to communicate only with that address. 
- One point of debugging: Always consider the new lines. As `nc` is sending the `\n` newlines as well, which needs to trim new lines. Basically alway make to sure sanitize request and responses
- If you are looking into some documentation for structs make sure to go look if they can be created natively, or without manual `&` creation. For example. creating `UDPAddr` is not neccesary because you have `ResolveUDPAddr` to create that for you. 

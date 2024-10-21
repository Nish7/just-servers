# 2: Means to an End
https://protohackers.com/problem/2

- Each client tracks the price of a different asset. Clients send messages to the server that either insert or query the prices.
- Each connection from a client is a separate session. Each session's data represents a different asset, so each session can only query the data supplied by itself.

Message Format:
- 9 bytes long - 72 bits
- 1st byte is indentifier  - in ASCII endcoding
- 2nd and 3rd segmetns of each 4bytes is in two complement 32 bits integer conversion

Takeaways:
- Using `buf:=make([]byte, 2048)` vs `bufio.NewReader`. While in the PROTO-1 i used bufio.NewReader, the decison on to use simple buf was made due the nature of the data being received. While in PROTO-1 it was json data
with new line and general delimeter capabilities. Thus always be mindful of type or struct to use.
    - Low-level Reading: This method provides more control but requires you to handle all the details—like splitting the incoming data into meaningful chunks (e.g., splitting it into lines, messages, etc.).

- Sending binary data is much more complicated rather than processing it. Binary data needs to be send in a 
byte form rather than ASCII binary strings. Which would happen if you add "100101" in the stdin on the nc 
console.
    - Sending the binary strings would simply mean we would be sending the ASCII codes for those digits
    - Ex. "01000001" -> would be sending 0x30, 0x31, 0x30, 0x30, 0x30, 0x30, 0x30, 0x31
    - Rather than -> 0x41
    - So How do we send byte form "binary data": `printf '\x41\x00\x01\xE2\x40\xFF\xF3\xF6\x1C' | nc <server-ip> <port>`

- The two's complement system is used in computing because it simplifies the representation and arithmetic operations on signed integers 
    - postive numbers are converted normally
    - For negative number are converted and inverted. then add 1
    - addition circuitlry can be reused for both subtraction

- Convertion to twos complement is really intresting 
    -  firstly, it assumes the data is in big-endian Format. Which mean MSB comes first
    -  In Big-Endian, the most significant byte (0x12) is stored at the lowest memory address (first = 0)
    - We have currently used binary package to handle the combining the 4 bytes into a 32 bit integer, internallly, it is handled by shifting the bits by its respective position in 32bit size and then using an 
    OR operator on those.
    - type casting `int32` on the 4 bytes would handle the conversion for twos complement


- Considering Data Race conditions could happen if using a shared map is being accessed by bunch of go routines
    - we make it safe by the use of mutex and sync.map
    - However, i opted for more secure by design approach
    - Rather making a shared map and accessing those, will create a new map per connection and keep track of that.
    - It would be a connection specific state managed, whenever we drop the go routine we clean up
    database as well. Which is the intended feature as well

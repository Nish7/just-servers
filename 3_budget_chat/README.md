# 2: Budget Chat
https://protohackers.com/problem/3

## Design

- Design consideration is big things:
- There 3 major design consideration that were taken
1. Having a simple handle connection design, where each request and connection is handled wihin the handleConnection handler
    - pros: simple
    - cons: have to manage client map and use of mutex
2. Having a event based strucutre. There are 3 types of event. Join, Leave and Message. There is a single channgel to handle the each events and an event dispatcher will route to a specific handler to handle each type of event
    - pros: can be expanded with more events; 
    - cons: little bit overkill for a simple application; use of  mutex
3. Even more aggresive, event based structure where each client has its own channel, each client its own consumption and writing. There is a common broadcaster and broadcast channel which handles routing and deciding whihch channel/user it should go to.
    - pros: no more handling of mutex; secure by design 
    - cons: way much more complicated to manage. 


![General-Architecture-Means](general-architecture.png)

## Local Test Commands

### Useful Links


## Takeaways:

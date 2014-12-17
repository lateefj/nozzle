# nozzle

Control the flow of http requests.

## Problem

The basic problem is that servers are very bad at responding to varying response times. When a server performance begins to degrade it is generally not aware of it. It doesn't descriminate the requests it is getting it just keep on trying to serve them no matter how long they are expected to take.

The idea for _Nozzle_ is that a proxy with a small amount of intellegence and memory could potentially make better decisions about handling client requests. There are some basic assumptions:

* The servers that are downstream are not heterogenious
* Failing a client request quickly is better than sending it to an overloaded or slow server
* Response times can be bucketed (grouped) by hashing the request path and parameter names

## Hypothesis

Client request can be bucketed into queues based on average reponse time of the hashed request path and parameter names. Based on the average reponse time the request queue size would be based on a simple formula:

queue size = ( (1/reponse time) * response rate) / max reponse time

Once the number of client requests filled the queue the request would be dropped at the proxy until there was more room on the queue. If response time decreased then the queue size will increase. If the response time increases the queue size would decrease. The failures would expose backpressure dynamically.


## Implementation Ideas

### Vulcanproxy

Initial implementation could use https://www.vulcanproxy.com/library.html#limiter to support rate and connection limits. However the basic calculation of throughput would still need to be implemented in the vulcan middleware concept. 

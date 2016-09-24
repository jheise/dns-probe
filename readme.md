DNS-PROBE
---

usage: dns-probe ifacename

dns-probe sniffs for traffic on udp port 53 then makes the following fields available via a json blob served on a zeromq stream on port 7777

SrcIP - ip the packet came from
DstIP - ip the packet is heading to
Request - boolean value defining if this is a query or an answer
Timestamp - time the packet was captured
Query - the domain that was queried
Answers - the response to the query

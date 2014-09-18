# broadcast-graph

Broadcast-graph is a simple graph backend example for
[broadcast](http://github.com/nyxtom/broadcast). This backend implements
various commands useful for updating directed or undirected weighted edges
as well as weighted vertices. Given that this is an example
implementation, it uses an ephemeral in-memory store to keep all data
(broadcast-stats does the same thing in this particular case). 

## Commands

```
127.0.0.1:7331> CMDS
&E
 Returns the intersection of all edges between the set of vertices with the sum of the weights
 usage: &e vertex [vertex ...]

*E
 Returns a list of all edges from the specified vertices
 usage: *e vertex [vertex ...]

+
 Increments a given vertex's own weight
 usage: + weight vertex [weight vertex ...]

+>
 Increments the directed edge weight
 usage: +> weight from to [from to ...]

-
 Decrements a given vertex's own weight
 usage: - weight vertex [weight vertex ...]

->
 Decrements the directed edge weight
 usage: -> weight from to [from to ...]

<+>
 Increments the symmetric edge weight
 usage: <+> weight from to [from to ...]

<->
 Decrements the symmetric edge weight
 usage: <-> weight from to [from to ...]

<=>
 Sets the symmetric edge weight
 usage: <=> weight from to [from to ...]

=
 Sets a given vertex's own weight
 usage: = weight vertex [weight vertex ...]

=>
 Sets the directed edge weight
 usage: => weight from to [from to ...]

CMDS
 List of available commands supported by the server

ECHO
 Echos back a message sent
 usage: ECHO "hello world"

INFO
 Current server status and information

PING
 Pings the server for a response

127.0.0.1:7331>
```

## Build and Install

Installation can be done via make or by running the command below.

```
go get github.com/nyxtom/broadcast-graph
```

## License
The MIT License (MIT)

Copyright (c) 2014 Thomas Holloway

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

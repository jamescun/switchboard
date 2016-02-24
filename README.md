<img src="http://i.imgur.com/pmUVbZk.png" width="400" height="84" />

Switchboard uses TLS Server Name Indication (SNI) to provide service discovery and load balancing to almost any protocol or service.

Switchboard intercepts the public TLS CLIENT_HELLO packet, extracts the service name, and provides upstream proxying thereafter. It does not need to decrypt any traffic or need knowledge of the certificates being used on either side, additionally it does not interfere with client or server certificate authentication.


Why
---

Service discovery proxying for an HTTP service is easy, HTTP/1.1+ contains a header called `Host` which can be used to direct traffic to an upstream server without using many IP addresses or ports per service. The trouble is that this is an application specific key with its own implementation, to support other protocols you must also implement their application specific key, and some protocols don't even have this. To work around this, solutions have been devised using firewall rules, userspace networks and kernel modules for broader protocol support.

Introduced to the IETF in 2004, implemented in OpenSSL in 2004; [RFC 3546](https://tools.ietf.org/html/rfc3546) introduced a new field to the initial handshake of a TLS connection, the hostname the client is requesting. This meant that any protocol connecting over TLS could utilise functionality similar to that of HTTP.

If you like, TLS SNI is the Host header for the internet.


### Pros

  - **TLS ALL THE THINGS**: modern hardware has [minimised the cost of TLS](https://istlsfastyet.com/) and modern implementations make it trivial to deploy; leaving no reason not to secure your applications traffic, ever more vitally so on public cloud networks.
  - **Management and Performance**: userspace networking, the most compatible alternative, can add significant latency and restrict usable bandwidth by needing to maintain their own L2/L3 stack. Switchboard acts like any other socket server and can be optimised as such, benchmarks show the difference can often be [an order of magnitude](http://blog.backslasher.net/ssh-openvpn-tunneling.html).
  - **Service Discovery Support**: Switchboard does not implement any service discovery, it exposes a simple interface for other systems (such as Etcd, Consul, Marathon etc) to integrate with.

### Cons

  - **TLS ALL THE THINGS**: if you are not currently using TLS, deployment of Switchboard includes the additional overhead of a fresh TLS deployment. Some servers do not natively support TLS but clients often do meaning, with the addition of tools like stunnel, they can still be used with Switchboard (e.g. redis).
  - **No STARTTLS Support**: protocols implementing STARTTLS (or similar) are not supported by Switchboard directly as the first packet on the wire is not a TLS CLIENT_HELLO. This includes POP, IMAP and SMTP.


Configuration
-------------

### Backends

#### Static Backend

The static backend does not perform any service discovery, only proxying to a set of upstreams defined on startup. This is useful if routing is performed on another upstream. It can still be used in conjunction with load balancing functionality.

    switchboard server static --upstream=192.168.0.11:8080 --upstream=192.168.0.12:8080

    export SWITCHBOARD_STATIC_UPSTREAM=192.168.0.11:8080,192.168.0.12:8080
    switchboard server static


#### Marathon Backend

The Marathon framework for Mesos exports information about cluster tasks through a REST API. Switchboard can query this API, by default using the label `switchboard`, for which tasks provide a service for a hostname.

	switchboard server marathon --marathon-host=192.168.0.11:8080 --marathon-host=192.168.0.12:8080

	export SWITCHBOARD_MARATHON_HOST=192.168.0.11:8080,192.168.0.12:8080
	export SWITCHBOARD_MARATHON_LABEL=service_name # non-default task label
	switchboard server marathon


### Cache

Switchboard by default will cache the results of service discovery for 1 minute to lessen the load on service discovery and improve routing performance. This value can be changed or disabled like below:

    switchboard server --cache=10m

	export SWITCHBOARD_CACHE=0 # disable caching
	switchboard server


### Balancing

Switchboard will perform naive load balancing across the upstreams provided by service discovery. The load balancing schemes currently implemented are:

  - **random**: new connections will be assigned to a random upstream.
  - **consistent**: use O(1) consistent hashing of CRC32(hostname + client ip) to maintain a single upstream for a client for the lifetime of a particular service discovery state.
  - **rendezvous**: use O(log n) rendezvous hashing, also known as Highest Random Weight (HRW) hashing, to maintain a single upstream for a client for the lifetime of a particular service discovery state.

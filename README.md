<img src="http://i.imgur.com/pmUVbZk.png" width="400" height="84" />

Switchboard is a reverse proxy server that uses TLS Server Name Indication to select an upstream from a service discovery mechanism. It also implements basic load balancing and retrying.

Switchboard intercepts the public TLS CLIENT_HELLO packet, extracts the value of the SNI extension, and provides reverse proxying to an upstream thereafter. It does not need to decrypt any traffic or need knowledge of the certificates being used on either side, additionally it does not interfere with client or server certificate authentication.


Why
---

Introduced to the IETF in 2004, implemented in OpenSSL in 2004; [RFC 3546](https://tools.ietf.org/html/rfc3546) added a new extension to the initial handshake of a TLS connection, the hostname the client is requesting. This meant that any protocol connecting over TLS could utilise functionality similar to that of the HTTP `Host` header. If you like, TLS SNI is the Host header for the internet.


### Pros

  - **TLS ALL THE THINGS**: the very core of Switchboard is TLS Server Name Indication, thus it requires all services to be available over TLS.
  - **Performance**: while TLS does have an overhead versus raw TCP, the cost of this with modern hardware is [minimal](https://istlsfastyet.com/) particularly when compared to userspace networking alternatives (i.e. [OpenVPN vs SSH Tunneling](http://blog.backslasher.net/ssh-openvpn-tunneling.html)).
  - **Management**: Switchboard is simply a TCP server that understands part of the TLS specification; it can be managed and implemented similarly to how you might manage a web server.
  - **Pluggable Service Discovery**: rather than re-invent the wheel, Switchboard relies on existing service discovery tools such as Etcd, Consul and Marathon.


### Cons
  - **TLS ALL THE THINGS**: the very core of Switchboard is TLS Server Name Indication, thus it requires all services to be available over TLS.
  - **No STARTTLS Support**: protocols implementing STARTTLS (or similar) are not supported by Switchboard directly as the first packet on the wire is not a TLS CLIENT_HELLO. This includes POP, IMAP and SMTP.


Configuration
-------------

### Service Discovery

Switchboard does not implement service discovery itself, rather it relies on existing service discovery tools for upstream configuration.

#### Static

The static backend does not perform any service discovery, only proxying to a set of upstreams defined on startup. This is useful if routing is performed on another upstream. It can still be used in conjunction with load balancing functionality.

	switchboard server static --upstream=192.168.0.11:8080 --upstream=192.168.0.12:8080

	export SWITCHBOARD_STATIC_UPSTREAM=192.168.0.11:8080,192.168.0.12:8080
	switchboard server static


#### Marathon

The Marathon framework for Mesos exports information about cluster tasks through a REST API. Switchboard can query this API, by default using the label `switchboard`, for which tasks provide a service for a hostname.

	switchboard server marathon --marathon-host=192.168.0.11:8080 --marathon-host=192.168.0.12:8080

	export SWITCHBOARD_MARATHON_HOST=192.168.0.11:8080,192.168.0.12:8080
	export SWITCHBOARD_MARATHON_LABEL=non_standard_label
	switchboard server marathon


### Cache

Switchboard by default will cache the results of service discovery for 1 minute to lessen the load on service discovery and improve routing performance. Caching can be disabled by setting the ttl to 0.

    switchboard server --cache=10m

	export SWITCHBOARD_CACHE=10m
	switchboard server


### Balancing

Switchboard will perform naive load balancing across the upstreams provided by service discovery. The load balancing schemes currently implemented are:

  - **random**: new connections will be assigned to a random upstream.
  - **consistent**: use O(1) consistent hashing of CRC32(hostname + client ip) to maintain a single upstream for a client for the lifetime of a particular service discovery state.
  - **rendezvous**: use O(log n) rendezvous hashing, also known as Highest Random Weight (HRW) hashing, to maintain a single upstream for a client for the lifetime of a particular service discovery state.

	switchboard server --balance=consistent

	export SWITCHBOARD_BALANCE=rendezvous
	switchboard server


### Retrying

Switchboard can exponentially retry connections that fail during the establishment phase, this can help with clients that assume the network is reliable. Switchboard will never attempt to retry after failure in an in-flight connection.

	switchboard server --retry --retry-max=3 --retry-timeout=30s

	export SWITCHBOARD_RETRY=1
	export SWITCHBOARD_RETRY_MAX=3
	export SWITCHBOARD_RETRY_TIMEOUT=30s
	switchboard server


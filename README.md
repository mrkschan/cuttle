cuttle
======

[![Build Status](https://travis-ci.org/mrkschan/cuttle.svg?branch=master)](https://travis-ci.org/mrkschan/cuttle)

Cuttle is a HTTP forward proxy designed for HTTP clients who need to respect rate limit. Its primary use case is to serve as the centralized outbound rate limit controller for API clients.


Quickstart
----------

```
GOPATH=`pwd` go get github.com/mrkschan/cuttle
bin/cuttle -f src/github.com/mrkschan/cuttle/cuttle.yml
```

Read `cuttle.yml` for available configurations.


What is Cuttle designed for?
----------------------------

* Centralized rate limit control for API client running in multiple processes.
* Centralized rate limit control for multi-process web crawler.


When NOT to use Cuttle?
-----------------------

* Fully fledged HTTP forward proxy with access control, caching, etc., consider Squid instead. [1]
* API gateway of your service with authentication, inbound rate limiting, etc., consider Kong instead. [2]
* HTTP reverse proxy in front of the API service, consider Nginx instead. [3]

[1]: http://www.squid-cache.org/
[2]: https://getkong.org/
[3]: http://nginx.org/


Why Cuttle is born?
-------------------

There are quite a number of ways to respect rate limit of an API service as below:

* Serialize all API calls into a single process, use sleep statement to make pause between consecutive calls. [4]
* Host a RPC server / use a task queue to make API calls. The RPC server / queue manager has to rate limit the calls. [5]
* Centralize all API calls with a HTTP proxy where the proxy performs rate limiting.

HTTP proxy stands out as the most generic solution on the list. And, there are quite a few existing options in using HTTP proxy. For example:

* Using Nginx reverse proxy to wrap the API, use its limit module to perform simple rate limit or Lua plugin for more sophisticated control. [6]
* Using Squid forward proxy to perform simple rate limit by its delay pools. [7]

At first glance, the Nginx reverse proxy option looks superior since we can have sophisticated rate limit control deployed. Though, using such approach would need to change the URL to the API service in the API client so that traffic goes through Nginx. Or, we have to modify DNS configuration to route the traffic.

Thus, reverse proxy may not be a good solution for everyone. Cuttle, in contrast, positions as a forward proxy with a sole focus on delivering generic but sophisticated rate limit control capability.

[4]: https://github.com/benbjohnson/slowweb
[5]: http://product.reverb.com/2015/03/07/shopify-rate-limits-sidekiq-and-you/
[6]: http://codetunes.com/2011/outbound-api-rate-limits-the-nginx-way/
[7]: http://wiki.squid-cache.org/Features/DelayPools


Behind the scene
----------------

Each HTTP request going through Cuttle is handled by a goroutine. Cuttle would inspect the HTTP headers and pick a configured rate limit controller. The request is then blocked until the rate limit controller turns on the green light. A dedicated goroutine would be created to forward the pending request afterwards.

Each rate limit controller in Cuttle is a goroutine that has two go channels - pending channel and ready channel. The controller is put to sleep until it receives a signal from its pending channel where the signal is sent from the pending request. Then, the controller might be blocked for a certain amount of time according to its rate limit rule. Green light is sent to its ready channel after the controller is unblocked and so the pending request would be forwarded.

In order to inspect the HTTP headers for rate limiting, Cuttle has to terminate SSL/TLS connection of a HTTPS request. In such case, Cuttle would establish a dedicated SSL/TLS connection with the upstream when the pending request is forwarded. Since Cuttle terminates SSL/TLS connection, the HTTP client would need to either verify the certificate sent by Cuttle with a custom certificate authority or skip verifying it.


Suggested setup
---------------

*Option 1*:

API client <-- private network --> Cuttle <-- Internet --> Origin

In this setup, the connection between Cuttle and the API client is already secured. The API client can skip verifying the SSL/TLS certificate sent from Cuttle. On the other hand, connection over the Internet is secured by SSL/TLS between Cuttle and the origin.

Sample API client - `https_proxy='127.0.0.1:3128' curl -k https://www.example.com/api/`

*Option 2*:

API client <-- public network --> Cuttle <-- Internet --> Origin

In this setup, we have to secure the communication between the API client and Cuttle. We need to create a self-signed x509 certificate and configure Cuttle to use it for signing SSL/TLS certificate it sends to the API client. Besides, the API client should use the self-signed certificate as a trusted certificate authority so that a secure connection can be established to Cuttle. Again, the connection between Cuttle and the origin is secured by SSL/TLS connection over the Internet.

To create a self-signed x509 certificate with plaintext private key. Run,

`openssl req -x509 -nodes -sha1 -newkey rsa:2048 -out cacert.pem -outform PEM -days 1825`

Sample API client - `https_proxy='127.0.0.1:3128' curl --cacert cacert.pem https://www.example.com/api/`

.

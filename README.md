# XJPVPN
> This is an early, unfinished implementation. It is intended for developers only and should **NOT** be used in a production environment.

![](static/logo.webp)


## Roadmap
| Feature       | Status | Description                                                                                    |
| ------------- | ------ | ---------------------------------------------------------------------------------------------- |
| Init Protocol | DONE   | The protocol uses websocket between a SOCKS5 client and proxy server.                          |
| CI System     | TODO   | The CI should validate everything functions as intended. Security is a priority consideration. |
| Cipher System | TODO   | HTTPS is not enabled. A possible solution is to place NGINX in front of the proxy server.      |


## Design Decision
+ **Free.** The project aims to remain free, non-profit, and simple in order to minimize security risks and maintenance needs. The current maintainer is located outside of China.
+ **Minimal.** The project will keep a minimal and simplified design, without attempting compatibility with other protocols.
+ **HTTP.** Generic HTTP is used as the tunneling protocol, but HTTPS will not be built internally. Instead, NGINX could provide HTTPS capability.


## Workflow
| Terminology    | Description                                                                                                                           |
| -------------- | ------------------------------------------------------------------------------------------------------------------------------------- |
| User           | The end user that initiates network requests, typically via a web browser.                                                            |
| Target         | The endpoint or destination that the user wants to access, usually a website.                                                         |
| Client         | The program that runs on the user's device. It communicates with the browser through SOCKS5 protocol and forwards data to the server. |
| Server         | The proxy program deployed on a cloud server. It receives requests from the client and makes the actual connections to the target.    |
| Direct Access  | User (chrome browser) -> (request to visit google.com) -> Target (Google's webserver)                                                 |
| Proxied Access | User -> (request to visit google.com) -> Client -> (obfused network traffic data) -> Server -> Target (Google's webserver)            |


## Early Testing for Developers
+ Build and run
``` bash
# 1. build
go build

# 2. execute
./xjpvpn --test
```

+ The program will open `localhost:1081` for Client/SOCKS5 and `localhost:8081` for Server/websocket, now configure your browser with proxy `socks5h://localhost:1081`


## Real-world usage and configuration
+ To be done: (1) configure NGINX with `https` enabled, (2) make `xjpvpn` as an upstream server, (3) configure the server and set the password

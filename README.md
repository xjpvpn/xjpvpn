# XJPVPN

> This is an early, unfinished implementation. It is intended for developers only and should **NOT** be used in a production environment.

![](static/logo.webp)


## Design Decision

+ **Free.** The project aims to remain free, non-profit, and simple in order to minimize security risks and maintenance needs. The current maintainer is located outside of China.
+ **Minimal.** The project will keep a minimal and simplified design, without attempting compatibility with other protocols.
+ **HTTP.** Generic HTTP is used as the tunneling protocol, but HTTPS will not be built internally. Instead, NGINX could provide HTTPS capability.


## Roadmap

| Feature       | Status | Description                                                                                    |
| ------------- | ------ | ---------------------------------------------------------------------------------------------- |
| Init Protocol | DONE   | The protocol uses websocket between a SOCKS5 client and proxy server.                          |
| CI System     | TODO   | The CI should validate everything functions as intended. Security is a priority consideration. |
| Cipher System | TODO   | HTTPS is not enabled. A possible solution is to place NGINX in front of the proxy server.      |

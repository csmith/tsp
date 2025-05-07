# tsp - tailscale proxy

A very simple binary that connects to tailscale,
listens for incoming TCP connections, and then
proxies them to somewhere else.

Designed to quickly and easily expose services
in containers on a tailnet, without having to deal
with a full tailscale container with serve configs
and so on.

For HTTP proxying, consider [Centauri](https://github.com/csmith/centauri),
which can deal with HTTP headers correctly, and add tailscale auth information
for the upstreams to use.

## Example

```docker-compose
services:
  fun:
    image: my-image
    restart: always
    # No need to export any ports or anything here.
    
  tsp:
    image: ghcr.io/csmith/tsp:latest
    restart: always
    depends_on: 
      - fun
    
    environment:
      # The name of the machine to register on the tailscale network
      TAILSCALE_HOSTNAME: "fun"
      
      # Port to listen on for incoming connections from the tailnet
      TAILSCALE_PORT: "80"
           
      # Auth key to enable automatic authentication with tailscale. If you
      # don't supply an auth key, you'll need to do an interactive auth (check
      # the logs for the link to click).
      TAILSCALE_AUTH_KEY: "tsnet-....."
      
      # The address and port of the service to proxy requests to.
      UPSTREAM: "fun:8080"
      
      # For debugging purposes. Default level of "info" should be fine for
      # most use.
      LOG_LEVEL: "debug"
      
      # Logs are by default plain text, but you can switch them to JSON if
      # desired.
      LOG_FORMAT: "json"
    
    volumes:
      - tailscale:/config
        
volumes:
  tailscale:
```
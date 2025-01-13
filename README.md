# AccessControl Middleware Plugin for Traefik

AccessControl is a middleware plugin for [Traefik](https://traefik.io/) that restricts access to your services based on a whitelist of IP addresses or the presence of a specific HTTP header with a specified value.

## Features
- Restrict access to services using an IP whitelist.
- Allow access based on a specific HTTP header and its value.
- Flexible configuration for custom needs.

## Traefik Configuration

### 1. Enable the Plugin in `traefik.yml`
Add the following configuration to your Traefik static configuration file:
```yaml
experimental:
  plugins:
    accesscontrol:
      moduleName: "github.com/<your-username>/accesscontrol"
      version: "v1.0.0"
```

### 2. Dynamic Configuration in `dynamic.yml`
Define the middleware and apply it to your routers:
```yaml
http:
  middlewares:
    accessControlMiddleware:
      plugin:
        accesscontrol:
          whitelist:
            - "<your-whitelisted-ip>"
          headerKey: "<your-header-key>"
          headerValue: "<your-header-value>"

  routers:
    example-router:
      rule: "Host(`example.com`)"
      service: example-service
      middlewares:
        - accessControlMiddleware

  services:
    example-service:
      loadBalancer:
        servers:
          - url: "http://127.0.0.1:8080"
```

### 3. Restart Traefik
Restart Traefik to apply the new configuration:
```bash
docker-compose restart traefik
```

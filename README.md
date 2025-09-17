# Vidiemme - Traefik Plugin

A list of Traefik plugins

## AccessControl Middleware Plugin for Traefik

AccessControl is a middleware plugin for [Traefik](https://traefik.io/) that restricts access to your services based on a whitelist of IP addresses or the presence of a specific HTTP header with a specified value.

### Features
- Restrict access to services using an IP whitelist.
- Allow access based on a specific HTTP header and its value.
- Flexible configuration for custom needs.

### What it does

The plugin is a *middleware* for Traefik that restricts access by either:

* checking if the client IP is on a whitelist **OR**
* checking if a certain HTTP header has a specific value.

### How to configure

You need two parts of configuration: **static** and **dynamic**.

1. **Static configuration**
   Enable the plugin in Traefik’s static configuration under `experimental.plugins`. You need the module name and version. Something like:

   ```yaml
   experimental:
     plugins:
       accesscontrol-ip-or-header:
         moduleName: github.com/vidiemme/accesscontrol-ip-or-header
         version: <plugin-version>
   ```

2. **Dynamic configuration**
   Define the middleware (with that plugin) in your HTTP dynamic configuration. Example:

   ```yaml
   http:
     middlewares:
       my-access-control:
         plugin:
           accesscontrol-ip-or-header:
             allowedIPs:
               - "192.0.2.0/24"
               - "203.0.113.5"
             headerName: "X-My-Secret"
             headerValue: "some-value"
   ```

   Then attach that middleware to your router(s). E.g.:

   ```yaml
   http:
     routers:
       my-router:
         rule: Host(`example.com`)
         service: my-service
         middlewares:
           - my-access-control
   ```

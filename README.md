<p align="center">
  <img src="https://i.imgur.com/7dtHykL.png" />
</p>

# messaget - getting messages, it's that easy 
MessaGet is a simple self-hosted SFU intended to be used as a relay to drive internal event services or serve real time component updates to frontend clients over websocket. Clients connect to one public endpoint, and your backend service authenticates over another endpoint through one of our libraries, from where it can manage clients and forward/broadcast messages.
<p align="center">
  <img src="https://i.imgur.com/oILIy31.png" />
</p>

#### Protocol Implementations
 - [JavaScript Client/Controller](https://github.com/messaget/js-client)

#### Usage
messaget is deployed from a single binary, and only accepts one argument. By default, running `./messaget` generates a `config.yml` in your current directory and just starts. You can alternatively point to a config elsewhere on your system with the `-config /path/to/config.yml` flag. There's nothing else you need to do, it's really that easy!

#### Configuration
```yaml
server:
  # Port used for the main rest/ws endpionts
  listen_port: 443
  # Auto cert configuration, ignored when disabled
  public_url: messaget.example.com
  use_auto_cert: false
  cert_path: /var/www/.cache
  
auth:
  # Password used for the controller
  password: super-secure-password
  # Client connection rate limit
  connections_per_second: 2
  # Max string length of the namespace
  max_namespace_length: 50
  # If clients should authenticate using a password, and if so, what it should be
  use_client_password: false
  client_password: another-secret-password
  
logging:
  file: ""
  json: false
  production: true

```

#### Development Notice
the messaget server is still in active development so the protocol and implementations are still subject to change, but I'm hoping to get those ready in a few months. The end goal is to make the development of real time applications as easy as possible without needing to worry about client managment or routing.
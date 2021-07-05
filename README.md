<p align="center">
  <img src="https://i.imgur.com/7dtHykL.png" />
</p>

# MessaGet - Yes, it does sound stupid.
MessaGet is a simple self-hosted SFU intended to be used as a relay to drive internal event services or serve real time component updates to frontend clients over websocket. Clients connect to one public endpoint, and your backend service authenticates over another endpoint through one of our libraries, from where it can manage clients and forward/broadcast messages
<p align="center">
  <img src="https://i.imgur.com/oILIy31.png" />
</p>

#### Protocol Implementations
 - [JavaScript Client/Controller](https://github.com/messaget/js-client)

#### Development Notice
the messaget server is still in active development so the protocol and implementations are still subject to change, but I'm hoping to get those ready in a few months. The end goal is to make the development of real time applications as easy as possible without needing to worry about client managment or routing.
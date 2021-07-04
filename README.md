<p align="center">
  <img src="https://i.imgur.com/7dtHykL.png" />
</p>
MessaGet is a simple self-hosted SFU intended to be used as a relay to drive internal event services or serve real time component updates to frontend clients over websocket. Clients connect to one public endpoint, and your backend service authenticates over another endpoint through one of our libraries, from where it can manage clients and forward/broadcast messages
<p align="center">
  <img src="https://i.imgur.com/oILIy31.png" />
</p>

the messaget server is still in active development and there aren't any client/backend libraries available yet, but I'm hoping to get those ready in a few months. The end goal is to make the development of real time applications as easy as possible without needing to worry about client managment or routing.
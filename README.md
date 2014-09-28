#Go DDP Client

###A golang ddp client for meteor

This work is in early development and is not ready for production.

####To run

1. Make sure your GOPATH includes ddp package
2. go run main.go http://mysite.com 3000 websocket
   
####What works

1. connects to remote host and gets session id
2. ping-pong

####Development checklist

1. subscribe/unsubscribe, handling of client side collections
2. login/logout
3. method call

Based on https://github.com/meteor/meteor/blob/devel/packages/ddp/DDP.md and 
https://github.com/arunoda/meteor-ddp-analyzer
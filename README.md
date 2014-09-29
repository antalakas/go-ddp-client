#Go DDP Client

###A golang ddp client for meteor

This work is in early development and is not ready for production.

####To run

1. Create a run.sh file next to main.go
2. chmod +x run.sh
3. Fill it it with following:
    - export GOPATH=$GOPATH:~/path/go-ddp-client
    - export MY_USERNAME='username'
    - export MY_PASSWORD='password'
    - go run main.go http://mysite.com 3000 websocket
   
####What works

1. connects to remote host and gets session id
2. login/logout
3. ping-pong

####Development checklist

1. subscribe/unsubscribe, handling of client side collections
2. method call

Based on https://github.com/meteor/meteor/blob/devel/packages/ddp/DDP.md and 
https://github.com/arunoda/meteor-ddp-analyzer
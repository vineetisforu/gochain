A simple code to implement Block-Chain in Go,

Pre-requisite:

Go 1.16 or higher
Unix/Linux terminal

How to run?

Just clone the project to your directory

1. git clone [gitlink]
2. go to the project leaf folder
3. run the following command,
	go run main.go

4. Once the commmand is run open few more terminal nodes and press the following command,
	nc localhost 9000
	
5. On the different terminal windows 
	enter the token amount : 1-100 (some random integer)
	enter the inventory counts : 1-1000 (some random integer) 

6. While you are entering the values on the #5, open a terminal or postman and try to connect to the app to get the validator status using,
	curl --location 'http://localhost:8080/getValidatorPer'
	curl --location 'http://localhost:8080/getActiveValidatorsPer'
	

While running the application you would get log something similar to this,

(main.Block) {
 Index: (int) 0,
 Timestamp: (string) (len=51) "2023-07-16 23:33:38.651341 +0530 IST m=+0.015777369",
 inventory: (int) 0,
 Hash: (string) (len=64) "96a296d224f285c67bee93c30f8a309157f0daa35dc5b87e410b78630a09cfc7",
 PrevHash: (string) "",
 Validator: (string) ""
}
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /getActiveValidatorsPer   --> main.getActiveValidatorsPer (3 handlers)
[GIN-debug] GET    /getValidatorPer          --> main.getValidatorPer (3 handlers)
[GIN-debug] [WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.
Please check https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies for details.
[GIN-debug] Listening and serving HTTP on localhost:8080
2023/07/16 23:33:38 TCP Server Listening on port : 9000
ActiveValidatorsPercent is  {0}
[GIN] 2023/07/16 - 23:33:41 | 200 |    1.123371ms |       127.0.0.1 | GET      "/getActiveValidatorsPer"
map[dabc73e82a43eeee6a6be89fba388312fc25c829aa1d44cc3f96b9ca2665ab1d:1]
map[dabc73e82a43eeee6a6be89fba388312fc25c829aa1d44cc3f96b9ca2665ab1d:1]
ActiveValidatorsPercent is  {100}
[GIN] 2023/07/16 - 23:34:49 | 200 |      31.765µs |       127.0.0.1 | GET      "/getActiveValidatorsPer"
map[c8600947b2ca714d362a74e6f16b1d728c110d0ba54821bfc23256c05a82d3e9:2 dabc73e82a43eeee6a6be89fba388312fc25c829aa1d44cc3f96b9ca2665ab1d:1]
ActiveValidatorsPercent is  {50}
[GIN] 2023/07/16 - 23:34:56 | 200 |      29.852µs |       127.0.0.1 | GET      "/getActiveValidatorsPer"
validatorsPercent is  {50}
[GIN] 2023/07/16 - 23:35:15 | 200 |      86.589µs |       127.0.0.1 | GET      "/getValidatorPer"
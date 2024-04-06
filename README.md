# pow-fiat-shamir

It is necessary for the client to prove to the server that he knows the secret of x during calculations using the Fiat-Shamir protocol.
[article](https://asecuritysite.com/golang/go_fiat2). At the same time, a random sequence of bytes arrives each time. 
Stages: 
1. The client sends a request - request service. The server sends a random 32 bytes and the points G and H for the elliptical curve in the response.
2. The client performs calculations and sends the result to the server cryptographic results with a payload.

The server does not need to check the mail or hash in the database as in hashcash. Only check random 32 bytes for originality. Perhaps using points on an elliptic curve, it is possible to generate different results for one sequence of bytes, which will eliminate the storage of all previously obtained hashes as hashcash, since it will be necessary to calculate for the same sequence of bytes but with different data for the elliptic curve in any case, otherwise the cryptographic results will not pass verification on the server.

At the moment, it is not clear how the complexity can be adjusted by choosing the points of the elliptic curve (G and H) and whether this is possible.

```
server - go run pow-fiat-shamir server --config ./example/config.yml
client - go run pow-fiat-shamir client --config ./example/config.yml
```

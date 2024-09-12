To run the TCP server, run the following command:

```bash
go build
./broadcast-server
```

To connect as a client, you can use the following command (you can open multiple terminals to simulate multiple clients):

```bash
nc localhost 3000
```

If you write something in the terminal where the client is connected, it will be broadcasted to all the clients connected to the server.

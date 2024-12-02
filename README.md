# just-servers 
Collection of Networking Server Implementations in Pure Go

## Challenges
[Protohacker Challenges](https://protohackers.co)

- [0: Smoke Test](https://github.com/Nish7/just-servers/tree/main/0_smoke_test)
- [1: Prime Time](https://github.com/Nish7/just-servers/tree/main/1_prime_time)
- [2: Means to an End](https://github.com/Nish7/just-servers/tree/main/2_means_to_an_end)
- [3: Budget Chat](https://github.com/Nish7/just-servers/tree/main/3_budget_chat)
- [4: Unusual Database Program](https://github.com/Nish7/just-servers/tree/main/4_unusual_database_program)
- [5: Mob in the Middle](https://github.com/Nish7/just-servers/tree/main/5_mob_in_the_middle)

### Deployment
Deployed a Droplet (1vCPU and 1GB Disk) on DigitalOcean

Steps:
```sh
# SSH into your droplet
$ ssh root@<your-droplet-ip>

# Download the Go binary
$ wget https://golang.org/dl/go1.20.3.linux-amd64.tar.gz

# Extract and install Go
$ sudo tar -C /usr/local -xzf go1.20.3.linux-amd64.tar.gz

# Add Go to the PATH
$ export PATH=$PATH:/usr/local/go/bin
$ source ~/.profile

# Verify the Go installation
$ go version

# Clone your TCP server repository
$ git clone https://github.com/username/your-tcp-server.git

# Build the TCP server
$ cd /path/to/tcp-server && go build -o tcp-server

# Run the TCP server in the background
$ ./tcp-server &
```

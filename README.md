# just-servers 
Low-Level Go Server Implementations

## Challenges
[Protohacker Challenges](https://protohackers.co)

- [0: Smoke Test](https://github.com/Nish7/just-servers/tree/main/0_smoke_test)

### Deployment
- Deployed a Droplet on DigitalOcean
- Specs: 1vCPU and 1GB Disk

Steps:
1. `ssh root@<your-droplet-ip>`
2. `wget https://golang.org/dl/go1.20.3.linux-amd64.tar.gz`
3. `sudo tar -C /usr/local -xzf go1.20.3.linux-amd64.tar.gz`
4. `export PATH=$PATH:/usr/local/go/bin`
5. `source ~/.profile`
6. `go version`
7. `git clone https://github.com/username/your-tcp-server.git`
8. `cd /path/to/tcp-server && go build -o tcp-server`
9. `./tcp-server &`

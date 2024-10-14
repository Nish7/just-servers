# just-servers
- Protohacker Challenges (https://protohackers.com/about)

# Deployment
- Deployed a Droplet on DigitalOcean
- Specs: 1vCPU and 1GB Disk
- Steps:
1.  `ssh root@<your-droplet-ip>`
2. ```
wget https://golang.org/dl/go1.20.3.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.20.3.linux-amd64.tar.gz
```
3. `export PATH=$PATH:/usr/local/go/bin`
4. `source ~/.profile`
5. `go version`
6. `git clone https://github.com/username/your-tcp-server.git`
7. ```
cd /path/to/tcp-server
go build -o tcp-server` 
```
8. `./tcp-server &`

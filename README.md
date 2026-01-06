# DNS UDP-to-TCP Bridge

A tool to help converting and acting as a bridge for DNS UDP packet to TCP packet format.

This case used if you just have a DNS on TCP listener like when using SSH Tunnel to forward the DNS TCP port to your main server but most DNS client request is using UDP packet format.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.
See deployment section for notes on how to deploy the project on a live system.

### Prerequisites

Prequisites packages:
* Go (Go Programming Language)
* GoReleaser (Go Automated Binaries Build)
* Make (Automated Execution using Makefile)

Optional packages:
* Docker (Application Containerization)

### Deployment

#### **Using Container**

1) Install Docker CE based on the [manual documentation](https://docs.docker.com/desktop/)

2) Run the following command on your Terminal or PowerShell
```sh
docker run -d \
  -p <UDP_PORT>:<UDP_PORT>
  --name dns-to-tcp \
  --rm dimaskiddo/dns-to-tcp:latest \
  dns-to-tcp -lh 0.0.0.0 -lp <UDP_PORT> -rh <REMOTE_TCP_IP> -rp <REMOTE_TCP_PORT> -p <CONNECTION_POOL_SIZE>

# Example of Usage

docker run -d \
  -p 53:53
  --name dns-to-tcp \
  --rm dimaskiddo/dns-to-tcp:latest \
  dns-to-tcp -lh 0.0.0.0 -lp 53 -rh 9.9.9.9 -rp 9953 -p 1024
```

#### **Using Pre-Build Binaries**

1) Download Pre-Build Binaries from the [release page](https://github.com/dimaskiddo/dns-to-tcp/releases)

2) Extract the zipped file

3) Run the pre-build binary
```sh
# MacOS / Linux
chmod 755 dns-to-tcp
# -- Example of Usage
# -- ./dns-to-tcp -lh 0.0.0.0 -lp <UDP_PORT> -rh <REMOTE_TCP_IP> -rp <REMOTE_TCP_PORT> -p <CONNECTION_POOL_SIZE>
./dns-to-tcp -lh 0.0.0.0 -lp 53 -rh 9.9.9.9 -rp 9953 -p 1024

# Windows
# You can double click it or using PowerShell
# -- Example of Usage
# -- .\dns-to-tcp.exe -lh 0.0.0.0 -lp <UDP_PORT> -rh <REMOTE_TCP_IP> -rp <REMOTE_TCP_PORT> -p <CONNECTION_POOL_SIZE>
.\dns-to-tcp.exe -lh 0.0.0.0 -lp 53 -rh 9.9.9.9 -rp 9953 -p 1024
```

#### **Build From Source**

Below is the instructions to make this source code running:

1) Create a Go Workspace directory and export it as the extended GOPATH directory
```sh
cd <your_go_workspace_directory>
export GOPATH=$GOPATH:"`pwd`"
```

2) Under the Go Workspace directory create a source directory
```sh
mkdir -p src/github.com/dimaskiddo/dns-to-tcp
```

3) Move to the created directory and pull codebase
```sh
cd src/github.com/dimaskiddo/dns-to-tcp
git clone -b master https://github.com/dimaskiddo/dns-to-tcp.git .
```

4) Run following command to pull vendor packages
```sh
make vendor
```

5) Until this step you already can run this code by using this command
```sh
make run
```

6) *(Optional)* Use following command to build this code into binary spesific platform
```sh
make build
```

7) *(Optional)* To make mass binaries distribution you can use following command
```sh
make release
```

### Running The Tests

Currently the test is not ready yet :)

## Built With

* [Go](https://golang.org/) - Go Programming Languange
* [GoReleaser](https://github.com/goreleaser/goreleaser) - Go Automated Binaries Build
* [Make](https://www.gnu.org/software/make/) - GNU Make Automated Execution
* [Docker](https://www.docker.com/) - Application Containerization

## Authors

* **Dimas Restu Hidayanto** - *Initial Work* - [DimasKiddo](https://github.com/dimaskiddo)

See also the list of [contributors](https://github.com/dimaskiddo/dns-to-tcp/contributors) who participated in this project

## Annotation

You can seek more information for the make command parameters in the [Makefile](https://github.com/dimaskiddo/dns-to-tcp/-/raw/master/Makefile)

## License

Copyright (C) 2026 Dimas Restu Hidayanto

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

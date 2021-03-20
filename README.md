# open-balance

<p align="center"> Open-Balance is a free and open-source implementation of balance system that serves the foundation for a payment system or a reward system. 
</p>
<p align="center">
</p>

<div align="center">
  <h3>
    <a href="https://www.opendigitalpay.io">
      Website
    </a>
    <span> | </span>
    <a href="https://docs.opendigitalpay.io">
      Documentation
    </a>
    <span> | </span>
    <a href="https://docs.opendigitalpay.io/API/balance/v1/">
      API
    </a>
    <span> | </span>
    <a href="https://docs.opendigitalpay.io/Contribute/">
      Contribute
    </a>
    <span> | </span>
  </h3>
</div>
<br/>

<p align="center">
  <a href="https://demo.opendigitalpay.io">View Demo</a>
  ·
  <a href="https://github.com/opendigitalpay-io/open-balance/issues">Report a bug</a>
  ·
  <a href="https://github.com/opendigitalpay-io/open-balance/discussions/new">Request a feature</a>
  ·
  <a href="https://docs.opendigitalpay.io/FAQ/">FAQ</a>
</p>

## 💼 Table of Contents

* [Use Cases](#-usecases)
* [Getting Started](#-getting-started)
* [Documentation](#-documentation)
* [Contributing](#-contributing)
* [Developing](#-developing)
    * [API](#-api)
* [Community](#-community)
* [License](#-license)

## 🎨 usecases
* use as a balance system
* use as a reward system

# Local Environment Setup
```shell
# (Optional) 0. make sure you have go installed
 brew install go 

# 1. enable go modules; Ref: https://blog.golang.org/using-go-modules
go env -w GO111MODULE=on

# 2. fallback option: pull modules from orgin of the source code (e.g., Github) when failing to pull from go modules
go env -w GOPROXY=direct 

# 3. manually fetch all external modules (dependencies) from go.mod
go mod download

# (Optional) 3.1 add a new dependency or change the required version of a dependency  
go get -u github.com/gin-gonic/gin

# (Optional) 3.2 implicitly add new dependencies to go.mod as needed
go build 
go test

# (Optional) 3.3 clean up unused dependencies:
go mod tidy

# 4. start the port locally
docker-compose up -d
*manually run sql queries in ./sql/schema.sql 
go run main.go

# 5. test local port
curl -X GET 127.0.0.1:8080/v1/user/1
curl -X POST 127.0.0.1:8080/v1/user -H "Content-Type: application/json" --data-raw '{"email": "xxx@gmail.com", "phone": "4166666666", "userName": "carlos"}'
```
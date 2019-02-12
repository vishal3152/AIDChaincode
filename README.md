[![Build Status](https://travis-ci.org/vishal3152/AIDChaincode.svg?branch=master)](https://travis-ci.org/vishal3152/AIDChaincode)

# AIDChaincode
Chaincode for managing life cycle of charity donations and spends.

### Project structure:
```sh
package-name
├── main
|  ├── app.go               --> chaincode entry point
|  ├── init.go              --> Chaincode Interface Init implementation 
|  ├── invoke.go            --> Chaincode Interface Invoke implementation 
|  ├── interfaces.go        --> AidAssetInterface interface 
|  ├── validate.go          --> Input arguments validations
|  ├── project.go           --> Project asset implements AidAssetInterface
|  ├── project_test.go      --> Unit tests for project asset
|  ├── item.go              --> Item asset implements AidAssetInterface
|  ├── item_test.go         --> Unit tests for item asset
|  ├── donation.go          --> Donation asset implements AidAssetInterface
|  ├── donation_test.go     --> Unit tests for donation asset
|  ├── spend.go             --> Spend asset implements AidAssetInterface
|  ├── spend_test.go        --> Unit tests for spend asset         
|  ├── util.go              --> Utility functions
   └── main_test.go         --> TestMain(m *testing.M) implementaion
```
### Prerequisites:
* [Golang](https://golang.org/dl/) - version go1.11.2
* [VSCode](https://code.visualstudio.com/download) - Or any other text editor of your choice
* [Govendor](https://github.com/kardianos/govendor) - Dependency vendoring tool


### Setting up Project:
Set GOPATH environment variables.
The GOPATH environment variable specifies the location of your workspace. Create two directory under GOPATH(workspace):

* src - contains Go source files, and
* bin - contains executable commands. 

Also,  add the workspace's bin subdirectory to your PATH.


Install govendor:
```sh
$ go get -u github.com/kardianos/govendor
```

```sh
$ cd GOPATH/src
```

Clone the chaincode project:
```sh
$ git clone https://github.com/vishal3152/AIDChaincode.git
```
```sh
$ cd AIDChaincode/aidcc
```

Download dependencies:
```sh
$ govendor add +external
$ govendor sync
```

### Unit test the chaincode:
```sh
$ go test
```

### Building the chaincode (go files ending with *_test are excluded from build)):
```sh
$ go build
```




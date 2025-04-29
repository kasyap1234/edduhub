# mockery install 
brew install mockery 
go install github.com/vektra/mockery/v3@v3.2.4
or use v2 if getting error --all 

# mockery update 
brew upgrade mockery 

# golangci-lint 
brew install golangci-lint
brew upgrade golangci-lint
or 
brew tap golangci/tap
brew install golangci/tap/golangci-lint

## Install task file 

brew install go-task/tap/go-



## gopfumpt
go install mvdan.cc/gofumpt@latest
export PATH="$PATH:$(go env GOPATH)/bin"

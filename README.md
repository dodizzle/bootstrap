1. Install go.

2. Add these lines in your~/.bash_profile.

`export PATH=$PATH:$(go env GOPATH)/bin`

`export GOPATH=$(go env GOPATH)`

3. After you edit the code run this command to build a new binary.

`go install`

The new binary should be in your path.

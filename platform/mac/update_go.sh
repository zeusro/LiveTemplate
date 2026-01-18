
    //备份
    mkdir -p $GOPATH/bin-old
    cp $GOPATH/bin/* $GOPATH/bin-old/
    
    //dlv
    mkdir -p $GOPATH/src/github.com/go-delve
    cd $GOPATH/src/github.com/go-delve/delve
    git pull
    cd cmd/dlv
    go install
    //golang
    x=$GOPATH/src/golang.org/x
    cd $x/tools
    git pull
    cd $x/tools/cmd/guru
    go install
    cd $x/tools/cmd/gorename
    go install
    cd $x/sys
    git pull
    cd $x/net
    git pull
    cd $x/time
    git pull
    cd $x/lint
    git pull
    go install
    cd $x/sync
    git pull
    go get -v golang.org/x/tools/cmd/goimports
    go get -v golang.org/x/tools/cmd/gopls
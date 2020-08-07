
brew install go

cat << EOF >>  ~/.zshrc
export GOPATH=$HOME/go
export GOPROXY=https://goproxy.io
EOF
source ~/.zshrc

cat << EOF >> ~/.bash_profile
export GOPATH=$HOME/go
export GOPROXY=https://goproxy.io
EOF
source ~/.bash_profile
echo $GOPATH
mkdir -p $GOPATH/src
mkdir -p $GOPATH/src/golang.org/x/
cd $GOPATH/src/golang.org/x/
xgo=$GOPATH/src/golang.org/x/
git clone git@github.com:golang/tools.git
cd $xgo/tools/cmd/goimports && go install
cd $xgo/tools/gopls && go install
git clone git@github.com:golang/sys.git
git clone git@github.com:golang/net.git
git clone git@github.com:golang/time.git
git clone git@github.com:golang/lint.git
cd $xgo/lint/golint && go install
git clone git@github.com:golang/sync.git
git clone git@github.com:golang/mod.git
git clone git@github.com:golang/xerrors.git



cd $GOPATH/src/github.com
gayhub=$GOPATH/src/github.com
# https://github.com/Microsoft/vscode-go/wiki/Go-tools-that-the-Go-extension-depends-on

cd $gayhub
git clone git@github.com:sqs/goreturns.git sqs/goreturns && cd $gayhub/sqs/goreturns && go install
# go get -u -v  github.com/sqs/goreturns

cd $gayhub
git clone git@github.com:go-delve/delve.git go-delve/delve && cd $gayhub/go-delve/delve/cmd/dlv && go install
# go get -u -v  github.com/go-delve/delve/cmd/dlv

cd $gayhub
git clone git@github.com:zmb3/gogetdoc.git zmb3/gogetdoc && cd $gayhub/zmb3/gogetdoc && go install
# go get -v github.com/zmb3/gogetdoc

cd $gayhub
git clone git@github.com:stamblerre/gocode.git  && cd $gayhub/stamblerre/gocode && go install && go build  && cp gocode ~/go/bin/gocode-gomod

# git clone -b bingo https://github.com/saibing/tools.git
# cd tools/cmd/gopls
# go install
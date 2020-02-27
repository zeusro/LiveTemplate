
# mac装机

```bash

brew() {
#brew
/usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
#替换brew.git:
cd "$(brew --repo)"
git remote set-url origin https://mirrors.ustc.edu.cn/brew.git
#替换homebrew-core.git:
cd "$(brew --repo)/Library/Taps/homebrew/homebrew-core"
git remote set-url origin https://mirrors.ustc.edu.cn/homebrew-core.git

echo 'export HOMEBREW_BOTTLE_DOMAIN=https://mirrors.ustc.edu.cn/homebrew-bottles' >> ~/.zshrc
source ~/.zshrc

brew install https://raw.githubusercontent.com/kadwanev/bigboybrew/master/Library/Formula/sshpass.rb


brew install  yarn gradle mpv wget kubernetes-cli maven node telnet git-lfs iproute2mac
brew install kubernetes-cli  kubectx
brew cask install java
}


zsh (){
sh -c "$(curl -fsSL https://raw.github.com/robbyrussell/oh-my-zsh/master/tools/install.sh)"
cd ~/.oh-my-zsh/custom/plugins
#zsh-autosuggestions
git clone git@github.com/zsh-users/zsh-autosuggestions ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autosuggestions
#zsh-syntax-highlighting
git clone git@github.com/zsh-users/zsh-syntax-highlighting.git ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-syntax-highlighting

source ~/.zshrc

# 中文设置
# https://blog.fazero.me/2015/09/04/Mac-iTerm2--chinese/
touch ~/.bash_profile
cat << EOF >>  ~/.zshrc
export LC_ALL=en_US.UTF-8  
export LANG=en_US.UTF-8
source ~/.bash_profile
EOF

}


updatebrew(){
    brew update
    brew upgrade
    brew cleanup
}

```

## 安装升级go

```


installgo (){

brew install go
cat << EOF >>  ~/.zshrc    
export GOPATH=$HOME/go
EOF

echo $GOPATH
cd $GOPATH/src
mkdir -p golang.org/x/
cd golang.org/x/
git clone git@github.com:golang/tools.git
git clone git@github.com:golang/sys.git
git clone git@github.com:golang/net.git
git clone git@github.com:golang/time.git
git clone git@github.com:golang/lint.git
git clone git@github.com:golang/sync.git

# https://github.com/Microsoft/vscode-go/wiki/Go-tools-that-the-Go-extension-depends-on
go get -u -v github.com/mdempsky/gocode
go get -u -v  github.com/ramya-rao-a/go-outline
go get -u -v  github.com/acroca/go-symbols
go get -u -v  github.com/stamblerre/gocode
go get -u -v  github.com/sqs/goreturns
go get -u -v  github.com/go-delve/delve/cmd/dlv

# 接着需要修改~/.bash_profile,配置3个变量
# export GOPATH=~/go
# export GOBIN=$GOPATH/bin
# export GOROOT=/usr/local/Cellar/go/1.12.4/libexec

go get -v golang.org/x/tools/cmd/gopls
go get -v golang.org/x/tools/cmd/goimports
go get -v github.com/zmb3/gogetdoc

# git clone -b bingo https://github.com/saibing/tools.git
# cd tools/cmd/gopls
# go install
}

updatego(){
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
    
}

```

 需要更新的工具有
  gocode
  gopkgs
  go-outline
  go-symbols
  guru
  gorename
  dlv
  gocode-gomod
  gogetdoc
  goimports
  golint
  gopls
等

## 修改host破解迅雷版权校验

```bash
# sudo vi /etc/hosts

127.0.0.1 hub5btmain.sandai.net 
127.0.0.1 hub5emu.sandai.net 
127.0.0.1 upgrade.xl9.xunlei.com
```


## 参考：

http://www.361way.com/cat-eof-cover-append/4298.html


brew uninstall --force tomcat
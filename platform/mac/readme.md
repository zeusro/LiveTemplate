

[Mac 命令行下编辑常用的快捷键](http://notes.11ten.net/mac-command-line-editing-commonly-used-shortcut-keys.html)

## 截屏

command+shift+3 三个键按下则抓取/截取全屏……

command+shift+4 然后用鼠标框选则抓取该区域的截图……

command+shift+4 然后按空格则抓取软件窗口。截图会自动保存到桌面。

## mac装机

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


brew install  yarn gradle mpv wget maven node telnet git-lfs iproute2mac
# cloud native
brew install kubernetes-cli  kubectx
brew cask install java
}


zsh (){
sh -c "$(curl -fsSL https://raw.github.com/robbyrussell/oh-my-zsh/master/tools/install.sh)"
cd ~/.oh-my-zsh/custom/plugins
# brew install https://raw.githubusercontent.com/Homebrew/homebrew-core/master/Formula/zsh-autosuggestions.rb

#zsh-autosuggestions 仓库已失效
# git clone git@github.com/zsh-users/zsh-autosuggestions ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autosuggestions


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

```bash


installgo (){

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
  gopkgs
  gocode-gomod
  gogetdoc
  golint
等

## 修改host破解迅雷版权校验

```bash
# sudo vi /etc/hosts

127.0.0.1 hub5btmain.sandai.net
127.0.0.1 hub5emu.sandai.net
127.0.0.1 upgrade.xl9.xunlei.com
```

## 安装Composer

```bash
curl -sS https://getcomposer.org/installer | php
php composer.phar --version
mv composer.phar /usr/local/bin/composer
composer selfupdate
```

## 参考链接:

1. [Mac中Composer的安装和使用](https://www.jianshu.com/p/fd1b53df3f4b)

## 参考

cat EOF追加与覆盖

http://www.361way.com/cat-eof-cover-append/4298.html


brew uninstall --force tomcat


## 其他

    #设置文件可执行
    chmod a+x /usr/local/bin/fly


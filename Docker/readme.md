

## Ubuntu 18在安装docker-compose 之后无法 docker login

```
sudo apt install gnupg2 pass 
gpg2 --full-generate-key
```

## Ubuntu 18 安装docker

```bash
sudo apt-get update -y
sudo apt-get remove docker docker-engine docker.io
sudo apt install -y docker.io
sudo systemctl start docker
sudo systemctl enable docker

```

## docker-compose

```bash
docker-compose rm --all &&
 docker-compose pull &&
 docker-compose build --no-cache &&
 docker-compose up -d --force-recreate &&
```

# dev
dev test

# prepare for Mac
```zsh
go install github.com/joho/godotenv/cmd/godotenv@latest
godotenv -f .env zsh
go run main.go
```

# install valkey for test
```zsh
brew install valkey
brew services start valkey
valkey-cli ping
```
### for docker
```zsh
docker run -d -p 6379:6379 --name valkey-test valkey/valkey:latest
docker run -v ./config/valkey:/usr/local/etc/valkey -d -p 6379:6379 --name valkey-test valkey/valkey:latest valkey-server /usr/local/etc/valkey/valkey.conf
```
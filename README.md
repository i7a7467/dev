# dev
dev test

# prepare for Mac
```
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
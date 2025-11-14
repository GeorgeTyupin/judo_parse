#!/bin/bash

# Читаем версию из .env
VERSION=$(grep "VERSION" configs/.env | cut -d '=' -f2 | tr -d ' ')

# Сборка для Windows 11
GOOS=windows GOARCH=amd64 go build -o judo_parse_v${VERSION}.exe ./cmd/parse

echo "Сборка завершена: judo_parse_v${VERSION}.exe"

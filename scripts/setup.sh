#!/bin/bash

# install dependencies
sudo apt install -y make

# install docker
echo "Installing Docker..."
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# setup .env
cp .env.sample .env

echo -n "DISCORD_BOT_TOKEN: "
read BOT_TOKEN
sed -i -e "/BOT_TOKEN/c BOT_TOKEN=$BOT_TOKEN" .env

echo -n "DISCORD_GUILD_ID: "
read GUILD_ID
sed -i -e "/GUILD_ID/c GUILD_ID=$GUILD_ID" .env

echo -n "DISCORD_README_MESSAGE_ID: "
read README_MESSAGE_ID
sed -i -e "/README_MESSAGE_ID/c README_MESSAGE_ID=$README_MESSAGE_ID" .env

echo -n "DISCORD_README_ROLE_ID: "
read README_ROLE_ID
sed -i -e "/README_ROLE_ID/c README_ROLE_ID=$README_ROLE_ID" .env

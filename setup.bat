@ECHO OFF
setlocal EnableDelayedExpansion
title setup frontend project

echo "current work dir: %~s0"

git clone https://github.com/215952619/GoProjectFrontend.git frontend

cd frontend

git pull

npm install

npm run build

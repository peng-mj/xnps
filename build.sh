#/bash/sh
cd $(dirname $0)

if [ -d "bin" ]; then
  rm -rf bin/
fi

export VERSION=0.26.16
export GOPROXY=direct
sudo apt-get update
sudo apt-get install gcc-mingw-w64-i686 gcc-multilib

mkdir bin



###################################-NPC-#########################################

echo "build binary to xnpc and xnps"
{

  {
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -extldflags -static -extldflags -static" -o bin/xnpc_linux_amd64 ./cmd/npc/npc.go
  } &
  {
    CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -ldflags "-s -w -extldflags -static -extldflags -static" -o bin/xnpc_linux_386 ./cmd/npc/npc.go
  } &
  {
    CGO_ENABLED=0 GOOS=freebsd GOARCH=386 go build -ldflags "-s -w -extldflags -static -extldflags -static" -o bin/xnpc_freebsd_386 ./cmd/npc/npc.go
  } &
  {
    CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 go build -ldflags "-s -w -extldflags -static -extldflags -static" -o bin/xnpc_freebsd_amd64 ./cmd/npc/npc.go
  } &
  {
    CGO_ENABLED=0 GOOS=freebsd GOARCH=arm go build -ldflags "-s -w -extldflags -static -extldflags -static" -o bin/xnpc_freebsd_arm ./cmd/npc/npc.go
  } &
  {
    CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -ldflags "-s -w -extldflags -static -extldflags -static" -o bin/xnpc_linux_arm_v7 ./cmd/npc/npc.go
  } &
  {
    CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -ldflags "-s -w -extldflags -static -extldflags -static" -o bin/xnpc_linux_arm_v6 ./cmd/npc/npc.go
  } &
  {
    CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=5 go build -ldflags "-s -w -extldflags -static -extldflags -static" -o bin/xnpc_linux_arm_v5 ./cmd/npc/npc.go
  } &
  {
    CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "-s -w -extldflags -static -extldflags -static" -o bin/xnpc_linux_arm64 ./cmd/npc/npc.go
  } &
  {
    CGO_ENABLED=0 GOOS=linux GOARCH=mips64 go build -ldflags "-s -w -extldflags -static -extldflags -static" -o bin/xnpc_linux_mips64 ./cmd/npc/npc.go
  } &
  {
    CGO_ENABLED=0 GOOS=linux GOARCH=mips64le go build -ldflags "-s -w -extldflags -static -extldflags -static" -o bin/xnpc_linux_mips64le ./cmd/npc/npc.go
  } &
  {
    CGO_ENABLED=0 GOOS=linux GOARCH=mipsle go build -ldflags "-s -w -extldflags -static -extldflags -static" -o bin/xnpc_linux_mipsle ./cmd/npc/npc.go
  } &
  {
    CGO_ENABLED=0 GOOS=linux GOARCH=mips go build -ldflags "-s -w -extldflags -static -extldflags -static" -o bin/xnpc_linux_mips ./cmd/npc/npc.go
  } &
  {
    CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -ldflags "-s -w -extldflags -static -extldflags -static" -o bin/xnpc_windows_386.exe ./cmd/npc/npc.go
  } &
  {
    CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w -extldflags -static -extldflags -static" -o bin/xnpc_windows_amd64.exe ./cmd/npc/npc.go
  } &
  {
    CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w -extldflags -static -extldflags -static" -o bin/xnpc_darwin_amd64 ./cmd/npc/npc.go
  } &
  {
    CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w -extldflags -static -extldflags -static" -o bin/xnpc_darwin_arm64 ./cmd/npc/npc.go
  } &
} &

###################################-NPS-#########################################

{
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -extldflags -static -extldflags -static" -o bin/xnps_linux_amd64 ./cmd/nps/nps.go
} &
{
  CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -ldflags "-s -w -extldflags -static -extldflags -static" -o bin/xnps_linux_386 ./cmd/nps/nps.go
} &
{
  CGO_ENABLED=0 GOOS=freebsd GOARCH=386 go build -ldflags "-s -w -extldflags -static -extldflags -static" -o bin/xnps_freebsd_386 ./cmd/nps/nps.go
} &
{
  CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 go build -ldflags "-s -w -extldflags -static -extldflags -static" -o bin/xnps_freebsd_amd64 ./cmd/nps/nps.go
} &
{
  CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w -extldflags -static -extldflags -static" -o bin/xnps_darwin_amd64 ./cmd/nps/nps.go
} &
{
  CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w -extldflags -static -extldflags -static" -o bin/xnps_darwin_arm64 ./cmd/nps/nps.go
} &
{
  CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w -extldflags -static -extldflags -static" -o bin/xnps_windows_amd64.exe ./cmd/nps/nps.go
} &
{
  CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -ldflags "-s -w -extldflags -static -extldflags -static" -o bin/xnps_windows_386.exe ./cmd/nps/nps.go
} &

wait

echo "make package and compress"

cp -r conf/ bin/conf
mkdir bin/web
cp -r web/static bin/web/
cp -r web/views bin/web/
cd bin/



{

  if [ -f "xnps_linux_amd64" ]; then
    mv xnps_linux_amd64 xnps
    tar -czvf ./xnps_linux_amd64.tar.gz conf/nps.conf conf/tasks.json conf/clients.json conf/hosts.json conf/server.key conf/server.pem web/views web/static xnps
  fi

  if [ -f "xnps_linux_386" ]; then
    mv xnps_linux_amd64 xnps
    tar -czvf ./xnps_linux_386.tar.gz conf/nps.conf conf/tasks.json conf/clients.json conf/hosts.json conf/server.key conf/server.pem web/views web/static xnps
  fi

  if [ -f "xnps_freebsd_386" ]; then
    mv xnps_linux_amd64 xnps
    tar -czvf ./xnps_freebsd_386.tar.gz conf/nps.conf conf/tasks.json conf/clients.json conf/hosts.json conf/server.key conf/server.pem web/views web/static xnps
  fi

  if [ -f "xnps_freebsd_amd64" ]; then
    mv xnps_linux_amd64 xnps
    tar -czvf ./xnps_freebsd_amd64.tar.gz conf/nps.conf conf/tasks.json conf/clients.json conf/hosts.json conf/server.key conf/server.pem web/views web/static xnps
  fi

  if [ -f "xnps_darwin_amd64" ]; then
    mv xnps_linux_amd64 xnps
    tar -czvf ./xnps_darwin_amd64.tar.gz conf/nps.conf conf/tasks.json conf/clients.json conf/hosts.json conf/server.key conf/server.pem web/views web/static xnps
  fi

  if [ -f "xnps_darwin_arm64" ]; then
    mv xnps_linux_amd64 xnps
    tar -czvf ./xnps_darwin_arm64.tar.gz conf/nps.conf conf/tasks.json conf/clients.json conf/hosts.json conf/server.key conf/server.pem web/views web/static xnps
  fi

  if [ -f "xnps_windows_amd64.exe" ]; then
    mv "xnps_linux_amd64.exe" "xnps.exe"
    tar -czvf ./xnps_windows_amd64.tar.gz conf/nps.conf conf/tasks.json conf/clients.json conf/hosts.json conf/server.key conf/server.pem web/views web/static xnps.exe
  fi

  if [ -f "xnps_windows_386.exe" ]; then
    mv "xnps_linux_amd64.exe" "xnps.exe"
    tar -czvf ./xnps_windows_386.tar.gz conf/nps.conf conf/tasks.json conf/clients.json conf/hosts.json conf/server.key conf/server.pem web/views web/static xnps.exe
  fi
  rm -r xnps xnps.exe
} &

{
  
  if [ -f "xnpc_linux_amd64" ]; then
    mv xnpc_linux_amd64 xnpc
    tar -czvf xnpc_linux_amd64.tar.gz xnpc conf/npc.conf
  fi

  if [ -f "xnpc_linux_386" ]; then
    mv xnpc_linux_386 xnpc
    tar -czvf xnpc_linux_386.tar.gz xnpc conf/npc.conf
  fi

  if [ -f "xnpc_freebsd_386" ]; then
    mv xnpc_freebsd_386 xnpc
    tar -czvf xnpc_freebsd_386.tar.gz xnpc conf/npc.conf
  fi

  if [ -f "xnpc_freebsd_amd64" ]; then
    mv xnpc_freebsd_amd64 xnpc
    tar -czvf xnpc_freebsd_amd64.tar.gz xnpc conf/npc.conf
  fi

  if [ -f "xnpc_freebsd_arm" ]; then
    mv xnpc_freebsd_arm xnpc
    tar -czvf xnpc_freebsd_arm.tar.gz xnpc conf/npc.conf
  fi

  if [ -f "xnpc_linux_arm_v7" ]; then
    mv xnpc_linux_arm_v7 xnpc
    tar -czvf xnpc_linux_arm_v7.tar.gz xnpc conf/npc.conf
  fi

  if [ -f "xnpc_linux_arm_v6" ]; then
    mv xnpc_linux_arm_v6 xnpc
    tar -czvf xnpc_linux_arm_v6.tar.gz xnpc conf/npc.conf
  fi

  if [ -f "xnpc_linux_arm_v5" ]; then
    mv xnpc_linux_arm_v5 xnpc
    tar -czvf xnpc_linux_arm_v5.tar.gz xnpc conf/npc.conf
  fi

  if [ -f "xnpc_linux_arm64" ]; then
    mv xnpc_linux_arm64 xnpc
    tar -czvf xnpc_linux_arm64.tar.gz xnpc conf/npc.conf
  fi

  if [ -f "xnpc_linux_mips64" ]; then
    mv xnpc_linux_mips64 xnpc
    tar -czvf xnpc_linux_mips64.tar.gz xnpc conf/npc.conf
  fi

  if [ -f "xnpc_linux_mips64le" ]; then
    mv xnpc_linux_mips64le xnpc
    tar -czvf xnpc_linux_mips64le.tar.gz xnpc conf/npc.conf
  fi

  if [ -f "xnpc_linux_mipsle" ]; then
    mv xnpc_linux_mipsle xnpc
    tar -czvf xnpc_linux_mipsle.tar.gz xnpc conf/npc.conf
  fi

  if [ -f "xnpc_linux_mips" ]; then
    mv xnpc_linux_mips xnpc
    tar -czvf xnpc_linux_mips.tar.gz xnpc conf/npc.conf
  fi

  if [ -f "xnpc_windows_386.exe" ]; then
    mv "xnpc_windows_386.exe" "xnpc.exe"
    tar -czvf xnpc_windows_386.tar.gz xnpc.exe conf/npc.conf
  fi

  if [ -f "xnpc_windows_amd64.exe" ]; then
    mv "xnpc_linux_amd64.exe" "xnpc.exe"
    tar -czvf xnpc_windows_amd64.tar.gz xnpc.exe conf/npc.conf
  fi

  if [ -f "xnpc_darwin_amd64" ]; then
    mv xnpc_darwin_amd64 xnpc
    tar -czvf xnpc_darwin_amd64.tar.gz xnpc conf/npc.conf
  fi

  if [ -f "xnpc_darwin_arm64" ]; then
    mv xnpc_darwin_arm64 xnpc
    tar -czvf xnpc_darwin_arm64.tar.gz xnpc conf/npc.conf
  fi
  rm xnpc "xnpc.exe"
} &

wait
echo "build succeed! The package in 'bin/*'"

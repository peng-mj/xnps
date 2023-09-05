#/bash/sh
cd $(dirname $0)
export VERSION=0.26.16
export GOPROXY=direct

sudo apt-get update
sudo apt-get install gcc-mingw-w64-i686 gcc-multilib
#env GOOS=windows GOARCH=386 CGO_ENABLED=1 CC=i686-w64-mingw32-gcc go build -ldflags "-s -w -extldflags -static -extldflags -static" -buildmode=c-shared -o bin/npc_sdk.dll cmd/npc/sdk.go
#env GOOS=linux GOARCH=386 CGO_ENABLED=1 CC=gcc go build -ldflags "-s -w -extldflags -static -extldflags -static" -buildmode=c-shared -o bin/npc_sdk.so cmd/npc/sdk.go
#tar -czvf bin/npc_sdk.tar.gz npc_sdk.dll npc_sdk.so npc_sdk.h


###################################-NPC-#########################################

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -extldflags -static -extldflags -static"  ./cmd/npc/npc.go
tar -czvf bin/linux_amd64_client.tar.gz npc conf/npc.conf conf/multi_account.conf


CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -ldflags "-s -w -extldflags -static -extldflags -static" ./cmd/npc/npc.go
tar -czvf bin/linux_386_client.tar.gz npc conf/npc.conf conf/multi_account.conf


CGO_ENABLED=0 GOOS=freebsd GOARCH=386 go build -ldflags "-s -w -extldflags -static -extldflags -static" ./cmd/npc/npc.go
tar -czvf bin/freebsd_386_client.tar.gz npc conf/npc.conf conf/multi_account.conf


CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 go build -ldflags "-s -w -extldflags -static -extldflags -static" ./cmd/npc/npc.go
tar -czvf bin/freebsd_amd64_client.tar.gz npc conf/npc.conf conf/multi_account.conf


CGO_ENABLED=0 GOOS=freebsd GOARCH=arm go build -ldflags "-s -w -extldflags -static -extldflags -static" ./cmd/npc/npc.go
tar -czvf bin/freebsd_arm_client.tar.gz npc conf/npc.conf conf/multi_account.conf


CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -ldflags "-s -w -extldflags -static -extldflags -static" ./cmd/npc/npc.go
tar -czvf bin/linux_arm_v7_client.tar.gz npc conf/npc.conf conf/multi_account.conf


CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -ldflags "-s -w -extldflags -static -extldflags -static" ./cmd/npc/npc.go
tar -czvf bin/linux_arm_v6_client.tar.gz npc conf/npc.conf conf/multi_account.conf


CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=5 go build -ldflags "-s -w -extldflags -static -extldflags -static" ./cmd/npc/npc.go
tar -czvf bin/linux_arm_v5_client.tar.gz npc conf/npc.conf conf/multi_account.conf


CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "-s -w -extldflags -static -extldflags -static" ./cmd/npc/npc.go
tar -czvf bin/linux_arm64_client.tar.gz npc conf/npc.conf conf/multi_account.conf


CGO_ENABLED=0 GOOS=linux GOARCH=mips64 go build -ldflags "-s -w -extldflags -static -extldflags -static" ./cmd/npc/npc.go
tar -czvf bin/linux_mips64_client.tar.gz npc conf/npc.conf conf/multi_account.conf


CGO_ENABLED=0 GOOS=linux GOARCH=mips64le go build -ldflags "-s -w -extldflags -static -extldflags -static" ./cmd/npc/npc.go
tar -czvf bin/linux_mips64le_client.tar.gz npc conf/npc.conf conf/multi_account.conf


CGO_ENABLED=0 GOOS=linux GOARCH=mipsle go build -ldflags "-s -w -extldflags -static -extldflags -static" ./cmd/npc/npc.go
tar -czvf bin/linux_mipsle_client.tar.gz npc conf/npc.conf conf/multi_account.conf


CGO_ENABLED=0 GOOS=linux GOARCH=mips go build -ldflags "-s -w -extldflags -static -extldflags -static" ./cmd/npc/npc.go
tar -czvf bin/linux_mips_client.tar.gz npc conf/npc.conf conf/multi_account.conf


CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -ldflags "-s -w -extldflags -static -extldflags -static" ./cmd/npc/npc.go
tar -czvf bin/windows_386_client.tar.gz npc.exe conf/npc.conf conf/multi_account.conf


CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w -extldflags -static -extldflags -static" ./cmd/npc/npc.go
tar -czvf bin/windows_amd64_client.tar.gz npc.exe conf/npc.conf conf/multi_account.conf


CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w -extldflags -static -extldflags -static" ./cmd/npc/npc.go
tar -czvf bin/darwin_amd64_client.tar.gz npc conf/npc.conf conf/multi_account.conf


CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w -extldflags -static -extldflags -static" ./cmd/npc/npc.go
tar -czvf bin/darwin_arm64_client.tar.gz npc conf/npc.conf conf/multi_account.conf


###################################-NPS-#########################################

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -extldflags -static -extldflags -static" ./cmd/nps/nps.go
tar -czvf bin/linux_amd64_server.tar.gz conf/nps.conf conf/tasks.json conf/clients.json conf/hosts.json conf/server.key  conf/server.pem web/views web/static nps

CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -ldflags "-s -w -extldflags -static -extldflags -static" ./cmd/nps/nps.go
tar -czvf bin/linux_386_server.tar.gz conf/nps.conf conf/tasks.json conf/clients.json conf/hosts.json conf/server.key  conf/server.pem web/views web/static nps

CGO_ENABLED=0 GOOS=freebsd GOARCH=386 go build -ldflags "-s -w -extldflags -static -extldflags -static" ./cmd/nps/nps.go
tar -czvf bin/freebsd_386_server.tar.gz conf/nps.conf conf/tasks.json conf/clients.json conf/hosts.json conf/server.key  conf/server.pem web/views web/static nps

CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 go build -ldflags "-s -w -extldflags -static -extldflags -static" ./cmd/nps/nps.go
tar -czvf bin/freebsd_amd64_server.tar.gz conf/nps.conf conf/tasks.json conf/clients.json conf/hosts.json conf/server.key  conf/server.pem web/views web/static nps

CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w -extldflags -static -extldflags -static" ./cmd/nps/nps.go
tar -czvf bin/darwin_amd64_server.tar.gz conf/nps.conf conf/tasks.json conf/clients.json conf/hosts.json conf/server.key  conf/server.pem web/views web/static nps

CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w -extldflags -static -extldflags -static" ./cmd/nps/nps.go
tar -czvf bin/darwin_arm64_server.tar.gz conf/nps.conf conf/tasks.json conf/clients.json conf/hosts.json conf/server.key  conf/server.pem web/views web/static nps

CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w -extldflags -static -extldflags -static" ./cmd/nps/nps.go
tar -czvf bin/windows_amd64_server.tar.gz conf/nps.conf conf/tasks.json conf/clients.json conf/hosts.json conf/server.key  conf/server.pem web/views web/static nps.exe

CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -ldflags "-s -w -extldflags -static -extldflags -static" ./cmd/nps/nps.go
tar -czvf bin/windows_386_server.tar.gz conf/nps.conf conf/tasks.json conf/clients.json conf/hosts.json conf/server.key  conf/server.pem web/views web/static nps.exe

rm -r nps nps.exe npc npc.exe

#!/bin/bash
unamestr=`uname`
SHA256='shasum -a 256'

VERSION=`date -u +%Y%m%d`
LDFLAGS="-X main.VERSION=$VERSION -s -w"
GCFLAGS=""

OSES=(linux darwin)
ARCHS=(amd64 386)
for os in ${OSES[@]}; do
	for arch in ${ARCHS[@]}; do
		suffix=""
        cgo_enabled=0
        if [ $os == "darwin" ]; then
            cgo_enabled=1
        fi
        env CGO_ENABLED=$cgo_enabled GOOS=$os GOARCH=$arch go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o tunnel_${os}_${arch}${suffix} github.com/ccsexyz/tunnel
		env CGO_ENABLED=$cgo_enabled GOOS=$os GOARCH=$arch go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -tags goprof -o tunnel_${os}_${arch}${suffix}_pprof github.com/ccsexyz/tunnel
		tar -zcf tunnel-${os}-${arch}-$VERSION.tar.gz tunnel_${os}_${arch}${suffix} tunnel_${os}_${arch}${suffix}_pprof
		$SHA256 tunnel-${os}-${arch}-$VERSION.tar.gz
	done
done

# ARM
ARMS=(5 6 7)
for v in ${ARMS[@]}; do
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=$v go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o tunnel_linux_arm${v}  github.com/ccsexyz/tunnel
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=$v go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -tags goprof -o tunnel_linux_arm${v}_pprof  github.com/ccsexyz/tunnel
done
env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o tunnel_linux_arm64  github.com/ccsexyz/tunnel
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -tags goprof -o tunnel_linux_arm64_pprof  github.com/ccsexyz/tunnel
tar -zcf tunnel-linux-arm-$VERSION.tar.gz tunnel_linux_arm* 
$SHA256 tunnel-linux-arm-$VERSION.tar.gz

#MIPS32LE
env CGO_ENABLED=0 GOOS=linux GOARCH=mipsle go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o tunnel_linux_mipsle github.com/ccsexyz/tunnel
env CGO_ENABLED=0 GOOS=linux GOARCH=mips go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o tunnel_linux_mips github.com/ccsexyz/tunnel
env CGO_ENABLED=0 GOOS=linux GOARCH=mipsle go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -tags goprof -o tunnel_linux_mipsle_pprof github.com/ccsexyz/tunnel
env CGO_ENABLED=0 GOOS=linux GOARCH=mips go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -tags goprof -o tunnel_linux_mips_pprof github.com/ccsexyz/tunnel

tar -zcf tunnel-linux-mipsle-$VERSION.tar.gz tunnel_linux_mipsle tunnel_linux_mipsle_pprof
tar -zcf tunnel-linux-mips-$VERSION.tar.gz tunnel_linux_mips tunnel_linux_mips_pprof
$SHA256 tunnel-linux-mipsle-$VERSION.tar.gz
$SHA256 tunnel-linux-mips-$VERSION.tar.gz

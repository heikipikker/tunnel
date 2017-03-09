#!/bin/bash
MD5='md5sum'
unamestr=`uname`
if [[ "$unamestr" == 'Darwin' ]]; then
	MD5='md5'
fi

UPX=false
if hash upx 2>/dev/null; then
	UPX=true
fi

VERSION=`date -u +%Y%m%d`
LDFLAGS="-X main.VERSION=$VERSION -s -w"
GCFLAGS=""

OSES=(linux darwin windows freebsd)
ARCHS=(amd64 386)
for os in ${OSES[@]}; do
	for arch in ${ARCHS[@]}; do
		suffix=""
        cgo_enabled=0
        env CGO_ENABLED=$cgo_enabled GOOS=$os GOARCH=$arch go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o tunnel_${os}_${arch}${suffix} github.com/ccsexyz/tunnel
		if $UPX; then upx -9 tunnel_${os}_${arch}${suffix} ;fi
		tar -zcf tunnel-${os}-${arch}-$VERSION.tar.gz tunnel_${os}_${arch}${suffix}
		$MD5 tunnel-${os}-${arch}-$VERSION.tar.gz
	done
done

# ARM
ARMS=(5 6 7)
for v in ${ARMS[@]}; do
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=$v go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o tunnel_linux_arm$v  github.com/ccsexyz/tunnel
done
if $UPX; then upx -9 tunnel_linux_arm*;fi
tar -zcf tunnel-linux-arm-$VERSION.tar.gz tunnel_linux_arm* 
$MD5 tunnel-linux-arm-$VERSION.tar.gz

#MIPS32LE
env CGO_ENABLED=0 GOOS=linux GOARCH=mipsle go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o tunnel_linux_mipsle github.com/ccsexyz/tunnel
env CGO_ENABLED=0 GOOS=linux GOARCH=mips go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o tunnel_linux_mips github.com/ccsexyz/tunnel

if $UPX; then upx -9 tunnel_linux_mips* server_linux_mips*;fi
tar -zcf tunnel-linux-mipsle-$VERSION.tar.gz tunnel_linux_mipsle
tar -zcf tunnel-linux-mips-$VERSION.tar.gz tunnel_linux_mips
$MD5 tunnel-linux-mipsle-$VERSION.tar.gz
$MD5 tunnel-linux-mips-$VERSION.tar.gz

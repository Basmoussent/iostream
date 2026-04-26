#!/usr/bin/env bash
# local-windows-ci.sh — emulate the Windows CI job locally via Docker + mingw.
#
# Runs in a one-shot golang:1.23 container, installs mingw-w64, downloads the
# libusb-1.0 Windows release, and cross-compiles iostream.exe (and tries the
# GUI). Outputs land in ./bin alongside libusb-1.0.dll. Faster than waiting
# for GitHub Actions.

set -euo pipefail

LIBUSB_VER="${LIBUSB_VER:-1.0.27}"
LIBUSB_URL="https://github.com/libusb/libusb/releases/download/v${LIBUSB_VER}/libusb-${LIBUSB_VER}.7z"

run_in_container() {
  cat <<'INNER' >/tmp/iostream-inner.sh
#!/usr/bin/env bash
set -euo pipefail

echo "==> apt: mingw-w64, p7zip-full, pkg-config"
apt-get update -qq
apt-get install -y --no-install-recommends \
  mingw-w64 p7zip-full pkg-config wget ca-certificates >/dev/null

echo "==> downloading libusb ${LIBUSB_VER}"
mkdir -p /opt/libusb && cd /opt/libusb
wget -q "${LIBUSB_URL}" -O libusb.7z
7z x -y libusb.7z >/dev/null

# libusb release ships a MinGW-compatible build under MinGW64/. Headers live
# at include/libusb.h (no libusb-1.0/ subdir) but pkg-config + gousb expect
# them at include/libusb-1.0/libusb.h, so we relocate.
SYSROOT=/usr/x86_64-w64-mingw32
cp -r MinGW64/static/* "${SYSROOT}/lib/" 2>/dev/null || true
cp -r MinGW64/dll/*    "${SYSROOT}/lib/"
mkdir -p "${SYSROOT}/include/libusb-1.0"
cp include/libusb.h    "${SYSROOT}/include/libusb-1.0/"

cat >/tmp/libusb-1.0.pc <<EOF
prefix=${SYSROOT}
exec_prefix=\${prefix}
libdir=\${prefix}/lib
includedir=\${prefix}/include

Name: libusb-1.0
Version: ${LIBUSB_VER}
Description: C API for USB device access
Libs: -L\${libdir} -lusb-1.0
Cflags: -I\${includedir}/libusb-1.0
EOF
mkdir -p "${SYSROOT}/lib/pkgconfig"
mv /tmp/libusb-1.0.pc "${SYSROOT}/lib/pkgconfig/"

echo "==> cross-compiling iostream.exe"
cd /src
mkdir -p bin
export CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc \
  PKG_CONFIG_PATH=${SYSROOT}/lib/pkgconfig
go build -buildvcs=false -o bin/iostream.exe ./cmd/iostream

echo "==> cross-compiling iostream-gui.exe"
# Webview cgo headers may not survive mingw; this is best-effort.
if go build -buildvcs=false -ldflags "-H windowsgui" -o bin/iostream-gui.exe ./cmd/iostream-gui 2>/tmp/gui.err; then
  echo "    GUI build: OK"
else
  echo "    GUI build: FAILED (mingw + WebView2 headers don't always cooperate)"
  echo "    See /tmp/gui.err — the GitHub Actions Windows job uses MSVC, which is the supported path."
  tail -5 /tmp/gui.err
fi

echo "==> bundling libusb-1.0.dll"
cp ${SYSROOT}/bin/libusb-1.0.dll bin/ 2>/dev/null \
  || cp /opt/libusb/MinGW64/dll/libusb-1.0.dll bin/

echo "==> done. bin/:"
ls -la bin/
INNER
  chmod +x /tmp/iostream-inner.sh

  docker run --rm \
    -e LIBUSB_VER="${LIBUSB_VER}" \
    -e LIBUSB_URL="${LIBUSB_URL}" \
    -v "$(pwd)":/src \
    -v /tmp/iostream-inner.sh:/tmp/iostream-inner.sh:ro \
    -w /src \
    golang:1.23 \
    bash /tmp/iostream-inner.sh
}

run_in_container

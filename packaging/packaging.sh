#!/usr/bin/env bash
PKG_VERSION=$1
BIN_NAME=ns
TDIR_NAME=notification-service
SETTINGS=config.toml

if [ "$PKG_VERSION" = "" ]; then
    echo "Ok"
    PKG_VERSION=manual
fi

VER=`cat VERSION`
PKG_VERSION="${VER}_${PKG_VERSION}"

# make build
echo $PKG_VERSION
mkdir -p pkg/opt/${TDIR_NAME}/log pkg/etc/logrotate.d pkg/usr/share/${TDIR_NAME}/initial_data
cp packaging/notification-service.logrotate pkg/etc/logrotate.d/notification-service
cp -r $BIN_NAME pkg/opt/${TDIR_NAME}
# cp -r fixtures/* pkg/usr/share/${TDIR_NAME}/initial_data
cp ${SETTINGS}.orig pkg/opt/${TDIR_NAME}/${SETTINGS}
rm *.deb
fpm -s dir -t deb --name notification-service --version 1.$1 -C pkg/ -m "Sergey Alipchenkov" --config-files /opt/${TDIR_NAME}/${SETTINGS} --before-install packaging/before-install.sh --before-upgrade packaging/before-upgrade.sh --before-remove packaging/before-remove.sh --after-install packaging/after-install.sh --after-remove packaging/after-remove.sh --after-upgrade packaging/after-upgrade.sh --deb-init packaging/init/deb/notification-service
rm -rf pkg

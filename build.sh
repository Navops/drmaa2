#!/bin/sh

# You need to source the settings.sh file (source /path/to/UGE/default/common/settings.sh)
# of your Univa Grid Engine installation before building. "default" is the CELL name, which can
# be different in your setup.

ARCH=`$SGE_ROOT/util/arch`
export CGO_LDFLAGS="-L$SGE_ROOT/lib/$ARCH/"
export CGO_CFLAGS="-I$SGE_ROOT/include"

go build -a
go install 


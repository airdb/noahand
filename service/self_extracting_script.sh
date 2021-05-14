#!/bin/bash
echo ""
echo "Self Extracting Installer"
echo ""
 
export TMPDIR=`mktemp -d /tmp/selfextract.XXXXXX`
 
ARCHIVE=`awk '/^__ARCHIVE_BELOW__/ {print NR + 1; exit 0; }' $0`
 
tail -n+$ARCHIVE $0 | tar xzv -C $TMPDIR
 
CDIR=`pwd`
cd $TMPDIR

# 改为压缩包中安装程序的地址
# sh ./install.sh   
sh service/install.sh
 
cd $CDIR
rm -rf $TMPDIR
 
exit 0
 
__ARCHIVE_BELOW__

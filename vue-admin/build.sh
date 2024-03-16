#!/bin/bash
time=$(date "+%Y%m%d%H%M")
tarName="${time}_dist.tar.gz" 
echo "${time}"

# 构建
npm run build
# 打包
tar -czvf ${tarName} dist/
#mv ${tarName} ./release
# 上传
#scp ./release/${tarName} root@10.100.214.8:/opt/magneto/web/release/

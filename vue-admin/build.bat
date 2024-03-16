

time=%time%
#tarName="${time}_dist.tar.gz" 
echo "$time"

# 构建
npm run build
# 打包
tar -czvf ${tarName} dist/

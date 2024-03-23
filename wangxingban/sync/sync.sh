cd /root/asm-deploy/lev
rm -rf tmp/
mkdir tmp

cp -r ./mongo_data ./tmp/mongo_data_tmp
cd tmp
zip -r mongo_data.zip ./mongo_data_tmp

cp -r ./mysql ./tmp/mysql_tmp
cd tmp
zip -r mysql_tmp.zip ./mysql_tmp
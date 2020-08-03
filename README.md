# mediawiki的安装

## Step1

解压installer.tar.gz，执行`./install`

## Step2

根据命令行提示向下安装， 直到显示连接到网页

## Step3

注意`statics/script`下的docker-compose文件，其中v1是为redhat6.7或centos6.7准备的，有如下几行需要重点注意
```yaml
MYSQL_DATABASE: my_wiki
MYSQL_USER: wikiuser
MYSQL_PASSWORD: example
```
分别表示wiki所使用的数据库名称， 数据库用户，数据库密码, 第一次安装不熟悉不建议修改，您将在进入网页后的第一个页面填入数据库的ip地址即您安装wiki的主机的地址， 然后按提示分别填入上述参数

## Step4

接下来的页面是wiki名称等，可以自己随意填写，勾选所有选择框

## Step5

进入wiki功能选择页面选择所有的选择框

## Step6

完成界面

## Step7

LocalSettings.php下载界面，您需要下载LocalSettings.php到您的安装目录下有个叫mediawikiRoot的目录下
修改docker-compose文件
```yaml
#    - ./LocalSettings.php:/var/www/html/LocalSettings.php
```
为
```yaml
    - ./LocalSettings.php:/var/www/html/LocalSettings.php
```
注意和上一行缩进一致，回到mediawikiRoot目录执行以下命令
```bash
docker-compose stop
docker-compose up -d
```


## 备注

1. 因centos6.7或redhat6.7的一些原因，仅仅能使用`docker-1.7.1 docker-compose-1.5.1`均是较老版本，而在docker早期两个容器在连接的时候会出现无法链接的bug，所以间接的导致mediawiki主程序没有办法和数据库进行通信，所以解决方案仅有将主机的80端口和3306进行占用， 这个bug在docker的后期版本中进行了修订，但是由于系统原因并不兼容高版本，所以在centos6.7版本和redhat6.7版本有且仅有这一种方式进行安装。具体问题参见`https://github.com/docker/compose/issues/2172`
2. 由于上述原因安装mediawiki的主机不能有apache或者nginx等http服务软件占用80端口， 以及不能有数据库如mysql和maridb等数据库软件占用3306端口
3. 由于一些软件较老的原因直接在主机上安装mediawiki需要大量的兼容环境，编译一系列软件，可能会间接导致系统软件兼容问题，故现在不予考虑。

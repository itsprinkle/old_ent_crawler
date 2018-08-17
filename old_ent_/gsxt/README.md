企业信息
=======


### 安装

1. 安装golang，设置`~/dev`为GOPATH

2. 安装外部依赖库
    ```
        go get -u github.com/jessevdk/go-flags
        go get -u github.com/PuerkitoBio/goquery
        go get -u github.com/levigross/grequests
    ```
3. 下载本地库

> 其中geetest下载到 GOPATH/src/gxst/geetest

> 其中link下载到 GOPATH/src/tools/link

> 其中credit下载到 GOPATH/src/gxst/credit

> 其中gxst下载到 GOPATH/src/gxst/gxst

3. 参考geetest,link两个二进制库的README，移动文件到对应的位置


4. 进去gxst，并且make


### 功能介绍

1. find 主要用来测试查询单个keyword

    find guangdong 前海翼联


2. v1 用来将当前的输出格式兼容老的版本

    v1 v2txt.txt v1txt.txt

3. grab

    dispatcher -h 抓取服务端，分配任务给worker，并且写入到本地文件，需要新建out/code, out/detail两个目录
    master -h http服务端，通过http请求分发任务给worker，然后再将内容返回给访问的人
    worker -h 真正干活抓取的

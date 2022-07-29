# 实现逻辑

```
目的：采集http上报的监控数据，并metrics出去。
遵照prometheus和opencensus的规范，扩展模块，使用prometheus自定义metrics并注册的方式，双port模式，一个port使用自定义url采集数据，一个prot使用metrics的url暴露采集信息。


```

### 兼容主要golang框架:

```
gin
fasthttp
gorouter

```

### 允许接驳监控数据规范 CNCF

```
prometheus

opencensus
```

### 代码章节

```
prometheus middle-ware for five kinds go-web-frame
mock test
metrics export
makefile golangci 
example
go-doc
单元测试
```

### SIGS
```
1、适配"数百万级别的QPS"
2、适配时序数据库 'TSDB' \ 'influxdb'
3、直接grafana展示qps及本程序运行开销
4、允许prometheus直接remote-read和remote-write
```
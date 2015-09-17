# schedule

### Quick Start
```sh
$ ./schedule -interval 1 -schedule [FILE]
```
```sh
$ ./schedule -h
#  -backup="/tmp/schedule.backup": 排程備份位置
#  -interval=10: 排程掃描間隔
#  -schedule="": 排程設定檔
```

### 排程設定檔格式

```
# 秒 分 時 日 月|次數|命令

# 每秒執行
* * * * *|0|echo hello world $(date) >> /tmp/schedule.test

# 每分鐘內的第5, 10, 20, 40秒...
5,10,20,40 * * * *|0|echo hello world $(date) >> /tmp/schedule.test
```

### 次數規則

- >= 1 限定次數
- == 0 不限次數
- == -1 不限次數,但可重疊執行 


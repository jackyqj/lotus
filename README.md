# lotus 操作
## 从log中查看封装时间
### PC1
PC1 可以用关键字 `seal > seal_pre_commit_phase1` 到日志里查询，然后用格式化输出
`grep "seal > seal_pre_commit_phase1" worker-3001-20210801-073156.log |awk -F' ' '{print $6 ", "$5 ", " $1}'`
PC1并行输出不大好看，但你可以抄到EXCEL为所欲为^_^

如果查单个Sector的话就简单用ID去查就好
```
grep "SectorId(5242)" worker-3001-20210801-073156.log
2021-08-12T07:58:04.323 INFO filecoin_proofs::api::seal > seal_pre_commit_phase1:start: SectorId(5242)
2021-08-12T11:23:08.407 INFO filecoin_proofs::api::seal > seal_pre_commit_phase1:finish: SectorId(5242)
```
### PC2
和PC1一样，只是关键字换成`seal > seal_pre_commit_phase2`
不过PC2的输出则看不到SECTOR ID


## 查看爆块信息
用下面命令可以在MINER日志查看爆块信息, 假设你的log文件为lotus-miner.log.
`grep "\"isWinner\": true" lotus-miner.log -B 1`
```
lotus-miner.log-2021-08-10T07:24:11.216+0800	INFO	miner	miner/miner.go:565	mined new block	{"cid": "bafy2bzacec7tnd34djdobimxxxxxxeunnw7cxwimatr7lfal6ooxqp4kw", "height": 1008169, "miner": "f010xxxxx", "parents": ["f0695014","f0154294","f0461791"], "parentTipset": "{bafy2bzacebr3yvk33nmwbf5vbb5yfiv4kslt63gvovzrzmuin53qioffdsg3c,bafy2bzacebpk6us6efd23tw6xtrdf5pakybshnvt6trdeoyevejjtmte6ifnw,bafy2bzaceb3n7qxher5p7zdxit2rn4r6hcuxzitai325fwpe6dfxhvzeh4gyu}", "took": 5.127382263}
lotus-miner.log:2021-08-10T07:24:11.216+0800	INFO	miner	miner/miner.go:475	completed mineOne	{"tookMilliseconds": 5127, "forRound": 1008169, "baseEpoch": 1008168, "baseDeltaSeconds": 6, "nullRounds": 0, "lateStart": false, "beaconEpoch": 1104013, "lookbackEpochs": 900, "networkPowerAtLookback": "10172966243152592896", "minerPowerAtLookback": "148330990534656", "isEligible": true, "isWinner": true, "error": null}
```
以上为成功，如果出现WARN则为爆块失败。

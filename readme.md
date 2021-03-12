# MultiVAC monitor

## 1. Function
1. Query whether there is an error.. log from the log folder of MultiVAC. If there is, judge whether there is an update. If the content is updated, prepare to send an email warning.
2. Take data from MultiVAC to analyze whether the block information of a miner in a certain segment can form a chain.
3. For a shard, check whether all miners in the shard are at the same height and whether the block information is the same.
4. Check whether a miner doesn't get out of the block and whether it is offline.

## 2. Usage

Need to use with MultiVAC, this program runs on the monitoring server.

   `go run main.go` 

Before using, you need to check whether the configuration information in config/config and adjust it to the corresponding configuration.


# Journal-based Blockchain

- [Running](#Running)
- [Configuration](#Configuration)
- [Simulation](#Simulation)



## Running

- **Prerequisite**: Golang

1. Set configuration file (`conf.json`)
2. `$ ./full`     // Initialize full node
3. `$ ./light`   // Synchronize light node with full node



## Configuration

| Name           | Description                        | Default   |
| -------------- | ---------------------------------- | --------- |
| Hostname       | [string] Host name                 | localhost |
| Port           | [string] Port number               | 8081      |
| BlockNumber    | [int64] The number of blocks       | 100       |
| Participants   | [int64] The number of participants | 100       |
| PrintMode      | [bool] Print blocks on console     | true      |
| MiningInterval | [int]  Mining Interval (sec)       | 0 sec     |



## Simulation

#### Simulation depending on the block number and participants number

1. Set `BlockNumber`, `Participants` in `network/conf.json`

2. Remove old DB directory, if it exists (or rename it)

   `$ rm -rf chaindata*`

3. Run full node

   `$ ./iot_full.go`

4. Synchronize light node with full node

   `$ ./iot_light.go`



#### Simulation depending on the network bandwidth and delay

- See `network/README.md`

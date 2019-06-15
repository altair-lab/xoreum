[reference] https://gist.github.com/ole1986/d9d6be5218affd41796610a35e3b069c

#### traffic-control.sh 툴
- 리눅스 시스템의 네트워크 환경을 통제 -> IoT 환경의 네트워크 상황을 시뮬레이션에 사용
- bandwidth (e.g. 2Mbps) 및 delay (=latency=RTT/2) (e.g. 10ms)를 조절할 수 있음

`sudo traffic-control.sh -o --uspeed=[BANDWIDTH] --delay=[DELAY] [TARGET_IP]`

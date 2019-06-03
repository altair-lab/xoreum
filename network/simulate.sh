for speed in '50Mbps 20Mbps 10Mbps 1Mbps 512kbps 256kbps 128kbps 56kbps'
do
	for delay in '100 200 300 400 500 600 700 800 900 1000'
	do
		sudo ./traffic-control.sh --dspeed $speed --uspeed $speed --delay $delay 147.46.123.249
		rm -rf chaindata-iot
		go run light/iot_light.go
	done
done


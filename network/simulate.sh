for speed in 50Mbps 20Mbps 10Mbps 1Mbps 512kbps 256kbps 128kbps 56kbps
do
	for delay in 1 2 4 8 16 32 64 128
	do
		delay = 500
#		echo $speed $delay
		sudo ./traffic-control.sh -o --uspeed=$speed --delay=$delay 147.46.123.249
		sudo ./traffic-control.sh -i --dspeed=$speed --delay=$delay 147.46.123.249
		rm -rf chaindata-iot
		go run light/iot_light.go
		
		exit 
	done
done


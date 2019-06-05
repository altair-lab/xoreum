for speed in 50Mbps 20Mbps 10Mbps 8Mbps 6Mbps 4Mbps 2Mbps 1Mbps  
do
	for delay in 1 2 4 8
	do
		if [ ! -f "$speed"_"$delay".txt ]; then
  			echo "$speed"_"$delay".txt
		sshpass -p ma55lab ssh -t -o StrictHostKeyChecking=no yyh@lynx.snu.ac.kr -p 2825 "echo ma55lab | sudo -S go/pkg/src/github.com/altair-lab/xoreum/network/traffic-control.sh -o --uspeed=$speed --delay=$delay 147.46.123.135"
		sleep 1
		sudo ./traffic-control.sh -o --uspeed=$speed --delay=$delay 147.46.123.249
		for iter in {1..5}
		do
			echo $iter
			rm -rf chaindata-iot
			timeout 60 go run light/iot_light.go>>"$speed"_"$delay".txt
			echo "">>"$speed"_"$delay".txt
			sleep 1
		done
		sleep 3
	
		fi
	 
	done
done



import os

# range(start block #, end block # + 1, interval)
for x in range(500, 10001, 500):
    #Cmd = "echo ------" + str(x) + "------- >> sizelog"
    #os.system(Cmd)
    for i in range(0, 3):
        # Remove old db directory
        Cmd = "rm -rf chaindata_" + str(x)
        os.system(Cmd)
    
        # Run the full node and initialize the chain
        Cmd = "go run iot_full_simulation.go 1000 " + str(x)
        os.system(Cmd)

        # Get db size
        Cmd = "go run iot_full_simulation.go 1000 " + str(x)
        os.system(Cmd)
        Cmd = "du -sc chaindata_" + str(x) + "/*.ldb | tail -n 1 | awk '{ print $1 }' >> sizelog"
        os.system(Cmd)

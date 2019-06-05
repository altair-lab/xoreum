import os

# range(start block #, end block # + 1, interval)
Cmd = "echo 1000 20000 1000 >> sizelog"
os.system(Cmd)

for x in range(1000, 20001, 1000):
    # Remove old db directory
    Cmd = "rm -rf chaindata_" + str(x)
    os.system(Cmd)
    
    # Run the full node and initialize the chain
    Cmd = "go run iot_full_simulation.go " + str(x)
    os.system(Cmd)

    # Get db size
    Cmd = "go run iot_full_simulation.go " + str(x)
    os.system(Cmd)
    Cmd = "du -sc chaindata_" + str(x) + "/*.ldb | tail -n 1 | awk '{ print $1 }' >> sizelog"
    os.system(Cmd)

    # Remove db directory
    # Cmd = "rm -rf chaindata_" + str(x)
    # os.system(Cmd)

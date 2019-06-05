import os

# range(start block #, end block # + 1, interval)
for x in range(100, 1001, 100):
    # Remove old db directory
    Cmd = "rm -rf chaindata_" + str(x)
    os.System(Cmd)
    
    # Run the full node and initialize the chain
    Cmd = "go run iot_full_simulation.go " + str(x)
    os.system(Cmd)

    # Get db size
    Cmd = "go run iot_full_simulation.go " + str(x)
    os.system(Cmd)
    Cmd = "du -sc chaindata_" + str(x) + "/*.ldb"
    os.system(Cmd)

    # Remove db directory
    Cmd = "rm -rf chaindata_" + str(x)
    os.System(Cmd)

import socket,os

sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)  
sock.bind(("lynx.snu.ac.kr", 8086))  
sock.listen(5)

while True:  
    # Connection is done after the blockchain is constructed
    connection, address = sock.accept()
    
    # Remove old db directory
    Cmd = "rm -rf chaindata-iot"
    os.system(Cmd)
    
    # Synchronize with full node
    Cmd = "go run iot_light.go"
    os.system(Cmd)
    
    # Get size from db
    Cmd = "go run iot_light.go"
    os.system(Cmd)
    Cmd = "du -sc chaindata-iot/*.ldb | tail -n 1 | awk '{ print $1 }' >> sizelog"
    os.system(Cmd)

    # Remove db directory
    Cmd = "rm -rf chaindata-iot"
    os.system(Cmd)

    connection.close()

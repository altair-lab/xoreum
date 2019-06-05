import socket,os

sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)  
sock.bind(("lynx.snu.ac.kr", 8086))  
sock.listen(5)

while True:  
    # Connection is done after the blockchain is constructed
    connection, address = sock.accept()
    
    # Synchronize with full node
    Cmd = "go run iot_light.go"
    os.system(Cmd)

    # Get size from db
    Cmd = "sh size.sh"
    os.system(Cmd)

    connection.close()

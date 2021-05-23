from websocket import create_connection

def Init():
  ws = create_connection('ws://localhost:8080/ws')
  print('Attempting to Connect to GoLang server on port 8080')
  return ws

def Send(ws, mes):
  ws.send(mes)
  print("Sent")
  # print("Receiving...")
  # result =  ws.recv()
  # print("Received '%s'" % result)
  # ws.close()
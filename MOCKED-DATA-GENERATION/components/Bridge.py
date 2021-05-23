from components import Titan

def Generate(signal, id_, msid, type_, typeVer, key, tic):
  
  generatedTitan = Titan.Generate(id_, msid, type_, typeVer, key, tic)

  return dict([
    ("bridge", generatedTitan),
    ("signal", signal)
  ])
from components import SensorData
from components import News2
from components import Bridge

def Generate(respiration, sp02, pulse, temperature, bloodPressure, motion, proximity, airNotOxygen,
              tempAggregateScore, pulseAggregateScore, resAggregateScore, bpAggregateScore, sp02AggregateScore1,
              overall, overallLevel, level, target,
              id_, msid, type_, typeVer, key, tic, 
              time, batteryLevel, signal):
  generatedBridge = []
  generatedSensorData = SensorData.Generate(respiration, sp02, pulse, temperature, bloodPressure, motion, proximity, overall, level, target, overallLevel)
  generatedNews2 = News2.Generate(resAggregateScore, sp02AggregateScore1, airNotOxygen, pulseAggregateScore, tempAggregateScore, bpAggregateScore, motion, overall)
  generatedBridge.append(Bridge.Generate(signal, id_, msid, type_, typeVer, key, tic))

  return dict([
    # ("id", id_),
    ("time", time),
    ("sensorData", generatedSensorData),
    ("news2", generatedNews2),
    ("batteryLevel", batteryLevel),
    ("location", generatedBridge)
  ])

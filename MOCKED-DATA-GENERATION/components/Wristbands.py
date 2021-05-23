# random
import random

# logic
from logic import Rand

# components
from components import Titan
from components import Data

def GenerateParameters(preassignedLevel, preassigned_next_levels_2, errorTarget):
  # errorTarget = ''
  # # Check if there is a space in the string
  # if ':' in preassignedLevel:
  #   tempLevelHolder = preassignedLevel.split(':')
  #   preassignedLevel = tempLevelHolder[0]
  #   errorTarget = tempLevelHolder[1]

  global inputted_level
  # Check for empty 2nd level case
  if preassigned_next_levels_2 == '':
    inputted_level = preassignedLevel
    # print('EMPTY 2nd CASE')
  else:
    inputted_level = preassigned_next_levels_2

  signal = random.randint(-70, -30)

  level = ""
  target = ""

  # vital scores
  if not inputted_level == 'unchanged':
    # Blood Pressure
    global bloodPressure
    if inputted_level == 'random' or inputted_level == 'error':
      bloodPressure = random.randint(90, 220)
    elif inputted_level == 'high':
      bloodPressure =  random.randint(60, 90) or random.randint(220, 250)
    elif inputted_level == 'medium':
      bloodPressure =  random.randint(91, 100)
    elif inputted_level == 'low-medium':
      bloodPressure =  random.randint(101, 110)
    elif inputted_level == 'low':
      bloodPressure =  random.randint(111, 219)

    global bpAggregateScore
    if bloodPressure <= 90 or bloodPressure >= 220:
      level = "high"
      target = "bloodPressure"
      bpAggregateScore = 3
    elif bloodPressure >= 91 and bloodPressure <= 100:
      if level == "medium" or level == "low-medium" or level == "low" or level == "":
        level = "medium"
        target = "bloodPressure"
      
      bpAggregateScore = 2
    elif bloodPressure >= 101 and bloodPressure <= 110:
      if level == "low-medium" or level == "low" or level == "":
        level = "low-medium"
        target = "bloodPressure"
    
      bpAggregateScore = 1
    elif bloodPressure >= 111 and bloodPressure <= 219:
      if level == "low" or level == "":
        level = "low"
        target = "bloodPressure"
      
      bpAggregateScore = 0

    # Motion
    global motion
    if inputted_level == 'random' or inputted_level == 'low-medium' or inputted_level == 'medium' or inputted_level == 'error':
      motion = Rand.Bool()
    elif inputted_level == 'high':
      motion =  False
    elif inputted_level == 'low':
      motion =  True

    # Proximity (False will throw an band data error)
    global proximity
    if inputted_level == 'error' and errorTarget == 'proximity':
      proximity = False
    elif inputted_level == 'error' and errorTarget == 'batteryLevel':
      proximity =  Rand.Bool()
    elif inputted_level == 'error' and not errorTarget == '':
      proximity = False
    elif inputted_level == 'random':
      proximity =  Rand.Bool()
    elif inputted_level == 'high' or inputted_level == 'medium' or inputted_level == 'low-medium':
      proximity =  True
    elif inputted_level == 'low':
      proximity =  True

    # Pulse
    global pulse
    if inputted_level == 'random' or inputted_level == 'error':
      pulse = random.randint(40, 131)
    elif inputted_level == 'high':
      pulse =  random.randint(20, 40) or random.randint(131, 150)
    elif inputted_level == 'medium':
      pulse =  random.randint(111, 130)
    elif inputted_level == 'low-medium':
      pulse =  random.randint(41, 50) or random.randint(91, 110)
    elif inputted_level == 'low':
      pulse =  random.randint(51, 90)

    global pulseAggregateScore
    if pulse <= 40 or pulse >= 131:
      level = "high"
      target = "pulse"
      pulseAggregateScore = 3
    elif pulse >= 111 and pulse <= 130:
      if level == "medium" or level == "low-medium" or level == "low":
        level = "medium"
        target = "pulse"
      
      pulseAggregateScore = 2
    elif pulse >= 41 and pulse <= 50 or pulse >= 91 and pulse <= 110:
      if level == "low-medium" or level == "low":
        level = "low-medium"
        target = "pulse"
    
      pulseAggregateScore = 1
    elif pulse >= 51 and pulse <= 90:
      if level == "low":
        level = "low"
        target = "pulse"
      
      pulseAggregateScore = 0

    # Respiration
    global respiration
    if inputted_level == 'random' or inputted_level == 'error':
      respiration = random.randint(8, 25)
    elif inputted_level == 'high':
      respiration =  random.randint(4, 8) or random.randint(25, 29)
    elif inputted_level == 'medium':
      respiration =  random.randint(21, 24)
    elif inputted_level == 'low-medium':
      respiration =  random.randint(9, 11)
    elif inputted_level == 'low':
      respiration =  random.randint(12, 20)

    global resAggregateScore
    if respiration <= 8 or respiration >= 25:
      level = "high"
      target = "respiration"
      resAggregateScore = 3
    elif respiration >= 21 and respiration <= 24:
      if level == "medium" or level == "low-medium" or level == "low":
        level = "medium"
        target = "respiration"
      
      resAggregateScore = 2
    elif respiration >= 9 and respiration <= 11:
      if level == "low-medium" or level == "low":
        level = "low-medium"
        target = "respiration"
    
      resAggregateScore = 1
    elif respiration >= 12 and respiration <= 20:
      if level == "low":
        level = "low"
        target = "respiration"
      
      resAggregateScore = 0

    # Air Nor Oxygen
    airNotOxygen = Rand.Bool()

    # SP02
    global sp02
    if inputted_level == 'random' or inputted_level == 'error':
      sp02 = random.randint(91, 99)
    elif inputted_level == 'high':
      sp02 =  random.randint(80, 91)
    elif inputted_level == 'medium':
      sp02 =  random.randint(92, 93)
    elif inputted_level == 'low-medium':
      sp02 =  random.randint(94, 95)
    elif inputted_level == 'low':
      sp02 =  random.randint(96, 99)

    global sp02AggregateScore1
    if sp02 <= 91:
      level = "high"
      target = "sp02"
      sp02AggregateScore1 = 3
    elif sp02 >= 92 and sp02 <= 93:
      if level == "medium" or level == "low-medium" or level == "low":
        level = "medium"
        target = "sp02"
      
      sp02AggregateScore1 = 2
    elif sp02 >= 94 and sp02 <= 95:
      if level == "low-medium" or level == "low":
        level = "low-medium"
        target = "sp02"
    
      sp02AggregateScore1 = 1
    else:
      if level == "low":
        level = "low"
        target = "sp02"
      
      sp02AggregateScore1 = 0


    # Temperature
    global temperature
    if inputted_level == 'random' or inputted_level == 'error':
      temperature = random.randint(35, 39)
    elif inputted_level == 'high':
      temperature = random.randint(28, 35)
    elif inputted_level == 'medium':
      temperature = round(random.uniform(39.1, 42.0), 1)
    elif inputted_level == 'low-medium':
      temperature = round(random.uniform(35.1, 36.0), 1) or round(random.uniform(38.1, 39.0), 1)
    elif inputted_level == 'low':
      temperature = round(random.uniform(36.1, 38.0), 1)

    global tempAggregateScore
    if temperature <= 35:
      level = "high"
      target = "temperature"
      tempAggregateScore = 3
    elif temperature >= 39.1:
      if level == "medium" or level == "low-medium" or level == "low":
        level = "medium"
        target = "temperature"
      
      tempAggregateScore = 2
    elif temperature >= 35.1 and temperature <= 36 or temperature >= 38.1 and temperature <= 39:
      if level == "low-medium" or level == "low":
        level = "low-medium"
        target = "temperature"
    
      tempAggregateScore = 1
    elif temperature >= 36.1 and temperature <= 38:
      if level == "low":
        level = "low"
        target = "temperature"
      
      tempAggregateScore = 0

    global batteryLevel
    if inputted_level == 'error':
      batteryLevel = random.randint(0, 20)
    elif inputted_level == 'random':
      batteryLevel =  random.randint(0, 100)
    elif inputted_level == 'high' or inputted_level == 'medium' or inputted_level == 'low-medium':
      batteryLevel =  random.randint(92, 93)
    elif inputted_level == 'low':
      batteryLevel =  random.randint(80, 100)

    # THE ONLY REASON WHY MOTION AND PROXIMITY IS NOT ORDERED ALPHABETICALLY WITH THE ABOVE PARAMETERS IS BECAUSE THEY ARE BOOLEAN VARIABLES
    # AND THEY WILL THROW ERROR INSTEAD OF AN ALERT STATE, ERROR STATE WILL BE PRIORITISED BEFORE ALERT
    # SO THEY WILL BE ORDERED ALPHABETICALLY DESCENDING SEPARATELY
    if not motion:
      level = "high"
      target = "motion"
    
    if not proximity:
      level = "error"
      target = "proximity"

    overall = tempAggregateScore + pulseAggregateScore + resAggregateScore + bpAggregateScore + sp02AggregateScore1

    # Overall Level
    global overallLevel
    if overall >= 0 and overall <= 4:
      overallLevel = "low"
    elif level == "high" and overall <= 4:
      overallLevel = "low-medium"
    elif overall >= 5 and overall <= 6:
      overallLevel = "medium"
    elif overall >= 7:
      overallLevel = "high"

    global unchangedParams
    unchangedParams = dict([
      ("respiration", respiration),
      ("sp02", sp02),
      ("temperature", temperature),
      ("bloodPressure", bloodPressure),
      ("motion", motion),
      ("proximity", proximity),
      ("pulse", pulse),
      ("onOxygen", airNotOxygen),
      ("tempAggregateScore", tempAggregateScore),
      ("pulseAggregateScore", pulseAggregateScore),
      ("resAggregateScore", resAggregateScore),
      ("bpAggregateScore", bpAggregateScore),
      ("sp02AggregateScore1", sp02AggregateScore1),
      ("batteryLevel", batteryLevel),
      ("signal", signal),
      ("overall", overall),
      ("overallLevel", overallLevel),
      ("level", level),
      ("target", target),
    ])
  else:
    respiration = unchangedParams["respiration"]
    sp02 = unchangedParams["sp02"]
    temperature = unchangedParams["temperature"]
    bloodPressure = unchangedParams["bloodPressure"]
    motion = unchangedParams["motion"]
    proximity = unchangedParams["proximity"]
    pulse = unchangedParams["pulse"]
    airNotOxygen = unchangedParams["onOxygen"]
    tempAggregateScore = unchangedParams["tempAggregateScore"]
    pulseAggregateScore = unchangedParams["pulseAggregateScore"]
    resAggregateScore = unchangedParams["resAggregateScore"]
    bpAggregateScore = unchangedParams["bpAggregateScore"]
    sp02AggregateScore1 = unchangedParams["sp02AggregateScore1"]
    batteryLevel = unchangedParams["batteryLevel"]
    signal = unchangedParams["signal"]
    overall = unchangedParams["overall"]
    overallLevel = unchangedParams["overallLevel"]
    level = unchangedParams["level"]
    target = unchangedParams["target"]

  return respiration, sp02, temperature, bloodPressure, motion, proximity, pulse, airNotOxygen, tempAggregateScore, pulseAggregateScore, resAggregateScore, bpAggregateScore, sp02AggregateScore1, batteryLevel, signal, overall, overallLevel, level, target

def Generate(level, next_level_2, errorTarget, active, activated, deactivated, 
              id_, msid, type_, typeVer, key, tic, time):
              # respiration, sp02, temperature, bloodPressure, motion, proximity, pulse, airNotOxygen,
              # tempAggregateScore, pulseAggregateScore, resAggregateScore, bpAggregateScore, sp02AggregateScore1,
              # overall, overallLevel, level, target
  # Generate Data
  respiration, sp02, temperature, bloodPressure, motion, proximity, pulse, airNotOxygen, tempAggregateScore, pulseAggregateScore, resAggregateScore, bpAggregateScore, sp02AggregateScore1, batteryLevel, signal, overall, overallLevel, level, target = GenerateParameters(level, next_level_2, errorTarget)

  generatedTitan = Titan.Generate(id_, msid, type_, typeVer, key, tic)
  generatedData = Data.Generate(respiration, sp02, pulse, temperature, bloodPressure, motion, proximity, airNotOxygen,
              tempAggregateScore, pulseAggregateScore, resAggregateScore, bpAggregateScore, sp02AggregateScore1,
              overall, overallLevel, level, target,
              id_, msid, type_, typeVer, key, tic,
              time, batteryLevel, signal)

  return generatedData

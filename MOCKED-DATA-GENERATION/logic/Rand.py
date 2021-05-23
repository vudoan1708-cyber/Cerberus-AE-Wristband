# random
import random

def String(length, isString):
  text = ''
  possible = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'

  for _ in range(length):
    if isString is True:
      text += possible[random.randint(0, len(possible) - 1)]
    else:
      text += possible[random.randint(52, len(possible) - 1)]
  return text

def Bool():
  if random.randrange(0, 1) < 0.5:
    return True
  else: return False

def DayMonth():
  day = random.randint(1, 29)
  month = random.randint(1, 12)
 
  return day, month

def Time():
  return random.randint(0, 23), random.randint(0, 59), random.randint(0, 59)

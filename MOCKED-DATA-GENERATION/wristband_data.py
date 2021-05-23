# datetime
import datetime

import time

# to read and form json obj
import json

# to access to the operating system
import os

import asyncio

# logic
from logic import Rand

# components
from components import Wristbands

# websockets
# from sockets import Sockets

# setup the GUI
try:                      # In order to be able to import tkinter for
  import tkinter as tk    # either in python 2 or in python 3
except ImportError:
  import Tkinter as tk
# from tkinter import *
root = tk.Tk()

tick = 5
offsetTimeBetweenBands = [2, 10]
levels = []
next_levels_2 = []
errorTargets = []
size = 0
wb_entries = list()
wb_entries_next_level = list()

def toJSON(json_object, i, where):

  # writing to json files
  with open(f"../internal/core/unittest_data/TestWristbandFactory_NewWristband/{where}/wbData_{i+1:03d}.json", "w") as outfile: 
    outfile.write(json_object)

def makeWristbandData(num, size, levels, next_levels_2, errorTargets):

  # Number of Bands
  for i in range(num):
    wbData = []
    # Size of Each Band
    for j in range(size):
      id_ = Rand.String(9, True)
      msid = int(Rand.String(9, False))
      type_ = int(Rand.String(10, False))
      typeVer = int(Rand.String(1, False))
      key = Rand.String(5, True)
      tic = Rand.String(10, True)
      active = True

      # activated day, month
      aday, amonth = Rand.DayMonth()

      # deactivated day, month
      dday, dmonth = Rand.DayMonth()

      # hour, minute, second
      hour, minute, second = Rand.Time()
      
      # activated and deactivated dates
      activated = str(datetime.datetime(2020, amonth, aday, hour, minute, second))

      hour, minute, second = Rand.Time()
      if active is False:
        deactivated = str(datetime.datetime(2020, dmonth, dday, hour, minute, second))
      else: deactivated = "None"

      next_level_el = next_levels_2[i] if j >= (size / 10) else ''
      errorTarget_el = errorTargets[i] if j >= (size / 10) else ''
      # instantiate Wristband object
      wristbands = Wristbands.Generate(levels[i], next_level_el, errorTarget_el, active, activated, deactivated,
                                      id_, msid, type_, typeVer, key, tic, str(datetime.datetime.now()))
      wbData.append(wristbands)

    toJSON(json.dumps(wbData, indent=2), i, "Nurse")

    # print(json.dumps(wbData, indent=2))

      # if not i == num - 1:
      #   time.sleep(random.randint(offsetTimeBetweenBands[0], offsetTimeBetweenBands[1]))

  # dumps json object
  # return json.dumps(wbData, indent=4)
  return

# init socket connection
# ws = Sockets.Init()

# Start The Data Generation
def startGeneratingData():
  # Number of Bands
  input = 1
  if entry.get() == '':
    pass
  else:
    input = int(entry.get())
  numOfBands = input

  # Size of Each Band Data
  if size_entry.get() == '':
    size = 50
  else:
    size = int(size_entry.get())

  # Level of Each Band Data
  for e in (wb_entries):
    level = ''
    next_level = ''
    if e.get() == '' or e.get()[0] == ':':
      level = 'low'
    else:
      level = e.get()

    errorTarget = ''
    # Check if there is a colon in the string (error case)
    # EX: error:proximity
    if ':' in level:
      tempLevelHolder = level.split(':')
      level = tempLevelHolder[0]
      errorTarget = tempLevelHolder[1]
      # Check if there is another level
      # EX: error:proximity->low
      if '->' in tempLevelHolder[1]:
        splitTempLevelHolder = tempLevelHolder[1].split('->')
        # level is error, error target is proximity, next_level is low
        errorTarget = splitTempLevelHolder[0]
        next_level = splitTempLevelHolder[1]
      # EX: low->error:proximity
      elif '->' in level:
        splitTempLevelHolder = level.split('->')
        # level is low, error target is proximity, next_level is error
        level = splitTempLevelHolder[0]
        errorTarget = tempLevelHolder[1]
        next_level = splitTempLevelHolder[1]

    # If no error case
    # EX: low->high
    elif '->' in level:
      tempLevelHolder = level.split('->')
      level = tempLevelHolder[0]
      next_level = tempLevelHolder[1]

    # Otherwise, no-level-changing case
    # EX: low-medium
    levels.append(level)
    next_levels_2.append(next_level)
    errorTargets.append(errorTarget)
  print(f"1st Individual Levels {levels}")
  print(f"2nd Individual Levels {next_levels_2}")

  # spawn a JSON obj with a length of any number provided as a parameter
  makeWristbandData(numOfBands, size, levels, next_levels_2, errorTargets)

  # Sockets.Send(ws, obj)
  # print(obj)
  # time.sleep(random.randint(2, 10))



#############################
# TKINTER

# Add Wristbands
def addWBLevel():
  # Get the text number in the entry, and loop through it
  for _ in range(int(entry.get())):
    wb_entries.append(tk.Entry(root))
    wb_entries[-1].pack()
  print(len(wb_entries))

## Add Next Levels
def addWBNextLevels():
  # Get the text number in the entry, and loop through it
  for _ in range(int(entry.get())):
    wb_entries_next_level.append(tk.Entry(root))
    wb_entries_next_level[-1].pack(side='bottom')
  print(len(wb_entries_next_level))

# change windows title
root.title('Wristband Data Generation')
# create canvas
width = 500
height = 500
colour = '#263d42'
canvas = tk.Canvas(root, width = width, height = height, bg = colour)

# attach canvas to the root
canvas.pack()

# add a frame (similar to adding a html tag)
frame = tk.Frame(root, bg = 'white')

label = tk.Label(frame, text="Add Number of Wristbands", fg='black', bg='white', font='none 9')
# using place method we can set the position of label
label.place(relx = 0.0,
            rely = 0.1,
            anchor ='sw')
# label.pack()

# add textbox entry
entry = tk.Entry(root)
canvas.create_window(width / 2, 110, window=entry)
# Add Wristbands Levels button
add_wb_btn = tk.Button(frame, text = "Add Wristbands' Levels",
                            padx = 10, pady = 5, fg = 'white', bg = colour, command = addWBLevel)
add_wb_btn.place(x = width / 2,
                 y = 100)

# SIZE
size_label = tk.Label(frame, text="Add A Wristband Data Size", fg='black', bg='white', font='none 9')
# using place method we can set the position of label
size_label.place(relx = 0.0,
                 rely = 0.4,
                 anchor ='sw')
# size_label.pack()
# add textbox entry for a wristband data file size
size_entry = tk.Entry(root)
canvas.create_window(width / 2, height / 2, window=size_entry)

# Add Wristbands Next Levels button
# add_next_wb_btn = tk.Button(frame, text = "Add Wristbands' Next Levels",
#                             padx = 10, pady = 5, fg = 'white', bg = colour, command = addWBNextLevels)
# add_next_wb_btn.place(x = width / 2,
#                  y = 250)

# attach frame to the root, set its width, height, x, y
# rel means relative
frame.place(relwidth = 0.8, relheight = 0.8, relx = 0.1, rely = 0.1)

# add buttons
# fg means foreground (applying colour to the text)
start_btn = tk.Button(root, text = 'Start',
                            padx = 10, pady = 5, fg = 'white', bg = colour, command = startGeneratingData)
start_btn.pack()

root.mainloop()

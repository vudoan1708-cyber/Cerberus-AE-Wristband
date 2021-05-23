import React from 'react'

const NurseDescription = ({data}) => {

    const target = data.updateWristbandDataAlert.target;
    let value = data.updateWristbandDataAlert.sensorData[data.updateWristbandDataAlert.target]

    let measurement = ""
    let name = "error"
    
    if (data.updateWristbandDataAlert.trend === "unchanged") {// if there are no changes in the vitals, display no rapid changes
      name = "No rapid changes"
      value = ""
    } else {
      switch(target) {
        case "bloodPressure":
          name = "blood pressure";
          measurement = "mmHg";
          break;
        case "motion":
          value = ""
          // No change in name, "No Rapid changes"
          measurement = "";
          break;
        case "proximity": 
          value = ""
          name = "No skin contact"
          measurement = "";
          break;
        case "pulse": 
          name = "Pulse rate"
          // default value
          measurement = "p/m";
          break;
        case "respiration":
          name = "Respiration rate";
          // default value
          measurement = "bpm";
          break;
        case "sp02":
          name = "Respiration rate";
          // default value
          measurement = "p/m";
          break;
        case "temperature":
          name = "Temperature";
          // default value
          measurement = "\u00B0C";// Unicode values are used to represent special characters inside strings
          break;
        case "batteryLevel":
          name = "Low Battery";
          value = ""
          // default value
          measurement = "";// Unicode values are used to represent special characters inside strings
          break;
        case "noConnection":
          name = "No connection";
          value = ""
          // default value
          measurement = "";// Unicode values are used to represent special characters inside strings
          break;
        default: 
          name = `error in target value > switch case ${target} + ${value}`;
          value = "";
          measurement = ""
      }
    }
    

    return <span>{name} {value + measurement}</span>
    
    
  }

export default NurseDescription
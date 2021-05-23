import React, { useEffect, useState } from "react";

//Packages
import { gql, useQuery, useSubscription } from "@apollo/client";
import { Link } from 'react-router-dom';
import { animated, useSpring } from 'react-spring';

//Components
import Arrow from "./Arrow";
import MovePatient from './options/MovePatient';
import NurseDescription from './NurseDescription';

//SVG
import BatteryIcon from "../assets/svg/battery.svg";
import BatteryLow from "../assets/svg/battery-low.svg";
import BatteryUnknown from "../assets/svg/battery-unknown.svg"

import Transfer from "../assets/svg/transfer.svg"
import ConnectionIcon from "../assets/svg/connection.svg"
import LostConnection from "../assets/svg/lost-connection.svg"

//GraphQL
const UPDATE_DATA_ALERTS = gql`
    subscription updateWristbandDataAlert($id: ID) {
        updateWristbandDataAlert(id: $id) {
            sensorData {
            respiration,
            sp02,
            pulse,
            temperature,
            bloodPressure
            motion,
            proximity,
            },
            level,
            trend,
            target,
            overall,
            overallLevel,
        }
    }
`;

const DATA_SUBSCRIPTION = gql`
  subscription OnDataAdded($id: ID) {
    updateWristbandData(id: $id) {
      batteryLevel
      sensorData {
        proximity,
      },
      news2 {
        overall
      }
    }

    
  }
`;

const ACTIVE_SUBSCRIPTION = gql`
  subscription updateWristbandActive($id: ID) {
    updateWristbandActive(id: $id) {
      id,
      active,
    }
  }
`

// const MULTIPLE_WRISTBAND_DATA = gql`
//   query($id: ID!, $howMany: Int!, $start: Int!, $end: Int!) {
//     getMultipleWristbandData(
//       id: $id
//       howMany: $howMany
//       start: $start
//       end: $end
//     ) {
//       batteryLevel
//     }
//   }
// `;

const WristbandCard = ({ refetch, wristband, setWristbandID, setShowDetail, nurseView, props, important, showcaseBatteryLevel }) => {
  
  //Hooks
  const [ importantBand, setImportantBand ] = useState();

  // Subscriptions
  const {loading, error, data } = useSubscription(DATA_SUBSCRIPTION, { variables: { id: wristband.id }});


  // eslint-disable-next-line no-unused-expressions
  wristband.name === "Kid Putowski" ? console.log(`${wristband.name} active: ${wristband.active}`) : undefined;

  const handleClick = () => {
    if (nurseView) {
      setWristbandID(wristband.id);
    } else {
      console.log("SetShowDetail: ")
      setShowDetail(true);
      console.log(wristband.id);
      setWristbandID(wristband.id);
    }  
  }

  // Important Band check
  useEffect(()=> {
    if (data !== undefined) {
      if (data.updateWristbandData.news2.overall >= 4 || !data.updateWristbandData.sensorData.proximity || data.updateWristbandData.batteryLevel < 30 || wristband.pregnant || wristband.child) {
          setImportantBand(true);
      } else setImportantBand(false)
    } else setImportantBand(false);
  },[important, data])

  //  if (loading) return 'Loading...';
  if(error) return("Error in Wristbands: " + error.message);


  // if (wristband.name === "Vu Doan") {
  //   console.log("Vu Doan data: ")
  //   console.log(data);
  // }

  // if (wristband.name === "Low Battery") {
  //   console.log("Low Battery data: ")
  //   console.log(data);
  // }

    if (wristband.active && wristband.department === "waiting-room") {// Checks if the wirstband is active and in the right department. At the moment there is only one department so it is hardcoded.
      
      if (important && importantBand || !important && !importantBand || !important && importantBand) {
        if (nurseView) {
          return (
          <NurseWristband refetch={refetch} wristbandData={data} wristband={wristband} props={props} handleClick={handleClick}/>
        )} else { 
            return (//Admin Wristband
          <animated.figure style={props} onClick={handleClick} className={data ? cardBoxShadow(data.updateWristbandData.sensorData): "card"}>
            <header>
              <Connection data={data}/>
              {showcaseBatteryLevel === true & wristband.name === "Low Battery" ? <img alt="battery level" src={BatteryIcon}/>
                :<Battery data={data}/>}
            </header>
            <h3>{adminDescription(wristband, data)}</h3>
            <h2>{wristband.name}</h2>
          </animated.figure>
        )}

      } else return null
    } else return null
  }

  const adminDescription = (wristband, data) => {

    let message ="No issues to report"

      if (wristband.child) {
        message ="Under 16"
      }
      if (wristband.pregnant) {
        message ="Pregnant"
      }

      if (data) {
        if (data.updateWristbandData) {
          let sensorData = data.updateWristbandData.sensorData
        
          if (sensorData.batteryLevel < 30) {
            message = "Low Battery"
          }
  
          if (sensorData.proximity === false) {
            message = "No skin contact"
          }
        }
        
      }
      

    return message
  }
  
  const NurseWristband = ({ refetch, wristband, wristbandData, props, handleClick}) => {
    
    const { loading, error, data } = useSubscription(UPDATE_DATA_ALERTS, {variables: {id: wristband.id}})
    const [transfer, setTransfer] = useState(false);
    
    // if (data) {
    //   console.log(wristband.name);
    //   console.log("trend is " + data.updateWristbandDataAlert.trend);
    //   console.log("level is " + data.updateWristbandDataAlert.level);
    // }
    
    if (wristband.name === "Andrew" && data) {
      console.log(wristband.name + ": Nurse card");
      console.log(wristbandData);
      console.log("temperature: " + data.updateWristbandDataAlert.sensorData.temperature)
      console.log("target: " + data.updateWristbandDataAlert.target)
      console.log("trend: " + data.updateWristbandDataAlert.trend)
      console.log("overallLevel: " + data.updateWristbandDataAlert.overallLevel)
      console.log(data);
      
    }
    
    
       return (
         <>
        
        <animated.article style={props, {borderLeftColor: data && data.updateWristbandDataAlert.sensorData.proximity ? bandColour(data.updateWristbandDataAlert.overallLevel) : "var(--dark-navy)"}} onClick={handleClick} className={wristbandData ? cardBoxShadow(wristbandData.updateWristbandData.sensorData): "card"}>
        <Link className="right" to={`patient-cart`}>
          <header>
            <Connection data={wristbandData}/>
            <Battery data={wristbandData}/>
          </header>
            <h2>{wristband.name}</h2>
            {data ? 
            <div>
              {
                wristband.name == "Donata Lesiak" && data && data.updateWristbandDataAlert.target =="respiration" ?
                <svg style={{transform: "rotate(180deg)"}} width="12" height="13" viewBox="0 0 12 13" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M0.0205962 5.7934C0.0628321 5.89465 0.16139 5.96072 0.270762 5.96072L3.23755 5.96072L3.23755 12.7293C3.23755 12.8787 3.35885 13 3.50831 13L7.8402 13C7.98965 13 8.11095 12.8787 8.11095 12.7293L8.11095 5.9607L11.0891 5.9607C11.1985 5.9607 11.297 5.89463 11.3393 5.79391C11.381 5.69266 11.3582 5.57623 11.2808 5.49879L5.87948 0.0796037C5.82859 0.0287125 5.7598 5.23247e-06 5.68779 5.22618e-06C5.61579 5.21988e-06 5.547 0.0287124 5.49611 0.0790705L0.0790765 5.49826C0.00163553 5.57573 -0.0216397 5.69213 0.0205962 5.7934Z" fill="var(--red)"/>
                </svg>
            : 
              wristband.name == "Andrew" && data && data.updateWristbandDataAlert.target =="temperature" ? //temporary showcase changes for Andrew Card
               <svg style={{transform: "rotate(0deg)"}} width="12" height="13" viewBox="0 0 12 13" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M0.0205962 5.7934C0.0628321 5.89465 0.16139 5.96072 0.270762 5.96072L3.23755 5.96072L3.23755 12.7293C3.23755 12.8787 3.35885 13 3.50831 13L7.8402 13C7.98965 13 8.11095 12.8787 8.11095 12.7293L8.11095 5.9607L11.0891 5.9607C11.1985 5.9607 11.297 5.89463 11.3393 5.79391C11.381 5.69266 11.3582 5.57623 11.2808 5.49879L5.87948 0.0796037C5.82859 0.0287125 5.7598 5.23247e-06 5.68779 5.22618e-06C5.61579 5.21988e-06 5.547 0.0287124 5.49611 0.0790705L0.0790765 5.49826C0.00163553 5.57573 -0.0216397 5.69213 0.0205962 5.7934Z" fill="var(--yellow)"/>
                </svg> : wristband.name == "Andrew" && data && data.updateWristbandDataAlert.target =="pulse" ? 
                <svg style={{transform: "rotate(0deg)"}} width="12" height="13" viewBox="0 0 12 13" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M0.0205962 5.7934C0.0628321 5.89465 0.16139 5.96072 0.270762 5.96072L3.23755 5.96072L3.23755 12.7293C3.23755 12.8787 3.35885 13 3.50831 13L7.8402 13C7.98965 13 8.11095 12.8787 8.11095 12.7293L8.11095 5.9607L11.0891 5.9607C11.1985 5.9607 11.297 5.89463 11.3393 5.79391C11.381 5.69266 11.3582 5.57623 11.2808 5.49879L5.87948 0.0796037C5.82859 0.0287125 5.7598 5.23247e-06 5.68779 5.22618e-06C5.61579 5.21988e-06 5.547 0.0287124 5.49611 0.0790705L0.0790765 5.49826C0.00163553 5.57573 -0.0216397 5.69213 0.0205962 5.7934Z" fill="var(--red)"/>
                </svg>
                : <Arrow pregnant={wristband.pregnant} child={wristband.child} level={data.updateWristbandDataAlert.level} trend={data.updateWristbandDataAlert.trend}/> 
              }
              <h3>{!wristbandData ? "data loading..." :
              
                wristband.name == "Donata Lesiak" && data && data.updateWristbandDataAlert.target =="respiration" ?
                `Respiration rate ${data.updateWristbandDataAlert.sensorData.respiration}bpm`
              : wristband.name == "Andrew" && data && data.updateWristbandDataAlert.target =="temperature" ? //temporary showcase changes for Andrew Card
              `Temperature ${data.updateWristbandDataAlert.sensorData.temperature} \u00B0C` :
              wristband.name == "Andrew" && data && data.updateWristbandDataAlert.target =="pulse" ?
              "Pulse rate 140 p/m" :
              wristbandData.updateWristbandData.batteryLevel > 1 ?  <NurseDescription data={data} />  : "No connection"}</h3>
            </div>
        
            : "loading data..."}
        </Link>
          <footer><img onClick={()=> setTransfer(true)} alt="transfer" src={Transfer}/></footer>
        </animated.article>
      
      { transfer ? <MovePatient refetch={refetch} wristbandID={wristband.id} exit={()=> setTransfer(false)}/> : null}
        </>
      )
}
  
//js functions

  const Connection = ({ data }) => {
    
    let batteryLevel;
    if (data !== undefined) {
      if (data.getMultipleWristbandData !== undefined) {
        batteryLevel = data.getMultipleWristbandData[0].batteryLevel;
      } else batteryLevel = data.updateWristbandData.batteryLevel;
    }

    let proximity
    if (data !== undefined) {
      if (data.getMultipleWristbandData !== undefined) {
        proximity = data.getMultipleWristbandData[0].proximity;
        console.log(data.getMultipleWristbandData[0]);
      } else proximity = data.updateWristbandData.sensorData.proximity;
    }
    

    if (proximity && batteryLevel) return (
      <img alt="connected" src={ConnectionIcon}/>
    ); else return <img alt="lost connection" src={LostConnection}/>
  }

  const Battery = ({ data }) => {

    let batteryLevel;
    if (data !== undefined) {
      if (data.getMultipleWristbandData !== undefined) {
        batteryLevel = data.getMultipleWristbandData[0].batteryLevel;
      } else batteryLevel = data.updateWristbandData.batteryLevel;
    }
    
    if (batteryLevel) {
      if (batteryLevel > 30) {
        return <img alt="battery level" src={BatteryIcon}/>;
      } else if (batteryLevel > 1) {
        return <img alt="battery level" src={BatteryLow}/>
      } else return <img alt="battery level" src={BatteryUnknown}/>;
    } else return <img alt="battery level" src={BatteryUnknown}/>;
  }

  

  const bandColour = (overallLevel) => {
    
    switch(overallLevel) {
      case "low":
        return "var(--grey-low)";
      case "low-medium":
        return "var(--yellow)"
      case "medium":
        return "var(--orange)"
      case "high":
        return "var(--red)"
      case "error":
        return "var(--dark-navy)"
      default: 
        return "white"
    }
  }

  const cardBoxShadow = (sensorData) => {
    
    if (sensorData) {
      if (sensorData.proximity === false) {
        return "card red"
      } else 
      return "card"
    } else 
    return "card"
    
  }
  

export default WristbandCard;
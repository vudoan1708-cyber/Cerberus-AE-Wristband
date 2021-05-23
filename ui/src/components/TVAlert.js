import React, { useState } from 'react';

//Packages
import {useSubscription, gql} from '@apollo/client';
import { Link } from 'react-router-dom';

//Components
import Arrow from "./Arrow";
import NurseDescription from './NurseDescription';
import MovePatient from './options/MovePatient';

//SVG
import Transfer from "../assets/svg/transfer.svg"
import ConnectionIcon from "../assets/svg/connection.svg"
import BatteryIcon from "../assets/svg/battery.svg";

//GraphQL

const UPDATE_DATA_ALERTS = gql`
    subscription updateWristbandDataAlert($id: ID!) {
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
            overallLevel,
            overall
        }
    }
`

const TVAlert = ({ name, id, child, pregnant, active, department, props, TVFilter }) => {

    const [transfer, setTransfer] = useState(false)

    const { loading, error, data } = useSubscription(UPDATE_DATA_ALERTS, {variables: {id: id}})

    console.log(id);
    console.log(TVFilter)
    console.log(data);
    console.log(active);
    console.log(department);


    if (data && active && department === "waiting-room") {// Checks if the wirstband is active and in the right department. At the moment there is only one department so it is hardcoded.
        const level = data.updateWristbandDataAlert.level;
        const trend = data.updateWristbandDataAlert.trend;
        const overall = data.updateWristbandDataAlert.overall;
        const overallLevel = data.updateWristbandDataAlert.overallLevel;

        const proximity = data.updateWristbandDataAlert.sensorData.proximity;
    
        console.log("Alert: " + name);
        console.log("filter: " + TVFilter);
        console.log("overallLevel: " + overallLevel);
        console.log("special: " + pregnant + child);
    
                if ((TVFilter === overallLevel && proximity && !child && !pregnant) ||  (TVFilter === "special" && pregnant) ||  (TVFilter === "special" && child)|| (TVFilter === "error" && !proximity)) {
                    console.log("Alert worked: " + name);
                    return (
                    <>
                    <article className={data ? alertBoxShadow(data.updateWristbandDataAlert.sensorData) : "alert"}>
                        <aside>
                            <OverallNumber child={child} pregnant={pregnant} level={level} overallLevel={overallLevel} overall={overall}/>
                            <p>{child ? "Under 16" : pregnant ?  "Pregnant" : overallLevel}</p>
                        </aside>
                        <section>
                                <div style={{justifyContent: "flex-end"}}>
                                    <img alt="connected" src={ConnectionIcon}/>
                                    <img alt="battery level" src={BatteryIcon}/>
                                </div>
                            <header>
                                <h2>{name}</h2>
                            </header>
                            <div>
                            {
                                name ==="Donata Lesiak" && data && data.updateWristbandDataAlert.target =="respiration" ?
                                <svg style={{transform: "rotate(180deg)"}} width="12" height="13" viewBox="0 0 12 13" fill="none" xmlns="http://www.w3.org/2000/svg">
                                    <path d="M0.0205962 5.7934C0.0628321 5.89465 0.16139 5.96072 0.270762 5.96072L3.23755 5.96072L3.23755 12.7293C3.23755 12.8787 3.35885 13 3.50831 13L7.8402 13C7.98965 13 8.11095 12.8787 8.11095 12.7293L8.11095 5.9607L11.0891 5.9607C11.1985 5.9607 11.297 5.89463 11.3393 5.79391C11.381 5.69266 11.3582 5.57623 11.2808 5.49879L5.87948 0.0796037C5.82859 0.0287125 5.7598 5.23247e-06 5.68779 5.22618e-06C5.61579 5.21988e-06 5.547 0.0287124 5.49611 0.0790705L0.0790765 5.49826C0.00163553 5.57573 -0.0216397 5.69213 0.0205962 5.7934Z" fill="var(--red)"/>
                                </svg>
                            : name === "Andrew" && data && data.updateWristbandDataAlert.target ==="pulse" ?  //Temporary showcase changes
                            <svg style={{transform: "rotate(0deg)"}} width="12" height="13" viewBox="0 0 12 13" fill="none" xmlns="http://www.w3.org/2000/svg">
                                <path d="M0.0205962 5.7934C0.0628321 5.89465 0.16139 5.96072 0.270762 5.96072L3.23755 5.96072L3.23755 12.7293C3.23755 12.8787 3.35885 13 3.50831 13L7.8402 13C7.98965 13 8.11095 12.8787 8.11095 12.7293L8.11095 5.9607L11.0891 5.9607C11.1985 5.9607 11.297 5.89463 11.3393 5.79391C11.381 5.69266 11.3582 5.57623 11.2808 5.49879L5.87948 0.0796037C5.82859 0.0287125 5.7598 5.23247e-06 5.68779 5.22618e-06C5.61579 5.21988e-06 5.547 0.0287124 5.49611 0.0790705L0.0790765 5.49826C0.00163553 5.57573 -0.0216397 5.69213 0.0205962 5.7934Z" fill="var(--red)"/>
                            </svg>
                            :<Arrow level={level} trend={trend}/> 
                            }
                                { 
                                data && data.updateWristbandDataAlert.target === "batteryLevel" ?  //Temporary showcase changes
                                <span>No connection</span>
                                :name === "Donata Lesiak" && data && data.updateWristbandDataAlert.target =="respiration" ?  //Temporary showcase changes
                                        <span>Respiration rate {data.updateWristbandDataAlert.sensorData.respiration} bpm</span>
                                        :
                                    name === "Andrew" && data && data.updateWristbandDataAlert.target =="pulse" ?  //Temporary showcase changes
                                    <span>Pulse rate 140 p/m</span>
                                    :<NurseDescription data={data}/>
                                }
                            </div>
                            <footer><img onClick={()=> setTransfer(true)} alt="transfer" src={Transfer}/></footer>
                        </section>
                    </article>
                    { transfer ? <MovePatient wristbandID={id} exit={()=> setTransfer(false)}/> : null}
                    </>
                    )
                } else return null
    } else return null
}


const OverallNumber = ({ overallLevel, overall, level, child, pregnant }) => {

    let colour = "black"

    if (child || pregnant) {
        return <h1 style={{color: colour}}>N/A</h1>
    } else {
        if (overallLevel === "high") {
            colour = "var(--red)"
        } else if (overallLevel === "medium") {
            colour = "var(--orange)"
        } else if (overallLevel === "low-medium") {
            colour = "var(--yellow)"
        } else if (overallLevel === "low") {
            colour = "var(--grey)"
        }
    
        if (level !== "error") return (
            <h1 style={{color: colour}}>{overall}</h1>
        ); else return <h1>Err</h1>
    }
    
}

const alertBoxShadow = (sensorData) => {
    
    if (sensorData) {
      if (sensorData.proximity === false) {
        return "alert red"
      } else 
      return "alert"
    } else 
    return "alert"
    
  }


export default TVAlert
import React, { useState } from 'react';
import '../../App.css';

//Components
import MovePatient from './MovePatient';
import EditDetails from './EditDetails';
import ReAssignPatient from './ReAssignPatient';
import RemovePatient from './RemovePatient';


const Options = ({ patientCart, refetch, wristbandID, editDetails, reAssign, RemoveBand, movePatient, setShowcaseBatteryLevel}) => {

    const [expanded, setExpanded] = useState(false);

    const [option, setOption] = useState(0);

    const handleClick = () => {
        if (expanded) {
            setExpanded(false);
        } else {
            setExpanded(true);
        }
    }
    
    // console.log("setShowcaseBatteryLevel: " + setShowcaseBatteryLevel);

    const shadow = "0px 0px 4px rgba(130, 130, 130, 0.25)";
    const expandedShadow = "-6px 6px 15px 0px rgba(0, 0, 0, 0.25)";

        return (
        <>
        <div className="options">
            <figure style={{WebkitBoxShadow: expanded ? expandedShadow : shadow, MozBoxShadow: expanded ? expandedShadow : shadow, boxShadow: expanded ? expandedShadow : shadow}}>
            <button onClick={handleClick}>
                <h2>...</h2>
                <h3>Options</h3>
                <h2>{expanded ? ">" : "<"}</h2>
            </button>
            {expanded ? <div className="expanded">
                {movePatient ? <button onClick={()=> setOption(1)}>Move patient</button> : null}
                {editDetails  ? <button onClick={()=> setOption(2)}>Edit details</button> : null}
                {reAssign ? <button onClick={()=> setOption(3)}>Re-assign to new band</button> : null}
                {RemoveBand ? <button onClick={()=> setOption(4)}>Remove band</button> : null}
            </div>: null}
            </figure>
        </div>

        <Selector patientCart={patientCart} refetch={refetch} setShowcaseBatteryLevel={setShowcaseBatteryLevel} wristbandID={wristbandID} setOption={setOption} option={option}/>
        </>
    
    )
}

const Selector = ({patientCart, refetch, option, setOption, wristbandID, setShowcaseBatteryLevel}) => {

    // console.log("setShowcaseBatteryLevel: " + setShowcaseBatteryLevel);

    switch(option) {
        case 1:
            return <MovePatient patientCart={patientCart} refetch={refetch} wristbandID={wristbandID} exit={()=> setOption(0)}/>
        case 2:
            return <EditDetails patientCart={patientCart} refetch={refetch} wristbandID={wristbandID} exit={()=> setOption(0)}/>
        case 3:
            return <ReAssignPatient patientCart={patientCart} refetch={refetch} setShowcaseBatteryLevel={setShowcaseBatteryLevel} wristbandID={wristbandID} exit={()=> setOption(0)}/>
        case 4:
            return <RemovePatient patientCart={patientCart} refetch={refetch} wristbandID={wristbandID} exit={()=> setOption(0)}/>
        default:
            return null
    }
}





export default Options
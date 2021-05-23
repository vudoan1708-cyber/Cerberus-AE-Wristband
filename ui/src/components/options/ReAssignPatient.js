import React, { useEffect, useState } from 'react';
import '../../App.css';

// Package
import {gql, useMutation } from '@apollo/client';
import { useHistory } from 'react-router-dom';

//Components
import Confirmation from '../Confirmation';

//GraphQL

    const REASSIGN_NEW_WRISTBAND = gql`
            mutation ($oldWristband: DeactivateWristbandInput!) {
            reassignNewWristband(oldWristband: $oldWristband) {
                id,
                tic,
                active,
                name,
                onOxygen,
                pregnant,
                child,
                key,
                department,
            }
        }
    `

const ReAssignPatient = ({ refetch, patientCart, exit, wristbandID, setShowcaseBatteryLevel }) => {

    // console.log("setShowcaseBatteryLevel: " + setShowcaseBatteryLevel);
    const [confirm, setConfirm] = useState(false);

    //Router
    let history = useHistory();

    const [reAssignPatient] = useMutation(REASSIGN_NEW_WRISTBAND, {
        onCompleted: ()=> {
            if (patientCart !== undefined) {
                console.log("history")
                history.push("/nurse/default-view");
            } else {
                refetch();
            }
        }
    });


 
    useEffect (()=> {
        if (confirm) {
            //console.log("Wristband " + wristbandID + " Re-assigned to " + (wristbandID + 1));
            //setShowcaseBatteryLevel(true);
            reAssignPatient({variables: {oldWristband: {id: wristbandID}}})
            // console.log("setShowcaseBatteryLevel: " + setShowcaseBatteryLevel);
        }

    }, [confirm, wristbandID])
 
    return (
        <>
        {!confirm ? <div className="pop-up ">
            <header>
                <h1>Re-assign patient to new band</h1>
            </header>
            <article>
                <h2 className="message">All data will be automatically transferred to the new band</h2>
                <article className="message">
                    <p>New band ID</p>
                    <h2>AE7849249</h2>
                </article>
            </article>
            <footer>
                <button className="empty" onClick={exit}>Cancel</button>
                <button className="navy" onClick={()=> setConfirm(true)}>Confirm</button>
            </footer>
        </div> : <Confirmation message="Patient band updated successfully" exit={exit}/>}
        <div className="background"></div>
        </>
    )
}

export default ReAssignPatient
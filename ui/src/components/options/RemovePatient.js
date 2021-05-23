import React, { useEffect, useState } from 'react';
import '../../App.css';

//Packages
import {gql, useMutation } from '@apollo/client';
import { useHistory } from 'react-router-dom';

//Components
import Confirmation from '../Confirmation';


//GraphQL
const DEACTIVATE_WRISTBAND = gql`
mutation deactivateWristband ($input: DeactivateWristbandInput!){
    deactivateWristband (input: $input) {
        id
    }
}
`

const RemovePatient = ({ patientCart, exit, wristbandID }) => {

    const [confirm, setConfirm] = useState(false);

    //Router
    let history = useHistory(); 
      
    const [deactivateWristband] = useMutation(DEACTIVATE_WRISTBAND, {
        onCompleted: ()=> {
            if (patientCart !== undefined) {
                console.log("history")
                history.push("/nurse/default-view");
            }
        }
    });

    
 
    useEffect (()=> {
        if (confirm) {
            console.log("Wristband " + (wristbandID) + " Deactivating");
            deactivateWristband({variables: {input: {id: wristbandID}}});
        }

    }, [confirm, deactivateWristband, wristbandID])

    return (
        <>
        {!confirm ? <div className="pop-up">
            <header>
                <h1>Remove band</h1>
                <h2 className="exit" onClick={exit}>+</h2>
            </header>
            <div className="message">
                <h2>Are you sure you want to remove the band?<br/>
                All recorded data will be permanently deleted</h2>
            </div>
            <footer>
                <button className="empty" onClick={exit}>Cancel</button>
                <button className="navy" onClick={() =>setConfirm(true)}>Confirm</button>
              </footer>
        </div> : <Confirmation message="Patient band removed successfully" exit={exit}/>}
        <div className="background"></div>
        </>
    )
    
}

export default RemovePatient
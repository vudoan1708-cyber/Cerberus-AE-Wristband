import React, { useState, useEffect } from 'react';
import '../../App.css';

//Packages
import {gql, useMutation } from '@apollo/client';
import { useHistory } from 'react-router-dom';
import Confirmation from '../Confirmation';

const MovePatient = ({ refetch, patientCart, exit, wristbandID }) => {

    const [confirm, setConfirm] = useState(false);
    const [department, setDepartment] = useState(null);// Customise these locations to your liking below

    //Router
    let history = useHistory();

    const MOVE_WRISTBAND = gql`
        mutation ($id:ID!, $value: String!){
            resetDepartment (id: $id, value: $value) {id}
        }
    `

    const [resetDepartment] = useMutation(MOVE_WRISTBAND, {
        onCompleted: ()=> {
            if (patientCart !== undefined) {
                console.log("history")
                history.push("/nurse/default-view");
            }

            if (refetch !== undefined) {
                refetch()
            }
        }
    });

    const departmentArray = ["Majors", "Ambulatory", "Main Hospital", "Discharged"];
 
    useEffect (()=> {
        if (confirm) {
            console.log("Wristband " + (wristbandID) + " Moved to " + departmentArray[department]);
            resetDepartment({variables: {id: wristbandID, value: departmentArray[department]}})
        }

    }, [confirm, wristbandID])

    const handleClick = () => {
        console.log(department);
        if (department >= 0 && department < departmentArray.length) {
            console.log("department");
            setConfirm(true)
        }
    }
 
    return (
        <>
        {!confirm ? <section className="pop-up ">
            <header>
                <h1>Move Patient</h1>
                <h1 className="exit">+</h1>
            </header>

            <section>
                <article className="message">
                    <p>Current location area</p>
                    <h2>Waiting Room</h2>
                </article>
                <article className="message">
                    <p>Move to</p>
                    <article className="transfer">{/*Locations are customisable for individual hopsitals*/}
                        <h2 className={department === 0 ? "highlight" : null} onClick={()=> setDepartment(0)}>Majors</h2>
                        <h2 className={department === 1 ? "highlight" : null} onClick={()=> setDepartment(1)}>Ambulatory</h2>
                        <h2 className={department === 2 ? "highlight" : null} onClick={()=> setDepartment(2)}>Main Hospital</h2>
                        <h2 className={department === 3 ? "highlight" : null} onClick={()=> setDepartment(3)}>Discharged</h2>
                    </article>
                </article>
            </section>

            <footer>
                <button className="empty" onClick={exit}>Cancel</button>
                <button className="navy" onClick={handleClick}>Confirm</button>
            </footer>
        </section> : <Confirmation message="Patient band moved successfully" exit={exit}/>}
        <div className="background"></div>
        </>
    )
}

export default MovePatient
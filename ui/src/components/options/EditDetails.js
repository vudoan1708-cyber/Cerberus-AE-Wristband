import React, { useEffect, useState } from 'react';
import '../../App.css';
import {gql, useMutation, useQuery } from '@apollo/client';

//Components
import Confirmation from '../Confirmation';
import Info from '../Info';

const EditDetails = ({refetch, exit, wristbandID }) => {

    const [confirm, setConfirm] = useState(false);

    const [name, setName] = useState('N/A');
    const [pregnant , setPregnant] = useState(true);
    const [child, setChild] = useState();
    //const [dateOfBirth, setDateOfBirth] = useState()

    //Edit Wristband Mutation
      
      
  //GraphQL
    const MULTIPLE_WRISTBAND_DATA = gql `
        query getMultipleWristbandData ($id: ID!, $howMany: Int!, $start: Int!, $end: Int!) {

            getWristband (id: $id) {
            name,
            dateOfBirth,
            key,
            child,
            pregnant,
            }

            getMultipleWristbandData (id: $id, howMany: $howMany, start: $start, end: $end) {
                sensorData {
                respiration,
                sp02,
                motion,
                temperature
                },
                batteryLevel
            }
        }
        `

    const {loading, error, data} = useQuery(MULTIPLE_WRISTBAND_DATA, {variables: {id: wristbandID, howMany: 1, start: 0, end: 0}})
  
  
    const EDIT_WRISTBAND_DETAILS = gql`
      mutation editWristband ($id: ID!, $options: [String!]!, $values: [String!]!){
        resetMultipleFields (id: $id, options: $options, values: $values)
            {
                id,
                tic,
                name,
                onOxygen,
                pregnant,
                child,
                key,
                department
            }
        }
    `

    const [editWristband] = useMutation(EDIT_WRISTBAND_DETAILS);
      
      const handleInputChange = (e) => {
        pregnant ? setPregnant(false) : setPregnant(true);
        console.log(pregnant);
      }
    
      const handleDate = (e) => {
          console.log("left focus");
          console.log(child);
          if (e) {
            console.log(e.target.value);
            const dob = new Date(e.target.value);
            console.log(typeof(dob));
            const d = new Date();
            console.log(typeof(d));
    
            const diffTime = Math.abs(d - dob);// milliseconds
            console.log(diffTime);
            const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));// days
            console.log(diffDays);
            const diffYears = diffDays / 365;
    
            console.log(diffYears);
            if (diffYears < 16) {
              console.log("this is a child");
              setChild(true);
            } else {
              setChild(false);
              console.log("legal adult");
            }
          }
    
        }

    useEffect(()=> {// Setting the Existing Input Values 
        if (data) {
            setPregnant(data.getWristband.pregnant);
            setChild(data.getWristband.child);
            setName(data.getWristband.name);
        }
    }, [data])

    useEffect (()=> {
        if (confirm) {
            console.log("Wristband Details" + wristbandID + " edited");
            editWristband({
                variables: {
                    id: wristbandID,
                    options: [
                        "Name",
                        "Pregnant",
                        "Child",
                        //"dateOfBirth"
                ],
                values: [
                    name.toString(),
                    pregnant.toString(),
                    child.toString(),
                    //dateOfBirth,
                ]
                }
              });
        }

    }, [confirm, name, child, pregnant, wristbandID])


    if (loading) return 'Loading...';
    if(error) return("Error in  Edit Details: " + error.message);

        if (data) return (
        <>
        {!confirm ? <div className="pop-up new-patient">
            <header>
                <h1>Edit Patient details</h1>
                <h2 className="exit" onClick={exit}>+</h2>
            </header>
            <section>
                <article>
                  <p>Patient name</p>
                  <input className="text-field" value={name} onChange={e=> (setName(e.target.value))} type="text"/>
                  <div><p>Date of birth</p><Info message="The system will automatically calculate the age. NEWS system does not apply to patients under 16 years old "/></div>
                  <input value={data ? data.getWristband.dateOfBirth: null} onBlur={(e)=> handleDate(e)} className="text-field" type="date"/>
                  <p>NHS number</p>
                    <input className="text-field" type="text"></input>
                </article>
                <article>
                  <p>Band ID</p>
                  <h2>AE784729</h2>
                  <div style={{display:"flex", alignItems:"center"}}><p>Pregnancy</p><Info message="NEWS system does not apply to pregnant patients"/></div>
                  <input checked={pregnant} onChange={()=> handleInputChange()} className="checkbox" type="checkbox"/>
                </article>
            </section>
            <footer>
                <button className="empty" onClick={exit}>Cancel</button>
                <button className="navy" onClick={() =>setConfirm(true)}>Save changes</button>
            </footer>
        </div>: <Confirmation message="Patient details updated successfully" exit={exit}/>}
        <div className="background"></div>
        </>
    )
}

export default EditDetails
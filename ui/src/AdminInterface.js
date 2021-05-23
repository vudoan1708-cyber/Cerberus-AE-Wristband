import React, {useEffect, useState, useContext} from 'react';
import './App.css';

//Packages
import {gql, useMutation, useQuery, useSubscription } from '@apollo/client';
import { animated, useTransition } from 'react-spring';


//Components
import DateAndTime from "./components/DateAndTime";
import Info from './components/Info';
import Options from './components/options/Options';
import Confirmation from './components/Confirmation';
import WristbandCard from './components/WristbandCard';
import Searchbar from './components/Searchbar';

import { StoreContext } from './components/utils/store'

// svg
import ConnectedCare from "./assets/svg/connected-care.svg";
import Battery from "./assets/svg/battery.svg";
import BatteryLow from "./assets/svg/battery-low.svg";
import Wifi from "./assets/svg/wifi.svg";
import alert from './assets/svg/alert.svg';

//GraphQL
  
// Get the data for all wristbands
const ALL_WRISTBANDS = gql`
  query {
    getWristbands {
      id,
      name,
      tic,
      active,
      onOxygen,
      pregnant,
      child,
      department,
      key,
    }
  }
`

//Subscription to add on new wristbands
const ADD_NEW_WRISTBAND = gql`
    subscription {
      updateWristbandAdded {
        id
        tic
        active
        name
        dateOfBirth
        onOxygen
        pregnant
        child
        key
        department
      }
  }
`

//Subscription to generate Important Bands
const UPDATE_IMPORTANT_WRISTBANDS = gql`
  subscription {
    updateImportantBands {
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

const AdminInterface = () => {

//Hooks

  //useContext
  const [wristbandID, setWristbandID] = useContext(StoreContext);

  //useState
  const [important, setImportant] = useState(false);
  const [showDetail, setShowDetail] = useState(false);

  // Temporary Measure for showcase
  const [showcaseBatteryLevel, setShowcaseBatteryLevel] = useState(false);
  


  const {subscribeToMore, loading, error, data, refetch } = useQuery(ALL_WRISTBANDS);

  console.log(data);


  // Refetch the database on mount to prevent glitches
  useEffect(()=> {
    refetch();
  },[])

  //const newWristbandSub = useSubscription(ADD_NEW_WRISTBAND)
  //if (!newWristbandSub.error && !newWristbandSub.loading) console.log(newWristbandSub.data);
  subscribeToMore({
      document: ADD_NEW_WRISTBAND,
      updateQuery: (prev, { subscriptionData }) => {
        console.log(subscriptionData);
        if (!subscriptionData.data) return prev;
        console.log(prev);
        const newFeedItem = subscriptionData.data.updateWristbandAdded;
        console.log(newFeedItem);
        return Object.assign({}, prev, {
          getWristbands: prev.getWristbands ? [newFeedItem, ...prev.getWristbands] : [newFeedItem]
        });
        
      },
    });

  console.log(data);

  const [wristbands, setWristbands] = useState([])
  const [activeWristbands, setActiveWristbands] = useState();
  
  useEffect(()=> {
    let isMounted = true; // note this flag denote mount status
    let activeWristbandsArray = [];
    
    if (data) {
      if (data.getWristbands) {
        
        data.getWristbands.forEach(wristband => {
          console.log(wristband.name + ": " + wristband.active);
          if (wristband.active) {
            activeWristbandsArray.push(wristband);
            
            
          }
        })
      }
    
    if (isMounted) setActiveWristbands(activeWristbandsArray);
    };
    return () => { isMounted = false }; // use effect cleanup to set flag false, if unmounted
    
  }, [data])

  

  if (loading) return 'Loading...';
  if(error) return("Error in Wristbands: " + error.message);

    return (
        <div className="admin-interface">
          <div className="main">
          <header>
            <div>
              <img onClick={()=>refetch()} className="connected-care left" alt="ConnectedCare Logo" src={ConnectedCare}/>
              <h1 className="left">Dashboard</h1>
            </div>
            <div>
              <AddWristbandButton className="right"/>
              <DateAndTime/>
            </div> 
          </header>
        
          <header>
            <div>
              <h1 onClick={() => setImportant(true)} className={important ? "left double-underline": "left light-text"}>Important Bands</h1>
              <h1 onClick={() =>  setImportant(false)} className={!important ? "left double-underline": "left light-text"}>All Bands</h1>
            </div>
            <Searchbar setShowDetail={setShowDetail} setWristbandID={setWristbandID}  setWristbands={setWristbands} activeWristbands={activeWristbands} className="right"/>
          </header>
          <section className="grid-5">
            <Wristbands showcaseBatteryLevel={showcaseBatteryLevel} important={important} wristbands={wristbands} setWristbandID={setWristbandID} setShowDetail={setShowDetail}/>
          </section>
          </div>
          <PatientDetail refetch={refetch} setShowcaseBatteryLevel={setShowcaseBatteryLevel} wristbandID={wristbandID} showDetail={showDetail} setShowDetail={setShowDetail}/>
      </div>
    ) 
}


const PatientDetail = ({refetch, showDetail, setShowDetail, wristbandID, setShowcaseBatteryLevel}) => {

  //GraphQL


  console.log("setShowcaseBatteryLevel: " + setShowcaseBatteryLevel);
  //leave parent as multipleWristbandData 
  const MULTIPLE_WRISTBAND_DATA = gql `
    query getMultipleWristbandData ($id: ID!, $howMany: Int!, $start: Int!, $end: Int!) {

      getWristband (id: $id) {
        id,
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

     

  if (error) console.log(error);
  if (loading) console.log(loading);
  
  //Hooks
  const [display, setDisplay] = useState();

  useEffect(()=> {
    if (data && showDetail) {
      setDisplay(true)
    } else {
      setDisplay(false)
    }
  },[data, showDetail])

  //React-Spring
  const transitions = useTransition(display, null, {
    from: { opacity: 0, right: "-100%" },
    enter: { opacity: 1, right: "0%" },
    leave: { opacity: 0, right: "-100%" }
  });
  
  console.log("Patient Detail");

  // if (loading) return 'Loading...';
  // if(error) return("Error in  Patient Detail: " + error.message);
  let latestWristbandData = 0;
        if (data && showDetail) {
          console.log(data)
          latestWristbandData = data.getMultipleWristbandData[0];
        }
       
      return (
        <>
        {transitions.map(
          ({item, key, props}) =>
        item && (
          <animated.div key={key}>
          <animated.aside style={props} className="patient-detail">

            <header>
              <div>
                <h2 className="exit left" onClick={()=> setShowDetail(false)}>+</h2>
                <h1 className="left">Patient</h1>
              </div>
              <div></div>
            </header>

            <section>
              <article>
                <p>patient name</p>
                <h3>{data ? data.getWristband.name: "Err"}</h3>
              </article>
              <article>
                <p>date of birth</p>
                <h3>{data ? data.getWristband.dateOfBirth: "Err"}</h3>
              </article>
              <article>
                <p>NHS number</p>
                <h3>450 557 7104</h3>
              </article>
              <article>
                <p>Additional details</p>
                <h3>N/A</h3>
              </article>
            </section>
  
            <section>
              <article>
                <p>Band ID</p>
                <h3>AE7593834</h3>
              </article>
              <article>
                <p>System connection</p>
                <div><img alt="Connection Icon" src={Wifi}/><h3>Connection</h3></div>
              </article>
              <article>
                <p>Battery Level</p>
                <div>
                {data ? latestWristbandData.batteryLevel > 20 ? <img alt="Battery level" src={Battery}/>: <img alt="Battery level" src={BatteryLow}/>: "Err"}
                   <h3>{data ? latestWristbandData.batteryLevel: null}%</h3> 
                  {latestWristbandData && latestWristbandData.batteryLevel < 20 ? <img alt="Low Battery" src={alert}/>: null}
                </div>
              </article>
              <article>
                <p>Skin contact</p>
                <h3>{latestWristbandData ? latestWristbandData.sensorData.proximity ? "Yes": "No": "Err"}</h3>
              </article>
              <Options refetch={refetch} setShowcaseBatteryLevel={setShowcaseBatteryLevel} wristbandID={wristbandID} editDetails={true} reAssign={true} RemoveBand={true}/>
            </section>
  
            <section className="location">
              <hgroup>
                <h1>Location</h1>
                <p>Last logged: 30 seconds ago</p>
              </hgroup>
              <canvas>
                
              </canvas>
              <aside>
                <button><h1>+</h1></button>
                <button><h1>-</h1></button>
              </aside>
            </section>
          </animated.aside>
        <animated.div style={props.opacity} onClick={()=> setShowDetail(false)} className="background"></animated.div>
      </animated.div>
        )
        )}
      </>
      )
  }
  
  const AddWristbandButton = () => {
    
    //Hooks
      const [display, setDisplay] = useState(false);
      const [confirm, setConfirm] = useState(false)
    
    
    //React-Spring
    const transitions = useTransition(display, null, {
      from: { opacity: 0, top: "0%" },
      enter: { opacity: 1, top: "50%" },
      leave: { opacity: 0, top: "0%" }
    });


    
      return (
        <>
        {transitions.map(
          ({item, key, props}) =>
          item && (
            <NewPatient
            key={key}
            setDisplay={setDisplay}
            setConfirm={setConfirm}
            props={props}
            />
          )
        )}
          <button className="navy right" onClick={()=> setDisplay(true)}><h1 className="plus">+</h1><h3> New patient</h3></button>
          {confirm ? 
            <Confirmation exit={()=> setConfirm(false)} message="Patient band added succesfully"/>
          : null}
        </>
      )
    }

  const NewPatient = ({ setDisplay, setConfirm, props}) => {

    //useState
    const [name, setName] = useState('N/A');
    const [pregnant , setPregnant] = useState(false);
    const [child, setChild] = useState();
    const [dateOfBirth, setDateOfBirth] = useState(null);
    
    //GraphQL
    const ADD_NEW_WRISTBAND = gql`
    mutation addWristband ($input: AddWristbandInput!) {
      addWristband(input: $input) {
        id
      }
    }
  `

  const [addWristband] = useMutation(ADD_NEW_WRISTBAND);

  const handleClick = () => {
    if (isValidDate(dateOfBirth) === true) {// Check if the date is valid
        setDisplay(false);
        addWristband({
          variables: {
            input: {
              tic: "2",
              name: name,
              dateOfBirth: dateOfBirth,
              onOxygen: false,
              pregnant: pregnant,
              child: child,
              department: "waiting-room",
              // data: {
              //   news2 : {
              //     temperature: "7",
              //   }
              // }
            }
          }
        });
        setConfirm(true);
    } else {
      console.log("Date Format Incorrect: Please Enter a Valid Date of Birth");
    }
    
  }

  function isValidDate(dateString) {
    if (dateString) {
      console.log(dateString);
    var regEx = /^\d{4}-\d{2}-\d{2}$/;
    if(!dateString.match(regEx)) return false;  // Invalid format
    var d = new Date(dateString);
    var dNum = d.getTime();
    if(!dNum && dNum !== 0) return false; // NaN value, Invalid date
    return d.toISOString().slice(0,10) === dateString;
    }
  }

  const handleDate = (e) => {
    console.log("left focus");
    console.log(child);
    if (e) {
      console.log(e.target.value);
      setDateOfBirth(e.target.value);
      console.log();
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
    
    return (
      <>
      <animated.div style={props} className="pop-up new-patient">
        <header>
          <h1>New patient</h1>
          <h2 className="exit" onClick={()=> setDisplay(false)}>+</h2>
        </header>
        <section>
          <article>
            <p>Patient name</p>
            <input className="text-field" value={name} onChange={e=> (setName(e.target.value))} type="text"/>
            <div><p>Date of birth</p><Info message="The system will automatically calculate the age. NEWS system does not apply to patients under 16 years old "/></div>
            <input onBlur={(e)=> handleDate(e)} className="text-field" type="date"/>
            <p>NHS number</p>
              <input className="text-field" type="text"></input>
          </article>
          <article>
            <p>Band ID</p>
            <h2>AE784729</h2>
            <div style={{display:"flex", alignItems:"center"}}><p>Pregnancy</p><Info message="NEWS system does not apply to pregnant patients"/></div>
            <input checked={pregnant} onChange={()=> setPregnant(prevState => !prevState)} className="checkbox" type="checkbox"/>
          </article>
        </section>
        <footer>
          <button className="empty" onClick={()=> setDisplay(false)}>Cancel</button>
          <button className="navy" onClick={() =>handleClick()}>Add patient</button>
        </footer>
      </animated.div>
      <animated.div style={props.opacity} className="background"></animated.div>
      </>
    )
  }
  
  const Wristbands = ({ setShowDetail, setWristbandID, wristbands, important, showcaseBatteryLevel}) => {

    const transitions = useTransition(wristbands ? wristbands: [], wristband => wristband.id, {
      from: { opacity: 0, transform:"translate(0px, -50px)" },
      enter: { opacity: 1, transform:"translate(0px, 0px)" },
      leave: { opacity: 0, transform:"translate(0px, -50px)" },
    })

    //console.log("showcaseBatteryLevel: " + showcaseBatteryLevel);

      if (wristbands) {
      return (
        <>
           {transitions.map(({item, key, props})=> (
             item ? (
           <WristbandCard key={parseInt(item.id)}
              showcaseBatteryLevel={showcaseBatteryLevel}
              important={important}
              props={props}
              nurseView={false}
              setShowDetail={setShowDetail}
              setWristbandID={setWristbandID} 
              wristband={item}/> ) : null
            ))}
        </>
      )} else return <div>No Wristbands</div>
    
  }

  const ImportantWristbands = ({ setShowDetail, setWristbandID, important}) => {

    const { loading, error, data } = useSubscription(UPDATE_IMPORTANT_WRISTBANDS);

    const transitions = useTransition(data ? data.updateImportantBands: [], ImportantWristband => ImportantWristband.id, {
      from: { opacity: 0, transform:"translate(0px, -50px)" },
      enter: { opacity: 1, transform:"translate(0px, 0px)" },
      leave: { opacity: 0, transform:"translate(0px, -50px)" },
    })

      if (important) {
      return (
        <>
          {transitions.map(({item, key, props})=> (
            item ? (
          <WristbandCard key={parseInt(item.id)}
            props={props}
            nurseView={false}
            setShowDetail={setShowDetail}
            setWristbandID={setWristbandID} 
            wristband={item}/> ) : null
          ))}
        </>
      )} else return <div>No Important Wristbands</div>
    
  }
  

export default AdminInterface
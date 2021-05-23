import React, { useEffect, useRef, useState, useContext } from 'react';
import './App.css'

//React Context
import { StoreContext }  from './components/utils/store'


//Packages
import {gql, useQuery} from '@apollo/client';
import { useHistory } from 'react-router-dom';
import { animated,  useSpring } from 'react-spring';

//Components
import BarChart from './components/BarChart';
import LineGraph from './components/LineGraph';
import Options from './components/options/Options';

//Icons
import Battery from "./assets/svg/battery.svg";
import BatteryLow from "./assets/svg/battery-low.svg";
import Connection from './assets/svg/connection.svg'
//Location Map
import LocationMap from './assets/svg/location.svg' 

//GraphQL

  // Query to first get the multiple wristband data
  const MULTIPLE_WRISTBAND_DATA = gql `
    query multipleWristbandData ($id: ID!, $howMany: Int!, $start: Int!, $end: Int!) {

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
        pulse,
        temperature,
        bloodPressure
        motion,
        proximity
        },
        news2 {
          sp02,
          pulse,
          temperature,
          bloodPressure
          motion,
          onOxygen,
          overall
        },
        batteryLevel
      }
    }
  `
  // Subscription updates the Wirstband Data
  const UPDATE_WRISTBAND_DATA = gql `
  subscription updateWristbandData ($id: ID!) {
    updateWristbandData (id: $id) {
      sensorData {
        respiration,
        sp02,
        pulse,
        temperature,
        bloodPressure
        motion,
        proximity
      },
      news2 {
        sp02,
        pulse,
        temperature,
        bloodPressure
        motion,
        onOxygen,
        overall
      },
      batteryLevel
    }
  }
  `

const PatientCart = () => {


//Hooks

  //useContext
  const [wristbandID] = useContext(StoreContext)
  

  //UseState
  const [currentStat, setCurrentStat] = useState(0)
  const [cartWristband, setCartWristband] = useState();

  const props = useSpring({from: {transform:"translate(0px, -50px)", opacity: 0}, to:{transform:"translate(0px, 0px)", opacity: 1}})




  // Harcoded data for D3
  const chartData = [
    {year: 1980, efficiency: 24.3, sales: 8},
    {year: 1985, efficiency: 27.6, sales: 10},
    {year: 1990, efficiency: 28, sales: 9},
    {year: 1991, efficiency: 28.4, sales: 8},
    {year: 1992, efficiency: 27.9, sales: 8},
    {year: 1993, efficiency: 28.4, sales: 80},
    {year: 1994, efficiency: 28.3, sales: 89},
    {year: 1995, efficiency: 28.6, sales: 86},
    {year: 1996, efficiency: 28.5, sales: 8},
    {year: 1997, efficiency: 28.7, sales: 82},
    {year: 1998, efficiency: 28.8, sales: 8},
    {year: 1999, efficiency: 28.3, sales: 8},
    {year: 2000, efficiency: 28.5, sales: 7},
    {year: 2001, efficiency: 28.8, sales: 9},
    {year: 2002, efficiency: 29, sales: 80}
  ]


  useEffect(()=> {
    console.log("Patient Cart View: " + wristbandID);
  },[wristbandID])
    
  const {subscribeToMore, loading, error, data, refetch} = useQuery(MULTIPLE_WRISTBAND_DATA, {variables: {id: wristbandID, howMany: 1, start: 0, end: 0}})
  
  subscribeToMore({
    document: UPDATE_WRISTBAND_DATA,
    variables: { id: wristbandID },
    updateQuery: (prev, { subscriptionData }) => {
      console.log(subscriptionData);
      if (!subscriptionData.data) return prev;
      console.log(prev);
      const newFeedItem = subscriptionData.data.updateWristbandData;
      console.log(newFeedItem);
      console.log(prev.getMultipleWristbandData);
      return Object.assign({}, prev, {
        getMultipleWristbandData: [newFeedItem, ...prev.getMultipleWristbandData]
      });
      
    },
  });

  if (loading) return 'Loading...';
  if(error) return("Error in  Patient Cart: " + error.message);
  console.log(data);
  let overall = data.getMultipleWristbandData[data.getMultipleWristbandData.length-1].news2.overall;
  console.log(overall);
  console.log(data.getMultipleWristbandData.length-1);
  let sensorData = data.getMultipleWristbandData[data.getMultipleWristbandData.length-1].sensorData;
    if (data) return (
        <animated.div className="patient-cart" style={props} style={{display:"grid",gridTemplateColumns:"9fr 5fr", maxWidth:"100%"}}>
          <div className="main">
            <section className="grid-4">
              <article>
                <p>Patient Cart</p>
                <h3>{data.getWristband.name}</h3>
              </article>
              <article>
                <p>Date of birth</p>
                <h3>{data.getWristband.dateOfBirth}</h3>
              </article>
              <article>
                <p>NHS Number</p>
                <h3>450 557 7104</h3>
              </article>
              <article>
                <p>Addtional details</p>
                {additionalDetails(data)}
              </article>
            </section>

            <section className="grid-5">
              <article>
                <p>Band ID</p>
                <h3>AE78435</h3>
              </article>
              <article>
                <p>System Connection</p>
                <div><img alt="Connection Icon" src={Connection}/><h3>Connection</h3></div>
              </article>
              <article>
                <p>Battery Level</p>
                <div>{battery(data)}<h3>{data && data.getMultipleWristbandData ? data.getMultipleWristbandData[data.getMultipleWristbandData.length-1].batteryLevel : undefined}%</h3></div>
              </article>
              <article>
                <p>Band ID</p>
                <h3>AE78435</h3>
              </article>
                <Options patientCart={true} refetch={refetch} wristbandID={wristbandID} movePatient={true} editDetails={true} reAssign={true} RemoveBand={true}/>
            </section>
            <div style={{margin: "1em"}} className="grid-7 news">
              <OverallNEWSScore value={overall} onClick={()=>setCurrentStat(0)} currentStat={currentStat} />
              <Temperature value={data.getMultipleWristbandData[data.getMultipleWristbandData.length-1].sensorData.temperature} onClick={()=>setCurrentStat(1)} currentStat={currentStat}  />
              <PulseRate value={sensorData.pulse} onClick={()=>setCurrentStat(2)} currentStat={currentStat}  />
              <BloodOxygenSaturation value={sensorData.sp02} onClick={()=>setCurrentStat(3)} currentStat={currentStat}  />
              <RespirationRate value={sensorData.respiration} onClick={()=>setCurrentStat(4)} currentStat={currentStat}  />
              <BloodPressure value={sensorData.bloodPressure} onClick={()=>setCurrentStat(5)} currentStat={currentStat}  />
              <MovementDetection value={sensorData.motion} onClick={()=>setCurrentStat(6)} currentStat={currentStat}  />
            </div>

            <section>
              <header>
                <h1>NEWS score</h1>
                <div>
                  <button className="right"><h2>&lt;</h2></button>
                  <button className="right"><h2>&gt;</h2></button>
                </div>
              </header>
              
              { currentStat < 1 ? <BarChart data={chartData}/> : <LineGraph data={chartData} />}
              
              
              
            </section>

        </div>
        <Location />
        </animated.div>
    )
}


const Location = ({ }) => {

  let location = useRef();


  useEffect(()=> {
    if (location.current) {

    }
  },[location])
  
  return (
    <aside className="main patient-cart location" >
      <section>
        <header>
          <h1>Location</h1>
        </header>
        <p>Last logged: 30 seconds ago</p>
        <figure ref={location}>
        <img src={LocationMap} alt="Patient Location"/>
        <aside>
          <button><h1>+</h1></button>
          <button><h1>-</h1></button>
        </aside> 
        </figure>
      </section>
    </aside> 
  )
}


//Javascript Formatting Functions

const additionalDetails = (data) => {

  if (data.getWristband.child) {
    return (
      <div style={{display: "flex", alignItems: "center"}}>
      <h3>Under 16</h3>
      <p className="exclamation">!</p>
      </div>
    )
  } else {
    if (data.getWristband.pregnant) {
      return (
        <div style={{display: "flex", alignItems: "center"}}>
          <h3>Pregnant</h3>
          <p className="exclamation">!</p>
        </div>
        )
    } else {
      return (
        <h3>N/A</h3>
        )
    }
  }
}

const battery = (data) => {
  let batteryLevel;
  if (data) {
    batteryLevel = data.getMultipleWristbandData[data.getMultipleWristbandData.length-1].batteryLevel
  }
  
  if (batteryLevel) {
    if (batteryLevel > 30) {
      return <img alt="battery level" src={Battery}/>
    } else return <img alt="battery level" src={BatteryLow}/>
  } else return "pending"
}

// News Scores

const OverallNEWSScore = ({value, currentStat, onClick}) => {
  
  const [colour, setColour] = useState('black');

  useEffect(()=> {

    console.log(value)
    if (340 >= value) {
      console.log("smaller than 35");
        setColour("var(--red)");
        console.log(colour)
    } else if (345 >= value && value > 340) {
        setColour("var(--yellow)");
        console.log(colour)
    } else if (380 >= value && value > 345) {
        setColour("var(--grey-low)");
        console.log(colour)
    } else if (420 >= value && value > 380) {
        setColour("var(--yellow)");
        console.log(colour)
    } else if (value > 420) {
      console.log("bigger than 39");
        setColour("var(--orange)");
        console.log(colour)
    }
  }, [value])
  
  
  return (
    <figure 
    style={{boxShadow: currentStat === 0 ? "0px 0px 7px rgba(0, 1, 73, 0.63)": null,
     WebkitBoxShadow: currentStat === 0 ? "0px 0px 7px rgba(0, 1, 73, 0.63)": null}} 
     onClick={onClick}>
      <h2>NEWS score</h2>
      <h1 style={{color:colour}}>{value !== "undefined" ? value: "Err"}</h1>
    </figure>
  )

}

const Temperature = ({value, currentStat, onClick}) => {
  
  const [colour, setColour] = useState('black')

  useEffect(()=> {
    if (35 >= value) {
      console.log("smaller than 35");
        setColour("var(--red)")
        console.log(colour)
    } else if (36 >= value && value > 35) {
        setColour("var(--yellow)");
        console.log(colour)
    } else if (38 >= value && value > 36) {
        setColour("var(--grey-low)");
        console.log(colour)
    } else if (39 >= value && value > 38) {
        setColour("var(--yellow)");
        console.log(colour)
    } else if (value > 39) {
      console.log("bigger than 39");
        setColour("var(--orange)");
        console.log(colour)
    }
  }, [value])

  return (
    <figure 
    style={{boxShadow: currentStat === 1 ? "0px 0px 7px rgba(0, 1, 73, 0.63)": null,
     WebkitBoxShadow: currentStat === 1 ? "0px 0px 7px rgba(0, 1, 73, 0.63)": null}} 
     onClick={onClick}>
      <h2>Temperature (&deg;C)</h2>
      <h1 style={{color:colour}}>{value}</h1>
    </figure>
  )

}

const PulseRate = ({value, currentStat, onClick}) => {
  
  const [colour, setColour] = useState('black')

  useEffect(()=> {
    if (40 >= value) {
        setColour("var(--red)")
    } else if (50 >= value && value > 40) {
        setColour("var(--yellow)")
    } else if (90 >= value && value > 50) {
        setColour("var(--grey-low)")
    } else if (110 >= value && value > 90) {
        setColour("var(--yellow)")
    } else if (130 >= value && value > 110) {
      setColour("var(--orange)")
    }
     else if (value > 130) {
        setColour("var(--red)")
    }
  }, [value])
  
  return (
    <figure 
    style={{boxShadow: currentStat === 2 ? "0px 0px 7px rgba(0, 1, 73, 0.63)": null,
     WebkitBoxShadow: currentStat === 2 ? "0px 0px 7px rgba(0, 1, 73, 0.63)": null}} 
     onClick={onClick}>
      <h2>Pulse rate</h2>
      <h1 style={{color:colour}}>{value}</h1>
    </figure>
  )

}

const BloodOxygenSaturation = ({value, currentStat, onClick}) => {

  const [colour, setColour] = useState('black')

  useEffect(()=> {
    if (83 >= value) {
        setColour("var(--red)")
    } else if (85 >= value && value > 83) {
        setColour("var(--orange)")
    } else if (87 >= value && value > 85) {
        setColour("var(--yellow)")
    } else if (92 >= value && value > 87) {
        setColour("var(--grey-low)")
    } else if (94 >= value && value > 92) {
      setColour("var(--yellow)")
    } else if (96 >= value && value > 94) {
      setColour("var(--orange)")
    }
     else if (value > 96) {
        setColour("var(--red)")
    }
  }, [value])
  
  return (
    <figure 
    style={{boxShadow: currentStat === 3 ? "0px 0px 7px rgba(0, 1, 73, 0.63)": null,
     WebkitBoxShadow: currentStat === 3 ? "0px 0px 7px rgba(0, 1, 73, 0.63)": null}} 
     onClick={onClick}>
      <h2>Blood oxygen saturation (SPO<sup>2</sup>)</h2>
      <h1 style={{color:colour}}>{value}</h1>
    </figure>
  )

}

const RespirationRate = ({value, currentStat, onClick}) => {
  
  const [colour, setColour] = useState('black')

  useEffect(()=> {
    if (8 >= value) {
        setColour("var(--red)")
    } else if (11 >= value && value > 8) {
        setColour("var(--yellow)")
    } else if (20 >= value && value > 11) {
        setColour("var(--grey-low)")
    } else if (24 >= value && value > 20) {
        setColour("var(--orange)")
    } else if (value > 25) {
        setColour("var(--red)")
    }
  }, [value])
  
  return (
    <figure 
    style={{boxShadow: currentStat === 4 ? "0px 0px 7px rgba(0, 1, 73, 0.63)": null,
     WebkitBoxShadow: currentStat === 4 ? "0px 0px 7px rgba(0, 1, 73, 0.63)": null}} 
     onClick={onClick}>
      <h2>Respiration rate (p/m)</h2>
      <h1 style={{color:colour}}>{value}</h1>
    </figure>
  )

}

const BloodPressure = ({value, currentStat, onClick}) => {
  
  const [colour, setColour] = useState('black')

  useEffect(()=> {
    if (90 >= value) {
        setColour("var(--red)")
    } else if (100 >= value && value > 90) {
        setColour("var(--yellow)")
    } else if (110 >= value && value > 100) {
        setColour("var(--grey-low)")
    } else if (219 >= value && value > 110) {
        setColour("var(--orange)")
    } else if (value > 219) {
        setColour("var(--red)")
    }
  }, [value])
  
  return (
    <figure 
    style={{boxShadow: currentStat === 5 ? "0px 0px 7px rgba(0, 1, 73, 0.63)": null,
     WebkitBoxShadow: currentStat === 5 ? "0px 0px 7px rgba(0, 1, 73, 0.63)": null}} 
     onClick={onClick}>
      <h2>Systolic blood pressure (mmHg)</h2>
      <h1 style={{color:colour}}>{value}</h1>
    </figure>
  )

}

const MovementDetection = ({value, currentStat, onClick}) => {
  
  const [motion, setMotion] = useState()

  useEffect(()=> {
    if (value > 0) {
        setMotion(true);
    } else {
      setMotion(false);
    }
  }, [value])
  
  return (
    <figure 
    style={{boxShadow: currentStat === 6 ? "0px 0px 7px rgba(0, 1, 73, 0.63)": null,
     WebkitBoxShadow: currentStat === 6 ? "0px 0px 7px rgba(0, 1, 73, 0.63)": null}} 
     onClick={onClick}>
      <h2>Movement detection</h2>
      <h1>{motion ? "Yes" : "No"}</h1>
    </figure>
  )

}


export default PatientCart
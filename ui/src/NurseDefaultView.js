import React, {useContext, useEffect, useState} from 'react';

//Packages
import {useQuery, gql, useSubscription} from '@apollo/client';
import {  useTransition } from 'react-spring';

//Components
import WristbandCard from './components/WristbandCard';
import Searchbar from './components/Searchbar';
import MovePatient from './components/options/MovePatient';

//React Context
import { StoreContext } from './components/utils/store'


//GraphQL

  // Get the data for all wristbands
  const ALL_WRISTBANDS = gql`
    query {
      getWristbands {
        id,
        tic,
        active,
        name,
        dateOfBirth,
        onOxygen,
        pregnant,
        child,
        key,
        department,
      }
    }
  `

  //Subscription to add on new wristbands
  const ADD_NEW_WRISTBAND = gql`
  subscription {
    updateWristbandAdded {
      id,
      tic,
      active,
      name,
      dateOfBirth,
      onOxygen,
      pregnant,
      child,
      key,
      department,
    }
  }
  `

  //Subscription to update and populate the Summary
  const UPDATE_SUMMARY_DATA = gql`
    subscription updateSummary {
      updateSummary {
        high,
        medium,
        lowMedium,
        low,
        other
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

const NurseDefaultView = () => {

//Hooks

  //useContext
   const [wristbandID, setWristbandID] = useContext(StoreContext)

   console.log(wristbandID);

  //UI
  const [important, setImportant] = useState(false);

  //Wristbands
  const [wristbands, setWristbands] = useState()
  const [activeWristbands, setActiveWristbands] = useState();

  const {subscribeToMore, loading, error, data, refetch } = useQuery(ALL_WRISTBANDS);

  // Refetch the database on mount to prevent glitches
  useEffect(()=> {
    refetch();
  },[])

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

  // Filter out inactive wristbands and wristband in departments other that Waiting room, then send the data to the searchbar
  useEffect(()=> {
    let activeWristbandsArray = [];
    
    if (data) {
      if (data.getWristbands) {
        data.getWristbands.forEach(wristband => {
          if (wristband.active) {
            activeWristbandsArray.push(wristband);

            
          }
        })
      }

    setActiveWristbands(activeWristbandsArray);
    };

    
  }, [data])

    if (loading) return 'Loading...';
    if(error) return("Error in Wristbands: " + error.message);

    return (
    <>
        <header onClick={()=>refetch()}>
            <Summary/>
            <Searchbar nurseView setWristbandID={setWristbandID} activeWristbands={activeWristbands} setWristbands={setWristbands}/>
        </header>

        <header>
        <div>
              <h1 onClick={() => setImportant(true)} className={important ? "left double-underline": "left light-text"}>Important Bands</h1>
              <h1 onClick={() =>  setImportant(false)} className={!important ? "left double-underline": "left light-text"}>All Bands</h1>
            </div>
        </header>

        <section className="grid-4">
          
            <Wristbands refetch={refetch} important={important} wristbands={wristbands} setWristbandID={setWristbandID}/>
        </section>
</>
)
}

const Wristbands = ({ refetch, important, setWristbandID,  wristbands }) => {



  //React-Spring
  const transitions = useTransition(wristbands ? wristbands: [], wristband => wristband.id, {
    from: { opacity: 0, transform:"translate(0px, -50px)" },
    enter: { opacity: 1, transform:"translate(0px, 0px)" },
    leave: { opacity: 0, transform:"translate(0px, -50px)" },
  })
  

    if (wristbands) {
    return (
      <>
         {transitions.map(({item, key, props})=> (
           item ? (
         <WristbandCard key={parseInt(item.id)}
            refetch={refetch}
            important={important}
            props={props}
            nurseView={true}
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
            nurseView={true}
            setShowDetail={setShowDetail}
            setWristbandID={setWristbandID} 
            wristband={item}/> ) : null
          ))}
      </>
    )} else return <div>No Important Wristbands</div>
  
}

const Summary = ({}) => {

  
  const { error, loading, data } = useSubscription(UPDATE_SUMMARY_DATA);

  return (
    <div className="summary">
                <figure>
                    <h3>High Score</h3>
                    <h1 style={{color: "var(--red)"}}>{data  ? data.updateSummary.high : "N/A"}</h1>
                </figure>
                <figure>
                    <h3>Medium Score</h3>
                    <h1 style={{color: "var(--orange)"}}>{data  ? data.updateSummary.medium : "N/A"}</h1>
                </figure>
                <figure>
                    <h3>Low-mid Score</h3>
                    <h1 style={{color: "var(--yellow)"}}>{data  ? data.updateSummary.lowMedium : "N/A"}</h1>
                </figure>
                <figure>
                    <h3>Low Score</h3>
                    <h1 style={{color: "var(--grey-low)"}}>{data ? data.updateSummary.low : "N/A"}</h1>
                </figure>
                <figure>
                    <h3>Other</h3>
                    <h1 style={{color: "black"}}>{data ? data.updateSummary.other : "N/A"}</h1>
                </figure>
            </div>
  )
}


export default NurseDefaultView
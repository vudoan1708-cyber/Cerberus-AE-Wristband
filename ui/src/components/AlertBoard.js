import React, { useState, useContext, useEffect } from 'react';

import { StoreContext } from './utils/store'

//Packages
import {useSubscription, useQuery, gql} from '@apollo/client';

//Components
import Alert from './Alert';



//GraphQL

    //Query retreives all active wristbands
    const GET_ALL_BANDS = gql`
    query {
        getWristbands {
            id,
            name,
            pregnant,
            child,
            active,
            department,
        }
    }
    `

    //Subscription to add on new wristbands, to check for alerts
    const ADD_NEW_WRISTBAND = gql`
    subscription {
    updateWristbandAdded {
        id,
        name,
        pregnant,
        child,
        active,
        department,
        }
    }
    `

const AlertBoard = () => {

    useEffect(()=>{
        console.log("refresh");
    },[])

    //useContext
    const [wristbandID, setWristbandID] = useContext(StoreContext);

    //useState
    const [alertFilter, setAlertFilter] = useState("all")

    //Apollo
    const { subscribeToMore, loading, error, data } = useQuery(GET_ALL_BANDS);

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
    
    // console.log(loading);
    // console.log(error);
    
    return (
        <aside className="alerts-container">
            <h1>Vitals alerts</h1>
            <nav>
                <a className={alertFilter === "all" ? "active": null} onClick={()=> setAlertFilter("all")}>All alerts</a>
                <a className={alertFilter === "high" ? "active": null} onClick={()=> setAlertFilter("high")}>High level</a>
                <a className={alertFilter === "medium" ? "active": null} onClick={()=> setAlertFilter("medium")}>Medium level</a>
                <a className={alertFilter === "low-medium" ? "active": null} onClick={()=> setAlertFilter("low-medium")}>Low-Medium level</a>
                <a className={alertFilter === "other" ? "active": null} onClick={()=> setAlertFilter("other")}>Other</a>
            </nav>
            {data && data.getWristbands ?
                data.getWristbands.map(( alert, index )=> (
                <Alert setWristbandID={setWristbandID} alertFilter={alertFilter} key={index} name={alert.name} id={alert.id} child={alert.child} pregnant={alert.pregnant} active={alert.active} department={alert.department}/>
                )) : null
            }
        </aside>
    )
}


export default AlertBoard
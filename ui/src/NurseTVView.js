import React, {useEffect, useState, useContext} from 'react';
import {useQuery, gql} from '@apollo/client';

//Components
import TVAlert from './components/TVAlert'

//GraphQL

  // Get the data for all wristbands
  const GET_ALL_BANDS = gql`
    query {
        getWristbands {
          id,
          active,
          name,
          pregnant,
          child,
          department,
        }
    }
  `
  //Subscription to add on new wristbands
    const ADD_NEW_WRISTBAND = gql`
    subscription {
      updateWristbandAdded {
          id,
          active,
          name,
          pregnant,
          child,
          department,
      }
    }
    `


const NurseTVView = () => {

//Hooks

  const {loading, error, data, subscribeToMore } = useQuery(GET_ALL_BANDS)// this refetch is temporary, to have to TV/Column View update

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

    if (loading) return 'Loading...';
    if(error) return("Error in Wristbands: " + error.message);

    return (
        <div className="tv-view">
            <div className="main">
                <section className="grid-5">{/*Temporary refetch to update*/}
                    <article>
                        <header style={{borderBottom:"4px solid var(--red)"}}><h2>High</h2></header>
                        {data && data.getWristbands ?
                          data.getWristbands.map(( alert, index )=> (
                            <TVAlert TVFilter="high" key={index}  name={alert.name} id={alert.id} child={alert.child} pregnant={alert.pregnant} active={alert.active} department={alert.department}/>
                        )) : null
                        }
                    </article>
                    <article>
                        <header style={{borderBottom:"4px solid var(--orange)"}}><h2>Medium</h2></header>
                        {data && data.getWristbands ?
                          data.getWristbands.map(( alert, index )=> (
                            <TVAlert TVFilter="medium" key={index}  name={alert.name} id={alert.id} child={alert.child} pregnant={alert.pregnant} active={alert.active} department={alert.department}/>
                        )) : null
                        }
                    </article>
                    <article>
                        <header style={{borderBottom:"4px solid var(--yellow)"}}><h2>Low-Medium</h2></header>
                        {data && data.getWristbands ?
                          data.getWristbands.map(( alert, index )=> {
                            console.log(alert)
                            return (
                            <TVAlert TVFilter="low-medium" key={index}  name={alert.name} id={alert.id} child={alert.child} pregnant={alert.pregnant} active={alert.active} department={alert.department}/>
                        )}) : null
                        }
                    </article>
                    <article>
                        <header style={{borderBottom:"4px double var(--dark-navy)"}}><h2>Error</h2></header>
                        {data && data.getWristbands ?
                          data.getWristbands.map(( alert, index )=> (
                            <TVAlert TVFilter="error" key={index}  name={alert.name} id={alert.id} child={alert.child} pregnant={alert.pregnant} active={alert.active} department={alert.department}/>
                        )) : null
                        }
                    </article>
                    <article>
                        <header style={{borderBottom:"4px solid black"}}><h2>Special</h2></header>
                        {data && data.getWristbands ?
                          data.getWristbands.map(( alert, index )=> (
                            <TVAlert TVFilter="special" key={index}  name={alert.name} id={alert.id} child={alert.child} pregnant={alert.pregnant} active={alert.active} department={alert.department}/>
                        )) : null
                        }
                    </article>
                </section>
            </div>
        </div>
    )
}



export default NurseTVView
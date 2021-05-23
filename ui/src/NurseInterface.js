import React, {useEffect, useState} from 'react';
import { Redirect, Link, Switch, Route, useRouteMatch, useLocation } from 'react-router-dom';
import './App.css';

import {useSubscription, useQuery, gql} from '@apollo/client';

//Components
import DateAndTime from "./components/DateAndTime";

//Icons
import ConnectedCare from "./assets/svg/connected-care.svg";
import GridView from "./assets/svg/grid-view.svg";
import GridViewNavy from "./assets/svg/grid-view-navy.svg";
import RowView from "./assets/svg/row-view.svg";
import RowViewWhite from './assets/svg/row-view-white.svg'

//Pages
import PatientCart from './PatientCart';
import NurseDefaultView from './NurseDefaultView';
import NurseTVView from './NurseTVView';
import AlertBoard from './components/AlertBoard';

const NurseInterface = () => {
    let location = useLocation();
    const {path, url} = useRouteMatch();

    const [gridView, setGridView] = useState(true);

    const [patientView, setPatientView] = useState(false);

    useEffect(()=> {
        console.log("refresh");
        if (location.pathname === "/nurse/patient-cart") {
            setPatientView(true);
        } else {
            setPatientView(false);
        }

    }, [location])

    return (
        <div className="nurse-interface">
            <div className="main">
                <header>
                    <div>
                    <img className="connected-care left" alt="ConnectedCare Logo" src={ConnectedCare}/>
                        <hgroup className="left">
                        <Link to={`${url}/default-view`}>
                            <h1 className= {patientView ? "light-text" : ""}>Dashboard</h1>
                            <p>Waiting Room</p>
                        </Link>
                        </hgroup>
                        {patientView ? 
                        <>
                            <h1 className="left light-text">&gt;</h1>
                            <h1 className="left">Patient cart</h1>
                        </> : null}
                    </div>

                    <div>
                        <h3 className="right">Change view</h3>
                        <Link className="right" to={`${url}/default-view`}><button className={ location.pathname === `/nurse/tv-view`  ? "empty ": "navy"}><img src={location.pathname === `/nurse/tv-view` ? GridViewNavy : GridView} alt="grid view" /></button></Link>
                        <Link className="right" to={`${url}/tv-view`}><button className={ location.pathname === `/nurse/tv-view`  ? "navy": "empty"}><img src={location.pathname === `/nurse/tv-view` ? RowViewWhite : RowView} alt="row view" /></button></Link>
                        <DateAndTime/>
                    </div>
                </header>
                <Switch>
                    <Route
                     exact path={`${path}/default-view`} 
                     render={(props)=> (
                         <NurseDefaultView {...props}/>// Alternative way of Routing Components, allowing you to pass props to it 
                     )}
                    />
                    <Route
                     exact path={`${path}/tv-view`} 
                     render={(props)=> (
                         <NurseTVView {...props} gridView={gridView}/>// Alternative way of Routing Components, allowing you to pass props to it 
                     )}
                    />
                    <Route
                     exact path={`${path}/patient-cart`} 
                     render={(props)=> (
                         <PatientCart {...props} />// Alternative way of Routing Components, allowing you to pass props to it 
                     )}
                    />

                    <Redirect exact from={`${path}`} to="/nurse/default-view" />
                </Switch>
            </div>
                {location.pathname !== `/nurse/tv-view` ? <AlertBoard location={location}/> : null}
        </div>
    )
}



export default NurseInterface
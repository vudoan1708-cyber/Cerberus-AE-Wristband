import React, { useState } from 'react';
import { animated, useTransition } from 'react-spring'
import './App.css';

//Libraries
import { HashRouter as Router, Switch, Route, Redirect, useLocation } from 'react-router-dom';

//Components
import AdminInterface from './AdminInterface';
import NurseInterface from './NurseInterface';


function App() {

  return (
    <div className="App">
      <Router>
        <Location/>
      </Router>
    </div>
  );
    
  
}

const Location = ()=> {

  const location = useLocation()

  const transitions = useTransition(location, location => location.pathname, { from: { position:"absolute"  , top:0, left:0, opacity: 0 },
      enter: { position:"relative", opacity: 1 },
      leave: { position:"absolute", top:0, left:0, opacity: 0 },
  })
  return transitions.map(({ item, props, key }) => (
    <animated.div key={key} style={props}>
    <Switch location={item}>
      <Route path={`/nurse`} component={NurseInterface}/>
      <Route path={`/admin`} component={AdminInterface}/>
      <Redirect exact from="/" to="/admin" />
    </Switch>
  </animated.div>
  ))
}




export default App;

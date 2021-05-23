import React from 'react';

//Packages
import {animated, useSpring } from 'react-spring'

const Confirmation = ({message, exit}) => {

    return (
        <>
            <animated.div className="pop-up confirm">
                <header>
                  <div></div>
                  <h2 className="exit" onClick={exit}>+</h2>
                </header>
                <div>
                  <h2>{message}</h2>
                </div>
              </animated.div>
            <animated.div className="background"></animated.div>
            </>
    )
}

export default Confirmation
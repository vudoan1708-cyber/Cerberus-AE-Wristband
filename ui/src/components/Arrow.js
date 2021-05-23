import React from 'react';

//Icon
import EmptyCircle from '../assets/svg/empty-circle.svg'

const Arrow = ({level, trend, child, pregnant}) => {

    let rotation = "rotate(45deg)";
    
    // console.log("rising " + trend);
    // console.log("level: " + level);

    if (trend === "rising") {
        rotation = "rotate(0deg)"
    } else if (trend === "falling") {
        rotation = "rotate(180deg)"
    }
    
    if (level == "error" || trend == "unchanged" || trend == "none") {
        return (// The empty circle
                <img src={EmptyCircle} alt="none"/>
        )
    } else {
        return (// The arrow
        <svg style={{transform: rotation}} width="12" height="13" viewBox="0 0 12 13" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M0.0205962 5.7934C0.0628321 5.89465 0.16139 5.96072 0.270762 5.96072L3.23755 5.96072L3.23755 12.7293C3.23755 12.8787 3.35885 13 3.50831 13L7.8402 13C7.98965 13 8.11095 12.8787 8.11095 12.7293L8.11095 5.9607L11.0891 5.9607C11.1985 5.9607 11.297 5.89463 11.3393 5.79391C11.381 5.69266 11.3582 5.57623 11.2808 5.49879L5.87948 0.0796037C5.82859 0.0287125 5.7598 5.23247e-06 5.68779 5.22618e-06C5.61579 5.21988e-06 5.547 0.0287124 5.49611 0.0790705L0.0790765 5.49826C0.00163553 5.57573 -0.0216397 5.69213 0.0205962 5.7934Z" fill={colour(level, child, pregnant )}/>
        </svg> 
        )
    }
    
}

const colour = (level, child, pregnant) => {

    if (child || pregnant) {// first of all, check if the patient is pregnant or under 16, as NEWS2 doesnt apply to them
        return "black"
    } else {
        if (level === "high") {
            return "var(--red)"
        } else if (level === "medium") {
            return "var(--orange)"
        } else if (level === "low-medium") {
            return "var(--yellow)"
        } else if (level === "low") {
            return "var(--grey-low)"
        } else if (level === "error") {
            return "black"
        }
    }
    
}

export default Arrow
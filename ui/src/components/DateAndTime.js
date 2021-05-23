import React, {useState, useEffect} from 'react';

const DateAndTime = () => {

  const [date, setDate] = useState("N/A")

  const GetDateAndTime = () => {

      let date = new Date();
    
      let today, hours, minutes, day;
      
      today = date.getDate();
      hours = date.getHours();
      minutes = date.getMinutes();
      if (minutes.toString().length < 2) {
        minutes = (minutes+"").padStart(2,"0");
      }
      let dayNumber = date.getDay();
      day = 0;
      
      switch (dayNumber) {
        case 0:
          day = "Mon"
          break;
        case 1:
          day = "Tue"
          break;
        case 2:
          day = "Wed"
          break;
        case 3:
          day = "Thurs"
          break;
        case 4:
          day = "Fri"
          break;
        case 5:
          day = "Sat"
          break;
        case 6:
          day = "Sun"
          break;
        default:
          day = "Unknown"
      }
      setDate(`${day} ${today} - ${hours}:${minutes}`)
    }

      useEffect (() => {
        const interval = setInterval(GetDateAndTime
        , 1000);

        return () => clearInterval(interval);
      }, [])

    return (
        <>
        <h3 id="date" className="right">{date}</h3>
        </>
    )
}

export default DateAndTime
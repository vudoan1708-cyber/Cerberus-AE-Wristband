import React, { useEffect, useState } from 'react'

//Packages
import { Link } from 'react-router-dom';

//Components
import Search from "../assets/svg/search.svg";



const Searchbar = ({ nurseView, activeWristbands, setWristbands, setShowDetail, setWristbandID })=> {
    //followed this tutorial https://hackernoon.com/how-to-build-a-search-bar-in-react-with-react-hooks-o77l3yl7

    const [word, setWord] = useState("");
    const [filteredWristbands, setFilteredWristbands] = useState(activeWristbands);

    useEffect(()=> {// Set the inital Value of wristbands on Component Mount when activeWristbands are calculated
      if (activeWristbands) {
      setWristbands(activeWristbands);
      }

    },[activeWristbands])

    const handleChange = e => {
      console.log(e);
      console.log(activeWristbands);
      let oldList = activeWristbands.map((wristband, index) => {

        return {
          name: activeWristbands[index].name.toLowerCase(),
          id: activeWristbands[index].id,
        }
      })
      console.log(oldList);
      
      if(e !== "") {
        let newList = [];

        setWord(e);

        newList = oldList.filter(wristband => wristband.name.includes(word.toLowerCase()));
        console.log(newList);
        setFilteredWristbands(newList);
      } else {
        setWord("");
        setFilteredWristbands(null)// If the Searchbar is empty, then just set Filter Display to all the entries
      }
    }

    
    useEffect(()=> {
      // check if there are any searched Wristbands
      let wristbandsArray = [];

      if (filteredWristbands) {
        console.log("Render Cards")
        
        
        activeWristbands.forEach(wristband => {
          console.log("Filtered Wristbands:" + "length: " + filteredWristbands.length)
          
          for(let i = 0; i < filteredWristbands.length; i++) {
            console.log("for loop");
            console.log(wristband.id);
            console.log(i);
            console.log(filteredWristbands[i].id);
            if(wristband.id === filteredWristbands[i].id) {
              console.log(wristband.name);
              wristbandsArray.push(wristband);

            }
          }
        })
        setWristbands(wristbandsArray);
        console.log(wristbandsArray);

      } else { // check if there are any searched Wristbands, if not, return the value of the default Active Wristbands
        setWristbands(activeWristbands)
      }
    }, [filteredWristbands])

    const handleClick = (id) => {
      if (nurseView) {
        //Patient Detail doesnt exist in the Nurse View, so the function setShowDetail wont
        setWristbandID(id)
      } else {
        setShowDetail(true)
        setWristbandID(id)
      }
      
    }

    return (
      <div className="searchbar">
        <div className="text-field">
          <img alt="Calendar Icon by Evil Icons" src={Search}/>
          <input value={word} onChange={e => handleChange(e.target.value)} className="text-field" type="text" placeholder="Search"/>
        </div>
        <aside>
        {filteredWristbands ? filteredWristbands.map((result, id)=> {
          if (nurseView) return(
            <Link key={id} to={`patient-cart`}>
              <article onClick={()=>handleClick(result.id)} >
                <h2>{result.name}</h2>
              </article>
            </Link>
            
          ); else return (
            <article  key={id} onClick={()=>handleClick(result.id)} >
              <h2>{result.name}</h2>
            </article>
          )
          }): null}
        </aside>
      </div>
    )
  }

export default Searchbar
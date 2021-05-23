import React, { useState } from 'react'

const Info = ({message})=> {

    const [messageBox, setMessageBox] = useState(false)

    return (
        <div className="info">
        {messageBox ? <aside><p>{message}</p></aside>: null}
        <p onMouseOver={()=> setMessageBox(true)} onMouseLeave={()=> setMessageBox(false)}>i</p>
        </div>
    )
}

export default Info
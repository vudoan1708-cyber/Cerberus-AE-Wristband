import React, { createContext, useState } from 'react';

export const StoreContext = createContext(null)

export default ({ children }) => {
    const wristbandID = useState();

    return <StoreContext.Provider value={wristbandID}>{children}</StoreContext.Provider>
}
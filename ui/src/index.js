import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';

//useContext
import StoreProvider from './components/utils/store'

//GraphQL
import {ApolloClient, ApolloProvider, split, HttpLink, InMemoryCache} from '@apollo/client';
import { getMainDefinition } from '@apollo/client/utilities';
import { WebSocketLink } from '@apollo/client/link/ws';

//Components
import App from './App';
import reportWebVitals from './reportWebVitals';







// Link for HTTP Requests
const httpLink = new HttpLink({
  uri: 'http://localhost:8080/api'
});

// Link for Websocket Links
const wsLink = new WebSocketLink({
  uri: 'ws://localhost:8080/api',
  options: {
    reconnect: true,
  }
});

// Split Function takes the operation to execute, and reuturns the Websocket Link or HTTP Link depending on a boolean value
const splitLink = split(
  ({ query }) => {
    const definition = getMainDefinition(query);
    return (
      definition.kind === 'OperationDefinition' &&
      definition.operation === 'subscription'
    );
  },
  wsLink,
  httpLink,
);

const client = new ApolloClient({
  connectToDevTools: true,
  cache: new InMemoryCache(),
  link: splitLink,
})

ReactDOM.render(
  <React.StrictMode>
    <ApolloProvider client={client}>
    <StoreProvider>
        <App/>
      </StoreProvider>
    </ApolloProvider>
  </React.StrictMode>,
  document.getElementById('root')
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();

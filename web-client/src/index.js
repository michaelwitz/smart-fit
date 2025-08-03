import React from 'react';
import ReactDOM from 'react-dom/client';
import { MantineProvider } from '@mantine/core';
import { Notifications } from '@mantine/notifications';
import App from './App';
import '@mantine/core/styles.css';
import '@mantine/notifications/styles.css';
import './index.css';

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <MantineProvider
      theme={{
        colorScheme: 'dark',
        primaryColor: 'pink',
        defaultRadius: 'md',
        fontFamily: 'Inter, system-ui, sans-serif',
      }}
    >
      <Notifications />
      <App />
    </MantineProvider>
  </React.StrictMode>
);

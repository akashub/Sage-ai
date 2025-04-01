// utils/StateDebugger.js
import React, { useEffect } from 'react';

// A higher-order component that wraps your component and logs state changes
export const withStateDebugger = (WrappedComponent, debugName = 'Component') => {
  return function WithStateDebugger(props) {
    // Log when the component renders
    console.log(`[${debugName}] Rendering with props:`, props);
    
    // Effect to log when the component mounts, updates, or unmounts
    useEffect(() => {
      console.log(`[${debugName}] Component mounted`);
      
      return () => {
        console.log(`[${debugName}] Component unmounting`);
      };
    }, []);
    
    // For each important prop change, add a separate effect
    useEffect(() => {
      console.log(`[${debugName}] selectedChat prop changed:`, props.selectedChat);
    }, [props.selectedChat]);
    
    // Return the wrapped component with original props
    return <WrappedComponent {...props} />;
  };
};

// This hook can be added to individual components to debug state
export const useStateDebugger = (stateName, value, deps = []) => {
  useEffect(() => {
    console.log(`State '${stateName}' changed:`, value);
  }, [value, stateName, ...deps]);
  
  return value;
};

export default withStateDebugger;
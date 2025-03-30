"use client"

import { useState } from "react";
import { Box, CssBaseline, ThemeProvider, useMediaQuery } from "@mui/material";
import { createTheme } from "@mui/material/styles";
import ChatSidebar from "../components/chat/ChatSidebar";
import ChatWindow from "../components/chat/ChatWindow";

const ChatInterface = () => {
  const [selectedChat, setSelectedChat] = useState(null);
  const prefersDarkMode = useMediaQuery('(prefers-color-scheme: dark)');
  
  const theme = createTheme({
    palette: {
      mode: "dark",
      primary: {
        main: "#5865F2",
      },
      background: {
        default: "#202225",
        paper: "#2F3136",
      },
    },
    typography: {
      fontFamily: "'Roboto', sans-serif",
    },
    components: {
      MuiButton: {
        styleOverrides: {
          root: {
            textTransform: 'none',
          },
        },
      },
    },
  });

  // Define sidebar width
  const drawerWidth = 300;

  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Box sx={{ display: "flex", height: "100vh", width: "100vw", overflow: "hidden" }}>
        {/* First, the sidebar component */}
        <ChatSidebar selectedChat={selectedChat} setSelectedChat={setSelectedChat} />
        
        {/* Main content area - explicitly set to take remaining width */}
        <Box 
          component="main" 
          sx={{ 
            flexGrow: 1, 
            width: { sm: `calc(100% - ${drawerWidth}px)` },
            ml: { sm: `${drawerWidth}px` }
          }}
        >
          <ChatWindow selectedChat={selectedChat} />
        </Box>
      </Box>
    </ThemeProvider>
  );
};

export default ChatInterface;
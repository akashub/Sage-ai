"use client"

import { useState } from "react";
import { Box, CssBaseline, ThemeProvider } from "@mui/material";
import { createTheme } from "@mui/material/styles";
import ChatSidebar from "../components/chat/ChatSidebar";
import ChatWindow from "../components/chat/ChatWindow";

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
});

const ChatInterface = () => {
  const [selectedChat, setSelectedChat] = useState(null);

  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Box sx={{ display: "flex", height: "100vh", width: "100vw", overflow: "hidden" }}>
        <ChatSidebar selectedChat={selectedChat} setSelectedChat={setSelectedChat} />
        <ChatWindow selectedChat={selectedChat} />
      </Box>
    </ThemeProvider>
  );
};

export default ChatInterface;

// "use client"

// import { useState } from "react";
// import { Box, CssBaseline, ThemeProvider, useMediaQuery } from "@mui/material";
// import { createTheme } from "@mui/material/styles";
// import ChatSidebar from "../components/chat/ChatSidebar";
// import ChatWindow from "../components/chat/ChatWindow";

// const ChatInterface = () => {
//   const [selectedChat, setSelectedChat] = useState(null);
//   const prefersDarkMode = useMediaQuery('(prefers-color-scheme: dark)');
  
//   const theme = createTheme({
//     palette: {
//       mode: "dark",
//       primary: {
//         main: "#5865F2",
//       },
//       background: {
//         default: "#202225",
//         paper: "#2F3136",
//       },
//     },
//     typography: {
//       fontFamily: "'Roboto', sans-serif",
//     },
//     components: {
//       MuiButton: {
//         styleOverrides: {
//           root: {
//             textTransform: 'none',
//           },
//         },
//       },
//     },
//   });

//   // Define sidebar width
//   const drawerWidth = 300;

//   return (
//     <ThemeProvider theme={theme}>
//       <CssBaseline />
//       <Box sx={{ display: "flex", height: "100vh", width: "100vw", overflow: "hidden" }}>
//         {/* First, the sidebar component */}
//         <ChatSidebar selectedChat={selectedChat} setSelectedChat={setSelectedChat} />
        
//         {/* Main content area - explicitly set to take remaining width */}
//         <Box 
//           component="main" 
//           sx={{ 
//             flexGrow: 1, 
//             width: { sm: `calc(100% - ${drawerWidth}px)` },
//             ml: { sm: `${drawerWidth}px` }
//           }}
//         >
//           <ChatWindow selectedChat={selectedChat} />
//         </Box>
//       </Box>
//     </ThemeProvider>
//   );
// };

// export default ChatInterface;

import { useState, useEffect } from "react";
import { Box, CssBaseline, ThemeProvider, useMediaQuery } from "@mui/material";
import { createTheme } from "@mui/material/styles";
import ChatSidebar from "../components/chat/ChatSidebar";
import ChatWindow from "../components/chat/ChatWindow";
import { createChat, fetchChatById } from "../utils/api";

const ChatInterface = () => {
  const [selectedChat, setSelectedChat] = useState(null);
  const [sidebarOpen, setSidebarOpen] = useState(true);
  const [loading, setLoading] = useState(false);
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
  const collapsedDrawerWidth = 60;
  
  // Handle new chat creation
  const handleNewChat = async () => {
    try {
      setLoading(true);
      const newChat = await createChat();
      setSelectedChat(newChat);
    } catch (error) {
      console.error("Error creating new chat:", error);
    } finally {
      setLoading(false);
    }
  };
  
  // Handle sidebar toggle
  const handleSidebarToggle = () => {
    setSidebarOpen(!sidebarOpen);
  };
  
  // Load initial chat or create a new one
  useEffect(() => {
    const loadInitialChat = async () => {
      try {
        // Check URL for chat ID parameter
        const urlParams = new URLSearchParams(window.location.search);
        const chatId = urlParams.get('chat');
        
        if (chatId) {
          // Try to load the specified chat
          setLoading(true);
          try {
            const chat = await fetchChatById(chatId);
            setSelectedChat(chat);
          } catch (error) {
            console.error("Error loading chat from URL parameter:", error);
            // If failed, create a new chat
            await handleNewChat();
          }
        } else {
          // No chat specified, create a new one
          await handleNewChat();
        }
      } catch (error) {
        console.error("Error in initial chat setup:", error);
      } finally {
        setLoading(false);
      }
    };
    
    loadInitialChat();
  }, []);

  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Box sx={{ display: "flex", height: "100vh", width: "100vw", overflow: "hidden" }}>
        {/* Sidebar component with toggle capability */}
        <ChatSidebar 
          selectedChat={selectedChat} 
          setSelectedChat={setSelectedChat} 
          onNewChat={handleNewChat}
          sidebarOpen={sidebarOpen}
          onSidebarToggle={handleSidebarToggle}
        />
        
        {/* Main content area - dynamically adjust width based on sidebar state */}
        <Box 
          component="main" 
          sx={{ 
            flexGrow: 1, 
            width: { 
              sm: `calc(100% - ${sidebarOpen ? drawerWidth : collapsedDrawerWidth}px)` 
            },
            ml: { 
              sm: `${sidebarOpen ? drawerWidth : collapsedDrawerWidth}px` 
            },
            transition: theme => theme.transitions.create(['margin', 'width'], {
              easing: theme.transitions.easing.sharp,
              duration: theme.transitions.duration.leavingScreen,
            }),
          }}
        >
          <ChatWindow 
            selectedChat={selectedChat}
            loading={loading}
          />
        </Box>
      </Box>
    </ThemeProvider>
  );
};

export default ChatInterface;
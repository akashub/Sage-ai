// // "use client"

// // import { useState } from "react";
// // import { Box, CssBaseline, ThemeProvider, useMediaQuery } from "@mui/material";
// // import { createTheme } from "@mui/material/styles";
// // import ChatSidebar from "../components/chat/ChatSidebar";
// // import ChatWindow from "../components/chat/ChatWindow";

// // const ChatInterface = () => {
// //   const [selectedChat, setSelectedChat] = useState(null);
// //   const prefersDarkMode = useMediaQuery('(prefers-color-scheme: dark)');
  
// //   const theme = createTheme({
// //     palette: {
// //       mode: "dark",
// //       primary: {
// //         main: "#5865F2",
// //       },
// //       background: {
// //         default: "#202225",
// //         paper: "#2F3136",
// //       },
// //     },
// //     typography: {
// //       fontFamily: "'Roboto', sans-serif",
// //     },
// //     components: {
// //       MuiButton: {
// //         styleOverrides: {
// //           root: {
// //             textTransform: 'none',
// //           },
// //         },
// //       },
// //     },
// //   });

// //   // Define sidebar width
// //   const drawerWidth = 300;

// //   return (
// //     <ThemeProvider theme={theme}>
// //       <CssBaseline />
// //       <Box sx={{ display: "flex", height: "100vh", width: "100vw", overflow: "hidden" }}>
// //         {/* First, the sidebar component */}
// //         <ChatSidebar selectedChat={selectedChat} setSelectedChat={setSelectedChat} />
        
// //         {/* Main content area - explicitly set to take remaining width */}
// //         <Box 
// //           component="main" 
// //           sx={{ 
// //             flexGrow: 1, 
// //             width: { sm: `calc(100% - ${drawerWidth}px)` },
// //             ml: { sm: `${drawerWidth}px` }
// //           }}
// //         >
// //           <ChatWindow selectedChat={selectedChat} />
// //         </Box>
// //       </Box>
// //     </ThemeProvider>
// //   );
// // };

// // export default ChatInterface;

// // import { useState, useEffect } from "react";
// // import { Box, CssBaseline, ThemeProvider, useMediaQuery } from "@mui/material";
// // import { createTheme } from "@mui/material/styles";
// // import ChatSidebar from "../components/chat/ChatSidebar";
// // import ChatWindow from "../components/chat/ChatWindow";
// // import { createChat, fetchChatById } from "../utils/api";

// // const ChatInterface = () => {
// //   const [selectedChat, setSelectedChat] = useState(null);
// //   const [sidebarOpen, setSidebarOpen] = useState(true);
// //   const [loading, setLoading] = useState(false);
// //   const prefersDarkMode = useMediaQuery('(prefers-color-scheme: dark)');
  
// //   const theme = createTheme({
// //     palette: {
// //       mode: "dark",
// //       primary: {
// //         main: "#5865F2",
// //       },
// //       background: {
// //         default: "#202225",
// //         paper: "#2F3136",
// //       },
// //     },
// //     typography: {
// //       fontFamily: "'Roboto', sans-serif",
// //     },
// //     components: {
// //       MuiButton: {
// //         styleOverrides: {
// //           root: {
// //             textTransform: 'none',
// //           },
// //         },
// //       },
// //     },
// //   });

// //   // Define sidebar width
// //   const drawerWidth = 300;
// //   const collapsedDrawerWidth = 60;
  
// //   // Handle new chat creation
// //   const handleNewChat = async () => {
// //     try {
// //       setLoading(true);
// //       const newChat = await createChat();
// //       setSelectedChat(newChat);
// //     } catch (error) {
// //       console.error("Error creating new chat:", error);
// //     } finally {
// //       setLoading(false);
// //     }
// //   };
  
// //   // Handle sidebar toggle
// //   const handleSidebarToggle = () => {
// //     setSidebarOpen(!sidebarOpen);
// //   };
  
// //   // Load initial chat or create a new one
// //   useEffect(() => {
// //     const loadInitialChat = async () => {
// //       try {
// //         // Check URL for chat ID parameter
// //         const urlParams = new URLSearchParams(window.location.search);
// //         const chatId = urlParams.get('chat');
        
// //         if (chatId) {
// //           // Try to load the specified chat
// //           setLoading(true);
// //           try {
// //             const chat = await fetchChatById(chatId);
// //             setSelectedChat(chat);
// //           } catch (error) {
// //             console.error("Error loading chat from URL parameter:", error);
// //             // If failed, create a new chat
// //             await handleNewChat();
// //           }
// //         } else {
// //           // No chat specified, create a new one
// //           await handleNewChat();
// //         }
// //       } catch (error) {
// //         console.error("Error in initial chat setup:", error);
// //       } finally {
// //         setLoading(false);
// //       }
// //     };
    
// //     loadInitialChat();
// //   }, []);

// //   return (
// //     <ThemeProvider theme={theme}>
// //       <CssBaseline />
// //       <Box sx={{ display: "flex", height: "100vh", width: "100vw", overflow: "hidden" }}>
// //         {/* Sidebar component with toggle capability */}
// //         <ChatSidebar 
// //           selectedChat={selectedChat} 
// //           setSelectedChat={setSelectedChat} 
// //           onNewChat={handleNewChat}
// //           sidebarOpen={sidebarOpen}
// //           onSidebarToggle={handleSidebarToggle}
// //         />
        
// //         {/* Main content area - dynamically adjust width based on sidebar state */}
// //         <Box 
// //           component="main" 
// //           sx={{ 
// //             flexGrow: 1, 
// //             width: { 
// //               sm: `calc(100% - ${sidebarOpen ? drawerWidth : collapsedDrawerWidth}px)` 
// //             },
// //             ml: { 
// //               sm: `${sidebarOpen ? drawerWidth : collapsedDrawerWidth}px` 
// //             },
// //             transition: theme => theme.transitions.create(['margin', 'width'], {
// //               easing: theme.transitions.easing.sharp,
// //               duration: theme.transitions.duration.leavingScreen,
// //             }),
// //           }}
// //         >
// //           <ChatWindow 
// //             selectedChat={selectedChat}
// //             loading={loading}
// //           />
// //         </Box>
// //       </Box>
// //     </ThemeProvider>
// //   );
// // };

// // export default ChatInterface;

// import { useState, useEffect, useRef } from "react";
// import { Box, CssBaseline, ThemeProvider, useMediaQuery, CircularProgress } from "@mui/material";
// import { createTheme } from "@mui/material/styles";
// import ChatSidebar from "../components/chat/ChatSidebar";
// import ChatWindow from "../components/chat/ChatWindow";
// import { createChat, fetchChatById } from "../utils/api";

// // Create a utility function to help with debugging
// const debugLog = (message, ...data) => {
//   console.log(`[ChatInterface] ${message}`, ...data);
// };

// const ChatInterface = () => {
//   // Debug logging for component
//   debugLog("Component rendering");
  
  
//   const [selectedChat, setSelectedChat] = useState(null);
//   const [sidebarOpen, setSidebarOpen] = useState(true);
//   const [loading, setLoading] = useState(false);
//   const [initComplete, setInitComplete] = useState(false);
//   const initialLoadAttempted = useRef(false);
  
//   const prefersDarkMode = useMediaQuery('(prefers-color-scheme: dark)');
  
//   // Create theme
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
//   const collapsedDrawerWidth = 60;
  
//   // Debug log state changes
//   useEffect(() => {
//     debugLog("selectedChat changed:", selectedChat);
//   }, [selectedChat]);
  
//   useEffect(() => {
//     debugLog("loading changed:", loading);
//   }, [loading]);
  
//   // Handle new chat creation
//   const handleNewChat = async () => {
//     try {
//       debugLog("Creating new chat...");
//       setLoading(true);
      
//       // Create a new chat
//       const newChat = await createChat();
//       debugLog("New chat created:", newChat);
      
//       // Important: Set to null first to force reset
//       setSelectedChat(null);
      
//       // Wait a moment to ensure state is updated
//       setTimeout(() => {
//         setSelectedChat(newChat);
//       }, 50);
      
//     } catch (error) {
//       console.error("Error creating new chat:", error);
//     } finally {
//       setLoading(false);
//     }
//   };
  
//   // Handle sidebar toggle
//   const handleSidebarToggle = () => {
//     setSidebarOpen(!sidebarOpen);
//   };
  
//   // Load initial chat or create a new one
//   useEffect(() => {
//     const loadInitialChat = async () => {
//       // Prevent multiple attempts
//       if (initialLoadAttempted.current) {
//         return;
//       }
      
//       initialLoadAttempted.current = true;
      
//       try {
//         debugLog("Attempting to load initial chat");
        
//         // Check URL for chat ID parameter
//         const urlParams = new URLSearchParams(window.location.search);
//         const chatId = urlParams.get('chat');
        
//         if (chatId) {
//           // Try to load the specified chat
//           debugLog("Chat ID found in URL:", chatId);
//           setLoading(true);
//           try {
//             const chat = await fetchChatById(chatId);
//             debugLog("Fetched chat by ID:", chat);
//             setSelectedChat(chat);
//           } catch (error) {
//             console.error("Error loading chat from URL parameter:", error);
//             // If failed, create a new chat
//             await handleNewChat();
//           }
//         } else {
//           // No chat specified, create a new one
//           debugLog("No chat ID in URL, creating new chat");
//           await handleNewChat();
//         }
        
//         setInitComplete(true);
//       } catch (error) {
//         console.error("Error in initial chat setup:", error);
//         setInitComplete(true);
//       } finally {
//         setLoading(false);
//       }
//     };
    
//     loadInitialChat();
//   }, []);  // Empty dependency array ensures this runs only once

//   return (
//     <ThemeProvider theme={theme}>
//       <CssBaseline />
//       {!initComplete ? (
//         <Box 
//           sx={{ 
//             display: 'flex', 
//             justifyContent: 'center', 
//             alignItems: 'center', 
//             height: '100vh' 
//           }}
//         >
//           <CircularProgress />
//         </Box>
//       ) : (
//         <Box sx={{ display: "flex", height: "100vh", width: "100vw", overflow: "hidden" }}>
//           {/* Sidebar component with toggle capability */}
//           <ChatSidebar 
//             selectedChat={selectedChat} 
//             setSelectedChat={setSelectedChat} 
//             onNewChat={handleNewChat}
//             sidebarOpen={sidebarOpen}
//             onSidebarToggle={handleSidebarToggle}
//           />
          
//           {/* Main content area - dynamically adjust width based on sidebar state */}
//           <Box 
//             component="main" 
//             sx={{ 
//               flexGrow: 1, 
//               width: { 
//                 sm: `calc(100% - ${sidebarOpen ? drawerWidth : collapsedDrawerWidth}px)` 
//               },
//               ml: { 
//                 sm: `${sidebarOpen ? drawerWidth : collapsedDrawerWidth}px` 
//               },
//               transition: theme => theme.transitions.create (['margin', 'width'], {
//                 easing: theme.transitions.easing.sharp,
//                 duration: theme.transitions.duration.leavingScreen,
//               }),
//             }}
//           >
//             {/* Use the ChatWindow component with a key to force remount when selectedChat changes */}
//             <ChatWindow 
//               key={selectedChat ? selectedChat.id : 'no-chat'} // Force remount when chat changes
//               selectedChat={selectedChat}
//               loading={loading}
//             />
//           </Box>
//         </Box>
//       )}
//     </ThemeProvider>
//   );
// };

// export default ChatInterface;

import React, { useState, useEffect, useRef } from "react";
import { Box, CssBaseline, ThemeProvider, useMediaQuery, CircularProgress } from "@mui/material";
import { createTheme } from "@mui/material/styles";
import ChatSidebar from "../components/chat/ChatSidebar";
import ChatWindow from "../components/chat/ChatWindow";
import { useChat } from "../context/ChatContext";

// Create a utility function to help with debugging
const debugLog = (message, ...data) => {
  console.log(`[ChatInterface] ${message}`, ...data);
};

const ChatInterface = () => {
  // Debug logging for component
  debugLog("Component rendering");
  
  // Use the chat context
  const { 
    currentChat, 
    loading, 
    loadChat, 
    newChat 
  } = useChat();
  
  const [sidebarOpen, setSidebarOpen] = useState(true);
  const [initComplete, setInitComplete] = useState(false);
  const initialLoadAttempted = useRef(false);
  
  const prefersDarkMode = useMediaQuery('(prefers-color-scheme: dark)');
  
  // Create theme
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
  
  // Debug log state changes
  useEffect(() => {
    debugLog("currentChat changed:", currentChat);
  }, [currentChat]);
  
  useEffect(() => {
    debugLog("loading changed:", loading);
  }, [loading]);
  
  useEffect(() => {
    debugLog("sidebarOpen changed:", sidebarOpen);
  }, [sidebarOpen]);
  
  // Handle new chat creation
  const handleNewChat = async () => {
    try {
      debugLog("Creating new chat...");
      newChat(); // Use the context function
    } catch (error) {
      console.error("Error creating new chat:", error);
    }
  };
  
  // Handle sidebar toggle
  const handleSidebarToggle = () => {
    debugLog(`Toggling sidebar from ${sidebarOpen} to ${!sidebarOpen}`);
    setSidebarOpen(!sidebarOpen);
  };
  
  // Handle selecting a chat
  const handleSelectChat = (chat) => {
    if (chat?.id) {
      loadChat(chat.id);
    }
  };
  
  // Load initial chat or create a new one
  useEffect(() => {
    const loadInitialChat = async () => {
      // Prevent multiple attempts
      if (initialLoadAttempted.current) {
        return;
      }
      
      initialLoadAttempted.current = true;
      
      try {
        debugLog("Attempting to load initial chat");
        
        // Check URL for chat ID parameter
        const urlParams = new URLSearchParams(window.location.search);
        const chatId = urlParams.get('chat');
        
        if (chatId) {
          // Try to load the specified chat
          debugLog("Chat ID found in URL:", chatId);
          loadChat(chatId);
        } else {
          // No chat specified, create a new one
          debugLog("No chat ID in URL, creating new chat");
          await handleNewChat();
        }
        
        setInitComplete(true);
      } catch (error) {
        console.error("Error in initial chat setup:", error);
        setInitComplete(true);
      }
    };
    
    loadInitialChat();
  }, [loadChat, newChat]);  // Include dependencies

  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      {!initComplete || loading ? (
        <Box 
          sx={{ 
            display: 'flex', 
            justifyContent: 'center', 
            alignItems: 'center', 
            height: '100vh' 
          }}
        >
          <CircularProgress />
        </Box>
      ) : (
        <Box sx={{ display: "flex", height: "100vh", width: "100vw", overflow: "hidden" }}>
          {/* Sidebar component with toggle capability */}
          <ChatSidebar 
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
            {/* Use the ChatWindow component with a key to force remount when selectedChat changes */}
            <ChatWindow 
              key={currentChat ? currentChat.id : 'no-chat'} // Force remount when chat changes
              selectedChat={currentChat}
              loading={loading}
            />
          </Box>
        </Box>
      )}
    </ThemeProvider>
  );
};

export default ChatInterface;
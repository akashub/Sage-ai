// import { Box, Drawer, List, ListItem, ListItemButton, ListItemText, Typography, IconButton } from "@mui/material";
// import { Add as AddIcon } from "@mui/icons-material";

// const drawerWidth = 300;

// const ChatSidebar = ({ selectedChat, setSelectedChat }) => {
//   const chatHistory = [
//     { id: 1, title: "SQL Query Optimization" },
//     { id: 2, title: "Database Schema Design" },
//     { id: 3, title: "Data Analysis Query" },
//   ];

//   return (
//     <Drawer
//       variant="permanent"
//       sx={{
//         width: drawerWidth,
//         flexShrink: 0,
//         "& .MuiDrawer-paper": {
//           width: drawerWidth,
//           boxSizing: "border-box",
//           backgroundColor: "#2F3136",
//           color: "white",
//         },
//       }}
//     >
//       <Box sx={{ p: 2, display: "flex", alignItems: "center", justifyContent: "space-between" }}>
//         <Typography variant="h6">Chat History</Typography>
//         <IconButton color="primary" aria-label="new chat">
//           <AddIcon />
//         </IconButton>
//       </Box>
//       <List>
//         {chatHistory.map((chat) => (
//           <ListItem key={chat.id} disablePadding>
//             <ListItemButton selected={selectedChat === chat.id} onClick={() => setSelectedChat(chat.id)}>
//               <ListItemText primary={chat.title} />
//             </ListItemButton>
//           </ListItem>
//         ))}
//       </List>
//     </Drawer>
//   );
// };

// export default ChatSidebar;

import { useState, useEffect } from "react";
import { 
  Box, 
  Drawer, 
  List, 
  ListItem, 
  ListItemButton, 
  ListItemIcon, 
  ListItemText, 
  Typography, 
  IconButton, 
  Divider,
  Button,
  useMediaQuery,
  useTheme,
  Tooltip
} from "@mui/material";
import { 
  Add as AddIcon, 
  History as HistoryIcon, 
  Folder as FolderIcon,
  TableChart as TableChartIcon,
  Menu as MenuIcon,
  Close as CloseIcon,
  Home as HomeIcon,
  ChevronLeft as ChevronLeftIcon,
  ChevronRight as ChevronRightIcon
} from "@mui/icons-material";
import { Link } from "react-router-dom";

const drawerWidth = 300;
const collapsedDrawerWidth = 65;

const ChatSidebar = ({ selectedChat, setSelectedChat }) => {
  const [chatHistory, setChatHistory] = useState([
    { id: 1, title: "Sales Analysis 2024", timestamp: "2025-03-01", file: "sales_2024.csv" },
    { id: 2, title: "Customer Segmentation", timestamp: "2025-02-28", file: "customers.csv" },
    { id: 3, title: "Movies Database Query", timestamp: "2025-02-25", file: "movies.csv" },
  ]);
  
  const [mobileOpen, setMobileOpen] = useState(false);
  const [isCollapsed, setIsCollapsed] = useState(false);
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down("md"));

  const handleDrawerToggle = () => {
    setMobileOpen(!mobileOpen);
  };

  const handleCollapseToggle = () => {
    setIsCollapsed(!isCollapsed);
  };

  const handleNewChat = () => {
    setSelectedChat(null);
    if (isMobile) {
      setMobileOpen(false);
    }
  };

  const handleSelectChat = (chat) => {
    setSelectedChat(chat);
    if (isMobile) {
      setMobileOpen(false);
    }
  };

  const fetchChatHistory = async () => {
    // This would be replaced with actual API call in production
    // const response = await fetch('/api/chat-history');
    // const data = await response.json();
    // setChatHistory(data);
  };

  useEffect(() => {
    fetchChatHistory();
  }, []);

  const drawerContent = (
    <>
      <Box sx={{ 
        p: 2, 
        display: "flex", 
        flexDirection: "column",
        height: "100%"
      }}>
        {!isCollapsed && (
          <>
            <Button
              variant="outlined"
              startIcon={<AddIcon />}
              onClick={handleNewChat}
              sx={{ 
                mb: 2,
                borderColor: 'rgba(255,255,255,0.2)',
                color: 'white',
                justifyContent: "flex-start",
                '&:hover': {
                  borderColor: 'white',
                  backgroundColor: 'rgba(255,255,255,0.1)'
                }
              }}
            >
              New chat
            </Button>

            <Box sx={{ mb: 4 }}>
              <Typography 
                variant="overline" 
                sx={{ 
                  color: 'rgba(255,255,255,0.7)',
                  fontSize: '0.7rem',
                  fontWeight: 500,
                  letterSpacing: '0.1em',
                  pl: 1
                }}
              >
                Yesterday
              </Typography>
              <List dense>
                {chatHistory.slice(0, 2).map((chat) => (
                  <ListItem 
                    disablePadding 
                    key={chat.id}
                  >
                    <ListItemButton
                      selected={selectedChat?.id === chat.id}
                      onClick={() => handleSelectChat(chat)}
                      sx={{ 
                        borderRadius: 1,
                        py: 0.5,
                        '&.Mui-selected': {
                          backgroundColor: 'rgba(255,255,255,0.1)',
                        },
                        '&:hover': {
                          backgroundColor: 'rgba(255,255,255,0.05)',
                        }
                      }}
                    >
                      <ListItemIcon sx={{ minWidth: 36 }}>
                        <TableChartIcon fontSize="small" sx={{ color: 'rgba(255,255,255,0.7)' }} />
                      </ListItemIcon>
                      <ListItemText 
                        primary={chat.title}
                        primaryTypographyProps={{ 
                          noWrap: true,
                          fontSize: '0.875rem'
                        }}
                      />
                    </ListItemButton>
                  </ListItem>
                ))}
              </List>

              <Typography 
                variant="overline" 
                sx={{ 
                  color: 'rgba(255,255,255,0.7)',
                  fontSize: '0.7rem',
                  fontWeight: 500,
                  letterSpacing: '0.1em',
                  pl: 1,
                  mt: 2,
                  display: 'block'
                }}
              >
                Previous 30 Days
              </Typography>
              <List dense>
                {chatHistory.slice(2).map((chat) => (
                  <ListItem 
                    disablePadding 
                    key={chat.id}
                  >
                    <ListItemButton
                      selected={selectedChat?.id === chat.id}
                      onClick={() => handleSelectChat(chat)}
                      sx={{ 
                        borderRadius: 1,
                        py: 0.5,
                        '&.Mui-selected': {
                          backgroundColor: 'rgba(255,255,255,0.1)',
                        },
                        '&:hover': {
                          backgroundColor: 'rgba(255,255,255,0.05)',
                        }
                      }}
                    >
                      <ListItemIcon sx={{ minWidth: 36 }}>
                        <TableChartIcon fontSize="small" sx={{ color: 'rgba(255,255,255,0.7)' }} />
                      </ListItemIcon>
                      <ListItemText 
                        primary={chat.title}
                        primaryTypographyProps={{ 
                          noWrap: true,
                          fontSize: '0.875rem'
                        }}
                      />
                    </ListItemButton>
                  </ListItem>
                ))}
              </List>
            </Box>

            <Box sx={{ mt: 'auto' }}>
              <Button
                component={Link}
                to="/"
                variant="text"
                fullWidth
                startIcon={<HomeIcon />}
                sx={{ 
                  justifyContent: "flex-start", 
                  textTransform: "none",
                  color: 'rgba(255,255,255,0.8)',
                  '&:hover': {
                    backgroundColor: 'rgba(255,255,255,0.1)',
                    color: 'white'
                  }
                }}
              >
                Back to Home
              </Button>
            </Box>
          </>
        )}
      </Box>
    </>
  );

  return (
    <>
      {/* Mobile hamburger menu button */}
      {isMobile && (
        <IconButton
          color="inherit"
          aria-label="open drawer"
          edge="start"
          onClick={handleDrawerToggle}
          sx={{ 
            position: "absolute", 
            top: 10, 
            left: 10, 
            zIndex: 1100, 
            backgroundColor: "background.paper",
            boxShadow: 2
          }}
        >
          <MenuIcon />
        </IconButton>
      )}

      {/* Sidebar drawer - responsive behavior */}
      <Box
        component="nav"
        sx={{ 
          width: { md: isCollapsed ? collapsedDrawerWidth : drawerWidth }, 
          flexShrink: { md: 0 },
          transition: theme.transitions.create(['width'], {
            easing: theme.transitions.easing.sharp,
            duration: theme.transitions.duration.enteringScreen,
          }),
        }}
      >
        {isMobile ? (
          <Drawer
            variant="temporary"
            open={mobileOpen}
            onClose={handleDrawerToggle}
            ModalProps={{ keepMounted: true }}
            sx={{
              "& .MuiDrawer-paper": { 
                width: drawerWidth,
                boxSizing: "border-box",
                backgroundColor: "#2F3136"
              }
            }}
          >
            {drawerContent}
          </Drawer>
        ) : (
          <Drawer
            variant="permanent"
            sx={{
              width: isCollapsed ? collapsedDrawerWidth : drawerWidth,
              flexShrink: 0,
              "& .MuiDrawer-paper": {
                width: isCollapsed ? collapsedDrawerWidth : drawerWidth,
                boxSizing: "border-box",
                backgroundColor: "#2F3136",
                color: "white",
                transition: theme.transitions.create(['width'], {
                  easing: theme.transitions.easing.sharp,
                  duration: theme.transitions.duration.enteringScreen,
                }),
              },
            }}
            open
          >
            {drawerContent}
          </Drawer>
        )}
      </Box>
    </>
  );
};

export default ChatSidebar;
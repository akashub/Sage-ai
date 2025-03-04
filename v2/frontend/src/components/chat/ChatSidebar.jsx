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
  useTheme
} from "@mui/material";
import { 
  Add as AddIcon, 
  History as HistoryIcon, 
  Folder as FolderIcon,
  TableChart as TableChartIcon,
  Menu as MenuIcon,
  Close as CloseIcon,
  Home as HomeIcon
} from "@mui/icons-material";
import { Link } from "react-router-dom";

const drawerWidth = 300;

const ChatSidebar = ({ selectedChat, setSelectedChat }) => {
  // Mock chat history - in a real app, you'd fetch this from your backend
  const [chatHistory, setChatHistory] = useState([
    { id: 1, title: "Sales Analysis 2024", timestamp: "2025-03-01", file: "sales_2024.csv" },
    { id: 2, title: "Customer Segmentation", timestamp: "2025-02-28", file: "customers.csv" },
    { id: 3, title: "Movies Database Query", timestamp: "2025-02-25", file: "movies.csv" },
  ]);
  
  const [mobileOpen, setMobileOpen] = useState(false);
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down("md"));

  const handleDrawerToggle = () => {
    setMobileOpen(!mobileOpen);
  };

  const handleNewChat = () => {
    // Logic to start a new chat
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

  // Add a function to fetch chat history from backend
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
      <Box sx={{ p: 2, display: "flex", alignItems: "center", justifyContent: "space-between" }}>
        <Typography variant="h6" sx={{ fontWeight: "bold" }}>Sage Chat</Typography>
        {isMobile && (
          <IconButton onClick={handleDrawerToggle}>
            <CloseIcon />
          </IconButton>
        )}
      </Box>

      <Button
        variant="contained"
        startIcon={<AddIcon />}
        fullWidth
        sx={{ mx: 2, mb: 2 }}
        onClick={handleNewChat}
      >
        New Analysis
      </Button>

      <Box sx={{ px: 2, mb: 1 }}>
        <Button
          component={Link}
          to="/"
          variant="text"
          fullWidth
          startIcon={<HomeIcon />}
          sx={{ justifyContent: "flex-start", textTransform: "none" }}
        >
          Back to Home
        </Button>
      </Box>

      <Divider />
      
      <Typography variant="subtitle2" sx={{ px: 2, py: 1, color: "text.secondary" }}>
        <HistoryIcon fontSize="small" sx={{ verticalAlign: "middle", mr: 1 }} />
        Recent Sessions
      </Typography>
      
      <List>
        {chatHistory.map((chat) => (
          <ListItem 
            disablePadding 
            key={chat.id}
            secondaryAction={
              <Typography variant="caption" color="text.secondary">
                {new Date(chat.timestamp).toLocaleDateString()}
              </Typography>
            }
          >
            <ListItemButton
              selected={selectedChat?.id === chat.id}
              onClick={() => handleSelectChat(chat)}
              sx={{ borderRadius: 1, mx: 1 }}
            >
              <ListItemIcon sx={{ minWidth: 36 }}>
                <TableChartIcon fontSize="small" />
              </ListItemIcon>
              <ListItemText 
                primary={chat.title}
                secondary={chat.file}
                primaryTypographyProps={{ noWrap: true }}
                secondaryTypographyProps={{ noWrap: true, fontSize: '0.75rem' }}
              />
            </ListItemButton>
          </ListItem>
        ))}
      </List>
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
        sx={{ width: { md: drawerWidth }, flexShrink: { md: 0 } }}
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
              width: drawerWidth,
              flexShrink: 0,
              "& .MuiDrawer-paper": {
                width: drawerWidth,
                boxSizing: "border-box",
                backgroundColor: "#2F3136",
                color: "white",
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
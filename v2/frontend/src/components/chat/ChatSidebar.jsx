import React, { useState, useEffect } from "react";
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
  Collapse,
  Badge,
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
  ExpandLess,
  ExpandMore,
  Code as CodeIcon,
  Description as DescriptionIcon,
  DataObject as DataObjectIcon,
  Refresh as RefreshIcon,
  Upload as UploadIcon,
  ChevronLeft as ChevronLeftIcon
} from "@mui/icons-material";
import { Link } from "react-router-dom";
import TrainingDataSection from "./TrainingDataSection";
import { fetchChatHistory, createChat, deleteChat } from "../../utils/api";

// Fixed width for the drawer
const drawerWidth = 300;

const ChatSidebar = ({ selectedChat, setSelectedChat, onNewChat }) => {
  const [chatHistory, setChatHistory] = useState([]);
  const [mobileOpen, setMobileOpen] = useState(false);
  const [loading, setLoading] = useState(false);
  const [sidebarOpen, setSidebarOpen] = useState(true);
  
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down("md"));

  const handleDrawerToggle = () => {
    setMobileOpen(!mobileOpen);
  };

  const handleSidebarToggle = () => {
    setSidebarOpen(!sidebarOpen);
  };

  const handleNewChat = async () => {
    try {
      setLoading(true);
      
      // Reset the selected chat first - this should trigger UI reset
      setSelectedChat(null);
      
      // Create a new chat on the server
      const newChat = await createChat();
      
      // Update local state with the new chat
      setChatHistory(prev => [newChat, ...prev]);
      
      // Don't select the new chat immediately - keep UI in welcome state
      // setSelectedChat(newChat); <-- Comment out or remove this line
      
      // Notify parent component
      if (onNewChat) onNewChat(null); // Pass null instead of newChat
      
      if (isMobile) {
        setMobileOpen(false);
      }
    } catch (error) {
      console.error("Error creating new chat:", error);
    } finally {
      setLoading(false);
    }
  };

  const handleSelectChat = (chat) => {
    setSelectedChat(chat);
    if (isMobile) {
      setMobileOpen(false);
    }
  };

  const handleDeleteChat = async (chatId, event) => {
    event.stopPropagation();
    try {
      setLoading(true);
      console.log(`Attempting to delete chat: ${chatId}`);
      
      // Optimistically update UI first
      setChatHistory(prev => prev.filter(chat => chat.id !== chatId));
      
      // If the deleted chat was selected, clear selection
      if (selectedChat && selectedChat.id === chatId) {
        setSelectedChat(null);
      }
      
      // Then actually delete on server
      const result = await deleteChat(chatId);
      console.log("Delete result:", result);
      
      if (!result.success) {
        console.error(`Server failed to delete chat ${chatId}:`, result.error);
        // Optionally, you could add the chat back to the history here if the server delete failed
        // But most users won't notice if it's just removed from the UI
      }
    } catch (error) {
      console.error(`Error in handleDeleteChat for chat ${chatId}:`, error);
      // Show error to user
      alert(`Failed to delete chat: ${error.message}`);
      // Reload chat history to ensure UI consistency
      loadChatHistory();
    } finally {
      setLoading(false);
    }
  };

  // Load chat history
  const loadChatHistory = async () => {
    setLoading(true);
    try {
      const data = await fetchChatHistory();
      setChatHistory(data || []);
    } catch (error) {
      console.error("Error fetching chat history:", error);
      setChatHistory([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadChatHistory();
  }, []);

  const drawerContent = (
    <>
      <Box sx={{ 
        p: 2, 
        display: "flex", 
        alignItems: "center", 
        justifyContent: "space-between" 
      }}>
        <Typography variant="h6" sx={{ fontWeight: "bold" }}>Sage Chat</Typography>
        {isMobile ? (
          <IconButton onClick={handleDrawerToggle}>
            <CloseIcon />
          </IconButton>
        ) : (
          <IconButton onClick={handleSidebarToggle}>
            <ChevronLeftIcon />
          </IconButton>
        )}
      </Box>

      <Button
        variant="contained"
        startIcon={<AddIcon />}
        fullWidth
        sx={{ mx: 2, mb: 2 }}
        onClick={handleNewChat}
        disabled={loading}
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
      
      {/* Training Data Section */}
      <TrainingDataSection />
      
      <Divider />
      
      <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', px: 2, py: 1 }}>
        <Typography variant="subtitle2" sx={{ color: "text.secondary" }}>
          <HistoryIcon fontSize="small" sx={{ verticalAlign: "middle", mr: 1 }} />
          Recent Sessions
        </Typography>
        <Tooltip title="Refresh">
          <IconButton size="small" onClick={loadChatHistory} disabled={loading}>
            <RefreshIcon fontSize="small" />
          </IconButton>
        </Tooltip>
      </Box>
      
      <List sx={{ maxHeight: '40vh', overflow: 'auto' }}>
        {chatHistory.length === 0 ? (
          <ListItem>
            <ListItemText
              primary="No recent chats"
              primaryTypographyProps={{ 
                variant: 'body2', 
                color: 'text.secondary',
                align: 'center'
              }}
            />
          </ListItem>
        ) : (
          chatHistory.map((chat) => (
            <ListItem 
              disablePadding 
              key={chat.id}
              secondaryAction={
                <IconButton 
                  edge="end" 
                  size="small" 
                  onClick={(e) => handleDeleteChat(chat.id, e)}
                  sx={{ opacity: 0, '&:hover': { opacity: 1 } }}
                >
                  <CloseIcon fontSize="small" />
                </IconButton>
              }
              sx={{ 
                '&:hover .MuiIconButton-root': { 
                  opacity: 0.7 
                } 
              }}
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
                  primary={chat.title || "Untitled Chat"}
                  secondary={
                    <Box component="span" sx={{ display: 'flex', justifyContent: 'space-between' }}>
                      <Typography variant="caption" component="span">
                        {chat.file || "No file"}
                      </Typography>
                      <Typography variant="caption" component="span" sx={{ ml: 1 }}>
                        {new Date(chat.timestamp).toLocaleDateString()}
                      </Typography>
                    </Box>
                  }
                  primaryTypographyProps={{ noWrap: true }}
                  secondaryTypographyProps={{ noWrap: true }}
                />
              </ListItemButton>
            </ListItem>
          ))
        )}
      </List>
    </>
  );

  // Collapsed sidebar content
  const collapsedDrawerContent = (
    <Box sx={{ py: 2, display: "flex", flexDirection: "column", alignItems: "center" }}>
      <IconButton onClick={handleSidebarToggle} sx={{ mb: 2 }}>
        <MenuIcon />
      </IconButton>
      
      <Tooltip title="New Analysis">
        <IconButton 
          onClick={handleNewChat} 
          sx={{ 
            my: 1, 
            backgroundColor: "primary.main", 
            color: "white",
            '&:hover': { backgroundColor: "primary.dark" }
          }}
        >
          <AddIcon />
        </IconButton>
      </Tooltip>
      
      <Tooltip title="Back to Home">
        <IconButton 
          component={Link} 
          to="/" 
          sx={{ my: 1 }}
        >
          <HomeIcon />
        </IconButton>
      </Tooltip>
      
      <Divider sx={{ width: '80%', my: 2 }} />
      
      <Tooltip title="Chat History">
        <IconButton sx={{ my: 1 }}>
          <Badge badgeContent={chatHistory.length} color="primary">
            <HistoryIcon />
          </Badge>
        </IconButton>
      </Tooltip>
      
      <Tooltip title="Knowledge Base">
        <IconButton sx={{ my: 1 }}>
          <DataObjectIcon />
        </IconButton>
      </Tooltip>
    </Box>
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
          width: { 
            md: sidebarOpen ? drawerWidth : 60 
          },
          flexShrink: { md: 0 } 
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
              width: sidebarOpen ? drawerWidth : 60,
              flexShrink: 0,
              "& .MuiDrawer-paper": {
                width: sidebarOpen ? drawerWidth : 60,
                boxSizing: "border-box",
                backgroundColor: "#2F3136",
                color: "white",
                border: "none",
                transition: theme => theme.transitions.create('width', {
                  easing: theme.transitions.easing.sharp,
                  duration: theme.transitions.duration.enteringScreen,
                }),
                overflowX: 'hidden'
              },
            }}
            open
          >
            {sidebarOpen ? drawerContent : collapsedDrawerContent}
          </Drawer>
        )}
      </Box>
    </>
  );
};

export default ChatSidebar;
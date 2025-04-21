"use client"

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
  Tooltip,
  Avatar
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
  ChevronLeft as ChevronLeftIcon,
  Person as PersonIcon
} from "@mui/icons-material";
import { Link } from "react-router-dom";
import TrainingDataSection from "./TrainingDataSection";
import { fetchChatHistory, createChat, deleteChat } from "../../utils/api";
import { useAuth } from "../auth/AuthContext";

// Fixed width for the drawer
const drawerWidth = 300;

const ChatSidebar = ({ selectedChat, setSelectedChat, onNewChat }) => {
  const { user } = useAuth();
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
    <Box 
      sx={{ 
        height: '100%', 
        display: 'flex', 
        flexDirection: 'column',
        bgcolor: '#2F3136'
      }}
    >
      {/* Header */}
      <Box sx={{ 
        p: 2, 
        display: "flex", 
        alignItems: "center", 
        justifyContent: "space-between",
        borderBottom: '1px solid rgba(255, 255, 255, 0.1)'
      }}>
        <Typography variant="h6" sx={{ fontWeight: "bold", color: 'white' }}>
          Sage Chat
        </Typography>
        {isMobile ? (
          <IconButton onClick={handleDrawerToggle} sx={{ color: 'white' }}>
            <CloseIcon />
          </IconButton>
        ) : (
          <IconButton onClick={handleSidebarToggle} sx={{ color: 'white' }}>
            <ChevronLeftIcon />
          </IconButton>
        )}
      </Box>

      {/* Main Content Area */}
      <Box sx={{ flex: 1, overflow: 'hidden', display: 'flex', flexDirection: 'column' }}>
        {/* Action Buttons */}
        <Box sx={{ p: 2 }}>
          <Button
            variant="contained"
            startIcon={<AddIcon />}
            fullWidth
            sx={{ 
              mb: 2,
              bgcolor: theme.palette.primary.main,
              '&:hover': {
                bgcolor: theme.palette.primary.dark,
              }
            }}
            onClick={handleNewChat}
            disabled={loading}
          >
            New Analysis
          </Button>

          <Button
            component={Link}
            to="/"
            variant="text"
            fullWidth
            startIcon={<HomeIcon />}
            sx={{ 
              justifyContent: "flex-start", 
              textTransform: "none",
              color: 'white',
              '&:hover': {
                bgcolor: 'rgba(255, 255, 255, 0.1)'
              }
            }}
          >
            Back to Home
          </Button>
        </Box>

        <Divider sx={{ bgcolor: 'rgba(255, 255, 255, 0.1)' }} />
        
        {/* Training Data Section */}
        <TrainingDataSection />
        
        <Divider sx={{ bgcolor: 'rgba(255, 255, 255, 0.1)' }} />
        
        {/* Chat History Section */}
        <Box sx={{ flex: 1, display: 'flex', flexDirection: 'column', minHeight: 0 }}>
          <Box sx={{ 
            display: 'flex', 
            alignItems: 'center', 
            justifyContent: 'space-between', 
            px: 2, 
            py: 1 
          }}>
            <Typography variant="subtitle2" sx={{ color: "text.secondary" }}>
              <HistoryIcon fontSize="small" sx={{ verticalAlign: "middle", mr: 1 }} />
              Recent Sessions
            </Typography>
            <Tooltip title="Refresh">
              <IconButton 
                size="small" 
                onClick={loadChatHistory} 
                disabled={loading}
                sx={{ color: 'white' }}
              >
                <RefreshIcon fontSize="small" />
              </IconButton>
            </Tooltip>
          </Box>
          
          <List sx={{ 
            flex: 1, 
            overflow: 'auto',
            px: 1
          }}>
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
                      sx={{ 
                        opacity: 0, 
                        '&:hover': { opacity: 1 },
                        color: 'white'
                      }}
                    >
                      <CloseIcon fontSize="small" />
                    </IconButton>
                  }
                  sx={{ 
                    '&:hover .MuiIconButton-root': { 
                      opacity: 0.7 
                    },
                    mb: 0.5
                  }}
                >
                  <ListItemButton
                    selected={selectedChat?.id === chat.id}
                    onClick={() => handleSelectChat(chat)}
                    sx={{ 
                      borderRadius: 1,
                      '&.Mui-selected': {
                        bgcolor: 'rgba(88, 101, 242, 0.3)',
                        '&:hover': {
                          bgcolor: 'rgba(88, 101, 242, 0.4)',
                        }
                      },
                      '&:hover': {
                        bgcolor: 'rgba(255, 255, 255, 0.1)'
                      }
                    }}
                  >
                    <ListItemIcon sx={{ minWidth: 36, color: 'white' }}>
                      <TableChartIcon fontSize="small" />
                    </ListItemIcon>
                    <ListItemText 
                      primary={chat.title || "Untitled Chat"}
                      secondary={
                        <Box component="span" sx={{ display: 'flex', justifyContent: 'space-between' }}>
                          <Typography variant="caption" component="span" sx={{ color: 'rgba(255, 255, 255, 0.5)' }}>
                            {chat.file || "No file"}
                          </Typography>
                          <Typography variant="caption" component="span" sx={{ ml: 1, color: 'rgba(255, 255, 255, 0.5)' }}>
                            {new Date(chat.timestamp).toLocaleDateString()}
                          </Typography>
                        </Box>
                      }
                      primaryTypographyProps={{ 
                        noWrap: true,
                        sx: { color: 'white' }
                      }}
                      secondaryTypographyProps={{ noWrap: true }}
                    />
                  </ListItemButton>
                </ListItem>
              ))
            )}
          </List>
        </Box>
      </Box>

      {/* Profile Button */}
      <Box 
        sx={{ 
          mt: 'auto',
          borderTop: '1px solid rgba(255, 255, 255, 0.1)',
          p: 1
        }}
      >
        <Button
          fullWidth
          variant="text"
          component={Link}
          to="/profile"
          sx={{
            justifyContent: 'flex-start',
            color: 'white',
            p: 1,
            borderRadius: 1,
            transition: theme => theme.transitions.create(['background-color', 'transform'], {
              duration: theme.transitions.duration.shorter,
            }),
            '&:hover': {
              bgcolor: 'rgba(255, 255, 255, 0.1)',
              transform: 'translateX(4px)'
            }
          }}
        >
          <Avatar 
            sx={{ 
              width: 32, 
              height: 32, 
              mr: 1,
              bgcolor: theme.palette.primary.main,
              transition: 'transform 0.2s',
              '&:hover': {
                transform: 'scale(1.1)'
              }
            }}
          >
            {user?.name?.[0]?.toUpperCase() || 'U'}
          </Avatar>
          <Box sx={{ flex: 1 }}>
            <Typography variant="body2" sx={{ fontWeight: 500 }}>
              {user?.name || 'User'}
            </Typography>
            <Typography variant="caption" sx={{ color: 'rgba(255, 255, 255, 0.5)' }}>
              {user?.plan || 'Free Plan'}
            </Typography>
          </Box>
        </Button>
      </Box>
    </Box>
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
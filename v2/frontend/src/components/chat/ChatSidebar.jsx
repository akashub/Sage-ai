import { Box, Drawer, List, ListItem, ListItemButton, ListItemText, Typography, IconButton } from "@mui/material";
import { Add as AddIcon } from "@mui/icons-material";

const drawerWidth = 300;

const ChatSidebar = ({ selectedChat, setSelectedChat }) => {
  const chatHistory = [
    { id: 1, title: "SQL Query Optimization" },
    { id: 2, title: "Database Schema Design" },
    { id: 3, title: "Data Analysis Query" },
  ];

  return (
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
    >
      <Box sx={{ p: 2, display: "flex", alignItems: "center", justifyContent: "space-between" }}>
        <Typography variant="h6">Chat History</Typography>
        <IconButton color="primary" aria-label="new chat">
          <AddIcon />
        </IconButton>
      </Box>
      <List>
        {chatHistory.map((chat) => (
          <ListItem key={chat.id} disablePadding>
            <ListItemButton selected={selectedChat === chat.id} onClick={() => setSelectedChat(chat.id)}>
              <ListItemText primary={chat.title} />
            </ListItemButton>
          </ListItem>
        ))}
      </List>
    </Drawer>
  );
};

export default ChatSidebar;

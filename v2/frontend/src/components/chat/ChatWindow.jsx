"use client";

import { useState } from "react";
import {
  Box,
  TextField,
  Button,
  Typography,
  Paper,
  IconButton,
  InputAdornment,
} from "@mui/material";
import { Send as SendIcon, AttachFile as AttachFileIcon } from "@mui/icons-material";
import { keyframes } from "@mui/system";

const gradientAnimation = keyframes`
  0% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
  100% { background-position: 0% 50%; }
`;

const samplePrompts = [
  "How do I write a query to get users older than 30?",
  "Show me the latest 10 orders placed.",
  "What's the total revenue for last month?",
];

const ChatWindow = () => {
  const [input, setInput] = useState("");
  const [messages, setMessages] = useState([]);

  const handleSend = () => {
    if (input.trim()) {
      setMessages([...messages, { text: input, sender: "user" }]);
      setInput("");
      setTimeout(() => {
        setMessages((prev) => [
          ...prev,
          {
            text: "This is a sample AI response. Your SQL query might look like: SELECT * FROM table WHERE condition;",
            sender: "ai",
          },
        ]);
      }, 1000);
    }
  };

  return (
    <Box sx={{ flexGrow: 1, display: "flex", flexDirection: "column", height: "100vh", overflow: "hidden" }}>
      <Paper sx={{ flexGrow: 1, display: "flex", flexDirection: "column", p: 2, overflowY: "auto" }}>
        {messages.length === 0 ? (
          <Box
            sx={{
              display: "flex",
              flexDirection: "column",
              alignItems: "center",
              justifyContent: "center",
              height: "100%",
            }}
          >
            <Typography
              variant="h5"
              sx={{
                background: "linear-gradient(90deg, #00D4FF, #FF4D4D)",
                backgroundSize: "200% 200%",
                animation: `${gradientAnimation} 5s ease infinite`,
                WebkitBackgroundClip: "text",
                WebkitTextFillColor: "transparent",
                fontWeight: 700,
                mb: 2,
              }}
            >
              Welcome to Sage AI Chat!
            </Typography>
            <Typography variant="body1" sx={{ mb: 1 }}>
              Try one of these prompts:
            </Typography>
            {samplePrompts.map((prompt, index) => (
              <Typography
                key={index}
                variant="body2"
                sx={{
                  mb: 0.5,
                  background: "linear-gradient(90deg, #FF4D4D, #00D4FF)",
                  backgroundSize: "200% 200%",
                  animation: `${gradientAnimation} 5s ease infinite`,
                  WebkitBackgroundClip: "text",
                  WebkitTextFillColor: "transparent",
                }}
              >
                {prompt}
              </Typography>
            ))}
          </Box>
        ) : (
          messages.map((msg, index) => (
            <Box key={index} sx={{ display: "flex", justifyContent: msg.sender === "user" ? "flex-end" : "flex-start", mb: 1 }}>
              <Box
                sx={{
                  p: 1.5,
                  borderRadius: 2,
                  maxWidth: "70%",
                  backgroundColor: msg.sender === "user" ? "primary.main" : "background.default",
                  color: msg.sender === "user" ? "white" : "text.primary",
                }}
              >
                <Typography variant="body1">{msg.text}</Typography>
              </Box>
            </Box>
          ))
        )}
      </Paper>
      <Box sx={{ p: 2, borderTop: "1px solid rgba(255,255,255,0.1)", display: "flex", gap: 1 }}>
        <TextField
          fullWidth
          variant="outlined"
          placeholder="Type your message..."
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyPress={(e) => {
            if (e.key === "Enter") {
              e.preventDefault();
              handleSend();
            }
          }}
          InputProps={{
            endAdornment: (
              <InputAdornment position="end">
                <IconButton component="label">
                  <AttachFileIcon />
                  <input type="file" hidden />
                </IconButton>
              </InputAdornment>
            ),
          }}
        />
        <Button variant="contained" endIcon={<SendIcon />} onClick={handleSend}>
          Send
        </Button>
      </Box>
    </Box>
  );
};

export default ChatWindow;

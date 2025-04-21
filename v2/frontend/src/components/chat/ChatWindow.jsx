"use client"

import React, { useState, useRef, useEffect, useCallback } from "react";
import { uploadFile, queryData, fetchChatById, updateChat, createChat } from '../../utils/api';

import {
  Box,
  TextField,
  Button,
  Typography,
  Paper,
  IconButton,
  InputAdornment,
  CircularProgress,
  Divider,
  Chip,
  Tooltip,
  Switch,
  FormControlLabel,
  useTheme,
  Fade,
  Card,
  CardContent,
} from "@mui/material";
import { 
  Send as SendIcon, 
  AttachFile as AttachFileIcon, 
  TableChart as TableChartIcon,
  FileOpen as FileOpenIcon,
  Lightbulb as LightbulbIcon,
  DataObject as DataObjectIcon,
  Code as CodeIcon,
  Description as DescriptionIcon,
  Key as KeyIcon,
  AutoAwesome as AutoAwesomeIcon
} from "@mui/icons-material";
import { keyframes } from "@mui/system";

const gradientAnimation = keyframes`
  0% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
  100% { background-position: 0% 50%; }
`;

const pulseAnimation = keyframes`
  0% { transform: scale(1); }
  50% { transform: scale(1.05); }
  100% { transform: scale(1); }
`;

const samplePrompts = [
  "Show me the movies with a rating above 8.0",
  "What are the top 5 highest grossing movies?",
  "List all movies in the action genre sorted by release date",
];

const ChatWindow = React.memo(({ selectedChat, loading: externalLoading }) => {
  // Define state with useState hooks
  const [input, setInput] = useState("");
  const [messages, setMessages] = useState([]);
  const [loading, setLoading] = useState(false);
  const [csvFile, setCsvFile] = useState(null);
  const [csvHeaders, setCsvHeaders] = useState([]);
  const [sessionActive, setSessionActive] = useState(false);
  const [useKnowledgeBase, setUseKnowledgeBase] = useState(true);
  const [currentChatId, setCurrentChatId] = useState(null);
  const [chatTitle, setChatTitle] = useState("New Chat");
  const [llmProvider, setLLMProvider] = useState(null);
  const [apiKey, setApiKey] = useState(null);
  const [llmModel, setLLMModel] = useState(null);
  const [llmConfigDialog, setLLMConfigDialog] = useState(false);
  const [apiKeyError, setApiKeyError] = useState("");
  
  // References
  const fileInputRef = useRef(null);
  const messageEndRef = useRef(null);

  const theme = useTheme();

  // Model options for each provider
  const modelOptions = {
    gemini: ["gemini-1.5-flash", "gemini-1.5-pro", "gemini-pro"],
    openai: ["gpt-4-turbo-preview", "gpt-4", "gpt-3.5-turbo"],
    anthropic: ["claude-3-opus-20240229", "claude-3-sonnet-20240229", "claude-3-haiku-20240307"],
    mistral: ["mistral-large-latest", "mistral-medium", "mistral-small"]
  };

  // Define resetChat as a useCallback to avoid recreation on every render
  const resetChat = useCallback(() => {
    console.log("Resetting chat state...");
    setCsvFile(null);
    setCsvHeaders([]);
    setSessionActive(false);
    setMessages([]);
    setCurrentChatId(null);
    setChatTitle("New Chat");
    setLLMProvider(null);
    setApiKey(null);
    setLLMModel(null);
    
    // Clear the upload field if you have one
    if (fileInputRef.current) {
      fileInputRef.current.value = "";
    }
  }, []);

  // Load LLM config from the selected chat
  useEffect(() => {
    if (selectedChat && selectedChat.llmConfig) {
      setLLMProvider(selectedChat.llmConfig.provider);
      setApiKey(selectedChat.llmConfig.api_key);
      setLLMModel(selectedChat.llmConfig.model);
    } else {
      setLLMProvider(null);
      setApiKey(null);
      setLLMModel(null);
    }
  }, [selectedChat]);

  // Handle LLM configuration
  const handleLLMConfig = async () => {
    if (!llmProvider || !apiKey) {
      setApiKeyError("Please select a provider and enter an API key");
      return;
    }
    
    try {
      // Validate API key with backend
      const result = await validateApiKey(llmProvider, apiKey);
      
      if (result.valid) {
        // Update chat with LLM config
        if (currentChatId) {
          await updateChat(currentChatId, {
            llmConfig: {
              provider: llmProvider,
              api_key: apiKey,
              model: llmModel
            }
          });
        }
        setLLMConfigDialog(false);
        setApiKeyError("");
      } else {
        setApiKeyError("Invalid API key for selected provider");
      }
    } catch (error) {
      console.error("Error validating API key:", error);
      setApiKeyError(`Failed to validate API key: ${error.message}`);
    }
  };

  // Check if provider change is needed
  const handleProviderChange = (newProvider) => {
    if (sessionActive && llmProvider) {
      // Show warning dialog
      if (window.confirm("Changing LLM provider requires starting a new chat. Continue?")) {
        startNewSession();
        setLLMProvider(newProvider);
        setLLMModel(null); // Reset model when changing provider
      }
    } else {
      setLLMProvider(newProvider);
      setLLMModel(null); // Reset model when changing provider
    }
  };

  // Load selected chat data when selectedChat prop changes
  const loadSelectedChat = useCallback(async (chat) => {
    try {
      console.log("Loading selected chat:", chat);
      setLoading(true);
      
      // Set chat title
      setChatTitle(chat.title || "Untitled Chat");
      
      // Set current chat ID first
      if (chat.id) {
        console.log("Setting current chat ID:", chat.id);
        setCurrentChatId(chat.id);
      }
      
      // Fetch full chat data if only ID is provided
      let chatData = chat;
      if (chat.id && (!chat.messages || !chat.file)) {
        try {
          console.log("Fetching full chat data for ID:", chat.id);
          chatData = await fetchChatById(chat.id);
        } catch (error) {
          console.error("Error fetching chat data:", error);
          // Continue with what we have
        }
      }
      
      // Set session data
      setSessionActive(true);
      
      if (chatData.file) {
        console.log("Setting CSV file:", chatData.file);
        setCsvFile({ name: chatData.file, path: chatData.file });
      }
      
      if (chatData.headers && Array.isArray(chatData.headers)) {
        console.log("Setting CSV headers:", chatData.headers);
        setCsvHeaders(chatData.headers);
      }
      
      // Load LLM config if available
      if (chatData.llmConfig) {
        setLLMProvider(chatData.llmConfig.provider);
        setApiKey(chatData.llmConfig.api_key);
        setLLMModel(chatData.llmConfig.model);
      }
      
      // Load chat messages
      if (chatData.messages && Array.isArray(chatData.messages)) {
        console.log("Setting messages:", chatData.messages.length);
        setMessages(chatData.messages);
      } else {
        // Initialize with a system message if no messages exist
        const initialMessage = {
          type: "system",
          text: `Loaded session "${chatData.title || 'Untitled'}"${chatData.file ? ` with file: ${chatData.file}` : ''}`,
          timestamp: new Date(),
          file: chatData.file
        };
        console.log("Setting initial system message");
        setMessages([initialMessage]);
      }
    } catch (error) {
      console.error("Error loading chat:", error);
      
      // Initialize with an error message
      const errorMessage = {
        type: "error",
        text: `Error loading chat: ${error.message || "Unknown error"}`,
        timestamp: new Date()
      };
      console.log("Setting error message");
      setMessages([errorMessage]);
    } finally {
      setLoading(false);
    }
  }, []);

  // Scroll to bottom when messages change
  useEffect(() => {
    messageEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  // Set selected chat if passed from parent
  useEffect(() => {
    console.log("selectedChat changed:", selectedChat);
    
    if (selectedChat) {
      loadSelectedChat(selectedChat);
    } else {
      console.log("No selectedChat, resetting...");
      resetChat();
    }
  }, [selectedChat, loadSelectedChat, resetChat]);

  const handleFileUpload = async (event) => {
    const file = event.target.files[0];
    if (!file) return;
  
    setLoading(true);
    
    try {
      // Upload the CSV file to the backend
      const response = await uploadFile(file);
      console.log("Upload response:", response);
      
      // Store both the file name and the server-provided file path
      const filePath = response.FilePath || response.filePath;
      
      setCsvFile({
        name: file.name,
        path: filePath
      });
      
      // Get headers from the response
      const headers = response.Headers || response.headers || [];
      setCsvHeaders(headers);
      setSessionActive(true);
  
      // Create a system message
      const systemMessage = {
        type: "system",
        text: `CSV file "${file.name}" uploaded successfully. Found columns: ${headers.length > 0 ? headers.join(", ") : "none detected"}`,
        timestamp: new Date(),
        file: file.name
      };
      
      // Update messages state
      setMessages([systemMessage]);
      
      try {
        // Try to create a new chat if none exists
        if (!currentChatId) {
          const newChat = await createChat({
            file: file.name,
            filePath: filePath,
            headers,
            messages: [systemMessage],
            title: `${file.name} Analysis`,
            llmConfig: llmProvider && apiKey ? {
              provider: llmProvider,
              api_key: apiKey,
              model: llmModel
            } : null
          });
          setCurrentChatId(newChat.id);
          setChatTitle(newChat.title || "Untitled Chat");
        } else {
          // Update existing chat
          await updateChat(currentChatId, {
            file: file.name,
            filePath: filePath,
            headers,
            messages: [systemMessage],
            lastUpdated: new Date().toISOString()
          });
        }
      } catch (chatErr) {
        console.error("Error managing chat record:", chatErr);
        // Continue with the session anyway - the CSV is uploaded and ready to use
      }
      
      console.log("Stored file path:", filePath);
    } catch (error) {
      console.error("Error uploading CSV file:", error);
      setMessages([
        {
          type: "error",
          text: `Error uploading CSV file: ${error.message || "Unknown error"}`,
          timestamp: new Date()
        }
      ]);
    } finally {
      setLoading(false);
    }
  };

  const startNewSession = useCallback(() => {
    // Save current chat if needed
    if (currentChatId && messages.length > 0) {
      updateChat(currentChatId, {
        messages,
        lastUpdated: new Date().toISOString()
      }).catch(err => console.error("Error saving chat before new session:", err));
    }
    
    // Reset for new session
    resetChat();
  }, [currentChatId, messages, resetChat]);

  const handleSend = async () => {
    if (!input.trim() || !sessionActive) return;
    
    // Check if LLM is configured
    if (!llmProvider || !apiKey) {
      setLLMConfigDialog(true);
      return;
    }
  
    // Create a fresh user message
    const userMessage = {
      type: "user",
      text: input,
      sender: "user",
      timestamp: new Date()
    };
  
    // Clear input field immediately to prevent re-submission
    const currentInput = input;
    setInput("");
    
    // Show the user message right away
    const updatedMessages = [...messages, userMessage];
    setMessages(updatedMessages);
    
    // Show loading state
    setLoading(true);
  
    try {
      // Get the file path from csvFile
      let filePath = "";
      if (typeof csvFile === 'object' && csvFile !== null) {
        filePath = csvFile.path || csvFile.name || "";
      } else {
        filePath = csvFile || "";
      }
      
      console.log("Sending query:", currentInput, "with file path:", filePath);
      
      // Execute the query with LLM config
      const response = await queryData(currentInput, filePath, { 
        useKnowledgeBase: useKnowledgeBase, 
        timestamp: new Date().getTime(),
        chatId: currentChatId,
        llmConfig: {
          provider: llmProvider,
          api_key: apiKey,
          model: llmModel
        }
      });
      
      console.log("Query response:", response);
      
      // Create AI response message
      const aiMessage = {
        type: "ai",
        text: response.response || "Query processed successfully.",
        sender: "ai",
        timestamp: new Date(),
        results: response.results || [],
        generatedQuery: response.sql || response.generatedQuery || "",
        knowledgeContext: response.knowledgeContext || []
      };
      
      // Add AI message to chat
      const newMessages = [...updatedMessages, aiMessage];
      setMessages(newMessages);
      
      // Update chat in database
      if (currentChatId) {
        await updateChat(currentChatId, {
          messages: newMessages,
          lastUpdated: new Date().toISOString()
        });
      }
    } catch (error) {
      console.error("Error processing query:", error);
      
      // Add error message
      const errorMessage = {
        type: "error",
        text: `Error processing query: ${error.message || "Unknown error"}`,
        sender: "ai",
        timestamp: new Date()
      };
      
      const newMessages = [...updatedMessages, errorMessage];
      setMessages(newMessages);
      
      // Still update chat in database
      if (currentChatId) {
        await updateChat(currentChatId, {
          messages: newMessages,
          lastUpdated: new Date().toISOString()
        });
      }
    } finally {
      setLoading(false);
    }
  };

  const handlePromptClick = (prompt) => {
    setInput(prompt);
  };
  
  const updateChatTitle = async (newTitle) => {
    if (!currentChatId) return;
    
    try {
      setChatTitle(newTitle);
      await updateChat(currentChatId, { title: newTitle });
    } catch (error) {
      console.error("Error updating chat title:", error);
    }
  };

  return (
    <Box 
      sx={{ 
        flexGrow: 1, 
        display: "flex", 
        flexDirection: "column", 
        height: "100vh", 
        overflow: "hidden",
        bgcolor: 'background.default'
      }}
    >
      {/* Header */}
      <Paper 
        elevation={0}
        sx={{ 
          p: 2, 
          display: "flex", 
          justifyContent: "space-between", 
          alignItems: "center",
          borderBottom: '1px solid rgba(255, 255, 255, 0.1)',
          bgcolor: 'background.paper',
          backdropFilter: 'blur(10px)'
        }} 
      >
        <Box sx={{ display: 'flex', alignItems: 'center' }}>
          <AutoAwesomeIcon sx={{ mr: 1, color: 'primary.main' }} />
          <Typography variant="h6" sx={{ fontWeight: 600 }}>
          {sessionActive 
            ? `${chatTitle} - ${csvFile?.name || "CSV File"}`
            : "Sage AI - Upload a CSV file to begin"}
        </Typography>
        </Box>
        
        <Box sx={{ display: 'flex', gap: 2, alignItems: 'center' }}>
          {sessionActive && (
            <FormControlLabel
              control={
                <Switch
                  checked={useKnowledgeBase}
                  onChange={(e) => setUseKnowledgeBase(e.target.checked)}
                  size="small"
                />
              }
              label={
                <Box sx={{ display: 'flex', alignItems: 'center' }}>
                  <DataObjectIcon fontSize="small" sx={{ mr: 0.5 }} />
                  <Typography variant="body2">Knowledge Base</Typography>
                </Box>
              }
            />
          )}
          
          {sessionActive ? (
            <Button 
              variant="outlined" 
              size="small"
              onClick={startNewSession}
              startIcon={<FileOpenIcon />}
              sx={{
                borderRadius: 2,
                textTransform: 'none',
                '&:hover': {
                  transform: 'translateY(-1px)',
                  transition: 'transform 0.2s'
                }
              }}
            >
              New Session
            </Button>
          ) : (
            <Button
              variant="contained"
              onClick={() => fileInputRef.current?.click()}
              startIcon={<AttachFileIcon />}
              disabled={loading}
              sx={{
                borderRadius: 2,
                textTransform: 'none',
                animation: !sessionActive ? `${pulseAnimation} 2s infinite` : 'none',
                '&:hover': {
                  transform: 'translateY(-1px)',
                  transition: 'transform 0.2s'
                }
              }}
            >
              Upload CSV
            </Button>
          )}
        </Box>
        <input
          type="file"
          accept=".csv"
          ref={fileInputRef}
          style={{ display: "none" }}
          onChange={handleFileUpload}
        />
      </Paper>

      {/* Welcome Screen */}
      {!sessionActive && (
        <Fade in={!sessionActive}>
          <Box
            sx={{
              display: 'flex',
              flexDirection: 'column',
              alignItems: 'center',
              justifyContent: 'center',
              height: '100%',
              p: 4,
              textAlign: 'center',
              background: `linear-gradient(135deg, 
                ${theme.palette.background.default} 0%, 
                ${theme.palette.background.paper} 100%)`,
            }}
          >
            <AutoAwesomeIcon 
              sx={{
                fontSize: 48, 
                color: 'primary.main',
                mb: 2,
                animation: `${pulseAnimation} 2s infinite`
              }}
            />
            <Typography variant="h4" gutterBottom sx={{ fontWeight: 600 }}>
              Welcome to Sage AI Chat!
            </Typography>
            <Typography variant="subtitle1" color="text.secondary" sx={{ mb: 4, maxWidth: 600 }}>
              Upload a CSV file and try asking questions like:
            </Typography>
            <Box sx={{ display: 'flex', gap: 2, flexWrap: 'wrap', justifyContent: 'center' }}>
            {samplePrompts.map((prompt, index) => (
                <Chip
                key={index}
                  label={prompt}
                onClick={() => handlePromptClick(prompt)}
                sx={{
                    bgcolor: 'background.paper',
                    '&:hover': {
                      bgcolor: 'primary.dark',
                      transform: 'translateY(-2px)',
                      transition: 'all 0.2s'
                    },
                    transition: 'all 0.2s'
                  }}
                />
              ))}
          </Box>
                        </Box>
        </Fade>
      )}

      {/* Message Area */}
      {sessionActive && (
        <Box 
                  sx={{ 
            flex: 1,
            overflow: 'auto',
            p: 2,
            display: 'flex',
            flexDirection: 'column',
            gap: 2
          }}
        >
          {messages.map((msg, index) => (
            <Fade in key={index}>
              <Card
                elevation={0}
                    sx={{ 
                  alignSelf: msg.sender === "user" ? "flex-end" : "flex-start",
                  maxWidth: "70%",
                  bgcolor: msg.sender === "user" ? 'primary.dark' : 'background.paper',
                  borderRadius: 2,
                  position: 'relative',
                  '&:hover': {
                    transform: 'translateY(-1px)',
                    transition: 'transform 0.2s'
                  }
                }}
              >
                <CardContent sx={{ p: 2, '&:last-child': { pb: 2 } }}>
                  <Typography 
                    variant="body1" 
                  sx={{ 
                      color: msg.sender === "user" ? 'white' : 'text.primary',
                      whiteSpace: 'pre-wrap'
                    }}
                  >
                    {msg.text}
                          </Typography>
                  
                  {msg.generatedQuery && (
                    <Box sx={{ mt: 2 }}>
                      <Tooltip title="Generated SQL Query">
                              <Chip 
                          icon={<CodeIcon />}
                          label="View SQL"
                                size="small"
                          sx={{ 
                            bgcolor: 'background.paper',
                            '&:hover': { bgcolor: 'primary.dark' }
                          }}
                              />
                            </Tooltip>
                      </Box>
                    )}
                    
                  <Typography 
                    variant="caption" 
                    sx={{ 
                      display: 'block',
                      mt: 1,
                      color: msg.sender === "user" ? 'rgba(255,255,255,0.7)' : 'text.secondary'
                    }}
                  >
                    {new Date(msg.timestamp).toLocaleTimeString()}
                        </Typography>
                </CardContent>
              </Card>
            </Fade>
          ))}
          <div ref={messageEndRef} />
                      </Box>
                    )}

      {/* Input Area */}
      <Box 
        sx={{ 
          p: 2, 
          borderTop: '1px solid rgba(255, 255, 255, 0.1)',
          bgcolor: 'background.paper',
          backdropFilter: 'blur(10px)'
        }}
      >
        <Box sx={{ display: 'flex', gap: 1 }}>
        <TextField
          fullWidth
          variant="outlined"
          placeholder={sessionActive 
            ? "Ask a question about your data..." 
            : "Upload a CSV file to begin"
          }
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyPress={(e) => {
            if (e.key === "Enter" && !e.shiftKey) {
              e.preventDefault();
              handleSend();
            }
          }}
          disabled={!sessionActive || loading}
            multiline
            maxRows={4}
            sx={{
              '& .MuiOutlinedInput-root': {
                borderRadius: 2,
                bgcolor: 'background.paper',
                '&.Mui-focused': {
                  '& .MuiOutlinedInput-notchedOutline': {
                    borderColor: 'primary.main',
                    borderWidth: 2
                  }
                }
              }
            }}
          InputProps={{
            endAdornment: !sessionActive && (
              <InputAdornment position="end">
                  <IconButton 
                    component="label" 
                    onClick={() => fileInputRef.current?.click()}
                    sx={{
                      '&:hover': {
                        transform: 'scale(1.1)',
                        transition: 'transform 0.2s'
                      }
                    }}
                  >
                  <AttachFileIcon />
                </IconButton>
              </InputAdornment>
            ),
          }}
        />
        <Button 
          variant="contained" 
          endIcon={loading ? <CircularProgress size={20} color="inherit" /> : <SendIcon />} 
          onClick={handleSend}
          disabled={!sessionActive || loading || !input.trim()}
            sx={{
              borderRadius: 2,
              px: 3,
              height: 56,
              minWidth: 100,
              textTransform: 'none',
              '&:not(:disabled):hover': {
                transform: 'translateY(-1px)',
                transition: 'transform 0.2s'
              }
            }}
        >
          Send
        </Button>
      </Box>
        </Box>
    </Box>
  );
});

export default ChatWindow;
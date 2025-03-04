// "use client";

// import { useState } from "react";
// import {
//   Box,
//   TextField,
//   Button,
//   Typography,
//   Paper,
//   IconButton,
//   InputAdornment,
// } from "@mui/material";
// import { Send as SendIcon, AttachFile as AttachFileIcon } from "@mui/icons-material";
// import { keyframes } from "@mui/system";

// const gradientAnimation = keyframes`
//   0% { background-position: 0% 50%; }
//   50% { background-position: 100% 50%; }
//   100% { background-position: 0% 50%; }
// `;

// const samplePrompts = [
//   "How do I write a query to get users older than 30?",
//   "Show me the latest 10 orders placed.",
//   "What's the total revenue for last month?",
// ];

// const ChatWindow = () => {
//   const [input, setInput] = useState("");
//   const [messages, setMessages] = useState([]);

//   const handleSend = () => {
//     if (input.trim()) {
//       setMessages([...messages, { text: input, sender: "user" }]);
//       setInput("");
//       setTimeout(() => {
//         setMessages((prev) => [
//           ...prev,
//           {
//             text: "This is a sample AI response. Your SQL query might look like: SELECT * FROM table WHERE condition;",
//             sender: "ai",
//           },
//         ]);
//       }, 1000);
//     }
//   };

//   return (
//     <Box sx={{ flexGrow: 1, display: "flex", flexDirection: "column", height: "100vh", overflow: "hidden" }}>
//       <Paper sx={{ flexGrow: 1, display: "flex", flexDirection: "column", p: 2, overflowY: "auto" }}>
//         {messages.length === 0 ? (
//           <Box
//             sx={{
//               display: "flex",
//               flexDirection: "column",
//               alignItems: "center",
//               justifyContent: "center",
//               height: "100%",
//             }}
//           >
//             <Typography
//               variant="h5"
//               sx={{
//                 background: "linear-gradient(90deg, #00D4FF, #FF4D4D)",
//                 backgroundSize: "200% 200%",
//                 animation: `${gradientAnimation} 5s ease infinite`,
//                 WebkitBackgroundClip: "text",
//                 WebkitTextFillColor: "transparent",
//                 fontWeight: 700,
//                 mb: 2,
//               }}
//             >
//               Welcome to Sage AI Chat!
//             </Typography>
//             <Typography variant="body1" sx={{ mb: 1 }}>
//               Try one of these prompts:
//             </Typography>
//             {samplePrompts.map((prompt, index) => (
//               <Typography
//                 key={index}
//                 variant="body2"
//                 sx={{
//                   mb: 0.5,
//                   background: "linear-gradient(90deg, #FF4D4D, #00D4FF)",
//                   backgroundSize: "200% 200%",
//                   animation: `${gradientAnimation} 5s ease infinite`,
//                   WebkitBackgroundClip: "text",
//                   WebkitTextFillColor: "transparent",
//                 }}
//               >
//                 {prompt}
//               </Typography>
//             ))}
//           </Box>
//         ) : (
//           messages.map((msg, index) => (
//             <Box key={index} sx={{ display: "flex", justifyContent: msg.sender === "user" ? "flex-end" : "flex-start", mb: 1 }}>
//               <Box
//                 sx={{
//                   p: 1.5,
//                   borderRadius: 2,
//                   maxWidth: "70%",
//                   backgroundColor: msg.sender === "user" ? "primary.main" : "background.default",
//                   color: msg.sender === "user" ? "white" : "text.primary",
//                 }}
//               >
//                 <Typography variant="body1">{msg.text}</Typography>
//               </Box>
//             </Box>
//           ))
//         )}
//       </Paper>
//       <Box sx={{ p: 2, borderTop: "1px solid rgba(255,255,255,0.1)", display: "flex", gap: 1 }}>
//         <TextField
//           fullWidth
//           variant="outlined"
//           placeholder="Type your message..."
//           value={input}
//           onChange={(e) => setInput(e.target.value)}
//           onKeyPress={(e) => {
//             if (e.key === "Enter") {
//               e.preventDefault();
//               handleSend();
//             }
//           }}
//           InputProps={{
//             endAdornment: (
//               <InputAdornment position="end">
//                 <IconButton component="label">
//                   <AttachFileIcon />
//                   <input type="file" hidden />
//                 </IconButton>
//               </InputAdornment>
//             ),
//           }}
//         />
//         <Button variant="contained" endIcon={<SendIcon />} onClick={handleSend}>
//           Send
//         </Button>
//       </Box>
//     </Box>
//   );
// };

// export default ChatWindow;

"use client";

import { useState, useRef, useEffect } from "react";
import { uploadFile, queryData } from '../../utils/api';

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
}

from "@mui/material";
import { 
  Send as SendIcon, 
  AttachFile as AttachFileIcon, 
  TableChart as TableChartIcon,
  FileOpen as FileOpenIcon
} 

from "@mui/icons-material";
import { keyframes } from "@mui/system";

const gradientAnimation = keyframes`
  0% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
  100% { background-position: 0% 50%; }
`;

const samplePrompts = [
  "Show me the movies with a rating above 8.0",
  "What are the top 5 highest grossing movies?",
  "List all movies in the action genre sorted by release date",
];

const ChatWindow = () => {
  const [input, setInput] = useState("");
  const [messages, setMessages] = useState([]);
  const [loading, setLoading] = useState(false);
  const [csvFile, setCsvFile] = useState(null);
  const [csvHeaders, setCsvHeaders] = useState([]);
  const [sessionActive, setSessionActive] = useState(false);
  const fileInputRef = useRef(null);
  const messageEndRef = useRef(null);

  // Scroll to bottom when messages change
  useEffect(() => {
    messageEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  const handleFileUpload = async (event) => {
    const file = event.target.files[0];
    if (!file) return;
  
    setLoading(true);
    setCsvFile(file);
  
    try {
      // Upload the CSV file to the backend
      const data = await uploadFile(file);
      setCsvHeaders(data.headers);
      setSessionActive(true);
  
      // Add system message
      setMessages([
        ...messages,
        {
          type: "system",
          text: `CSV file "${file.name}" uploaded successfully. Found columns: ${data.headers.join(", ")}`,
          timestamp: new Date(),
          file: file.name
        }
      ]);
    } catch (error) {
      console.error("Error uploading CSV file:", error);
      setMessages([
        ...messages,
        {
          type: "error",
          text: `Error uploading CSV file: ${error.message}`,
          timestamp: new Date()
        }
      ]);
    } finally {
      setLoading(false);
    }
  };

  const startNewSession = () => {
    setCsvFile(null);
    setCsvHeaders([]);
    setSessionActive(false);
    setMessages([
      ...messages,
      {
        type: "system",
        text: "Started a new session. Please upload a CSV file.",
        timestamp: new Date()
      }
    ]);
  };

  const handleSend = async () => {
    if (!input.trim() || !sessionActive) return;
  
    const userMessage = {
      type: "user",
      text: input,
      sender: "user",
      timestamp: new Date()
    };
  
    setMessages([...messages, userMessage]);
    setInput("");
    setLoading(true);
  
    try {
      // Send the query to the backend
      const data = await queryData(input, csvFile ? csvFile.name : "");
      
      // Format and display the results
      setMessages(prevMessages => [
        ...prevMessages,
        {
          type: "ai",
          text: "Query processed successfully.",
          sender: "ai",
          timestamp: new Date(),
          results: data.results,
          generatedQuery: data.generatedQuery
        }
      ]);
    } catch (error) {
      console.error("Error processing query:", error);
      setMessages(prevMessages => [
        ...prevMessages,
        {
          type: "error",
          text: `Error processing query: ${error.message}`,
          sender: "ai",
          timestamp: new Date()
        }
      ]);
    } finally {
      setLoading(false);
    }
  };

  const handlePromptClick = (prompt) => {
    setInput(prompt);
  };

  return (
    <Box sx={{ flexGrow: 1, display: "flex", flexDirection: "column", height: "100vh", overflow: "hidden" }}>
      {/* Header */}
      <Paper 
        sx={{ 
          p: 2, 
          display: "flex", 
          justifyContent: "space-between", 
          alignItems: "center",
          borderRadius: 0
        }} 
        elevation={2}
      >
        <Typography variant="h6">
          {sessionActive 
            ? `Sage AI - Session with ${csvFile?.name || "CSV File"}`
            : "Sage AI - Upload a CSV file to begin"}
        </Typography>
        {sessionActive ? (
          <Button 
            variant="outlined" 
            size="small"
            onClick={startNewSession}
            startIcon={<FileOpenIcon />}
          >
            New Session
          </Button>
        ) : (
          <Button
            variant="contained"
            onClick={() => fileInputRef.current.click()}
            startIcon={<AttachFileIcon />}
            disabled={loading}
          >
            Upload CSV
          </Button>
        )}
        <input
          type="file"
          accept=".csv"
          ref={fileInputRef}
          style={{ display: "none" }}
          onChange={handleFileUpload}
        />
      </Paper>

      {/* Message Area */}
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
              Upload a CSV file and try asking questions like:
            </Typography>
            {samplePrompts.map((prompt, index) => (
              <Typography
                key={index}
                variant="body2"
                onClick={() => handlePromptClick(prompt)}
                sx={{
                  mb: 0.5,
                  background: "linear-gradient(90deg, #FF4D4D, #00D4FF)",
                  backgroundSize: "200% 200%",
                  animation: `${gradientAnimation} 5s ease infinite`,
                  WebkitBackgroundClip: "text",
                  WebkitTextFillColor: "transparent",
                  cursor: "pointer",
                  "&:hover": {
                    opacity: 0.8,
                  }
                }}
              >
                {prompt}
              </Typography>
            ))}
            <Button
              variant="contained"
              sx={{ mt: 3 }}
              onClick={() => fileInputRef.current.click()}
              startIcon={<AttachFileIcon />}
            >
              Upload CSV
            </Button>
          </Box>
        ) : (
          messages.map((msg, index) => {
            if (msg.type === "system") {
              return (
                <Box 
                  key={index}
                  sx={{ 
                    display: "flex", 
                    justifyContent: "center", 
                    mb: 2 
                  }}
                >
                  <Paper 
                    sx={{ 
                      p: 1, 
                      backgroundColor: "rgba(0,0,0,0.2)", 
                      color: "#dddddd", 
                      maxWidth: "90%",
                      borderRadius: 2
                    }}
                  >
                    <Typography variant="body2" align="center">
                      {msg.file && (
                        <Box sx={{ display: "flex", alignItems: "center", justifyContent: "center", mb: 1 }}>
                          <TableChartIcon fontSize="small" sx={{ mr: 1 }} />
                          <Typography variant="caption">{msg.file}</Typography>
                        </Box>
                      )}
                      {msg.text}
                    </Typography>
                  </Paper>
                </Box>
              );
            } else if (msg.type === "error") {
              return (
                <Box 
                  key={index}
                  sx={{ 
                    display: "flex", 
                    justifyContent: "center", 
                    mb: 2 
                  }}
                >
                  <Paper 
                    sx={{ 
                      p: 1, 
                      backgroundColor: "#770000", 
                      color: "#ffffff", 
                      maxWidth: "90%",
                      borderRadius: 2
                    }}
                  >
                    <Typography variant="body2" align="center">
                      {msg.text}
                    </Typography>
                  </Paper>
                </Box>
              );
            } else {
              return (
                <Box 
                  key={index} 
                  sx={{ 
                    display: "flex", 
                    justifyContent: msg.sender === "user" ? "flex-end" : "flex-start", 
                    mb: 2 
                  }}
                >
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
                    
                    {msg.generatedQuery && (
                      <Box sx={{ mt: 1, mb: 1, p: 1, bgcolor: "rgba(0,0,0,0.3)", borderRadius: 1 }}>
                        <Typography variant="caption" color="text.secondary">
                          Generated SQL:
                        </Typography>
                        <Typography variant="body2" component="pre" sx={{ 
                          overflowX: "auto", 
                          whiteSpace: "pre-wrap",
                          wordBreak: "break-word" 
                        }}>
                          {msg.generatedQuery}
                        </Typography>
                      </Box>
                    )}
                    
                    {msg.results && msg.results.length > 0 && (
                      <Box sx={{ mt: 1 }}>
                        <Typography variant="subtitle2" sx={{ mb: 1 }}>
                          Results:
                        </Typography>
                        <Paper variant="outlined" sx={{ p: 1, bgcolor: "rgba(0,0,0,0.2)" }}>
                          {msg.results.map((result, idx) => (
                            <Box key={idx} sx={{ mb: idx < msg.results.length - 1 ? 1 : 0 }}>
                              {Object.entries(result).map(([key, value]) => (
                                <Typography key={key} variant="body2">
                                  <strong>{key}:</strong> {value?.toString()}
                                </Typography>
                              ))}
                              {idx < msg.results.length - 1 && <Divider sx={{ my: 1 }} />}
                            </Box>
                          ))}
                        </Paper>
                      </Box>
                    )}
                    
                    <Typography variant="caption" sx={{ opacity: 0.7, display: "block", mt: 1 }}>
                      {new Date(msg.timestamp).toLocaleTimeString()}
                    </Typography>
                  </Box>
                </Box>
              );
            }
          })
        )}
        <div ref={messageEndRef} />
      </Paper>

      {/* Input Area */}
      <Box sx={{ p: 2, borderTop: "1px solid rgba(255,255,255,0.1)", display: "flex", gap: 1 }}>
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
          InputProps={{
            endAdornment: !sessionActive && (
              <InputAdornment position="end">
                <IconButton component="label" onClick={() => fileInputRef.current.click()}>
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
        >
          Send
        </Button>
      </Box>
      
      {/* CSV Headers */}
      {csvHeaders.length > 0 && (
        <Box sx={{ p: 1, borderTop: "1px solid rgba(255,255,255,0.1)", display: "flex", flexWrap: "wrap", gap: 0.5 }}>
          <Typography variant="caption" color="text.secondary" sx={{ mr: 1, alignSelf: "center" }}>
            Available columns:
          </Typography>
          {csvHeaders.map((header, idx) => (
            <Typography key={idx} variant="caption" sx={{ 
              bgcolor: "rgba(88, 101, 242, 0.2)", 
              px: 0.7, 
              py: 0.3, 
              borderRadius: 1,
              fontSize: "0.7rem"
            }}>
              {header}
            </Typography>
          ))}
        </Box>
      )}
    </Box>
  );
};

export default ChatWindow;
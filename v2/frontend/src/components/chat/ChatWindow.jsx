// // // // // "use client";

// // // // // import { useState, useRef, useEffect } from "react";
// // // // // import { uploadFile, queryData } from '../../utils/api';

// // // // // import {
// // // // //   Box,
// // // // //   TextField,
// // // // //   Button,
// // // // //   Typography,
// // // // //   Paper,
// // // // //   IconButton,
// // // // //   InputAdornment,
// // // // //   CircularProgress,
// // // // //   Divider,
// // // // // }

// // // // // from "@mui/material";
// // // // // import { 
// // // // //   Send as SendIcon, 
// // // // //   AttachFile as AttachFileIcon, 
// // // // //   TableChart as TableChartIcon,
// // // // //   FileOpen as FileOpenIcon
// // // // // } 

// // // // // from "@mui/icons-material";
// // // // // import { keyframes } from "@mui/system";

// // // // // const gradientAnimation = keyframes`
// // // // //   0% { background-position: 0% 50%; }
// // // // //   50% { background-position: 100% 50%; }
// // // // //   100% { background-position: 0% 50%; }
// // // // // `;

// // // // // const samplePrompts = [
// // // // //   "Show me the movies with a rating above 8.0",
// // // // //   "What are the top 5 highest grossing movies?",
// // // // //   "List all movies in the action genre sorted by release date",
// // // // // ];

// // // // // const ChatWindow = () => {
// // // // //   const [input, setInput] = useState("");
// // // // //   const [messages, setMessages] = useState([]);
// // // // //   const [loading, setLoading] = useState(false);
// // // // //   const [csvFile, setCsvFile] = useState(null);
// // // // //   const [csvHeaders, setCsvHeaders] = useState([]);
// // // // //   const [sessionActive, setSessionActive] = useState(false);
// // // // //   const fileInputRef = useRef(null);
// // // // //   const messageEndRef = useRef(null);

// // // // //   // Scroll to bottom when messages change
// // // // //   useEffect(() => {
// // // // //     messageEndRef.current?.scrollIntoView({ behavior: "smooth" });
// // // // //   }, [messages]);

// // // // //   const handleFileUpload = async (event) => {
// // // // //     const file = event.target.files[0];
// // // // //     if (!file) return;
  
// // // // //     setLoading(true);
// // // // //     setCsvFile(file);
  
// // // // //     try {
// // // // //       // Upload the CSV file to the backend
// // // // //       const data = await uploadFile(file);
// // // // //       setCsvHeaders(data.headers);
// // // // //       setSessionActive(true);
  
// // // // //       // Add system message
// // // // //       setMessages([
// // // // //         ...messages,
// // // // //         {
// // // // //           type: "system",
// // // // //           text: `CSV file "${file.name}" uploaded successfully. Found columns: ${data.headers.join(", ")}`,
// // // // //           timestamp: new Date(),
// // // // //           file: file.name
// // // // //         }
// // // // //       ]);
// // // // //     } catch (error) {
// // // // //       console.error("Error uploading CSV file:", error);
// // // // //       setMessages([
// // // // //         ...messages,
// // // // //         {
// // // // //           type: "error",
// // // // //           text: `Error uploading CSV file: ${error.message}`,
// // // // //           timestamp: new Date()
// // // // //         }
// // // // //       ]);
// // // // //     } finally {
// // // // //       setLoading(false);
// // // // //     }
// // // // //   };

// // // // //   const startNewSession = () => {
// // // // //     setCsvFile(null);
// // // // //     setCsvHeaders([]);
// // // // //     setSessionActive(false);
// // // // //     setMessages([
// // // // //       ...messages,
// // // // //       {
// // // // //         type: "system",
// // // // //         text: "Started a new session. Please upload a CSV file.",
// // // // //         timestamp: new Date()
// // // // //       }
// // // // //     ]);
// // // // //   };

// // // // //   const handleSend = async () => {
// // // // //     if (!input.trim() || !sessionActive) return;
  
// // // // //     const userMessage = {
// // // // //       type: "user",
// // // // //       text: input,
// // // // //       sender: "user",
// // // // //       timestamp: new Date()
// // // // //     };
  
// // // // //     setMessages([...messages, userMessage]);
// // // // //     setInput("");
// // // // //     setLoading(true);
  
// // // // //     try {
// // // // //       // Send the query to the backend
// // // // //       const data = await queryData(input, csvFile ? csvFile.name : "");
      
// // // // //       // Format and display the results
// // // // //       setMessages(prevMessages => [
// // // // //         ...prevMessages,
// // // // //         {
// // // // //           type: "ai",
// // // // //           text: "Query processed successfully.",
// // // // //           sender: "ai",
// // // // //           timestamp: new Date(),
// // // // //           results: data.results,
// // // // //           generatedQuery: data.generatedQuery
// // // // //         }
// // // // //       ]);
// // // // //     } catch (error) {
// // // // //       console.error("Error processing query:", error);
// // // // //       setMessages(prevMessages => [
// // // // //         ...prevMessages,
// // // // //         {
// // // // //           type: "error",
// // // // //           text: `Error processing query: ${error.message}`,
// // // // //           sender: "ai",
// // // // //           timestamp: new Date()
// // // // //         }
// // // // //       ]);
// // // // //     } finally {
// // // // //       setLoading(false);
// // // // //     }
// // // // //   };

// // // // //   const handlePromptClick = (prompt) => {
// // // // //     setInput(prompt);
// // // // //   };

// // // // //   return (
// // // // //     <Box sx={{ flexGrow: 1, display: "flex", flexDirection: "column", height: "100vh", overflow: "hidden" }}>
// // // // //       {/* Header */}
// // // // //       <Paper 
// // // // //         sx={{ 
// // // // //           p: 2, 
// // // // //           display: "flex", 
// // // // //           justifyContent: "space-between", 
// // // // //           alignItems: "center",
// // // // //           borderRadius: 0
// // // // //         }} 
// // // // //         elevation={2}
// // // // //       >
// // // // //         <Typography variant="h6">
// // // // //           {sessionActive 
// // // // //             ? `Sage AI - Session with ${csvFile?.name || "CSV File"}`
// // // // //             : "Sage AI - Upload a CSV file to begin"}
// // // // //         </Typography>
// // // // //         {sessionActive ? (
// // // // //           <Button 
// // // // //             variant="outlined" 
// // // // //             size="small"
// // // // //             onClick={startNewSession}
// // // // //             startIcon={<FileOpenIcon />}
// // // // //           >
// // // // //             New Session
// // // // //           </Button>
// // // // //         ) : (
// // // // //           <Button
// // // // //             variant="contained"
// // // // //             onClick={() => fileInputRef.current.click()}
// // // // //             startIcon={<AttachFileIcon />}
// // // // //             disabled={loading}
// // // // //           >
// // // // //             Upload CSV
// // // // //           </Button>
// // // // //         )}
// // // // //         <input
// // // // //           type="file"
// // // // //           accept=".csv"
// // // // //           ref={fileInputRef}
// // // // //           style={{ display: "none" }}
// // // // //           onChange={handleFileUpload}
// // // // //         />
// // // // //       </Paper>

// // // // //       {/* Message Area */}
// // // // //       <Paper sx={{ flexGrow: 1, display: "flex", flexDirection: "column", p: 2, overflowY: "auto" }}>
// // // // //         {messages.length === 0 ? (
// // // // //           <Box
// // // // //             sx={{
// // // // //               display: "flex",
// // // // //               flexDirection: "column",
// // // // //               alignItems: "center",
// // // // //               justifyContent: "center",
// // // // //               height: "100%",
// // // // //             }}
// // // // //           >
// // // // //             <Typography
// // // // //               variant="h5"
// // // // //               sx={{
// // // // //                 background: "linear-gradient(90deg, #00D4FF, #FF4D4D)",
// // // // //                 backgroundSize: "200% 200%",
// // // // //                 animation: `${gradientAnimation} 5s ease infinite`,
// // // // //                 WebkitBackgroundClip: "text",
// // // // //                 WebkitTextFillColor: "transparent",
// // // // //                 fontWeight: 700,
// // // // //                 mb: 2,
// // // // //               }}
// // // // //             >
// // // // //               Welcome to Sage AI Chat!
// // // // //             </Typography>
// // // // //             <Typography variant="body1" sx={{ mb: 1 }}>
// // // // //               Upload a CSV file and try asking questions like:
// // // // //             </Typography>
// // // // //             {samplePrompts.map((prompt, index) => (
// // // // //               <Typography
// // // // //                 key={index}
// // // // //                 variant="body2"
// // // // //                 onClick={() => handlePromptClick(prompt)}
// // // // //                 sx={{
// // // // //                   mb: 0.5,
// // // // //                   background: "linear-gradient(90deg, #FF4D4D, #00D4FF)",
// // // // //                   backgroundSize: "200% 200%",
// // // // //                   animation: `${gradientAnimation} 5s ease infinite`,
// // // // //                   WebkitBackgroundClip: "text",
// // // // //                   WebkitTextFillColor: "transparent",
// // // // //                   cursor: "pointer",
// // // // //                   "&:hover": {
// // // // //                     opacity: 0.8,
// // // // //                   }
// // // // //                 }}
// // // // //               >
// // // // //                 {prompt}
// // // // //               </Typography>
// // // // //             ))}
// // // // //             <Button
// // // // //               variant="contained"
// // // // //               sx={{ mt: 3 }}
// // // // //               onClick={() => fileInputRef.current.click()}
// // // // //               startIcon={<AttachFileIcon />}
// // // // //             >
// // // // //               Upload CSV
// // // // //             </Button>
// // // // //           </Box>
// // // // //         ) : (
// // // // //           messages.map((msg, index) => {
// // // // //             if (msg.type === "system") {
// // // // //               return (
// // // // //                 <Box 
// // // // //                   key={index}
// // // // //                   sx={{ 
// // // // //                     display: "flex", 
// // // // //                     justifyContent: "center", 
// // // // //                     mb: 2 
// // // // //                   }}
// // // // //                 >
// // // // //                   <Paper 
// // // // //                     sx={{ 
// // // // //                       p: 1, 
// // // // //                       backgroundColor: "rgba(0,0,0,0.2)", 
// // // // //                       color: "#dddddd", 
// // // // //                       maxWidth: "90%",
// // // // //                       borderRadius: 2
// // // // //                     }}
// // // // //                   >
// // // // //                     <Typography variant="body2" align="center">
// // // // //                       {msg.file && (
// // // // //                         <Box sx={{ display: "flex", alignItems: "center", justifyContent: "center", mb: 1 }}>
// // // // //                           <TableChartIcon fontSize="small" sx={{ mr: 1 }} />
// // // // //                           <Typography variant="caption">{msg.file}</Typography>
// // // // //                         </Box>
// // // // //                       )}
// // // // //                       {msg.text}
// // // // //                     </Typography>
// // // // //                   </Paper>
// // // // //                 </Box>
// // // // //               );
// // // // //             } else if (msg.type === "error") {
// // // // //               return (
// // // // //                 <Box 
// // // // //                   key={index}
// // // // //                   sx={{ 
// // // // //                     display: "flex", 
// // // // //                     justifyContent: "center", 
// // // // //                     mb: 2 
// // // // //                   }}
// // // // //                 >
// // // // //                   <Paper 
// // // // //                     sx={{ 
// // // // //                       p: 1, 
// // // // //                       backgroundColor: "#770000", 
// // // // //                       color: "#ffffff", 
// // // // //                       maxWidth: "90%",
// // // // //                       borderRadius: 2
// // // // //                     }}
// // // // //                   >
// // // // //                     <Typography variant="body2" align="center">
// // // // //                       {msg.text}
// // // // //                     </Typography>
// // // // //                   </Paper>
// // // // //                 </Box>
// // // // //               );
// // // // //             } else {
// // // // //               return (
// // // // //                 <Box 
// // // // //                   key={index} 
// // // // //                   sx={{ 
// // // // //                     display: "flex", 
// // // // //                     justifyContent: msg.sender === "user" ? "flex-end" : "flex-start", 
// // // // //                     mb: 2 
// // // // //                   }}
// // // // //                 >
// // // // //                   <Box
// // // // //                     sx={{
// // // // //                       p: 1.5,
// // // // //                       borderRadius: 2,
// // // // //                       maxWidth: "70%",
// // // // //                       backgroundColor: msg.sender === "user" ? "primary.main" : "background.default",
// // // // //                       color: msg.sender === "user" ? "white" : "text.primary",
// // // // //                     }}
// // // // //                   >
// // // // //                     <Typography variant="body1">{msg.text}</Typography>
                    
// // // // //                     {msg.generatedQuery && (
// // // // //                       <Box sx={{ mt: 1, mb: 1, p: 1, bgcolor: "rgba(0,0,0,0.3)", borderRadius: 1 }}>
// // // // //                         <Typography variant="caption" color="text.secondary">
// // // // //                           Generated SQL:
// // // // //                         </Typography>
// // // // //                         <Typography variant="body2" component="pre" sx={{ 
// // // // //                           overflowX: "auto", 
// // // // //                           whiteSpace: "pre-wrap",
// // // // //                           wordBreak: "break-word" 
// // // // //                         }}>
// // // // //                           {msg.generatedQuery}
// // // // //                         </Typography>
// // // // //                       </Box>
// // // // //                     )}
                    
// // // // //                     {msg.results && msg.results.length > 0 && (
// // // // //                       <Box sx={{ mt: 1 }}>
// // // // //                         <Typography variant="subtitle2" sx={{ mb: 1 }}>
// // // // //                           Results:
// // // // //                         </Typography>
// // // // //                         <Paper variant="outlined" sx={{ p: 1, bgcolor: "rgba(0,0,0,0.2)" }}>
// // // // //                           {msg.results.map((result, idx) => (
// // // // //                             <Box key={idx} sx={{ mb: idx < msg.results.length - 1 ? 1 : 0 }}>
// // // // //                               {Object.entries(result).map(([key, value]) => (
// // // // //                                 <Typography key={key} variant="body2">
// // // // //                                   <strong>{key}:</strong> {value?.toString()}
// // // // //                                 </Typography>
// // // // //                               ))}
// // // // //                               {idx < msg.results.length - 1 && <Divider sx={{ my: 1 }} />}
// // // // //                             </Box>
// // // // //                           ))}
// // // // //                         </Paper>
// // // // //                       </Box>
// // // // //                     )}
                    
// // // // //                     <Typography variant="caption" sx={{ opacity: 0.7, display: "block", mt: 1 }}>
// // // // //                       {new Date(msg.timestamp).toLocaleTimeString()}
// // // // //                     </Typography>
// // // // //                   </Box>
// // // // //                 </Box>
// // // // //               );
// // // // //             }
// // // // //           })
// // // // //         )}
// // // // //         <div ref={messageEndRef} />
// // // // //       </Paper>

// // // // //       {/* Input Area */}
// // // // //       <Box sx={{ p: 2, borderTop: "1px solid rgba(255,255,255,0.1)", display: "flex", gap: 1 }}>
// // // // //         <TextField
// // // // //           fullWidth
// // // // //           variant="outlined"
// // // // //           placeholder={sessionActive 
// // // // //             ? "Ask a question about your data..." 
// // // // //             : "Upload a CSV file to begin"
// // // // //           }
// // // // //           value={input}
// // // // //           onChange={(e) => setInput(e.target.value)}
// // // // //           onKeyPress={(e) => {
// // // // //             if (e.key === "Enter" && !e.shiftKey) {
// // // // //               e.preventDefault();
// // // // //               handleSend();
// // // // //             }
// // // // //           }}
// // // // //           disabled={!sessionActive || loading}
// // // // //           InputProps={{
// // // // //             endAdornment: !sessionActive && (
// // // // //               <InputAdornment position="end">
// // // // //                 <IconButton component="label" onClick={() => fileInputRef.current.click()}>
// // // // //                   <AttachFileIcon />
// // // // //                 </IconButton>
// // // // //               </InputAdornment>
// // // // //             ),
// // // // //           }}
// // // // //         />
// // // // //         <Button 
// // // // //           variant="contained" 
// // // // //           endIcon={loading ? <CircularProgress size={20} color="inherit" /> : <SendIcon />} 
// // // // //           onClick={handleSend}
// // // // //           disabled={!sessionActive || loading || !input.trim()}
// // // // //         >
// // // // //           Send
// // // // //         </Button>
// // // // //       </Box>
      
// // // // //       {/* CSV Headers */}
// // // // //       {csvHeaders.length > 0 && (
// // // // //         <Box sx={{ p: 1, borderTop: "1px solid rgba(255,255,255,0.1)", display: "flex", flexWrap: "wrap", gap: 0.5 }}>
// // // // //           <Typography variant="caption" color="text.secondary" sx={{ mr: 1, alignSelf: "center" }}>
// // // // //             Available columns:
// // // // //           </Typography>
// // // // //           {csvHeaders.map((header, idx) => (
// // // // //             <Typography key={idx} variant="caption" sx={{ 
// // // // //               bgcolor: "rgba(88, 101, 242, 0.2)", 
// // // // //               px: 0.7, 
// // // // //               py: 0.3, 
// // // // //               borderRadius: 1,
// // // // //               fontSize: "0.7rem"
// // // // //             }}>
// // // // //               {header}
// // // // //             </Typography>
// // // // //           ))}
// // // // //         </Box>
// // // // //       )}
// // // // //     </Box>
// // // // //   );
// // // // // };

// // // // // export default ChatWindow;

// // import React, { useState, useRef, useEffect, useCallback } from "react";
// // import { uploadFile, queryData, fetchChatById, updateChat, createChat } from '../../utils/api';

// // import {
// //   Box,
// //   TextField,
// //   Button,
// //   Typography,
// //   Paper,
// //   IconButton,
// //   InputAdornment,
// //   CircularProgress,
// //   Divider,
// //   Chip,
// //   Tooltip,
// //   Switch,
// //   FormControlLabel,
// // } from "@mui/material";
// // import { 
// //   Send as SendIcon, 
// //   AttachFile as AttachFileIcon, 
// //   TableChart as TableChartIcon,
// //   FileOpen as FileOpenIcon,
// //   Lightbulb as LightbulbIcon,
// //   DataObject as DataObjectIcon,
// //   Code as CodeIcon,
// //   Description as DescriptionIcon
// // } from "@mui/icons-material";
// // import { keyframes } from "@mui/system";

// // const gradientAnimation = keyframes`
// //   0% { background-position: 0% 50%; }
// //   50% { background-position: 100% 50%; }
// //   100% { background-position: 0% 50%; }
// // `;

// // const samplePrompts = [
// //   "Show me the movies with a rating above 8.0",
// //   "What are the top 5 highest grossing movies?",
// //   "List all movies in the action genre sorted by release date",
// // ];

// // const ChatWindow = React.memo(({ selectedChat, loading: externalLoading }) => {
// //   // Define state with useState hooks
// //   const [input, setInput] = useState("");
// //   const [messages, setMessages] = useState([]);
// //   const [loading, setLoading] = useState(false);
// //   const [csvFile, setCsvFile] = useState(null);
// //   const [csvHeaders, setCsvHeaders] = useState([]);
// //   const [sessionActive, setSessionActive] = useState(false);
// //   const [useKnowledgeBase, setUseKnowledgeBase] = useState(true);
// //   const [currentChatId, setCurrentChatId] = useState(null);
// //   const [chatTitle, setChatTitle] = useState("New Chat");
  
// //   // References
// //   const fileInputRef = useRef(null);
// //   const messageEndRef = useRef(null);

// //   // Define resetChat as a useCallback to avoid recreation on every render
// //   const resetChat = useCallback(() => {
// //     console.log("Resetting chat state...");
// //     setCsvFile(null);
// //     setCsvHeaders([]);
// //     setSessionActive(false);
// //     setMessages([]);
// //     setCurrentChatId(null);
// //     setChatTitle("New Chat");
    
// //     // Clear the upload field if you have one
// //     if (fileInputRef.current) {
// //       fileInputRef.current.value = "";
// //     }
// //   }, []);

// //   // Load selected chat data when selectedChat prop changes
// //   const loadSelectedChat = useCallback(async (chat) => {
// //     try {
// //       console.log("Loading selected chat:", chat);
// //       setLoading(true);
      
// //       // Set chat title
// //       setChatTitle(chat.title || "Untitled Chat");
      
// //       // Set current chat ID first
// //       if (chat.id) {
// //         console.log("Setting current chat ID:", chat.id);
// //         setCurrentChatId(chat.id);
// //       }
      
// //       // Fetch full chat data if only ID is provided
// //       let chatData = chat;
// //       if (chat.id && (!chat.messages || !chat.file)) {
// //         try {
// //           console.log("Fetching full chat data for ID:", chat.id);
// //           chatData = await fetchChatById(chat.id);
// //           // Don't set ID here again, already set above
// //         } catch (error) {
// //           console.error("Error fetching chat data:", error);
// //           // Continue with what we have
// //         }
// //       }
      
// //       // Set session data
// //       setSessionActive(true);
      
// //       if (chatData.file) {
// //         console.log("Setting CSV file:", chatData.file);
// //         setCsvFile({ name: chatData.file, path: chatData.file });
// //       }
      
// //       if (chatData.headers && Array.isArray(chatData.headers)) {
// //         console.log("Setting CSV headers:", chatData.headers);
// //         setCsvHeaders(chatData.headers);
// //       }
      
// //       // Load chat messages
// //       if (chatData.messages && Array.isArray(chatData.messages)) {
// //         console.log("Setting messages:", chatData.messages.length);
// //         setMessages(chatData.messages);
// //       } else {
// //         // Initialize with a system message if no messages exist
// //         const initialMessage = {
// //           type: "system",
// //           text: `Loaded session "${chatData.title || 'Untitled'}"${chatData.file ? ` with file: ${chatData.file}` : ''}`,
// //           timestamp: new Date(),
// //           file: chatData.file
// //         };
// //         console.log("Setting initial system message");
// //         setMessages([initialMessage]);
// //       }
// //     } catch (error) {
// //       console.error("Error loading chat:", error);
      
// //       // Initialize with an error message
// //       const errorMessage = {
// //         type: "error",
// //         text: `Error loading chat: ${error.message || "Unknown error"}`,
// //         timestamp: new Date()
// //       };
// //       console.log("Setting error message");
// //       setMessages([errorMessage]);
// //     } finally {
// //       setLoading(false);
// //     }
// //   }, []);

// //   // Scroll to bottom when messages change
// //   useEffect(() => {
// //     messageEndRef.current?.scrollIntoView({ behavior: "smooth" });
// //   }, [messages]);

// //   // Set selected chat if passed from parent
// //   useEffect(() => {
// //     console.log("selectedChat changed:", selectedChat);
    
// //     if (selectedChat) {
// //       loadSelectedChat(selectedChat);
// //     } else {
// //       console.log("No selectedChat, resetting...");
// //       resetChat();
// //     }
// //   }, [selectedChat, loadSelectedChat, resetChat]);

// //   const handleFileUpload = async (event) => {
// //     const file = event.target.files[0];
// //     if (!file) return;
  
// //     setLoading(true);
    
// //     try {
// //       // Upload the CSV file to the backend
// //       const response = await uploadFile(file);
// //       console.log("Upload response:", response);
      
// //       // Store both the file name and the server-provided file path
// //       // Handle different case formats (FilePath vs filePath)
// //       const filePath = response.FilePath || response.filePath;
      
// //       setCsvFile({
// //         name: file.name,
// //         path: filePath
// //       });
      
// //       // Get headers from the response
// //       const headers = response.Headers || response.headers || [];
// //       setCsvHeaders(headers);
// //       setSessionActive(true);
  
// //       // Create a system message
// //       const systemMessage = {
// //         type: "system",
// //         text: `CSV file "${file.name}" uploaded successfully. Found columns: ${headers.length > 0 ? headers.join(", ") : "none detected"}`,
// //         timestamp: new Date(),
// //         file: file.name
// //       };
      
// //       // Update messages state
// //       setMessages([systemMessage]);
      
// //       try {
// //         // Try to create a new chat if none exists
// //         if (!currentChatId) {
// //           const newChat = await createChat({
// //             file: file.name,
// //             filePath: filePath,
// //             headers,
// //             messages: [systemMessage],
// //             title: `${file.name} Analysis`
// //           });
// //           setCurrentChatId(newChat.id);
// //           setChatTitle(newChat.title || "Untitled Chat");
// //         } else {
// //           // Update existing chat
// //           await updateChat(currentChatId, {
// //             file: file.name,
// //             filePath: filePath,
// //             headers,
// //             messages: [systemMessage],
// //             lastUpdated: new Date().toISOString()
// //           });
// //         }
// //       } catch (chatErr) {
// //         console.error("Error managing chat record:", chatErr);
// //         // Continue with the session anyway - the CSV is uploaded and ready to use
// //       }
      
// //       console.log("Stored file path:", filePath);
// //     } catch (error) {
// //       console.error("Error uploading CSV file:", error);
// //       setMessages([
// //         {
// //           type: "error",
// //           text: `Error uploading CSV file: ${error.message || "Unknown error"}`,
// //           timestamp: new Date()
// //         }
// //       ]);
// //     } finally {
// //       setLoading(false);
// //     }
// //   };

// //   const startNewSession = useCallback(() => {
// //     // Save current chat if needed
// //     if (currentChatId && messages.length > 0) {
// //       updateChat(currentChatId, {
// //         messages,
// //         lastUpdated: new Date().toISOString()
// //       }).catch(err => console.error("Error saving chat before new session:", err));
// //     }
    
// //     // Reset for new session
// //     resetChat();
// //   }, [currentChatId, messages, resetChat]);

// //   const handleSend = async () => {
// //     if (!input.trim() || !sessionActive) return;
  
// //     // Create a fresh user message
// //     const userMessage = {
// //       type: "user",
// //       text: input,
// //       sender: "user",
// //       timestamp: new Date()
// //     };
  
// //     // Clear input field immediately to prevent re-submission
// //     const currentInput = input;
// //     setInput("");
    
// //     // Show the user message right away
// //     const updatedMessages = [...messages, userMessage];
// //     setMessages(updatedMessages);
    
// //     // Show loading state
// //     setLoading(true);
  
// //     try {
// //       // Get the file path from csvFile
// //       let filePath = "";
// //       if (typeof csvFile === 'object' && csvFile !== null) {
// //         filePath = csvFile.path || csvFile.name || "";
// //       } else {
// //         filePath = csvFile || "";
// //       }
      
// //       console.log("Sending query:", currentInput, "with file path:", filePath);
      
// //       // Execute the query with fresh options
// //       const response = await queryData(currentInput, filePath, { 
// //         useKnowledgeBase: useKnowledgeBase, 
// //         timestamp: new Date().getTime() // Add timestamp to prevent caching
// //       });
      
// //       console.log("Query response:", response);
      
// //       // Create AI response message
// //       const aiMessage = {
// //         type: "ai",
// //         text: response.response || "Query processed successfully.",
// //         sender: "ai",
// //         timestamp: new Date(),
// //         results: response.results || [],
// //         generatedQuery: response.sql || response.generatedQuery || "",
// //         knowledgeContext: response.knowledgeContext || []
// //       };
      
// //       // Add AI message to chat
// //       const newMessages = [...updatedMessages, aiMessage];
// //       setMessages(newMessages);
      
// //       // Update chat in database
// //       if (currentChatId) {
// //         await updateChat(currentChatId, {
// //           messages: newMessages,
// //           lastUpdated: new Date().toISOString()
// //         });
// //       }
// //     } catch (error) {
// //       console.error("Error processing query:", error);
      
// //       // Add error message
// //       const errorMessage = {
// //         type: "error",
// //         text: `Error processing query: ${error.message || "Unknown error"}`,
// //         sender: "ai",
// //         timestamp: new Date()
// //       };
      
// //       const newMessages = [...updatedMessages, errorMessage];
// //       setMessages(newMessages);
      
// //       // Still update chat in database
// //       if (currentChatId) {
// //         await updateChat(currentChatId, {
// //           messages: newMessages,
// //           lastUpdated: new Date().toISOString()
// //         });
// //       }
// //     } finally {
// //       setLoading(false);
// //     }
// //   };

// //   const handlePromptClick = (prompt) => {
// //     setInput(prompt);
// //   };
  
// //   const updateChatTitle = async (newTitle) => {
// //     if (!currentChatId) return;
    
// //     try {
// //       setChatTitle(newTitle);
// //       await updateChat(currentChatId, { title: newTitle });
// //     } catch (error) {
// //       console.error("Error updating chat title:", error);
// //     }
// //   };

// //   return (
// //     <Box sx={{ flexGrow: 1, display: "flex", flexDirection: "column", height: "100vh", overflow: "hidden" }}>
// //       {/* Header */}
// //       <Paper 
// //         sx={{ 
// //           p: 2, 
// //           display: "flex", 
// //           justifyContent: "space-between", 
// //           alignItems: "center",
// //           borderRadius: 0
// //         }} 
// //         elevation={2}
// //       >
// //         <Typography 
// //           variant="h6" 
// //           sx={{ cursor: 'pointer' }}
// //           onClick={() => {
// //             if (currentChatId) {
// //               const newTitle = prompt("Enter new chat title:", chatTitle);
// //               if (newTitle && newTitle.trim()) {
// //                 updateChatTitle(newTitle.trim());
// //               }
// //             }
// //           }}
// //         >
// //           {sessionActive 
// //             ? `${chatTitle} - ${csvFile?.name || "CSV File"}`
// //             : "Sage AI - Upload a CSV file to begin"}
// //         </Typography>
// //         <Box sx={{ display: 'flex', alignItems: 'center' }}>
// //           {sessionActive && (
// //             <FormControlLabel
// //               control={
// //                 <Switch
// //                   checked={useKnowledgeBase}
// //                   onChange={(e) => setUseKnowledgeBase(e.target.checked)}
// //                   size="small"
// //                 />
// //               }
// //               label={
// //                 <Box sx={{ display: 'flex', alignItems: 'center' }}>
// //                   <DataObjectIcon fontSize="small" sx={{ mr: 0.5 }} />
// //                   <Typography variant="body2">Use Knowledge Base</Typography>
// //                 </Box>
// //               }
// //               sx={{ mr: 2 }}
// //             />
// //           )}
          
// //           {sessionActive ? (
// //             <Button 
// //               variant="outlined" 
// //               size="small"
// //               onClick={startNewSession}
// //               startIcon={<FileOpenIcon />}
// //             >
// //               New Session
// //             </Button>
// //           ) : (
// //             <Button
// //               variant="contained"
// //               onClick={() => fileInputRef.current?.click()}
// //               startIcon={<AttachFileIcon />}
// //               disabled={loading}
// //             >
// //               Upload CSV
// //             </Button>
// //           )}
// //         </Box>
// //         <input
// //           type="file"
// //           accept=".csv"
// //           ref={fileInputRef}
// //           style={{ display: "none" }}
// //           onChange={handleFileUpload}
// //         />
// //       </Paper>

// //       {/* Message Area */}
// //       <Paper sx={{ flexGrow: 1, display: "flex", flexDirection: "column", p: 2, overflowY: "auto" }}>
// //         {messages.length === 0 ? (
// //           <Box
// //             sx={{
// //               display: "flex",
// //               flexDirection: "column",
// //               alignItems: "center",
// //               justifyContent: "center",
// //               height: "100%",
// //             }}
// //           >
// //             <Typography
// //               variant="h5"
// //               sx={{
// //                 background: "linear-gradient(90deg, #00D4FF, #FF4D4D)",
// //                 backgroundSize: "200% 200%",
// //                 animation: `${gradientAnimation} 5s ease infinite`,
// //                 WebkitBackgroundClip: "text",
// //                 WebkitTextFillColor: "transparent",
// //                 fontWeight: 700,
// //                 mb: 2,
// //               }}
// //             >
// //               Welcome to Sage AI Chat!
// //             </Typography>
// //             <Typography variant="body1" sx={{ mb: 1 }}>
// //               Upload a CSV file and try asking questions like:
// //             </Typography>
// //             {samplePrompts.map((prompt, index) => (
// //               <Typography
// //                 key={index}
// //                 variant="body2"
// //                 onClick={() => handlePromptClick(prompt)}
// //                 sx={{
// //                   mb: 0.5,
// //                   background: "linear-gradient(90deg, #FF4D4D, #00D4FF)",
// //                   backgroundSize: "200% 200%",
// //                   animation: `${gradientAnimation} 5s ease infinite`,
// //                   WebkitBackgroundClip: "text",
// //                   WebkitTextFillColor: "transparent",
// //                   cursor: "pointer",
// //                   "&:hover": {
// //                     opacity: 0.8,
// //                   }
// //                 }}
// //               >
// //                 {prompt}
// //               </Typography>
// //             ))}
// //             <Button
// //               variant="contained"
// //               sx={{ mt: 3 }}
// //               onClick={() => fileInputRef.current?.click()}
// //               startIcon={<AttachFileIcon />}
// //             >
// //               Upload CSV
// //             </Button>
// //           </Box>
// //         ) : (
// //           messages.map((msg, index) => {
// //             if (msg.type === "system") {
// //               return (
// //                 <Box 
// //                   key={index}
// //                   sx={{ 
// //                     display: "flex", 
// //                     justifyContent: "center", 
// //                     mb: 2 
// //                   }}
// //                 >
// //                   <Paper 
// //                     sx={{ 
// //                       p: 1, 
// //                       backgroundColor: "rgba(0,0,0,0.2)", 
// //                       color: "#dddddd", 
// //                       maxWidth: "90%",
// //                       borderRadius: 2
// //                     }}
// //                   >
// //                     <Typography variant="body2" align="center">
// //                       {msg.file && (
// //                         <Box component="span" sx={{ display: "flex", alignItems: "center", justifyContent: "center", mb: 1 }}>
// //                           <TableChartIcon fontSize="small" sx={{ mr: 1 }} />
// //                           <Typography component="span" variant="caption">{msg.file}</Typography>
// //                         </Box>
// //                       )}
// //                       {msg.text}
// //                     </Typography>
// //                   </Paper>
// //                 </Box>
// //               );
// //             } else if (msg.type === "error") {
// //               return (
// //                 <Box 
// //                   key={index}
// //                   sx={{ 
// //                     display: "flex", 
// //                     justifyContent: "center", 
// //                     mb: 2 
// //                   }}
// //                 >
// //                   <Paper 
// //                     sx={{ 
// //                       p: 1, 
// //                       backgroundColor: "#770000", 
// //                       color: "#ffffff", 
// //                       maxWidth: "90%",
// //                       borderRadius: 2
// //                     }}
// //                   >
// //                     <Typography variant="body2" align="center">
// //                       {msg.text}
// //                     </Typography>
// //                   </Paper>
// //                 </Box>
// //               );
// //             } else {
// //               return (
// //                 <Box 
// //                   key={index} 
// //                   sx={{ 
// //                     display: "flex", 
// //                     justifyContent: msg.sender === "user" ? "flex-end" : "flex-start", 
// //                     mb: 2 
// //                   }}
// //                 >
// //                   <Box
// //                     sx={{
// //                       p: 1.5,
// //                       borderRadius: 2,
// //                       maxWidth: "70%",
// //                       backgroundColor: msg.sender === "user" ? "primary.main" : "background.default",
// //                       color: msg.sender === "user" ? "white" : "text.primary",
// //                     }}
// //                   >
// //                     <Typography variant="body1">{msg.text}</Typography>
                    
// //                     {/* Knowledge Context Used */}
// //                     {msg.knowledgeContext && msg.knowledgeContext.length > 0 && (
// //                       <Box sx={{ mt: 1, mb: 1 }}>
// //                         <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
// //                           <LightbulbIcon fontSize="small" sx={{ mr: 0.5, color: 'primary.main' }} />
// //                           <Typography variant="caption" color="primary.main">
// //                             Knowledge Base Used:
// //                           </Typography>
// //                         </Box>
// //                         <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
// //                           {msg.knowledgeContext.map((item, i) => (
// //                             <Tooltip 
// //                               key={i} 
// //                               title={
// //                                 <Box>
// //                                   <Typography variant="caption" sx={{ fontWeight: 'bold' }}>
// //                                     {item.description}
// //                                   </Typography>
// //                                   <Typography variant="caption" sx={{ display: 'block', mt: 0.5 }}>
// //                                     {item.type === 'question_sql' 
// //                                       ? `Q: ${item.question} | SQL: ${item.sql}`
// //                                       : item.content?.substring(0, 200) + (item.content?.length > 200 ? '...' : '')
// //                                     }
// //                                   </Typography>
// //                                 </Box>
// //                               }
// //                             >
// //                               <Chip 
// //                                 size="small"
// //                                 label={item.description}
// //                                 icon={
// //                                   item.type === 'ddl' ? <CodeIcon fontSize="small" /> :
// //                                   item.type === 'documentation' ? <DescriptionIcon fontSize="small" /> :
// //                                   <DataObjectIcon fontSize="small" />
// //                                 }
// //                                 sx={{ bgcolor: 'rgba(88, 101, 242, 0.2)' }}
// //                               />
// //                             </Tooltip>
// //                           ))}
// //                         </Box>
// //                       </Box>
// //                     )}
                    
// //                     {msg.generatedQuery && (
// //                       <Box sx={{ mt: 1, mb: 1, p: 1, bgcolor: "rgba(0,0,0,0.3)", borderRadius: 1 }}>
// //                         <Typography variant="caption" color="text.secondary">
// //                           Generated SQL:
// //                         </Typography>
// //                         <Typography variant="body2" component="pre" sx={{ 
// //                           overflowX: "auto", 
// //                           whiteSpace: "pre-wrap",
// //                           wordBreak: "break-word" 
// //                         }}>
// //                           {msg.generatedQuery}
// //                         </Typography>
// //                       </Box>
// //                     )}
                    
// //                     {msg.results && msg.results.length > 0 && (
// //                       <Box sx={{ mt: 1 }}>
// //                         <Typography variant="subtitle2" sx={{ mb: 1 }}>
// //                           Results:
// //                         </Typography>
// //                         <Paper variant="outlined" sx={{ p: 1, bgcolor: "rgba(0,0,0,0.2)" }}>
// //                           {msg.results.map((result, idx) => (
// //                             <Box key={idx} sx={{ mb: idx < msg.results.length - 1 ? 1 : 0 }}>
// //                               {Object.entries(result).map(([key, value]) => (
// //                                 <Typography key={key} variant="body2">
// //                                   <strong>{key}:</strong> {value?.toString()}
// //                                 </Typography>
// //                               ))}
// //                               {idx < msg.results.length - 1 && <Divider sx={{ my: 1 }} />}
// //                             </Box>
// //                           ))}
// //                         </Paper>
// //                       </Box>
// //                     )}
                    
// //                     <Typography variant="caption" sx={{ opacity: 0.7, display: "block", mt: 1 }}>
// //                       {new Date(msg.timestamp).toLocaleTimeString()}
// //                     </Typography>
// //                   </Box>
// //                 </Box>
// //               );
// //             }
// //           })
// //         )}
// //         <div ref={messageEndRef} />
// //       </Paper>

// //       {/* Input Area */}
// //       <Box sx={{ p: 2, borderTop: "1px solid rgba(255,255,255,0.1)", display: "flex", gap: 1 }}>
// //         <TextField
// //           fullWidth
// //           variant="outlined"
// //           placeholder={sessionActive 
// //             ? "Ask a question about your data..." 
// //             : "Upload a CSV file to begin"
// //           }
// //           value={input}
// //           onChange={(e) => setInput(e.target.value)}
// //           onKeyPress={(e) => {
// //             if (e.key === "Enter" && !e.shiftKey) {
// //               e.preventDefault();
// //               handleSend();
// //             }
// //           }}
// //           disabled={!sessionActive || loading}
// //           InputProps={{
// //             endAdornment: !sessionActive && (
// //               <InputAdornment position="end">
// //                 <IconButton component="label" onClick={() => fileInputRef.current?.click()}>
// //                   <AttachFileIcon />
// //                 </IconButton>
// //               </InputAdornment>
// //             ),
// //           }}
// //         />
// //         <Button 
// //           variant="contained" 
// //           endIcon={loading ? <CircularProgress size={20} color="inherit" /> : <SendIcon />} 
// //           onClick={handleSend}
// //           disabled={!sessionActive || loading || !input.trim()}
// //         >
// //           Send
// //         </Button>
// //       </Box>
      
// //       {/* CSV Headers */}
// //       {csvHeaders.length > 0 && (
// //         <Box sx={{ p: 1, borderTop: "1px solid rgba(255,255,255,0.1)", display: "flex", flexWrap: "wrap", gap: 0.5 }}>
// //           <Typography variant="caption" color="text.secondary" sx={{ mr: 1, alignSelf: "center" }}>
// //             Available columns:
// //           </Typography>
// //           {csvHeaders.map((header, idx) => (
// //             <Typography key={idx} variant="caption" sx={{ 
// //               bgcolor: "rgba(88, 101, 242, 0.2)", 
// //               px: 0.7, 
// //               py: 0.3, 
// //               borderRadius: 1,
// //               fontSize: "0.7rem"
// //             }}>
// //               {header}
// //             </Typography>
// //           ))}
// //         </Box>
// //       )}
// //     </Box>
// //   );
// // });

// // export default ChatWindow;

// // frontend/src/components/chat/ChatWindow.jsx
// import React, { useState, useRef, useEffect, useCallback } from "react";
// import { uploadFile, queryData, fetchChatById, updateChat, createChat } from '../../utils/api';

// import {
//   Box,
//   TextField,
//   Button,
//   Typography,
//   Paper,
//   IconButton,
//   InputAdornment,
//   CircularProgress,
//   Divider,
//   Chip,
//   Tooltip,
//   Switch,
//   FormControlLabel,
//   FormControl,
//   InputLabel,
//   Select,
//   MenuItem,
//   Dialog,
//   DialogTitle,
//   DialogContent,
//   DialogActions,
//   Alert,
// } from "@mui/material";
// import { 
//   Send as SendIcon, 
//   AttachFile as AttachFileIcon, 
//   TableChart as TableChartIcon,
//   FileOpen as FileOpenIcon,
//   Lightbulb as LightbulbIcon,
//   DataObject as DataObjectIcon,
//   Code as CodeIcon,
//   Description as DescriptionIcon,
//   Key as KeyIcon,
// } from "@mui/icons-material";
// import { keyframes } from "@mui/system";

// const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

// const gradientAnimation = keyframes`
//   0% { background-position: 0% 50%; }
//   50% { background-position: 100% 50%; }
//   100% { background-position: 0% 50%; }
// `;

// const samplePrompts = [
//   "Show me the movies with a rating above 8.0",
//   "What are the top 5 highest grossing movies?",
//   "List all movies in the action genre sorted by release date",
// ];

// const ChatWindow = React.memo(({ selectedChat, loading: externalLoading }) => {
//   // Define state with useState hooks
//   const [input, setInput] = useState("");
//   const [messages, setMessages] = useState([]);
//   const [loading, setLoading] = useState(false);
//   const [csvFile, setCsvFile] = useState(null);
//   const [csvHeaders, setCsvHeaders] = useState([]);
//   const [sessionActive, setSessionActive] = useState(false);
//   const [useKnowledgeBase, setUseKnowledgeBase] = useState(true);
//   const [currentChatId, setCurrentChatId] = useState(null);
//   const [chatTitle, setChatTitle] = useState("New Chat");
//   const [llmProvider, setLLMProvider] = useState(null);
//   const [apiKey, setApiKey] = useState(null);
//   const [llmConfigDialog, setLLMConfigDialog] = useState(false);
//   const [apiKeyError, setApiKeyError] = useState("");
//   const [llmModel, setLLMModel] = useState(null);
  
//   // References
//   const fileInputRef = useRef(null);
//   const messageEndRef = useRef(null);

//   const modelOptions = {
//     gemini: ["gemini-1.5-flash", "gemini-1.5-pro"],
//     openai: ["gpt-4", "gpt-4-turbo", "gpt-3.5-turbo"],
//     anthropic: ["claude-3-opus-20240229", "claude-3-sonnet-20240229", "claude-3-haiku-20240307"],
//     mistral: ["mistral-large-latest", "mistral-medium"]
//   };

//   // Define resetChat as a useCallback to avoid recreation on every render
//   const resetChat = useCallback(() => {
//     console.log("Resetting chat state...");
//     setCsvFile(null);
//     setCsvHeaders([]);
//     setSessionActive(false);
//     setMessages([]);
//     setCurrentChatId(null);
//     setChatTitle("New Chat");
//     setLLMProvider(null);
//     setApiKey(null);
    
//     // Clear the upload field if you have one
//     if (fileInputRef.current) {
//       fileInputRef.current.value = "";
//     }
//   }, []);

//   // Load LLM config from the selected chat
//   useEffect(() => {
//     if (selectedChat && selectedChat.llmConfig) {
//       setLLMProvider(selectedChat.llmConfig.provider);
//       setApiKey(selectedChat.llmConfig.api_key);
//     } else {
//       setLLMProvider(null);
//       setApiKey(null);
//     }
//   }, [selectedChat]);

//   // Handle LLM configuration
//   const handleLLMConfig = async () => {
//     if (!llmProvider || !apiKey) {
//       setApiKeyError("Please select a provider and enter an API key");
//       return;
//     }
    
//     try {
//       console.log("Validating API key for provider:", llmProvider);
      
//       // Validate API key with backend
//       const response = await fetch(`${API_BASE_URL}/api/validate-api-key`, {
//         method: 'POST',
//         headers: {
//           'Content-Type': 'application/json',
//         },
//         body: JSON.stringify({
//           provider: llmProvider,
//           api_key: apiKey
//         }),
//       });
      
//       console.log("Validation response status:", response.status);
      
//       if (!response.ok) {
//         const errorText = await response.text();
//         console.error("Validation failed:", errorText);
//         setApiKeyError(`Validation failed: ${errorText}`);
//         return;
//       }
      
//       const result = await response.json();
//       console.log("Validation result:", result);
      
//       if (result.valid) {
//         // Update chat with LLM config
//         if (currentChatId) {
//           await updateChat(currentChatId, {
//             llmConfig: {
//               provider: llmProvider,
//               api_key: apiKey
//             }
//           });
//         }
//         setLLMConfigDialog(false);
//         setApiKeyError("");
//       } else {
//         setApiKeyError("Invalid API key for selected provider");
//       }
//     } catch (error) {
//       console.error("Error validating API key:", error);
//       setApiKeyError(`Failed to validate API key: ${error.message}`);
//     }
//   };

//   // Check if provider change is needed
//   const handleProviderChange = (newProvider) => {
//     if (sessionActive && llmProvider) {
//       // Show warning dialog
//       if (window.confirm("Changing LLM provider requires starting a new chat. Continue?")) {
//         startNewSession();
//         setLLMProvider(newProvider);
//       }
//     } else {
//       setLLMProvider(newProvider);
//     }
//   };

//   // Load selected chat data when selectedChat prop changes
//   const loadSelectedChat = useCallback(async (chat) => {
//     try {
//       console.log("Loading selected chat:", chat);
//       setLoading(true);
      
//       // Set chat title
//       setChatTitle(chat.title || "Untitled Chat");
      
//       // Set current chat ID first
//       if (chat.id) {
//         console.log("Setting current chat ID:", chat.id);
//         setCurrentChatId(chat.id);
//       }
      
//       // Fetch full chat data if only ID is provided
//       let chatData = chat;
//       if (chat.id && (!chat.messages || !chat.file)) {
//         try {
//           console.log("Fetching full chat data for ID:", chat.id);
//           chatData = await fetchChatById(chat.id);
//         } catch (error) {
//           console.error("Error fetching chat data:", error);
//           // Continue with what we have
//         }
//       }
      
//       // Set session data
//       setSessionActive(true);
      
//       if (chatData.file) {
//         console.log("Setting CSV file:", chatData.file);
//         setCsvFile({ name: chatData.file, path: chatData.file });
//       }
      
//       if (chatData.headers && Array.isArray(chatData.headers)) {
//         console.log("Setting CSV headers:", chatData.headers);
//         setCsvHeaders(chatData.headers);
//       }
      
//       // Load LLM config if available
//       if (chatData.llmConfig) {
//         setLLMProvider(chatData.llmConfig.provider);
//         setApiKey(chatData.llmConfig.api_key);
//       }
      
//       // Load chat messages
//       if (chatData.messages && Array.isArray(chatData.messages)) {
//         console.log("Setting messages:", chatData.messages.length);
//         setMessages(chatData.messages);
//       } else {
//         // Initialize with a system message if no messages exist
//         const initialMessage = {
//           type: "system",
//           text: `Loaded session "${chatData.title || 'Untitled'}"${chatData.file ? ` with file: ${chatData.file}` : ''}`,
//           timestamp: new Date(),
//           file: chatData.file
//         };
//         console.log("Setting initial system message");
//         setMessages([initialMessage]);
//       }
//     } catch (error) {
//       console.error("Error loading chat:", error);
      
//       // Initialize with an error message
//       const errorMessage = {
//         type: "error",
//         text: `Error loading chat: ${error.message || "Unknown error"}`,
//         timestamp: new Date()
//       };
//       console.log("Setting error message");
//       setMessages([errorMessage]);
//     } finally {
//       setLoading(false);
//     }
//   }, []);

//   // Scroll to bottom when messages change
//   useEffect(() => {
//     messageEndRef.current?.scrollIntoView({ behavior: "smooth" });
//   }, [messages]);

//   // Set selected chat if passed from parent
//   useEffect(() => {
//     console.log("selectedChat changed:", selectedChat);
    
//     if (selectedChat) {
//       loadSelectedChat(selectedChat);
//     } else {
//       console.log("No selectedChat, resetting...");
//       resetChat();
//     }
//   }, [selectedChat, loadSelectedChat, resetChat]);

//   const handleFileUpload = async (event) => {
//     const file = event.target.files[0];
//     if (!file) return;
  
//     setLoading(true);
    
//     try {
//       // Upload the CSV file to the backend
//       const response = await uploadFile(file);
//       console.log("Upload response:", response);
      
//       // Store both the file name and the server-provided file path
//       const filePath = response.FilePath || response.filePath;
      
//       setCsvFile({
//         name: file.name,
//         path: filePath
//       });
      
//       // Get headers from the response
//       const headers = response.Headers || response.headers || [];
//       setCsvHeaders(headers);
//       setSessionActive(true);
  
//       // Create a system message
//       const systemMessage = {
//         type: "system",
//         text: `CSV file "${file.name}" uploaded successfully. Found columns: ${headers.length > 0 ? headers.join(", ") : "none detected"}`,
//         timestamp: new Date(),
//         file: file.name
//       };
      
//       // Update messages state
//       setMessages([systemMessage]);
      
//       try {
//         // Try to create a new chat if none exists
//         if (!currentChatId) {
//           const newChat = await createChat({
//             file: file.name,
//             filePath: filePath,
//             headers,
//             messages: [systemMessage],
//             title: `${file.name} Analysis`,
//             llmConfig: llmProvider && apiKey ? {
//               provider: llmProvider,
//               api_key: apiKey
//             } : null
//           });
//           setCurrentChatId(newChat.id);
//           setChatTitle(newChat.title || "Untitled Chat");
//         } else {
//           // Update existing chat
//           await updateChat(currentChatId, {
//             file: file.name,
//             filePath: filePath,
//             headers,
//             messages: [systemMessage],
//             lastUpdated: new Date().toISOString()
//           });
//         }
//       } catch (chatErr) {
//         console.error("Error managing chat record:", chatErr);
//         // Continue with the session anyway - the CSV is uploaded and ready to use
//       }
      
//       console.log("Stored file path:", filePath);
//     } catch (error) {
//       console.error("Error uploading CSV file:", error);
//       setMessages([
//         {
//           type: "error",
//           text: `Error uploading CSV file: ${error.message || "Unknown error"}`,
//           timestamp: new Date()
//         }
//       ]);
//     } finally {
//       setLoading(false);
//     }
//   };

//   const startNewSession = useCallback(() => {
//     // Save current chat if needed
//     if (currentChatId && messages.length > 0) {
//       updateChat(currentChatId, {
//         messages,
//         lastUpdated: new Date().toISOString()
//       }).catch(err => console.error("Error saving chat before new session:", err));
//     }
    
//     // Reset for new session
//     resetChat();
//   }, [currentChatId, messages, resetChat]);

//   const handleSend = async () => {
//     if (!input.trim() || !sessionActive) return;
    
//     // Check if LLM is configured
//     if (!llmProvider || !apiKey) {
//       setLLMConfigDialog(true);
//       return;
//     }
  
//     // Create a fresh user message
//     const userMessage = {
//       type: "user",
//       text: input,
//       sender: "user",
//       timestamp: new Date()
//     };
  
//     // Clear input field immediately to prevent re-submission
//     const currentInput = input;
//     setInput("");
    
//     // Show the user message right away
//     const updatedMessages = [...messages, userMessage];
//     setMessages(updatedMessages);
    
//     // Show loading state
//     setLoading(true);
  
//     try {
//       // Get the file path from csvFile
//       let filePath = "";
//       if (typeof csvFile === 'object' && csvFile !== null) {
//         filePath = csvFile.path || csvFile.name || "";
//       } else {
//         filePath = csvFile || "";
//       }
      
//       console.log("Sending query:", currentInput, "with file path:", filePath);
      
//       // Execute the query with LLM config
//       const response = await queryData(currentInput, filePath, { 
//         useKnowledgeBase: useKnowledgeBase, 
//         timestamp: new Date().getTime(),
//         chatId: currentChatId,
//         llmConfig: {
//           provider: llmProvider,
//           api_key: apiKey
//         }
//       });
      
//       console.log("Query response:", response);
      
//       // Create AI response message
//       const aiMessage = {
//         type: "ai",
//         text: response.response || "Query processed successfully.",
//         sender: "ai",
//         timestamp: new Date(),
//         results: response.results || [],
//         generatedQuery: response.sql || response.generatedQuery || "",
//         knowledgeContext: response.knowledgeContext || []
//       };
      
//       // Add AI message to chat
//       const newMessages = [...updatedMessages, aiMessage];
//       setMessages(newMessages);
      
//       // Update chat in database
//       if (currentChatId) {
//         await updateChat(currentChatId, {
//           messages: newMessages,
//           lastUpdated: new Date().toISOString()
//         });
//       }
//     } catch (error) {
//       console.error("Error processing query:", error);
      
//       // Add error message
//       const errorMessage = {
//         type: "error",
//         text: `Error processing query: ${error.message || "Unknown error"}`,
//         sender: "ai",
//         timestamp: new Date()
//       };
      
//       const newMessages = [...updatedMessages, errorMessage];
//       setMessages(newMessages);
      
//       // Still update chat in database
//       if (currentChatId) {
//         await updateChat(currentChatId, {
//           messages: newMessages,
//           lastUpdated: new Date().toISOString()
//         });
//       }
//     } finally {
//       setLoading(false);
//     }
//   };

//   const handlePromptClick = (prompt) => {
//     setInput(prompt);
//   };
  
//   const updateChatTitle = async (newTitle) => {
//     if (!currentChatId) return;
    
//     try {
//       setChatTitle(newTitle);
//       await updateChat(currentChatId, { title: newTitle });
//     } catch (error) {
//       console.error("Error updating chat title:", error);
//     }
//   };

//   return (
//     <Box sx={{ flexGrow: 1, display: "flex", flexDirection: "column", height: "100vh", overflow: "hidden" }}>
//       {/* Header */}
//       <Paper 
//         sx={{ 
//           p: 2, 
//           display: "flex", 
//           justifyContent: "space-between", 
//           alignItems: "center",
//           borderRadius: 0
//         }} 
//         elevation={2}
//       >
//         <Typography 
//           variant="h6" 
//           sx={{ cursor: 'pointer' }}
//           onClick={() => {
//             if (currentChatId) {
//               const newTitle = prompt("Enter new chat title:", chatTitle);
//               if (newTitle && newTitle.trim()) {
//                 updateChatTitle(newTitle.trim());
//               }
//             }
//           }}
//         >
//           {sessionActive 
//             ? `${chatTitle} - ${csvFile?.name || "CSV File"}`
//             : "Sage AI - Upload a CSV file to begin"}
//         </Typography>
//         <Box sx={{ display: 'flex', alignItems: 'center' }}>
//           {sessionActive && (
//             <FormControlLabel
//               control={
//                 <Switch
//                   checked={useKnowledgeBase}
//                   onChange={(e) => setUseKnowledgeBase(e.target.checked)}
//                   size="small"
//                 />
//               }
//               label={
//                 <Box sx={{ display: 'flex', alignItems: 'center' }}>
//                   <DataObjectIcon fontSize="small" sx={{ mr: 0.5 }} />
//                   <Typography variant="body2">Use Knowledge Base</Typography>
//                 </Box>
//               }
//               sx={{ mr: 2 }}
//             />
//           )}
          
//           {sessionActive ? (
//             <Button 
//               variant="outlined" 
//               size="small"
//               onClick={startNewSession}
//               startIcon={<FileOpenIcon />}
//             >
//               New Session
//             </Button>
//           ) : (
//             <Button
//               variant="contained"
//               onClick={() => fileInputRef.current?.click()}
//               startIcon={<AttachFileIcon />}
//               disabled={loading}
//             >
//               Upload CSV
//             </Button>
//           )}
//         </Box>
//         <input
//           type="file"
//           accept=".csv"
//           ref={fileInputRef}
//           style={{ display: "none" }}
//           onChange={handleFileUpload}
//         />
//       </Paper>

//       {/* LLM Configuration Alert */}
//       {!llmProvider && sessionActive && (
//         <Paper
//           sx={{
//             p: 2,
//             m: 2,
//             backgroundColor: "rgba(88, 101, 242, 0.1)",
//             borderRadius: 2,
//             display: "flex",
//             alignItems: "center",
//             justifyContent: "space-between"
//           }}
//         >
//           <Box sx={{ display: "flex", alignItems: "center" }}>
//             <KeyIcon sx={{ mr: 1 }} />
//             <Typography>Configure LLM provider and API key to start chatting</Typography>
//           </Box>
//           <Button
//             variant="contained"
//             size="small"
//             onClick={() => setLLMConfigDialog(true)}
//           >
//             Configure
//           </Button>
//         </Paper>
//       )}

//       {/* LLM provider configured indicator */}
//       {llmProvider && sessionActive && (
//         <Chip
//           icon={<KeyIcon />}
//           label={`${llmProvider.toUpperCase()} configured`}
//           sx={{ m: 2 }}
//           color="primary"
//           variant="outlined"
//           onDelete={() => setLLMConfigDialog(true)}
//           deleteIcon={<Tooltip title="Change LLM Settings"><KeyIcon /></Tooltip>}
//         />
//       )}

//       {/* Message Area */}
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
//               Upload a CSV file and try asking questions like:
//             </Typography>
//             {samplePrompts.map((prompt, index) => (
//               <Typography
//                 key={index}
//                 variant="body2"
//                 onClick={() => handlePromptClick(prompt)}
//                 sx={{
//                   mb: 0.5,
//                   background: "linear-gradient(90deg, #FF4D4D, #00D4FF)",
//                   backgroundSize: "200% 200%",
//                   animation: `${gradientAnimation} 5s ease infinite`,
//                   WebkitBackgroundClip: "text",
//                   WebkitTextFillColor: "transparent",
//                   cursor: "pointer",
//                   "&:hover": {
//                     opacity: 0.8,
//                   }
//                 }}
//               >
//                 {prompt}
//               </Typography>
//             ))}
//             <Button
//               variant="contained"
//               sx={{ mt: 3 }}
//               onClick={() => fileInputRef.current?.click()}
//               startIcon={<AttachFileIcon />}
//             >
//               Upload CSV
//             </Button>
//           </Box>
//         ) : (
//           messages.map((msg, index) => {
//             if (msg.type === "system") {
//               return (
//                 <Box 
//                   key={index}
//                   sx={{ 
//                     display: "flex", 
//                     justifyContent: "center", 
//                     mb: 2 
//                   }}
//                 >
//                   <Paper 
//                     sx={{ 
//                       p: 1, 
//                       backgroundColor: "rgba(0,0,0,0.2)", 
//                       color: "#dddddd", 
//                       maxWidth: "90%",
//                       borderRadius: 2
//                     }}
//                   >
//                     <Typography variant="body2" align="center">
//                       {msg.file && (
//                         <Box component="span" sx={{ display: "flex", alignItems: "center", justifyContent: "center", mb: 1 }}>
//                           <TableChartIcon fontSize="small" sx={{ mr: 1 }} />
//                           <Typography component="span" variant="caption">{msg.file}</Typography>
//                         </Box>
//                       )}
//                       {msg.text}
//                     </Typography>
//                   </Paper>
//                 </Box>
//               );
//             } else if (msg.type === "error") {
//               return (
//                 <Box 
//                   key={index}
//                   sx={{ 
//                     display: "flex", 
//                     justifyContent: "center", 
//                     mb: 2 
//                   }}
//                 >
//                   <Paper 
//                     sx={{ 
//                       p: 1, 
//                       backgroundColor: "#770000", 
//                       color: "#ffffff", 
//                       maxWidth: "90%",
//                       borderRadius: 2
//                     }}
//                   >
//                     <Typography variant="body2" align="center">
//                       {msg.text}
//                     </Typography>
//                   </Paper>
//                 </Box>
//               );
//             } else {
//               return (
//                 <Box 
//                   key={index} 
//                   sx={{ 
//                     display: "flex", 
//                     justifyContent: msg.sender === "user" ? "flex-end" : "flex-start", 
//                     mb: 2 
//                   }}
//                 >
//                   <Box
//                     sx={{
//                       p: 1.5,
//                       borderRadius: 2,
//                       maxWidth: "70%",
//                       backgroundColor: msg.sender === "user" ? "primary.main" : "background.default",
//                       color: msg.sender === "user" ? "white" : "text.primary",
//                     }}
//                   >
//                     <Typography variant="body1">{msg.text}</Typography>
                    
//                     {/* Knowledge Context Used */}
//                     {msg.knowledgeContext && msg.knowledgeContext.length > 0 && (
//                       <Box sx={{ mt: 1, mb: 1 }}>
//                         <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
//                           <LightbulbIcon fontSize="small" sx={{ mr: 0.5, color: 'primary.main' }} />
//                           <Typography variant="caption" color="primary.main">
//                             Knowledge Base Used:
//                           </Typography>
//                         </Box>
//                         <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
//                           {msg.knowledgeContext.map((item, i) => (
//                             <Tooltip 
//                               key={i} 
//                               title={
//                                 <Box>
//                                   <Typography variant="caption" sx={{ fontWeight: 'bold' }}>
//                                     {item.description}
//                                   </Typography>
//                                   <Typography variant="caption" sx={{ display: 'block', mt: 0.5 }}>
//                                     {item.type === 'question_sql' 
//                                       ? `Q: ${item.question} | SQL: ${item.sql}`
//                                       : item.content?.substring(0, 200) + (item.content?.length > 200 ? '...' : '')
//                                     }
//                                   </Typography>
//                                 </Box>
//                               }
//                             >
//                               <Chip 
//                                 size="small"
//                                 label={item.description}
//                                 icon={
//                                   item.type === 'ddl' ? <CodeIcon fontSize="small" /> :
//                                   item.type === 'documentation' ? <DescriptionIcon fontSize="small" /> :
//                                   <DataObjectIcon fontSize="small" />
//                                 }
//                                 sx={{ bgcolor: 'rgba(88, 101, 242, 0.2)' }}
//                               />
//                             </Tooltip>
//                           ))}
//                         </Box>
//                       </Box>
//                     )}
                    
//                     {msg.generatedQuery && (
//                       <Box sx={{ mt: 1, mb: 1, p: 1, bgcolor: "rgba(0,0,0,0.3)", borderRadius: 1 }}>
//                         <Typography variant="caption" color="text.secondary">
//                           Generated SQL:
//                         </Typography>
//                         <Typography variant="body2" component="pre" sx={{ 
//                           overflowX: "auto", 
//                           whiteSpace: "pre-wrap",
//                           wordBreak: "break-word" 
//                         }}>
//                           {msg.generatedQuery}
//                         </Typography>
//                       </Box>
//                     )}
                    
//                     {msg.results && msg.results.length > 0 && (
//                       <Box sx={{ mt: 1 }}>
//                         <Typography variant="subtitle2" sx={{ mb: 1 }}>
//                           Results:
//                         </Typography>
//                         <Paper variant="outlined" sx={{ p: 1, bgcolor: "rgba(0,0,0,0.2)" }}>
//                           {msg.results.map((result, idx) => (
//                             <Box key={idx} sx={{ mb: idx < msg.results.length - 1 ? 1 : 0 }}>
//                               {Object.entries(result).map(([key, value]) => (
//                                 <Typography key={key} variant="body2">
//                                   <strong>{key}:</strong> {value?.toString()}
//                                 </Typography>
//                               ))}
//                               {idx < msg.results.length - 1 && <Divider sx={{ my: 1 }} />}
//                             </Box>
//                           ))}
//                         </Paper>
//                       </Box>
//                     )}
                    
//                     <Typography variant="caption" sx={{ opacity: 0.7, display: "block", mt: 1 }}>
//                       {new Date(msg.timestamp).toLocaleTimeString()}
//                     </Typography>
//                   </Box>
//                 </Box>
//               );
//             }
//           })
//         )}
//         <div ref={messageEndRef} />
//       </Paper>

//       {/* Input Area */}
//       <Box sx={{ p: 2, borderTop: "1px solid rgba(255,255,255,0.1)", display: "flex", gap: 1 }}>
//         <TextField
//           fullWidth
//           variant="outlined"
//           placeholder={sessionActive 
//             ? "Ask a question about your data..." 
//             : "Upload a CSV file to begin"
//           }
//           value={input}
//           onChange={(e) => setInput(e.target.value)}
//           onKeyPress={(e) => {
//             if (e.key === "Enter" && !e.shiftKey) {
//               e.preventDefault();
//               handleSend();
//             }
//           }}
//           disabled={!sessionActive || loading}
//           InputProps={{
//             endAdornment: !sessionActive && (
//               <InputAdornment position="end">
//                 <IconButton component="label" onClick={() => fileInputRef.current?.click()}>
//                   <AttachFileIcon />
//                 </IconButton>
//               </InputAdornment>
//             ),
//           }}
//         />
//         <Button 
//           variant="contained" 
//           endIcon={loading ? <CircularProgress size={20} color="inherit" /> : <SendIcon />} 
//           onClick={handleSend}
//           disabled={!sessionActive || loading || !input.trim()}
//         >
//           Send
//         </Button>
//       </Box>
      
//       {/* CSV Headers */}
//       {csvHeaders.length > 0 && (
//         <Box sx={{ p: 1, borderTop: "1px solid rgba(255,255,255,0.1)", display: "flex", flexWrap: "wrap", gap: 0.5 }}>
//           <Typography variant="caption" color="text.secondary" sx={{ mr: 1, alignSelf: "center" }}>
//             Available columns:
//           </Typography>
//           {csvHeaders.map((header, idx) => (
//             <Typography key={idx} variant="caption" sx={{ 
//               bgcolor: "rgba(88, 101, 242, 0.2)", 
//               px: 0.7, 
//               py: 0.3, 
//               borderRadius: 1,
//               fontSize: "0.7rem"
//             }}>
//               {header}
//             </Typography>
//           ))}
//         </Box>
//       )}

//       {/* LLM Configuration Dialog */}
//       <Dialog open={llmConfigDialog} onClose={() => setLLMConfigDialog(false)}>
//         <DialogTitle>Configure LLM Provider</DialogTitle>
//         <DialogContent>
//           {apiKeyError && (
//             <Alert severity="error" sx={{ mb: 2 }}>
//               {apiKeyError}
//             </Alert>
//           )}
//           <FormControl fullWidth sx={{ mt: 2 }}>
//             <InputLabel>LLM Provider</InputLabel>
//             <Select
//               value={llmProvider || ""}
//               onChange={(e) => handleProviderChange(e.target.value)}
//               label="LLM Provider"
//               disabled={sessionActive && !!llmProvider}
//             >
//               <MenuItem value="gemini">Google Gemini</MenuItem>
//               <MenuItem value="openai">OpenAI</MenuItem>
//               <MenuItem value="anthropic">Anthropic Claude</MenuItem>
//               <MenuItem value="mistral">Mistral</MenuItem>
//             </Select>
//           </FormControl>
//           {sessionActive && llmProvider && (
//             <Typography variant="caption" color="warning.main" sx={{ mt: 1, display: 'block' }}>
//               Note: Changing the provider requires starting a new chat session.
//             </Typography>
//           )}
//           <TextField
//             fullWidth
//             margin="normal"
//             label="API Key"
//             type="password"
//             value={apiKey || ""}
//             onChange={(e) => setApiKey(e.target.value)}
//             error={!!apiKeyError}
//             helperText={apiKeyError || "Your API key will be used for this chat session only"}
//           />
//         </DialogContent>
//         <DialogActions>
//           <Button onClick={() => setLLMConfigDialog(false)}>Cancel</Button>
//           <Button onClick={handleLLMConfig} variant="contained">
//             Configure
//           </Button>
//         </DialogActions>
//       </Dialog>
//     </Box>
//   );
// });

// export default ChatWindow;

import React, { useState, useRef, useEffect, useCallback } from "react";
import { uploadFile, queryData, fetchChatById, updateChat, createChat, validateApiKey } from '../../utils/api';

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
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Alert,
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
} from "@mui/icons-material";
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
        <Typography 
          variant="h6" 
          sx={{ cursor: 'pointer' }}
          onClick={() => {
            if (currentChatId) {
              const newTitle = prompt("Enter new chat title:", chatTitle);
              if (newTitle && newTitle.trim()) {
                updateChatTitle(newTitle.trim());
              }
            }
          }}
        >
          {sessionActive 
            ? `${chatTitle} - ${csvFile?.name || "CSV File"}`
            : "Sage AI - Upload a CSV file to begin"}
        </Typography>
        <Box sx={{ display: 'flex', alignItems: 'center' }}>
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
                  <Typography variant="body2">Use Knowledge Base</Typography>
                </Box>
              }
              sx={{ mr: 2 }}
            />
          )}
          
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
              onClick={() => fileInputRef.current?.click()}
              startIcon={<AttachFileIcon />}
              disabled={loading}
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

      {/* LLM Configuration Alert */}
      {!llmProvider && sessionActive && (
        <Paper
          sx={{
            p: 2,
            m: 2,
            backgroundColor: "rgba(88, 101, 242, 0.1)",
            borderRadius: 2,
            display: "flex",
            alignItems: "center",
            justifyContent: "space-between"
          }}
        >
          <Box sx={{ display: "flex", alignItems: "center" }}>
            <KeyIcon sx={{ mr: 1 }} />
            <Typography>Configure LLM provider and API key to start chatting</Typography>
          </Box>
          <Button
            variant="contained"
            size="small"
            onClick={() => setLLMConfigDialog(true)}
          >
            Configure
          </Button>
        </Paper>
      )}

      {/* LLM provider configured indicator */}
      {llmProvider && sessionActive && (
        <Chip
          icon={<KeyIcon />}
          label={`${llmProvider.toUpperCase()} (${llmModel || 'Default Model'})`}
          sx={{ m: 2 }}
          color="primary"
          variant="outlined"
          onDelete={() => setLLMConfigDialog(true)}
          deleteIcon={<Tooltip title="Change LLM Settings"><KeyIcon /></Tooltip>}
        />
      )}

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
              onClick={() => fileInputRef.current?.click()}
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
                        <Box component="span" sx={{ display: "flex", alignItems: "center", justifyContent: "center", mb: 1 }}>
                          <TableChartIcon fontSize="small" sx={{ mr: 1 }} />
                          <Typography component="span" variant="caption">{msg.file}</Typography>
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
                    
                    {/* Knowledge Context Used */}
                    {msg.knowledgeContext && msg.knowledgeContext.length > 0 && (
                      <Box sx={{ mt: 1, mb: 1 }}>
                        <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                          <LightbulbIcon fontSize="small" sx={{ mr: 0.5, color: 'primary.main' }} />
                          <Typography variant="caption" color="primary.main">
                            Knowledge Base Used:
                          </Typography>
                        </Box>
                        <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                          {msg.knowledgeContext.map((item, i) => (
                            <Tooltip 
                              key={i} 
                              title={
                                <Box>
                                  <Typography variant="caption" sx={{ fontWeight: 'bold' }}>
                                    {item.description}
                                  </Typography>
                                  <Typography variant="caption" sx={{ display: 'block', mt: 0.5 }}>
                                    {item.type === 'question_sql' 
                                      ? `Q: ${item.question} | SQL: ${item.sql}`
                                      : item.content?.substring(0, 200) + (item.content?.length > 200 ? '...' : '')
                                    }
                                  </Typography>
                                </Box>
                              }
                            >
                              <Chip 
                                size="small"
                                label={item.description}
                                icon={
                                  item.type === 'ddl' ? <CodeIcon fontSize="small" /> :
                                  item.type === 'documentation' ? <DescriptionIcon fontSize="small" /> :
                                  <DataObjectIcon fontSize="small" />
                                }
                                sx={{ bgcolor: 'rgba(88, 101, 242, 0.2)' }}
                              />
                            </Tooltip>
                          ))}
                        </Box>
                      </Box>
                    )}
                    
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
                <IconButton component="label" onClick={() => fileInputRef.current?.click()}>
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

      {/* LLM Configuration Dialog */}
      <Dialog open={llmConfigDialog} onClose={() => setLLMConfigDialog(false)}>
        <DialogTitle>Configure LLM Provider</DialogTitle>
        <DialogContent>
          {apiKeyError && (
            <Alert severity="error" sx={{ mb: 2 }}>
              {apiKeyError}
            </Alert>
          )}
          <FormControl fullWidth sx={{ mt: 2 }}>
            <InputLabel>LLM Provider</InputLabel>
            <Select
              value={llmProvider || ""}
              onChange={(e) => handleProviderChange(e.target.value)}
              label="LLM Provider"
              disabled={sessionActive && !!llmProvider}
            >
              <MenuItem value="gemini">Google Gemini</MenuItem>
              <MenuItem value="openai">OpenAI</MenuItem>
              <MenuItem value="anthropic">Anthropic Claude</MenuItem>
              <MenuItem value="mistral">Mistral</MenuItem>
            </Select>
          </FormControl>
          
          {llmProvider && (
            <FormControl fullWidth sx={{ mt: 2 }}>
              <InputLabel>Model</InputLabel>
              <Select
                value={llmModel || ""}
                onChange={(e) => setLLMModel(e.target.value)}
                label="Model"
              >
                {modelOptions[llmProvider]?.map((model) => (
                  <MenuItem key={model} value={model}>{model}</MenuItem>
                ))}
              </Select>
            </FormControl>
          )}
          
          {sessionActive && llmProvider && (
            <Typography variant="caption" color="warning.main" sx={{ mt: 1, display: 'block' }}>
              Note: Changing the provider requires starting a new chat session.
            </Typography>
          )}
          
          <TextField
            fullWidth
            margin="normal"
            label="API Key"
            type="password"
            value={apiKey || ""}
            onChange={(e) => setApiKey(e.target.value)}
            error={!!apiKeyError}
            helperText={apiKeyError || "Your API key will be used for this chat session only"}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setLLMConfigDialog(false)}>Cancel</Button>
          <Button onClick={handleLLMConfig} variant="contained">
            Configure
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
});

export default ChatWindow;
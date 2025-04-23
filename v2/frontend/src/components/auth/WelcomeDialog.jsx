// // import React, { useState } from "react";
// // import {
// //   Dialog,
// //   DialogTitle,
// //   DialogContent,
// //   DialogActions,
// //   TextField,
// //   Button,
// //   Typography,
// //   Box
// // } from "@mui/material";
// // import { useAuth } from "./AuthContext";

// // const WelcomeDialog = ({ open, onClose }) => {
// //   const [preferredName, setPreferredName] = useState("");
// //   const [error, setError] = useState("");
// //   const { user, updateUserName } = useAuth();

// //   const handleSubmit = async () => {
// //     if (!preferredName.trim()) {
// //       setError("Please enter a name");
// //       return;
// //     }

// //     try {
// //       await updateUserName(preferredName);
// //       onClose();
// //     } catch (err) {
// //       console.error("Failed to update user name:", err);
// //       setError("Failed to save your name. Please try again.");
// //     }
// //   };

// //   return (
// //     <Dialog 
// //       open={open} 
// //       onClose={onClose}
// //       maxWidth="sm"
// //       fullWidth
// //       PaperProps={{
// //         style: {
// //           backgroundColor: "#2F3136",
// //           borderRadius: "16px",
// //         },
// //       }}
// //     >
// //       <DialogTitle sx={{ textAlign: "center", pt: 4 }}>
// //         <Typography variant="h5" sx={{ fontWeight: 600, color: "white" }}>
// //           Welcome to Sage.AI!
// //         </Typography>
// //       </DialogTitle>
// //       <DialogContent>
// //         <Box sx={{ mb: 3, textAlign: "center" }}>
// //           <Typography variant="body1" sx={{ color: "rgba(255, 255, 255, 0.8)" }}>
// //             Thanks for joining us! What would you like to be called?
// //           </Typography>
// //         </Box>
// //         <TextField
// //           autoFocus
// //           fullWidth
// //           label="Your name"
// //           value={preferredName}
// //           onChange={(e) => setPreferredName(e.target.value)}
// //           error={!!error}
// //           helperText={error}
// //           sx={{
// //             mb: 2,
// //             "& .MuiOutlinedInput-root": {
// //               bgcolor: "rgba(0, 0, 0, 0.2)",
// //             },
// //           }}
// //         />
// //         <Typography variant="caption" sx={{ color: "rgba(255, 255, 255, 0.5)" }}>
// //           This name will be displayed on your profile and in your chats.
// //         </Typography>
// //       </DialogContent>
// //       <DialogActions sx={{ p: 3 }}>
// //         <Button 
// //           variant="outlined" 
// //           onClick={() => {
// //             setPreferredName(user?.email?.split('@')[0] || "User");
// //             onClose();
// //           }}
// //         >
// //           Skip
// //         </Button>
// //         <Button 
// //           variant="contained" 
// //           onClick={handleSubmit}
// //           sx={{ bgcolor: "#5865F2" }}
// //         >
// //           Continue
// //         </Button>
// //       </DialogActions>
// //     </Dialog>
// //   );
// // };

// // export default WelcomeDialog;

// import React, { useState } from "react";
// import {
//   Dialog,
//   DialogTitle,
//   DialogContent,
//   DialogActions,
//   TextField,
//   Button,
//   Typography,
//   Box,
//   CircularProgress
// } from "@mui/material";
// import { useAuth } from "./AuthContext";

// const WelcomeDialog = ({ open, onClose }) => {
//   const [preferredName, setPreferredName] = useState("");
//   const [error, setError] = useState("");
//   const [loading, setLoading] = useState(false);
//   const { user, updateUserName } = useAuth();

//   const handleSubmit = async () => {
//     if (!preferredName.trim()) {
//       setError("Please enter a name");
//       return;
//     }

//     try {
//       setLoading(true);
//       await updateUserName(preferredName);
//       setLoading(false);
//       onClose();
//     } catch (err) {
//       console.error("Failed to update user name:", err);
//       setError("Failed to save your name. Please try again.");
//       setLoading(false);
//     }
//   };

//   return (
//     <Dialog 
//       open={open} 
//       onClose={onClose}
//       maxWidth="sm"
//       fullWidth
//       PaperProps={{
//         style: {
//           backgroundColor: "#2F3136",
//           borderRadius: "16px",
//         },
//       }}
//     >
//       <DialogTitle sx={{ textAlign: "center", pt: 4 }}>
//         <Typography variant="h5" sx={{ fontWeight: 600, color: "white" }}>
//           Welcome to Sage.AI!
//         </Typography>
//       </DialogTitle>
//       <DialogContent>
//         <Box sx={{ mb: 3, textAlign: "center" }}>
//           <Typography variant="body1" sx={{ color: "rgba(255, 255, 255, 0.8)" }}>
//             Thanks for joining us! What would you like to be called?
//           </Typography>
//         </Box>
//         <TextField
//           autoFocus
//           fullWidth
//           label="Your name"
//           value={preferredName}
//           onChange={(e) => {
//             setPreferredName(e.target.value);
//             setError("");
//           }}
//           error={!!error}
//           helperText={error}
//           sx={{
//             mb: 2,
//             "& .MuiOutlinedInput-root": {
//               bgcolor: "rgba(0, 0, 0, 0.2)",
//             },
//           }}
//         />
//         <Typography variant="caption" sx={{ color: "rgba(255, 255, 255, 0.5)" }}>
//           This name will be displayed on your profile and in your chats.
//         </Typography>
//       </DialogContent>
//       <DialogActions sx={{ p: 3 }}>
//         <Button 
//           variant="outlined" 
//           onClick={() => {
//             // Set a default name using the first part of the email
//             const defaultName = user?.email?.split('@')[0] || "User";
//             // Auto-update with this default name
//             updateUserName(defaultName).then(() => {
//               onClose();
//             });
//           }}
//           disabled={loading}
//         >
//           Skip
//         </Button>
//         <Button 
//           variant="contained" 
//           onClick={handleSubmit}
//           sx={{ bgcolor: "#5865F2" }}
//           disabled={loading}
//         >
//           {loading ? <CircularProgress size={24} /> : "Continue"}
//         </Button>
//       </DialogActions>
//     </Dialog>
//   );
// };

// export default WelcomeDialog;

import React, { useState } from "react";
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Button,
  Typography,
  Box,
  CircularProgress,
  Alert
} from "@mui/material";
import { useAuth } from "./AuthContext";

const WelcomeDialog = ({ open, onClose }) => {
  const [preferredName, setPreferredName] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const { user } = useAuth();

  // Generate a suggested name from email if available
  const suggestedName = user?.email ? user.email.split('@')[0] : "";

  const handleSubmit = async () => {
    if (!preferredName.trim()) {
      setError("Please enter a name");
      return;
    }

    try {
      setLoading(true);
      setError("");

      // Make a direct API call to update the profile
      // This avoids going through the AuthContext which might have errors
      const response = await fetch("/api/profile", {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({ 
          name: preferredName,
          // Make sure we don't send undefined values
          profilePicUrl: user?.profilePicUrl || "" 
        }),
      });

      if (!response.ok) {
        console.error("Failed to update profile:", response.status);
        throw new Error("Failed to update your profile. Please try again.");
      }

      // Wait for profile update to complete
      await response.json();
      
      // Reset states
      setLoading(false);
      
      // Refresh the page to ensure updated data is displayed
      window.location.reload();
      
      // Close the dialog
      onClose();
    } catch (err) {
      console.error("Failed to update user name:", err);
      setError(err.message || "Failed to save your name. Please try again.");
      setLoading(false);
    }
  };

  const handleUseDefault = async () => {
    // Use email username as default name
    const defaultName = suggestedName || "User";
    setPreferredName(defaultName);
    
    try {
      setLoading(true);
      setError("");

      // Make a direct API call to update the profile
      const response = await fetch("/api/profile", {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({ 
          name: defaultName,
          profilePicUrl: user?.profilePicUrl || ""
        }),
      });

      if (!response.ok) {
        console.error("Failed to update profile with default name:", response.status);
        throw new Error("Failed to update your profile. Please try again.");
      }

      // Wait for profile update to complete
      await response.json();
      
      // Reset states
      setLoading(false);
      
      // Refresh the page to ensure updated data is displayed
      window.location.reload();
      
      // Close dialog
      onClose();
    } catch (err) {
      console.error("Error setting default name:", err);
      setError(err.message || "Failed to set default name. Please try entering a name manually.");
      setLoading(false);
    }
  };

  return (
    <Dialog 
      open={open} 
      onClose={(event, reason) => {
        // Only allow closing if not a generic user
        if (user && !['Yash Kishore', 'GitHub User', 'Google User', 'OAuth User'].includes(user.name)) {
          onClose();
        }
      }}
      maxWidth="sm"
      fullWidth
      PaperProps={{
        style: {
          backgroundColor: "#2F3136",
          borderRadius: "16px",
        },
      }}
    >
      <DialogTitle sx={{ textAlign: "center", pt: 4 }}>
        <Typography variant="h5" sx={{ fontWeight: 600, color: "white" }}>
          Welcome to Sage.AI!
        </Typography>
      </DialogTitle>
      <DialogContent>
        <Box sx={{ mb: 3, textAlign: "center" }}>
          <Typography variant="body1" sx={{ color: "rgba(255, 255, 255, 0.8)" }}>
            What would you like to be called?
          </Typography>
        </Box>
        
        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}
        
        <TextField
          autoFocus
          fullWidth
          label="Your name"
          placeholder={suggestedName || "Your preferred name"}
          value={preferredName}
          onChange={(e) => {
            setPreferredName(e.target.value);
            setError("");
          }}
          error={!!error}
          sx={{
            mb: 2,
            "& .MuiOutlinedInput-root": {
              bgcolor: "rgba(0, 0, 0, 0.2)",
            },
          }}
        />
        <Typography variant="caption" sx={{ color: "rgba(255, 255, 255, 0.5)" }}>
          This name will be displayed on your profile and in your chats.
        </Typography>
      </DialogContent>
      <DialogActions sx={{ p: 3 }}>
        <Button 
          variant="outlined" 
          onClick={handleUseDefault}
          disabled={loading}
        >
          {loading ? <CircularProgress size={20} /> : "Use Default"}
        </Button>
        <Button 
          variant="contained" 
          onClick={handleSubmit}
          sx={{ bgcolor: "#5865F2" }}
          disabled={loading || !preferredName.trim()}
        >
          {loading ? <CircularProgress size={24} color="inherit" /> : "Continue"}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default WelcomeDialog;
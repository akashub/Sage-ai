// // pages/OAuthCallback.jsx
// import { useEffect, useState } from "react";
// import { useLocation, useNavigate } from "react-router-dom";
// import { Box, CircularProgress, Typography, Paper } from "@mui/material";
// import { useAuth } from "./AuthContext";

// const OAuthCallback = () => {
//   const [status, setStatus] = useState("Processing your authentication...");
//   const [error, setError] = useState(null);
//   const location = useLocation();
//   const navigate = useNavigate();
//   const { oauthSignIn } = useAuth();

//   useEffect(() => {
//     const handleOAuthCallback = async () => {
//       try {
//         // Parse URL params
//         const searchParams = new URLSearchParams(location.search);
//         const code = searchParams.get("code");
//         const state = searchParams.get("state");
//         const error = searchParams.get("error");

//         if (error) {
//           throw new Error(`Authentication error: ${error}`);
//         }

//         if (!code) {
//           throw new Error("No authorization code received");
//         }

//         // Determine provider from state or referrer
//         let provider = "unknown";
//         if (document.referrer.includes("google")) {
//           provider = "google";
//         } else if (document.referrer.includes("github")) {
//           provider = "github";
//         }

//         setStatus(`Completing authentication with ${provider}...`);

//         // Complete OAuth flow
//         await oauthSignIn(provider, code, window.location.origin + "/oauth-callback");
        
//         // Redirect will happen automatically via AuthContext
//         navigate("/chat");
//       } catch (err) {
//         console.error("OAuth callback error:", err);
//         setError(err.message);
//         setStatus("Authentication failed");
//       }
//     };

//     handleOAuthCallback();
//   }, [location, oauthSignIn]);

//   if (error) {
//     // Redirect to home after a short delay on error
//     setTimeout(() => {
//       return <Navigate to="/" />;
//     }, 3000);
//   }

//   return (
//     <Box
//       sx={{
//         display: "flex",
//         flexDirection: "column",
//         alignItems: "center",
//         justifyContent: "center",
//         height: "100vh",
//         bgcolor: "background.default",
//       }}
//     >
//       <Paper
//         elevation={3}
//         sx={{
//           p: 4,
//           display: "flex",
//           flexDirection: "column",
//           alignItems: "center",
//           maxWidth: 500,
//         }}
//       >
//         {error ? (
//           <>
//             <Typography variant="h5" color="error" gutterBottom>
//               Authentication Failed
//             </Typography>
//             <Typography color="error">{error}</Typography>
//             <Typography variant="body2" sx={{ mt: 2 }}>
//               Redirecting you back to home...
//             </Typography>
//           </>
//         ) : (
//           <>
//             <CircularProgress size={60} thickness={4} sx={{ mb: 3 }} />
//             <Typography variant="h6">{status}</Typography>
//           </>
//         )}
//       </Paper>
//     </Box>
//   );
// };

// export default OAuthCallback;

// Import Navigate properly at the top
import { useEffect, useState } from "react";
import { useLocation, useNavigate} from "react-router-dom";
import { Box, CircularProgress, Typography, Paper, Button } from "@mui/material";
import { useAuth } from "./AuthContext"; // Note this path change to use the correct import

const OAuthCallback = () => {
  const [status, setStatus] = useState("Processing your authentication...");
  const [error, setError] = useState(null);
  const location = useLocation();
  const navigate = useNavigate();
  const { oauthSignIn } = useAuth();

  useEffect(() => {
    const handleOAuthCallback = async () => {
      try {
        // Parse URL params
        const searchParams = new URLSearchParams(location.search);
        const code = searchParams.get("code");
        const error = searchParams.get("error");

        if (error) {
          throw new Error(`Authentication error: ${error}`);
        }

        if (!code) {
          throw new Error("No authorization code received");
        }

        // Better provider detection
        const provider = detectProvider();
        console.log("Detected provider:", provider);
        
        // Don't proceed if we couldn't determine the provider
        if (provider === "unknown") {
          throw new Error("Unable to determine OAuth provider");
        }

        setStatus(`Completing authentication with ${provider}...`);

        // Complete OAuth flow
        const response = await oauthSignIn(provider, code, window.location.origin + "/oauth-callback");
        
        // Navigate to chat after successful authentication
        navigate("/chat");
      } 
      // catch (err) {
      //   console.error("OAuth callback error:", err);
      //   setError(err.message);
      //   setStatus("Authentication failed");
      catch (err) {
        console.error("OAuth callback error:", err);
        
        // Extract more meaningful error messages
        let errorMessage = err.message;
        if (err.response) {
          try {
            // Try to parse as JSON first
            const errorData = await err.response.json();
            errorMessage = errorData.message || errorMessage;
          } catch {
            // If not JSON, get as text
            const errorText = await err.response.text();
            errorMessage = errorText || errorMessage;
          }
        }
        
        setError(errorMessage);
        setStatus("Authentication failed");
      }
    };

    // Helper function to detect the provider
    const detectProvider = () => {
      // First checking sessionStorage (most reliable)
      const storedProvider = sessionStorage.getItem('oauth_provider');
      if (storedProvider) {
        return storedProvider;
      }
      // Try to get provider from URL state or params first
      const searchParams = new URLSearchParams(location.search);
      if (searchParams.get("provider")) {
        return searchParams.get("provider");
      }
      // Try URL path
      const pathSegments = location.pathname.split('/');
      if (pathSegments.includes("github")) return "github";
      if (pathSegments.includes("google")) return "google";
      
      // Try referrer
      if (document.referrer) {
        if (document.referrer.includes("github.com")) return "github";
        if (document.referrer.includes("google.com")) return "google";
      }
      
      // If we can't determine, return unknown
      return "unknown";
    };

    handleOAuthCallback();
  }, [location, oauthSignIn, navigate]);

  // If error occurred, redirect to home
  if (error) {
    return (
      <Box
        sx={{
          display: "flex",
          flexDirection: "column",
          alignItems: "center",
          justifyContent: "center",
          height: "100vh",
          color: "white",
        }}
      >
        <Typography variant="h5" color="error" gutterBottom>
          Authentication Failed
        </Typography>
        <Typography color="error">{error}</Typography>
        <Typography variant="body2" sx={{ mt: 2 }}>
          Redirecting you back to home...
        </Typography>
        {/* Use navigate in a useEffect, not Navigate component directly */}
        <Box sx={{ mt: 4 }}>
          <Button 
            variant="contained"
            onClick={() => navigate("/")}
            sx={{ mt: 2 }}
          >
            Return to Home
          </Button>
        </Box>
      </Box>
    );
  }

  return (
    <Box
      sx={{
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        justifyContent: "center",
        height: "100vh",
        bgcolor: "background.default",
        color: "white",
      }}
    >
      <Paper
        elevation={3}
        sx={{
          p: 4,
          display: "flex",
          flexDirection: "column",
          alignItems: "center",
          maxWidth: 500,
          bgcolor: "background.paper",
        }}
      >
        <CircularProgress size={60} thickness={4} sx={{ mb: 3 }} />
        <Typography variant="h6">{status}</Typography>
      </Paper>
    </Box>
  );
};

export default OAuthCallback;
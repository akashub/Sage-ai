// pages/OAuthCallback.jsx
import { useEffect, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { Box, CircularProgress, Typography, Paper } from "@mui/material";
import { useAuth } from "./AuthContext";

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
        const state = searchParams.get("state");
        const error = searchParams.get("error");

        if (error) {
          throw new Error(`Authentication error: ${error}`);
        }

        if (!code) {
          throw new Error("No authorization code received");
        }

        // Determine provider from state or referrer
        let provider = "unknown";
        if (document.referrer.includes("google")) {
          provider = "google";
        } else if (document.referrer.includes("github")) {
          provider = "github";
        }

        setStatus(`Completing authentication with ${provider}...`);

        // Complete OAuth flow
        await oauthSignIn(provider, code, window.location.origin + "/oauth-callback");
        
        // Redirect will happen automatically via AuthContext
        navigate("/chat");
      } catch (err) {
        console.error("OAuth callback error:", err);
        setError(err.message);
        setStatus("Authentication failed");
      }
    };

    handleOAuthCallback();
  }, [location, oauthSignIn]);

  if (error) {
    // Redirect to home after a short delay on error
    setTimeout(() => {
      return <Navigate to="/" />;
    }, 3000);
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
        }}
      >
        {error ? (
          <>
            <Typography variant="h5" color="error" gutterBottom>
              Authentication Failed
            </Typography>
            <Typography color="error">{error}</Typography>
            <Typography variant="body2" sx={{ mt: 2 }}>
              Redirecting you back to home...
            </Typography>
          </>
        ) : (
          <>
            <CircularProgress size={60} thickness={4} sx={{ mb: 3 }} />
            <Typography variant="h6">{status}</Typography>
          </>
        )}
      </Paper>
    </Box>
  );
};

export default OAuthCallback;
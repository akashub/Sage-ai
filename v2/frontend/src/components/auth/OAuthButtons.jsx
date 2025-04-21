// components/auth/OAuthButtons.jsx
"use client"
import { Button, Box } from "@mui/material"
import GoogleIcon from "@mui/icons-material/Google"
import GitHubIcon from "@mui/icons-material/GitHub"
import { motion } from "framer-motion"
import { useAuth } from "./AuthContext"
import { useState } from "react"

const OAuthButton = ({ icon, children, onClick, delay, disabled }) => (
  <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ duration: 0.5, delay }}>
    <Button
      fullWidth
      variant="outlined"
      startIcon={icon}
      onClick={onClick}
      disabled={disabled}
      sx={{
        mb: 2,
        py: 1.5,
        color: "white",
        borderColor: "rgba(255, 255, 255, 0.1)",
        backgroundColor: "rgba(255, 255, 255, 0.05)",
        textTransform: "none",
        fontSize: "0.9rem",
        fontWeight: 500,
        "&:hover": {
          backgroundColor: "rgba(255, 255, 255, 0.1)",
          borderColor: "rgba(255, 255, 255, 0.2)",
        },
      }}
    >
      {children}
    </Button>
  </motion.div>
)

const OAuthButtons = () => {
  const { getOAuthUrl } = useAuth();
  const [isLoading, setIsLoading] = useState(false);

  const handleOAuthSignIn = async (provider) => {
    try {
      setIsLoading(true);
      const redirectUri = `${window.location.origin}/oauth-callback`;
      console.log(`Requesting OAuth URL for ${provider} with redirect URI: ${redirectUri}`);
      
      // Store the provider in sessionStorage before redirecting
      sessionStorage.setItem('oauth_provider', provider);
      
      // Request the OAuth URL with proper error handling
      const response = await fetch(`/api/auth/oauth/url/${provider}?redirect_uri=${encodeURIComponent(redirectUri)}`);
      
      if (!response.ok) {
        console.error(`Error response from server: ${response.status} ${response.statusText}`);
        throw new Error(`Failed to get OAuth URL: ${response.status}`);
      }
      
      const data = await response.json();
      console.log(`Redirecting to: ${data.url}`);
      
      // Redirect to the OAuth provider
      window.location.href = data.url;
    } catch (err) {
      console.error(`OAuth sign in with ${provider} failed:`, err);
      alert(`Failed to authenticate with ${provider}. Please try again later.`);
      setIsLoading(false);
    }
  };
  
  return (
    <Box>
      <OAuthButton 
        icon={<GitHubIcon />} 
        onClick={() => handleOAuthSignIn('github')} 
        delay={0.1}
        disabled={isLoading}
      >
        Continue with GitHub
      </OAuthButton>

      <OAuthButton 
        icon={<GoogleIcon />} 
        onClick={() => handleOAuthSignIn('google')} 
        delay={0.2}
        disabled={isLoading}
      >
        Continue with Google
      </OAuthButton>
    </Box>
  )
}

export default OAuthButtons

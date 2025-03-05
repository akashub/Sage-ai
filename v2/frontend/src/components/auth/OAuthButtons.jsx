// "use client"
// import { Button, Box } from "@mui/material"
// import GoogleIcon from "@mui/icons-material/Google"
// import GitHubIcon from "@mui/icons-material/GitHub"
// import { motion } from "framer-motion"

// const OAuthButton = ({ icon, children, onClick, delay }) => (
//   <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ duration: 0.5, delay }}>
//     <Button
//       fullWidth
//       variant="outlined"
//       startIcon={icon}
//       onClick={onClick}
//       sx={{
//         mb: 2,
//         py: 1.5,
//         color: "white",
//         borderColor: "rgba(255, 255, 255, 0.1)",
//         backgroundColor: "rgba(255, 255, 255, 0.05)",
//         textTransform: "none",
//         fontSize: "0.9rem",
//         fontWeight: 500,
//         "&:hover": {
//           backgroundColor: "rgba(255, 255, 255, 0.1)",
//           borderColor: "rgba(255, 255, 255, 0.2)",
//         },
//       }}
//     >
//       {children}
//     </Button>
//   </motion.div>
// )

// const OAuthButtons = () => {
//   const handleGoogleSignIn = () => {
//     console.log("Google Sign In")
//   }

//   const handleGithubSignIn = () => {
//     console.log("GitHub Sign In")
//   }

//   return (
//     <Box>
//       <OAuthButton icon={<GitHubIcon />} onClick={handleGithubSignIn} delay={0.1}>
//         Continue with GitHub
//       </OAuthButton>

//       <OAuthButton icon={<GoogleIcon />} onClick={handleGoogleSignIn} delay={0.2}>
//         Continue with Google
//       </OAuthButton>
//     </Box>
//   )
// }

// export default OAuthButtons

// components/auth/OAuthButtons.jsx
"use client"
import { Button, Box } from "@mui/material"
import GoogleIcon from "@mui/icons-material/Google"
import GitHubIcon from "@mui/icons-material/GitHub"
import { motion } from "framer-motion"
import { useAuth } from "./AuthContext"

const OAuthButton = ({ icon, children, onClick, delay }) => (
  <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ duration: 0.5, delay }}>
    <Button
      fullWidth
      variant="outlined"
      startIcon={icon}
      onClick={onClick}
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

  // const handleOAuthSignIn = async (provider) => {
  //   try {
  //     const redirectUri = `${window.location.origin}/oauth-callback`;
  //     const authUrl = await getOAuthUrl(provider, redirectUri);
  //     window.location.href = authUrl;
  //   } catch (err) {
  //     console.error(`OAuth sign in with ${provider} failed:`, err);
  //   }
  // };
  const handleOAuthSignIn = async (provider) => {
    try {
      const redirectUri = `${window.location.origin}/oauth-callback`;
      console.log(`Requesting OAuth URL for ${provider} with redirect URI: ${redirectUri}`);
      
      // Pass the provider in the query parameter so we can retrieve it later
      const response = await fetch(`/api/auth/oauth/url/${provider}?redirect_uri=${encodeURIComponent(redirectUri)}&provider=${provider}`);
      
      if (!response.ok) {
        const errorText = await response.text();
        console.error(`Error response from server:`, errorText);
        throw new Error(`Failed to get OAuth URL: ${response.status} ${response.statusText}`);
      }
      
      const data = await response.json();
      console.log(`Redirecting to: ${data.url}`);
      
      // Store the provider in sessionStorage before redirecting
      sessionStorage.setItem('oauth_provider', provider);
      
      window.location.href = data.url;
    } catch (err) {
      console.error(`OAuth sign in with ${provider} failed:`, err);
    }
  };
  
  return (
    <Box>
      <OAuthButton 
        icon={<GitHubIcon />} 
        onClick={() => handleOAuthSignIn('github')} 
        delay={0.1}
      >
        Continue with GitHub
      </OAuthButton>

      <OAuthButton 
        icon={<GoogleIcon />} 
        onClick={() => handleOAuthSignIn('google')} 
        delay={0.2}
      >
        Continue with Google
      </OAuthButton>
    </Box>
  )
}

export default OAuthButtons
"use client"

import { useState } from "react"
import { Box, TextField, Button, Link, InputAdornment, IconButton, Typography, CircularProgress, Alert } from "@mui/material"
import { Visibility, VisibilityOff, Google as GoogleIcon, GitHub as GitHubIcon } from "@mui/icons-material"
import { useAuth } from "./AuthContext"

const SignInForm = ({ onSignUpClick, onForgotClick }) => {
  const [showPassword, setShowPassword] = useState(false)
  const [formData, setFormData] = useState({
    email: "",
    password: "",
  })
  const [isSubmitting, setIsSubmitting] = useState(false)
  const { signIn, oauthSignIn, getOAuthUrl, error } = useAuth()

  const handleSubmit = async (e) => {
    e.preventDefault()
    setIsSubmitting(true)
    
    try {
      await signIn(formData.email, formData.password)
    } catch (err) {
      console.error("Sign in failed:", err)
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleChange = (e) => {
    const { name, value } = e.target
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }))
  }

  const handleOAuthSignIn = async (provider) => {
    try {
      const redirectUri = `${window.location.origin}/oauth-callback`
      const authUrl = await getOAuthUrl(provider, redirectUri)
      window.location.href = authUrl
    } catch (err) {
      console.error(`OAuth sign in with ${provider} failed:`, err)
    }
  }

  return (
    <Box component="form" onSubmit={handleSubmit} sx={{ width: "100%" }}>
      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}
      
      <TextField
        fullWidth
        label="Email"
        name="email"
        type="email"
        value={formData.email}
        onChange={handleChange}
        margin="normal"
        required
        disabled={isSubmitting}
        sx={{
          "& .MuiOutlinedInput-root": {
            bgcolor: "background.paper",
          },
        }}
      />
      
      <TextField
        fullWidth
        label="Password"
        name="password"
        type={showPassword ? "text" : "password"}
        value={formData.password}
        onChange={handleChange}
        margin="normal"
        required
        disabled={isSubmitting}
        InputProps={{
          endAdornment: (
            <InputAdornment position="end">
              <IconButton onClick={() => setShowPassword(!showPassword)} edge="end" disabled={isSubmitting}>
                {showPassword ? <VisibilityOff /> : <Visibility />}
              </IconButton>
            </InputAdornment>
          ),
        }}
        sx={{
          "& .MuiOutlinedInput-root": {
            bgcolor: "background.paper",
          },
        }}
      />
      
      <Box sx={{ textAlign: "right", mb: 2 }}>
        <Link 
          component="button" 
          type="button" 
          onClick={onForgotClick} 
          disabled={isSubmitting}
          sx={{ color: "rgba(255, 255, 255, 0.7)" }}
        >
          Forgot password?
        </Link>
      </Box>
      
      <Button
        fullWidth
        type="submit"
        variant="contained"
        disabled={isSubmitting}
        sx={{
          mt: 2,
          mb: 2,
          bgcolor: "black",
          color: "white",
          "&:hover": {
            bgcolor: "rgba(0, 0, 0, 0.8)",
          },
        }}
      >
        {isSubmitting ? <CircularProgress size={24} color="inherit" /> : "Sign In"}
      </Button>
      
      <Box sx={{ textAlign: "center" }}>
        <Typography variant="body2" sx={{ color: "rgba(255, 255, 255, 0.7)" }}>
          Don't have an account?{" "}
          <Link 
            component="button" 
            type="button" 
            onClick={onSignUpClick} 
            disabled={isSubmitting}
            sx={{ color: "rgba(255, 255, 255, 1.0)" }}
          >
            Sign up
          </Link>
        </Typography>
      </Box>
    </Box>
  )
}

export default SignInForm
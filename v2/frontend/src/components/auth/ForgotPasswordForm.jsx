"use client"

import { useState } from "react"
import { Box, TextField, Button, Link, Alert, Typography } from "@mui/material"
import { useAuth } from "./AuthContext"

const ForgotPasswordForm = ({ onBackToSignIn }) => {
  const { resetPassword } = useAuth()
  const [email, setEmail] = useState("")
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState("")
  const [successMessage, setSuccessMessage] = useState("")

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError("")
    setSuccessMessage("")
    setLoading(true)
    
    try {
      const result = await resetPassword(email)
      if (result.success) {
        setSuccessMessage(result.message || "Password reset email sent successfully!")
      } else {
        setError(result.error || "Failed to send password reset email")
      }
    } catch (err) {
      setError("An unexpected error occurred")
      console.error(err)
    } finally {
      setLoading(false)
    }
  }

  return (
    <Box component="form" onSubmit={handleSubmit} sx={{ width: "100%" }}>
      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}
      
      {successMessage && (
        <Alert severity="success" sx={{ mb: 2 }}>
          {successMessage}
        </Alert>
      )}
      
      <Typography variant="body2" color="text.secondary" paragraph>
        Enter your email address and we'll send you a link to reset your password.
      </Typography>
      
      <TextField
        fullWidth
        label="Email"
        type="email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        margin="normal"
        required
        sx={{
          "& .MuiOutlinedInput-root": {
            bgcolor: "rgba(47, 49, 54, 0.6)",
          },
        }}
      />
      
      <Button
        fullWidth
        type="submit"
        variant="contained"
        color="primary"
        disabled={loading}
        sx={{
          mt: 3,
          mb: 2,
          height: '50px',
          borderRadius: '28px',
          textTransform: 'none',
          fontSize: '1rem',
          fontWeight: 500,
        }}
      >
        {loading ? "Sending..." : "Send Reset Link"}
      </Button>
      
      <Box sx={{ textAlign: "center" }}>
        <Link 
          component="button" 
          type="button" 
          onClick={onBackToSignIn} 
          sx={{ color: "secondary.main" }}
        >
          Back to Sign In
        </Link>
      </Box>
    </Box>
  )
}

export default ForgotPasswordForm
"use client"

import { useState } from "react"
import { Box, TextField, Button, Link, InputAdornment, IconButton, Alert } from "@mui/material"
import { Visibility, VisibilityOff } from "@mui/icons-material"
import { useAuth } from "./AuthContext"

const SignInForm = ({ onSignUpClick, onForgotClick }) => {
  const { signIn } = useAuth()
  const [showPassword, setShowPassword] = useState(false)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState("")
  const [formData, setFormData] = useState({
    email: "",
    password: "",
  })

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError("")
    setLoading(true)
    
    try {
      const result = await signIn(formData.email, formData.password)
      if (!result.success) {
        setError(result.error || "Failed to sign in")
      }
    } catch (err) {
      setError("An unexpected error occurred")
      console.error(err)
    } finally {
      setLoading(false)
    }
  }

  const handleChange = (e) => {
    const { name, value } = e.target
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }))
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
        sx={{
          "& .MuiOutlinedInput-root": {
            bgcolor: "rgba(47, 49, 54, 0.6)",
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
        InputProps={{
          endAdornment: (
            <InputAdornment position="end">
              <IconButton onClick={() => setShowPassword(!showPassword)} edge="end">
                {showPassword ? <VisibilityOff /> : <Visibility />}
              </IconButton>
            </InputAdornment>
          ),
        }}
        sx={{
          "& .MuiOutlinedInput-root": {
            bgcolor: "rgba(47, 49, 54, 0.6)",
          },
        }}
      />
      
      <Box sx={{ textAlign: "right", mt: 1 }}>
        <Link 
          component="button"
          type="button" 
          onClick={onForgotClick} 
          sx={{ color: "secondary.main" }}
        >
          Forgot password?
        </Link>
      </Box>
      
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
        {loading ? "Signing in..." : "Sign In"}
      </Button>
      
      <Box sx={{ textAlign: "center" }}>
        Don't have an account?{" "}
        <Link 
          component="button" 
          type="button" 
          onClick={onSignUpClick} 
          sx={{ color: "secondary.main" }}
        >
          Sign up
        </Link>
      </Box>
    </Box>
  )
}

export default SignInForm
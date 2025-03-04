"use client"

import { useState } from "react"
import { Box, TextField, Button, Link, InputAdornment, IconButton, Alert } from "@mui/material"
import { Visibility, VisibilityOff } from "@mui/icons-material"
import { useAuth } from "./AuthContext"

const SignUpForm = ({ onSignInClick }) => {
  const { signUp } = useAuth()
  const [showPassword, setShowPassword] = useState(false)
  const [showConfirmPassword, setShowConfirmPassword] = useState(false)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState("")
  const [formData, setFormData] = useState({
    email: "",
    password: "",
    confirmPassword: ""
  })

  // Validation
  const [validationErrors, setValidationErrors] = useState({
    email: "",
    password: "",
    confirmPassword: ""
  })

  const validateForm = () => {
    const errors = {
      email: "",
      password: "",
      confirmPassword: ""
    }

    // Email validation
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    if (!emailRegex.test(formData.email)) {
      errors.email = "Please enter a valid email address"
    }

    // Password validation
    if (formData.password.length < 8) {
      errors.password = "Password must be at least 8 characters"
    }

    // Confirm password validation
    if (formData.password !== formData.confirmPassword) {
      errors.confirmPassword = "Passwords do not match"
    }

    setValidationErrors(errors)
    return !Object.values(errors).some(error => error)
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError("")
    
    // Validate form
    if (!validateForm()) {
      return
    }
    
    setLoading(true)
    
    try {
      const result = await signUp(formData.email, formData.password)
      if (!result.success) {
        setError(result.error || "Failed to create account")
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
    
    // Clear validation error when user types
    setValidationErrors(prev => ({
      ...prev,
      [name]: ""
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
        error={!!validationErrors.email}
        helperText={validationErrors.email}
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
        error={!!validationErrors.password}
        helperText={validationErrors.password}
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
      
      <TextField
        fullWidth
        label="Confirm Password"
        name="confirmPassword"
        type={showConfirmPassword ? "text" : "password"}
        value={formData.confirmPassword}
        onChange={handleChange}
        margin="normal"
        required
        error={!!validationErrors.confirmPassword}
        helperText={validationErrors.confirmPassword}
        InputProps={{
          endAdornment: (
            <InputAdornment position="end">
              <IconButton onClick={() => setShowConfirmPassword(!showConfirmPassword)} edge="end">
                {showConfirmPassword ? <VisibilityOff /> : <Visibility />}
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
        {loading ? "Creating Account..." : "Create Account"}
      </Button>
      
      <Box sx={{ textAlign: "center" }}>
        Already have an account?{" "}
        <Link 
          component="button" 
          type="button" 
          onClick={onSignInClick} 
          sx={{ color: "secondary.main" }}
        >
          Sign in
        </Link>
      </Box>
    </Box>
  )
}

export default SignUpForm
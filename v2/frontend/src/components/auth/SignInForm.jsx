"use client"

import { useState } from "react"
import { Box, TextField, Button, Link, InputAdornment, IconButton } from "@mui/material"
import { Visibility, VisibilityOff } from "@mui/icons-material"
import { color } from "framer-motion"

const SignInForm = ({ onSignUpClick }) => {
  const [showPassword, setShowPassword] = useState(false)
  const [formData, setFormData] = useState({
    email: "",
    password: "",
  })

  const handleSubmit = (e) => {
    e.preventDefault()
    // Handle sign in logic here
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
            bgcolor: "background.paper",
          },
        }}
      />
      <Button
        fullWidth
        type="submit"
        variant="contained"
        sx={{
          mt: 3,
          mb: 2,
          bgcolor: "black",
          color: "white",
          "&:hover": {
            bgcolor: "rgba(0, 0, 0, 0.8)",
          },
        }}
      >
        Sign In
      </Button>
      <Box sx={{ textAlign: "center" }}>
        Don't have an account?{" "}
        <Link component="button" type="button" onClick={onSignUpClick} sx={{ color: "rgba(255, 255, 255, 1.0)" }}>
          Sign up
        </Link>
      </Box>
    </Box>
  )
}

export default SignInForm


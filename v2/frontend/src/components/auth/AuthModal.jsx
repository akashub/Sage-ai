"use client"

import { useState, useEffect } from "react"
import { Dialog, DialogContent, IconButton, Box, Typography, Fade } from "@mui/material"
import { Close as CloseIcon } from "@mui/icons-material"
import SignUpForm from "./SignUpForm"
import SignInForm from "./SignInForm"
import { motion, AnimatePresence } from "framer-motion"

const AuthModal = ({ open, onClose }) => {
  const [isSignUp, setIsSignUp] = useState(true)
  const [mounted, setMounted] = useState(false)

  useEffect(() => {
    setMounted(true)
  }, [])

  if (!mounted) return null

  return (
    <Dialog
      open={open}
      onClose={onClose}
      maxWidth="sm"
      fullWidth
      PaperProps={{
        sx: {
          bgcolor: "rgba(255, 255, 255, 0.1)",
          backdropFilter: "blur(20px)",
          borderRadius: "16px",
          boxShadow: "0 8px 32px rgba(0, 0, 0, 0.2)",
          overflow: "hidden",
        },
      }}
      sx={{
        "& .MuiBackdrop-root": {
          backgroundColor: "rgba(0, 0, 0, 0.1)",
          backdropFilter: "blur(10px)",
        },
      }}
    >
      <Fade in={open}>
        <DialogContent sx={{ p: 0, position: "relative" }}>
          <IconButton
            onClick={onClose}
            sx={{
              position: "absolute",
              right: 8,
              top: 8,
              color: "text.secondary",
              zIndex: 1,
            }}
          >
            <CloseIcon />
          </IconButton>
          <Box sx={{ p: 4 }}>
            <AnimatePresence mode="wait">
              <motion.div
                key={isSignUp ? "signup" : "signin"}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -20 }}
                transition={{ duration: 0.3 }}
              >
                <Box sx={{ textAlign: "center", mb: 4 }}>
                  <Typography variant="h4" component="h2" sx={{ mb: 1, color: "text.primary", fontWeight: "bold" }}>
                    {isSignUp ? "Create Account" : "Welcome Back"}
                  </Typography>
                  <Typography variant="body1" color="text.secondary">
                    {isSignUp ? "Start creating amazing SQL queries with AI" : "Sign in to continue your SQL journey"}
                  </Typography>
                </Box>

                {isSignUp ? (
                  <SignUpForm onSignInClick={() => setIsSignUp(false)} />
                ) : (
                  <SignInForm onSignUpClick={() => setIsSignUp(true)} />
                )}
              </motion.div>
            </AnimatePresence>
          </Box>
        </DialogContent>
      </Fade>
    </Dialog>
  )
}

export default AuthModal


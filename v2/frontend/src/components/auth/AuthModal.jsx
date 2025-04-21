"use client"

import { useState, useEffect } from "react"
import { Dialog, Box, IconButton, useTheme, Typography, Button } from "@mui/material"
import CloseIcon from "@mui/icons-material/Close"
import { useAuth } from "./AuthContext"
import SignInForm from "./SignInForm"
import SignUpForm from "./SignUpForm"
import ForgotPasswordForm from "./ForgotPasswordForm"
import OAuthButtons from "./OAuthButtons"
import { motion, AnimatePresence } from "framer-motion"

const ROTATING_TEXTS = ["Assistant", "Companion", "Generator", "Expert", "Partner"]

const AuthModal = ({ open, onClose, initialMode = "signin" }) => {
  const theme = useTheme()
  const { switchAuthMode, authMode } = useAuth()
  const [currentMode, setCurrentMode] = useState(initialMode)
  const [textIndex, setTextIndex] = useState(0)

  useEffect(() => {
    if (open) {
      setCurrentMode(initialMode)
    }
  }, [open, initialMode])

  useEffect(() => {
    const interval = setInterval(() => {
      setTextIndex((current) => (current + 1) % ROTATING_TEXTS.length)
    }, 2000)
    return () => clearInterval(interval)
  }, [])

  const handleModeSwitch = (mode) => {
    setCurrentMode(mode)
    switchAuthMode(mode)
  }

  return (
    <Dialog
      open={open}
      onClose={onClose}
      maxWidth="lg"
      fullWidth
      PaperProps={{
        style: {
          backgroundColor: "#1a1a1a",
          borderRadius: "16px",
          overflow: "hidden",
          maxWidth: "1000px",
          position: "relative", // Added for absolute positioning context
        },
      }}
    >
      {/* Close button positioned relative to Dialog */}
      <IconButton
        onClick={onClose}
        sx={{
          position: "absolute",
          right: "16px",
          top: "16px",
          color: "rgba(255, 255, 255, 0.5)",
          zIndex: 1,
          "&:hover": {
            backgroundColor: "rgba(255, 255, 255, 0.1)",
          },
        }}
      >
        <CloseIcon />
      </IconButton>

      <Box sx={{ display: "flex", minHeight: "600px" }}>
        {/* Left Section */}
        <Box
          sx={{
            flex: "1",
            p: 4,
            position: "relative",
          }}
        >
          <Box sx={{ maxWidth: "400px", mx: "auto", pt: 4 }}>
            <Box sx={{ display: "flex", justifyContent: "center", mb: 3 }}>
              <img src="/logo.png" alt="SAGE.AI Logo" height="60" />
            </Box>

            <Typography variant="h4" sx={{ mb: 1, fontWeight: 600, textAlign: "center" }}>
              {currentMode === "signin" ? "Welcome Back" : "Create Account"}
            </Typography>

            <OAuthButtons />

            <Box sx={{ my: 3, display: "flex", alignItems: "center" }}>
              <Box sx={{ flex: 1, height: "1px", bgcolor: "rgba(255, 255, 255, 0.1)" }} />
              <Typography sx={{ px: 2, color: "rgba(255, 255, 255, 0.5)" }}>OR</Typography>
              <Box sx={{ flex: 1, height: "1px", bgcolor: "rgba(255, 255, 255, 0.1)" }} />
            </Box>

            {currentMode === "signin" && (
              <SignInForm
                onSignUpClick={() => handleModeSwitch("signup")}
                onForgotClick={() => handleModeSwitch("forgot")}
              />
            )}

            {currentMode === "signup" && <SignUpForm onSignInClick={() => handleModeSwitch("signin")} />}

            {currentMode === "forgot" && <ForgotPasswordForm onBackToSignIn={() => handleModeSwitch("signin")} />}
          </Box>
        </Box>

        {/* Right Section */}
        <Box
          sx={{
            flex: "1",
            background: "linear-gradient(135deg, rgba(88, 101, 242, 0.3), rgba(235, 69, 158, 0.3))",
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            p: 6,
          }}
        >
          <Box sx={{ position: "relative", maxWidth: "400px" }}>
            <Typography
              variant="h3"
              sx={{
                fontWeight: 700,
                mb: 3,
                color: "white",
                textAlign: "center",
                display: "flex",
                flexDirection: "column",
                alignItems: "center",
                gap: 1
              }}
            >
              <span>Your AI-Powered SQL</span>
              <Box
                sx={{
                  position: "relative",
                  height: "60px",
                  width: "100%",
                  maxWidth: "400px",
                  overflow: "hidden"
                }}
              >
                <AnimatePresence mode="wait">
                  <motion.span
                    key={textIndex}
                    initial={{ y: 20, opacity: 0 }}
                    animate={{ y: 0, opacity: 1 }}
                    exit={{ y: -20, opacity: 0 }}
                    transition={{ duration: 0.5 }}
                    style={{
                      position: "absolute",
                      left: 0,
                      right: 0,
                      textAlign: "center",
                      background: "linear-gradient(135deg, #5865F2 0%, #EB459E 100%)",
                      WebkitBackgroundClip: "text",
                      WebkitTextFillColor: "transparent",
                      backgroundClip: "text",
                      display: "block",
                      fontWeight: 800,
                      fontSize: "0.9em"
                    }}
                  >
                    {ROTATING_TEXTS[textIndex]}
                  </motion.span>
                </AnimatePresence>
              </Box>
            </Typography>
            <Typography
              variant="h6"
              sx={{
                color: "rgba(255, 255, 255, 0.8)",
                textAlign: "center",
                lineHeight: 1.6,
              }}
            >
              Transform natural language into powerful SQL queries instantly with{" "}
              <Box component="span" sx={{ display: "inline-flex", alignItems: "center", verticalAlign: "middle" }}>
                <img src="/logo.png" alt="SAGE.AI Logo" height="24" style={{ marginLeft: "4px", marginRight: "4px" }} />
              </Box>
            </Typography>
          </Box>
        </Box>
      </Box>
    </Dialog>
  )
}

export default AuthModal

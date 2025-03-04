"use client"

import { useState, useEffect, useRef } from "react"
import { Dialog, Box, IconButton, useTheme, Typography } from "@mui/material"
import CloseIcon from "@mui/icons-material/Close"
import { useAuth } from "./AuthContext"
import SignInForm from "./SignInForm"
import SignUpForm from "./SignUpForm"
import ForgotPasswordForm from "./ForgotPasswordForm"
import OAuthButtons from "./OAuthButtons"
import { motion } from "framer-motion"

const ROTATING_TEXTS = ["Assistant", "Companion", "Generator"]

const AuthModal = ({ open, onClose, initialMode = "signin" }) => {
  const theme = useTheme()
  const { switchAuthMode, authMode } = useAuth()
  const [currentMode, setCurrentMode] = useState(initialMode)
  const [displayText, setDisplayText] = useState(ROTATING_TEXTS[0])
  const [isAnimating, setIsAnimating] = useState(false)
  const textIndexRef = useRef(0)
  const animationTimeoutRef = useRef(null)

  useEffect(() => {
    if (open) {
      setCurrentMode(initialMode)
    }
  }, [open, initialMode])

  useEffect(() => {
    const rotateText = () => {
      setIsAnimating(true)

      // After fade out, change the text
      animationTimeoutRef.current = setTimeout(() => {
        textIndexRef.current = (textIndexRef.current + 1) % ROTATING_TEXTS.length
        setDisplayText(ROTATING_TEXTS[textIndexRef.current])
        setIsAnimating(false)
      }, 500)
    }

    const interval = setInterval(rotateText, 2500)

    return () => {
      clearInterval(interval)
      if (animationTimeoutRef.current) {
        clearTimeout(animationTimeoutRef.current)
      }
    }
  }, [])

  

  const handleModeSwitch = (mode) => {
    setCurrentMode(mode)
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
          position: "relative",
        },
      }}
    >
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
              {currentMode === "signin" ? "Sign in" : "Create account"}
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
            position: "relative",
            overflow: "hidden",
          }}
        >
          <Box
            sx={{
              position: "relative",
              maxWidth: "400px",
              width: "100%",
            }}
          >
            <Box sx={{ textAlign: "center" }}>
              <Typography
                variant="h2"
                component="div"
                sx={{
                  fontWeight: 700,
                  mb: 3,
                  fontSize: { xs: "2rem", sm: "2.5rem", md: "3rem" },
                  lineHeight: 1.2,
                  background: "linear-gradient(to right, #fff, #e0e0e0)",
                  WebkitBackgroundClip: "text",
                  WebkitTextFillColor: "transparent",
                  display: "flex",
                  flexDirection: "column",
                  alignItems: "center",
                }}
              >
                <span>Your AI-Powered SQL</span>
                <motion.div
                  animate={{
                    opacity: isAnimating ? 0 : 1,
                    y: isAnimating ? -20 : 0,
                  }}
                  transition={{ duration: 0.5, ease: "easeInOut" }}
                  style={{
                    background: "linear-gradient(to right, #5865F2, #EB459E)",
                    WebkitBackgroundClip: "text",
                    WebkitTextFillColor: "transparent",
                    marginTop: "0.2em",
                    height: "1.2em",
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "center",
                  }}
                >
                  {displayText}
                </motion.div>
              </Typography>
            </Box>

            <Typography
              variant="h6"
              sx={{
                color: "rgba(255, 255, 255, 0.9)",
                textAlign: "center",
                lineHeight: 1.6,
                mt: 4,
                fontSize: { xs: "1rem", sm: "1.25rem" },
              }}
            >
              Transform natural language into powerful SQL queries instantly with{" "}
              <Box
                component="span"
                sx={{
                  display: "inline-flex",
                  alignItems: "center",
                  verticalAlign: "middle",
                  mx: 1,
                }}
              >
                <img
                  src="/logo.png"
                  alt="SAGE.AI Logo"
                  height="24"
                  style={{
                    filter: "brightness(1.2) contrast(1.1)",
                  }}
                />
              </Box>
            </Typography>
          </Box>
        </Box>
      </Box>
    </Dialog>
  )
}

export default AuthModal


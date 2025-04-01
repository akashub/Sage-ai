"use client"

import { Box, Typography, Container, Button, Avatar, AvatarGroup } from "@mui/material"
import { useInView } from "react-intersection-observer"
import { ArrowOutward } from "@mui/icons-material"
import { useAuth } from "../auth/AuthContext"

const EarlyAccessSection = ({ onAuthOpen }) => {
  const { ref, inView } = useInView({
    triggerOnce: true,
    threshold: 0.1,
  })

  const { openAuthModal } = useAuth()

  return (
    <Box
      sx={{
        py: 15,
        background: "linear-gradient(180deg, rgba(0,0,0,0.9) 0%, rgba(64,78,237,0.7) 100%)",
        borderRadius: "24px",
        mx: { xs: 2, md: 12 },
        my: 15,
        boxShadow: "0 8px 32px rgba(0, 0, 0, 0.1)",
        backdropFilter: "blur(10px)",
      }}
    >
      <Container maxWidth="lg">
        <Box
          ref={ref}
          sx={{
            display: "flex",
            flexDirection: { xs: "column", md: "row" },
            alignItems: "center",
            justifyContent: "space-between",
            gap: 4,
          }}
        >
          <Box sx={{ maxWidth: "600px" }}>
            <AvatarGroup
              sx={{
                mb: 3,
                opacity: inView ? 1 : 0,
                transform: inView ? "translateX(0)" : "translateX(-20px)",
                transition: "all 0.6s ease-out",
              }}
            >
              {[...Array(4)].map((_, i) => (
                <Avatar
                  key={i}
                  src={`/avatar-${i + 1}.png`}
                  alt={`User ${i + 1}`}
                  sx={{
                    width: 48,
                    height: 48,
                    border: "2px solid #404EED",
                  }}
                />
              ))}
            </AvatarGroup>
            <Typography
              variant="h2"
              sx={{
                fontWeight: 800,
                mb: 2,
                opacity: inView ? 1 : 0,
                transform: inView ? "translateY(0)" : "translateY(20px)",
                transition: "all 0.6s ease-out 0.2s",
              }}
            >
              Say goodbye to complex SQL queries
              <br />
              and hello to AI-powered simplicity
            </Typography>
            <Typography
              variant="h6"
              sx={{
                color: "rgba(255, 255, 255, 0.8)",
                mb: 4,
                opacity: inView ? 1 : 0,
                transform: inView ? "translateY(0)" : "translateY(20px)",
                transition: "all 0.6s ease-out 0.3s",
              }}
            >
              Join thousands of developers who are already simplifying their database workflows
            </Typography>
          </Box>
          <Box
            sx={{
              opacity: inView ? 1 : 0,
              transform: inView ? "translateY(0)" : "translateY(20px)",
              transition: "all 0.6s ease-out 0.4s",
            }}
          >
            <Button
              variant="contained"
              size="large"
              endIcon={<ArrowOutward />}
              onClick={openAuthModal}
              sx={{
                backgroundColor: "#5865F2",
                color: "white",
                px: 4,
                py: 2,
                fontSize: "1.2rem",
                "&:hover": {
                  backgroundColor: "#4752C4",
                },
              }}
            >
              Request Early Access
            </Button>
          </Box>
        </Box>
      </Container>
    </Box>
  )
}

export default EarlyAccessSection
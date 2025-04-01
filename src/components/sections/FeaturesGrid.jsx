"use client"

import { Box, Typography, Container, Grid, Paper } from "@mui/material"
import { useInView } from "react-intersection-observer"
import { Code, Storage, Speed, AutoFixHigh } from "@mui/icons-material"

const features = [
  {
    icon: <Code />,
    title: "Natural Language to SQL",
    description: "Convert plain English to perfect SQL queries instantly",
    color: "#FF4D4D",
  },
  {
    icon: <Storage />,
    title: "Multi-Database Support",
    description: "Works with PostgreSQL, MySQL, and more",
    color: "#00D4FF",
  },
  {
    icon: <Speed />,
    title: "Real-time Optimization",
    description: "Get the most efficient queries automatically",
    color: "#FFD700",
  },
  {
    icon: <AutoFixHigh />,
    title: "Smart Suggestions",
    description: "AI-powered query improvements and fixes",
    color: "#4CAF50",
  },
]

const FeaturesGrid = () => {
  const { ref, inView } = useInView({
    triggerOnce: true,
    threshold: 0.1,
  })

  return (
    <Box
      sx={{
        py: 15,
        background: "linear-gradient(180deg, rgba(64,78,237,0.9) 0%, rgba(64,78,237,0.7) 100%)",
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
            display: "grid",
            gridTemplateColumns: { xs: "1fr", md: "1fr 1fr" },
            gap: 6,
            alignItems: "center",
          }}
        >
          <Box>
            <Typography
              variant="h2"
              sx={{
                fontWeight: 800,
                mb: 3,
                opacity: inView ? 1 : 0,
                transform: inView ? "translateY(0)" : "translateY(20px)",
                transition: "all 0.6s ease-out",
              }}
            >
              POWERFUL
              <br />
              FEATURES
            </Typography>
            <Typography
              variant="h6"
              sx={{
                color: "rgba(255, 255, 255, 0.8)",
                mb: 4,
                opacity: inView ? 1 : 0,
                transform: inView ? "translateY(0)" : "translateY(20px)",
                transition: "all 0.6s ease-out 0.2s",
              }}
            >
              Everything you need to manage your database queries efficiently
            </Typography>
          </Box>
          <Grid container spacing={3}>
            {features.map((feature, index) => (
              <Grid item xs={12} sm={6} key={index}>
                <Paper
                  sx={{
                    p: 3,
                    height: "100%",
                    background: "rgba(255, 255, 255, 0.05)",
                    backdropFilter: "blur(10px)",
                    borderRadius: 4,
                    border: "1px solid rgba(255, 255, 255, 0.1)",
                    transition: "all 0.3s ease-in-out",
                    opacity: inView ? 1 : 0,
                    transform: inView ? "translateY(0)" : "translateY(20px)",
                    transitionDelay: `${index * 0.1}s`,
                    "&:hover": {
                      transform: "translateY(-8px)",
                      boxShadow: `0 12px 24px rgba(${feature.color}, 0.2)`,
                    },
                  }}
                >
                  <Box
                    sx={{
                      color: feature.color,
                      mb: 2,
                      "& > svg": {
                        fontSize: 40,
                      },
                    }}
                  >
                    {feature.icon}
                  </Box>
                  <Typography variant="h6" gutterBottom>
                    {feature.title}
                  </Typography>
                  <Typography variant="body2" sx={{ color: "rgba(255, 255, 255, 0.7)" }}>
                    {feature.description}
                  </Typography>
                </Paper>
              </Grid>
            ))}
          </Grid>
        </Box>
      </Container>
    </Box>
  )
}

export default FeaturesGrid
import { Box, Typography, Container } from "@mui/material"
import { useInView } from "react-intersection-observer"

const DemoSection = () => {
  const { ref, inView } = useInView({
    triggerOnce: true,
    threshold: 0.1,
  })

  return (
    <Box
      sx={{
        py: 15,
        background: "linear-gradient(180deg, rgba(64,78,237,0.7) 0%, rgba(64,78,237,0.5) 100%)",
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
            textAlign: "center",
            opacity: inView ? 1 : 0,
            transform: inView ? "translateY(0)" : "translateY(20px)",
            transition: "all 0.6s ease-out",
          }}
        >
          <Typography variant="h2" sx={{ fontWeight: 800, mb: 3 }}>
            Watch Demo
          </Typography>
          <Typography
            variant="h6"
            sx={{
              color: "rgba(255, 255, 255, 0.8)",
              mb: 6,
              maxWidth: "600px",
              mx: "auto",
            }}
          >
            See how easy it is to transform natural language into SQL queries
          </Typography>
          <Box
            sx={{
              position: "relative",
              width: "100%",
              maxWidth: "800px",
              mx: "auto",
              aspectRatio: "16/9",
              borderRadius: "16px",
              overflow: "hidden",
              background: "rgba(0, 0, 0, 0.2)",
              backdropFilter: "blur(10px)",
              border: "1px solid rgba(255, 255, 255, 0.1)",
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
            }}
          >
            <Typography
              variant="h3"
              sx={{
                fontWeight: 800,
                background: "linear-gradient(90deg, #00D4FF, #FF4D4D)",
                backgroundSize: "200% 200%",
                animation: "gradient 5s ease infinite",
                WebkitBackgroundClip: "text",
                WebkitTextFillColor: "transparent",
                "@keyframes gradient": {
                  "0%": { backgroundPosition: "0% 50%" },
                  "50%": { backgroundPosition: "100% 50%" },
                  "100%": { backgroundPosition: "0% 50%" },
                },
              }}
            >
              Coming Soon
            </Typography>
          </Box>
        </Box>
      </Container>
    </Box>
  )
}

export default DemoSection
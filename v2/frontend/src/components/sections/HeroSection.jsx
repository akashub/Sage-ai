import { useState, useEffect } from "react";
import { Box, Typography, Button, Container } from "@mui/material";
import { keyframes } from "@mui/system";
import { PlayArrow, Code } from "@mui/icons-material";

const gradientText = keyframes`
  0% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
  100% { background-position: 0% 50%; }
`;

const float = keyframes`
  0% { transform: translate(0, 0) rotate(0deg); }
  25% { transform: translate(10px, -10px) rotate(2deg); }
  50% { transform: translate(0, -20px) rotate(0deg); }
  75% { transform: translate(-10px, -10px) rotate(-2deg); }
  100% { transform: translate(0, 0) rotate(0deg); }
`;

const sqlSnippets = [
  {
    code: "SELECT * FROM users WHERE created_at > NOW() - INTERVAL '7 days'",
    position: { top: "8%", left: "5%" },
    rotation: "-15deg",
    delay: "0s",
  },
  {
    code: "INSERT INTO orders (user_id, product_id, quantity) VALUES (1, 100, 5)",
    position: { top: "8%", right: "2%" },
    rotation: "10deg",
    delay: "0.2s",
  },
  {
    code: "UPDATE products SET stock = stock - 1 WHERE id = 123",
    position: { bottom: "30%", left: "5%" },
    rotation: "5deg",
    delay: "0.4s",
  },
  {
    code: "DELETE FROM cart WHERE last_modified < NOW() - INTERVAL '24 hours'",
    position: { bottom: "30%", right: "5%" },
    rotation: "-8deg",
    delay: "0.6s",
  },
];

const FloatingCode = ({ code, position, rotation, delay }) => (
  <Box
    sx={{
      position: "absolute",
      ...position,
      padding: 2,
      backgroundColor: "rgba(0, 0, 0, 0.7)",
      borderRadius: 2,
      maxWidth: 300,
      animation: `${float} 6s ease-in-out infinite`,
      animationDelay: delay,
      transform: `rotate(${rotation})`,
      display: { xs: "none", md: "block" },
      zIndex: 1,
      backdropFilter: "blur(10px)",
      border: "1px solid rgba(255, 255, 255, 0.1)",
    }}
  >
    <Typography
      sx={{
        fontFamily: "monospace",
        color: "#FFD700",
        fontSize: "0.8rem",
        whiteSpace: "pre-wrap",
      }}
    >
      {code}
    </Typography>
  </Box>
);

const HeroSection = ({ onAuthOpen }) => {
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);

  return (
    <Box
      sx={{
        minHeight: "100vh",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        position: "relative",
        overflow: "hidden",
        pt: 12,
      }}
    >
      {sqlSnippets.map((snippet, index) => (
        <FloatingCode key={index} {...snippet} />
      ))}

      <Container maxWidth="lg" sx={{ position: "relative", zIndex: 2 }}>
        <Box
          sx={{
            textAlign: "center",
            opacity: mounted ? 1 : 0,
            transform: mounted ? "translateY(0)" : "translateY(20px)",
            transition: "all 0.6s ease-out",
          }}
        >
          <Typography
            variant="h1"
            sx={{
              fontSize: { xs: "4rem", md: "8rem" },
              fontWeight: 900,
              lineHeight: 1,
              mb: 2,
              background: "linear-gradient(90deg, #00D4FF, #FF4D4D)",
              backgroundSize: "200% 200%",
              animation: `${gradientText} 5s ease infinite`,
              WebkitBackgroundClip: "text",
              WebkitTextFillColor: "transparent",
              textTransform: "uppercase",
            }}
          >
            Generate
            <br />
            SQL
          </Typography>
          <Typography
            variant="h6"
            sx={{
              mb: 4,
              color: "#D3D3D3",
              maxWidth: "800px",
              mx: "auto",
            }}
          >
            Manage, generate SQL queries.
          </Typography>
          <Box sx={{ display: "flex", gap: 2, justifyContent: "center" }}>
            <Button
              variant="contained"
              size="large"
              startIcon={<Code />}
              onClick={onAuthOpen}
              sx={{
                backgroundColor: "white",
                color: "background.default",
                "&:hover": {
                  backgroundColor: "rgba(255, 255, 255, 0.9)",
                },
              }}
            >
              Try Generator
            </Button>
            <Button
              variant="outlined"
              size="large"
              startIcon={<PlayArrow />}
              sx={{
                borderColor: "white",
                color: "white",
                "&:hover": {
                  borderColor: "rgba(255, 255, 255, 0.9)",
                  backgroundColor: "rgba(255, 255, 255, 0.1)",
                },
              }}
            >
              Watch Demo
            </Button>
          </Box>
        </Box>
      </Container>
    </Box>
  );
};

export default HeroSection;

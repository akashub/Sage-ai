"use client"

import { Box, Typography, Container, Paper } from "@mui/material"
import { useInView } from "react-intersection-observer"
import { keyframes } from "@mui/system"
import CodeIcon from "@mui/icons-material/Code"
import SearchIcon from "@mui/icons-material/Search"
import BarChartIcon from "@mui/icons-material/BarChart"
import SpeedIcon from "@mui/icons-material/Speed"

const slideIn = keyframes`
  from {
    opacity: 0;
    transform: translateY(50px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
`

const float = keyframes`
  0%, 100% { transform: translateY(0px); }
  50% { transform: translateY(-10px); }
`

const pulse = keyframes`
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.03); }
`

const typeAndErase = keyframes`
  0%, 33% { width: 0; }
  16.5%, 49.5% { width: 100%; }
  66%, 100% { width: 0; }
`

const words = [
  {
    text: "WRITE",
    description: "Natural language to SQL",
    icon: <CodeIcon />,
    color: "#FF4D4D",
    commands: [
      "SELECT * FROM users WHERE created_at > DATE_SUB(NOW(), INTERVAL 7 DAY);",
      "INSERT INTO orders (user_id, product_id, quantity) VALUES (1, 100, 2);",
      "UPDATE products SET stock = stock - 1 WHERE id = 123;",
    ],
  },
  {
    text: "QUERY",
    description: "Instant database access",
    icon: <SearchIcon />,
    color: "#00D4FF",
    commands: [
      "SELECT p.name, c.category FROM products p JOIN categories c ON p.category_id = c.id;",
      "SELECT AVG(price) FROM products GROUP BY category_id;",
      "SELECT u.name, COUNT(o.id) FROM users u LEFT JOIN orders o ON u.id = o.user_id GROUP BY u.id;",
    ],
  },
  {
    text: "ANALYZE",
    description: "Smart data insights",
    icon: <BarChartIcon />,
    color: "#9C27B0",
    chart: true,
  },
  {
    text: "OPTIMIZE",
    description: "Performance tuning",
    icon: <SpeedIcon />,
    color: "#4CAF50",
    optimize: true,
  },
]

const ActionWords = () => {
  return (
    <Box
      sx={{
        py: 10,
        position: "relative",
        overflow: "hidden",
      }}
    >
      <Container maxWidth="lg">
        {words.map((word, index) => (
          <WordSection key={word.text} word={word} index={index} />
        ))}
      </Container>
    </Box>
  )
}

const WordSection = ({ word, index }) => {
  const { ref, inView } = useInView({
    threshold: 0.2,
    triggerOnce: true,
  })

  return (
    <Box
      ref={ref}
      sx={{
        minHeight: "40vh",
        display: "flex",
        alignItems: "center",
        justifyContent: "space-between",
        flexDirection: index % 2 === 0 ? "row" : "row-reverse",
        mb: 8,
        opacity: inView ? 1 : 0,
        animation: inView ? `${slideIn} 0.8s ease-out ${index * 0.2}s` : "none",
        "@media (max-width: 900px)": {
          flexDirection: "column",
          textAlign: "center",
          gap: 3,
        },
      }}
    >
      <Box sx={{ flex: 1 }}>
        <Typography
          variant="h2"
          sx={{
            fontSize: { xs: "2.5rem", md: "3.5rem" },
            fontWeight: 800,
            color: word.color,
            mb: 1,
            display: "flex",
            alignItems: "center",
            gap: 2,
          }}
        >
          {word.icon} {word.text}
        </Typography>
        <Typography
          variant="h5"
          sx={{
            color: "rgba(255, 255, 255, 0.8)",
            fontWeight: 500,
          }}
        >
          {word.description}
        </Typography>
      </Box>
      <Box
        sx={{
          flex: 1,
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
        }}
      >
        <FrameBox color={word.color}>
          {word.chart ? (
            <AnalyzeChart color={word.color} />
          ) : word.optimize ? (
            <OptimizeIllustration color={word.color} />
          ) : (
            <TerminalIllustration commands={word.commands} />
          )}
        </FrameBox>
      </Box>
    </Box>
  )
}

const FrameBox = ({ children, color }) => (
  <Box
    sx={{
      position: "relative",
      width: "100%",
      maxWidth: 400,
      height: 250,
      borderRadius: 4,
      overflow: "hidden",
      "&::before": {
        content: '""',
        position: "absolute",
        top: 0,
        left: 0,
        right: 0,
        bottom: 0,
        background: `linear-gradient(45deg, ${color}33, transparent)`,
        zIndex: -1,
      },
      "&::after": {
        content: '""',
        position: "absolute",
        top: 2,
        left: 2,
        right: 2,
        bottom: 2,
        borderRadius: 3,
        border: `2px solid ${color}66`,
        zIndex: 1,
        pointerEvents: "none",
      },
    }}
  >
    {children}
  </Box>
)

const TerminalIllustration = ({ commands }) => (
  <Paper
    elevation={3}
    sx={{
      background: "rgba(0, 0, 0, 0.7)",
      borderRadius: 2,
      p: 2,
      width: "100%",
      height: "100%",
      display: "flex",
      flexDirection: "column",
      overflow: "hidden",
    }}
  >
    <Box sx={{ mb: 2, display: "flex", gap: 1 }}>
      <Box sx={{ width: 8, height: 8, borderRadius: "50%", backgroundColor: "#FF5F56" }} />
      <Box sx={{ width: 8, height: 8, borderRadius: "50%", backgroundColor: "#FFBD2E" }} />
      <Box sx={{ width: 8, height: 8, borderRadius: "50%", backgroundColor: "#27C93F" }} />
    </Box>
    <Box sx={{ flexGrow: 1, position: "relative", overflow: "hidden" }}>
      {commands.map((command, index) => (
        <Typography
          key={index}
          variant="body2"
          sx={{
            fontFamily: "monospace",
            color: "#FFFFFF",
            position: "absolute",
            top: `${index * 33.33}%`,
            left: 0,
            whiteSpace: "nowrap",
            overflow: "hidden",
            fontSize: "0.8rem",
            animation: `${typeAndErase} 15s linear infinite`,
            animationDelay: `${index * 5}s`,
            "&::after": {
              content: '"|"',
              animation: `blink 0.7s infinite`,
            },
          }}
        >
          {command}
        </Typography>
      ))}
    </Box>
  </Paper>
)

const AnalyzeChart = ({ color }) => (
  <Box
    sx={{
      width: "100%",
      height: "100%",
      position: "relative",
      animation: `${float} 6s ease-in-out infinite`,
    }}
  >
    <Box
      sx={{
        position: "absolute",
        bottom: 20,
        left: 20,
        right: 20,
        height: "calc(100% - 40px)",
        display: "flex",
        alignItems: "flex-end",
        justifyContent: "space-around",
      }}
    >
      {[40, 70, 30, 90, 60, 80, 50].map((height, index) => (
        <Box
          key={index}
          sx={{
            width: "12%",
            height: `${height}%`,
            backgroundColor: color,
            opacity: 0.7,
            borderTopLeftRadius: 4,
            borderTopRightRadius: 4,
            transition: "height 0.5s ease-in-out",
            "&:hover": {
              height: `${height + 10}%`,
              opacity: 1,
            },
          }}
        />
      ))}
    </Box>
    <Box
      sx={{
        position: "absolute",
        top: 20,
        left: 20,
        right: 20,
        bottom: 20,
        backgroundImage: `
          linear-gradient(to right, ${color}4D 1px, transparent 1px),
          linear-gradient(to bottom, ${color}4D 1px, transparent 1px),
          radial-gradient(circle at center, ${color}4D 1px, transparent 1px)
        `,
        backgroundSize: "14.28% 20%, 14.28% 20%, 14.28% 20%",
      }}
    />
  </Box>
)

const OptimizeIllustration = ({ color }) => (
  <Box
    sx={{
      width: "100%",
      height: "100%",
      position: "relative",
      animation: `${pulse} 4s ease-in-out infinite`,
    }}
  >
    <Box
      sx={{
        position: "absolute",
        top: "50%",
        left: "50%",
        transform: "translate(-50%, -50%)",
        width: "80%",
        height: "80%",
        border: `3px solid ${color}`,
        borderRadius: "50%",
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
      }}
    >
      <Box
        sx={{
          width: "60%",
          height: "60%",
          border: `2px solid ${color}`,
          borderRadius: "50%",
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
        }}
      >
        <Box
          sx={{
            width: "40%",
            height: "40%",
            backgroundColor: color,
            borderRadius: "50%",
          }}
        />
      </Box>
    </Box>
    <Box
      sx={{
        position: "absolute",
        top: "10%",
        left: "50%",
        height: "40%",
        width: "3px",
        backgroundColor: color,
        transform: "translateX(-50%)",
        transformOrigin: "bottom",
        animation: `${rotate} 2s linear infinite`,
        "&::after": {
          content: '""',
          position: "absolute",
          top: 0,
          left: "-3px",
          width: 0,
          height: 0,
          borderLeft: "4px solid transparent",
          borderRight: "4px solid transparent",
          borderBottom: `8px solid ${color}`,
        },
      }}
    />
  </Box>
)

const rotate = keyframes`
  from { transform: translateX(-50%) rotate(0deg); }
  to { transform: translateX(-50%) rotate(360deg); }
`

export default ActionWords


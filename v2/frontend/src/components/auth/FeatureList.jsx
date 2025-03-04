"use client"
import { Box, Typography, useTheme } from "@mui/material"
import { motion } from "framer-motion"
import { keyframes } from "@mui/system"

const gradientAnimation = keyframes`
  0% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
  100% { background-position: 0% 50%; }
`

const features = [
  {
    icon: "⚡️",
    title: "Save Hours of Work",
    description: "Turn simple text prompts into ready-to-publish SQL queries in minutes, not days",
    gradient: "linear-gradient(45deg, #FF6B6B, #FFE66D)",
    shadowColor: "rgba(255, 107, 107, 0.2)",
  },
  {
    icon: "📊",
    title: "Tell Better Stories",
    description: "Create engaging database queries that help you understand your data better",
    gradient: "linear-gradient(45deg, #4ECDC4, #556270)",
    shadowColor: "rgba(78, 205, 196, 0.2)",
  },
  {
    icon: "🎯",
    title: "Stand Out on Social",
    description: "Get more insights with professional-looking queries that capture attention",
    gradient: "linear-gradient(45deg, #A8E6CF, #FFD3B6)",
    shadowColor: "rgba(168, 230, 207, 0.2)",
  },
]

const FeatureItem = ({ feature, index }) => {
  const theme = useTheme()

  return (
    <motion.div
      initial={{ opacity: 0, y: 20, scale: 0.95 }}
      animate={{ opacity: 1, y: 0, scale: 1 }}
      transition={{
        duration: 0.5,
        delay: index * 0.2,
        type: "spring",
        stiffness: 100,
      }}
      whileHover={{
        scale: 1.02,
        transition: { duration: 0.2 },
      }}
      whileTap={{ scale: 0.98 }}
    >
      <Box
        sx={{
          mb: 2,
          p: 3,
          borderRadius: "16px",
          background: `${feature.gradient}`,
          position: "relative",
          overflow: "hidden",
          boxShadow: `0 8px 32px ${feature.shadowColor}`,
          "&::before": {
            content: '""',
            position: "absolute",
            top: 0,
            left: 0,
            right: 0,
            bottom: 0,
            background: "rgba(255, 255, 255, 0.1)",
            backdropFilter: "blur(10px)",
            borderRadius: "inherit",
          },
          "&::after": {
            content: '""',
            position: "absolute",
            top: 0,
            left: 0,
            right: 0,
            bottom: 0,
            background: feature.gradient,
            opacity: 0.7,
            animation: `${gradientAnimation} 5s ease infinite`,
            backgroundSize: "200% 200%",
            mixBlendMode: "overlay",
            borderRadius: "inherit",
          },
        }}
      >
        <Box sx={{ position: "relative", zIndex: 1 }}>
          <motion.div
            initial={{ scale: 1 }}
            whileHover={{ scale: 1.2, rotate: [0, -10, 10, -10, 0] }}
            transition={{ duration: 0.5 }}
          >
            <Typography
              variant="h3"
              sx={{
                fontSize: "2.5rem",
                mb: 1,
                filter: "drop-shadow(0 2px 4px rgba(0,0,0,0.2))",
              }}
            >
              {feature.icon}
            </Typography>
          </motion.div>

          <Typography
            variant="h6"
            sx={{
              fontWeight: 700,
              mb: 1,
              color: "white",
              textShadow: "0 2px 4px rgba(0,0,0,0.2)",
              fontFamily: '"gg sans", "Noto Sans", "Helvetica Neue", Helvetica, Arial, sans-serif',
            }}
          >
            {feature.title}
          </Typography>

          <Typography
            variant="body1"
            sx={{
              color: "rgba(255, 255, 255, 0.9)",
              fontWeight: 500,
              lineHeight: 1.6,
              textShadow: "0 1px 2px rgba(0,0,0,0.1)",
              fontFamily: '"gg sans", "Noto Sans", "Helvetica Neue", Helvetica, Arial, sans-serif',
            }}
          >
            {feature.description}
          </Typography>
        </Box>
      </Box>
    </motion.div>
  )
}

export const FeatureList = () => {
  return (
    <Box
      sx={{
        p: 2,
        background: "rgba(0,0,0,0.2)",
        borderRadius: "20px",
        backdropFilter: "blur(10px)",
      }}
    >
      <Typography
        variant="h5"
        sx={{
          mb: 3,
          fontWeight: 700,
          background: "linear-gradient(45deg, #FF6B6B, #4ECDC4)",
          backgroundClip: "text",
          textFillColor: "transparent",
          WebkitBackgroundClip: "text",
          WebkitTextFillColor: "transparent",
          textAlign: "center",
          fontSize: "1.8rem",
          textShadow: "0 2px 4px rgba(0,0,0,0.1)",
          fontFamily: '"gg sans", "Noto Sans", "Helvetica Neue", Helvetica, Arial, sans-serif',
        }}
      >
        Why Choose SAGE.AI?
      </Typography>

      <Box sx={{ mt: 4 }}>
        {features.map((feature, index) => (
          <FeatureItem key={feature.title} feature={feature} index={index} />
        ))}
      </Box>
    </Box>
  )
}


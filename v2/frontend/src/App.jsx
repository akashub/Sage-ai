import { ThemeProvider, createTheme, CssBaseline } from "@mui/material"
import Navigation from "./components/layout/Navigation"
import HeroSection from "./components/sections/HeroSection"
import FeaturesGrid from "./components/sections/FeaturesGrid"
import ActionWords from "./components/sections/ActionWords"
import DemoSection from "./components/sections/DemoSection"
import SupportedPlatforms from "./components/sections/SupportedPlatforms"
import EarlyAccessSection from "./components/sections/EarlyAccessSection"
import Footer from "./components/layout/Footer"

const theme = createTheme({
  palette: {
    mode: "dark",
    primary: {
      main: "#5865F2", // Discord blue
    },
    secondary: {
      main: "#EB459E", // Discord pink
    },
    background: {
      default: "#404EED", // Discord background
      paper: "#2F3136", // Discord dark
    },
    text: {
      primary: "#FFFFFF",
      secondary: "rgba(255, 255, 255, 0.7)",
    },
  },
  typography: {
    fontFamily: '"gg sans", "Noto Sans", "Helvetica Neue", Helvetica, Arial, sans-serif',
    h1: {
      fontSize: "4rem",
      fontWeight: 800,
      letterSpacing: "-0.025em",
    },
    h2: {
      fontSize: "3rem",
      fontWeight: 700,
      letterSpacing: "-0.025em",
    },
  },
  components: {
    MuiButton: {
      styleOverrides: {
        root: {
          borderRadius: "28px",
          padding: "16px 32px",
          fontSize: "1rem",
          textTransform: "none",
          fontWeight: 500,
        },
      },
    },
  },
})

function App() {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <div
        style={{
          background: "linear-gradient(180deg, rgba(5,3,55,1) 0%, rgba(57,7,112,1) 50%, rgba(3,33,98,1) 100%)",
          minHeight: "100vh",
          display: "flex",
          flexDirection: "column",
          position: "relative",
        }}
      >
        <Navigation />
        <main style={{ flex: 1, position: "relative", zIndex: 1 }}>
          <div id="home">
            <HeroSection />
          </div>
          <div id="features">
            <FeaturesGrid />
          </div>
          <ActionWords />
          <div id="demo">
            <DemoSection />
          </div>
          <div id="platforms">
            <SupportedPlatforms />
          </div>
          <div id="early-access">
            <EarlyAccessSection />
          </div>
        </main>
        <Footer />
      </div>
    </ThemeProvider>
  )
}

export default App


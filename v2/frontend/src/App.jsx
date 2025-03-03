import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import { ThemeProvider, createTheme, CssBaseline } from "@mui/material";
import Navigation from "./components/layout/Navigation";
import HeroSection from "./components/sections/HeroSection";
import FeaturesGrid from "./components/sections/FeaturesGrid";
import ActionWords from "./components/sections/ActionWords";
import DemoSection from "./components/sections/DemoSection";
import SupportedPlatforms from "./components/sections/SupportedPlatforms";
import EarlyAccessSection from "./components/sections/EarlyAccessSection";
import Footer from "./components/layout/Footer";
import { AuthProvider } from "./components/auth/AuthContext";
import AuthModalWrapper from "./components/auth/AuthModalWrapper";
import ChatInterface from "./pages/ChatInterface";


const theme = createTheme({
  palette: {
    mode: "dark",
    primary: {
      main: "#90caf9",
    },
    secondary: {
      main: "#f48fb1",
    },
    background: {
      default: "#050337",
      paper: "#390770",
    },
  },
  typography: {
    fontFamily: "'Roboto', sans-serif",
  },
});

function App() {
  return (
    <Router>
      <AuthProvider>
        <ThemeProvider theme={theme}>
          <CssBaseline />
          <Routes>
            <Route
              path="/"
              element={
                <>
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
                  <AuthModalWrapper />
                </>
              }
            />
            <Route path="/chat" element={<ChatInterface />} />
          </Routes>
        </ThemeProvider>
      </AuthProvider>
    </Router>
  );
}

export default App;

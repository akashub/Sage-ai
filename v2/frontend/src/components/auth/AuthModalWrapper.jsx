import { useAuth } from "./AuthContext"
import AuthModal from "./AuthModal"
import { Box } from "@mui/material"

const AuthModalWrapper = () => {
  const { authModalOpen, closeAuthModal, authMode } = useAuth()

  return (
    <>
      <Box
        sx={{
          position: "fixed",
          top: 0,
          left: 0,
          right: 0,
          bottom: 0,
          zIndex: 1299,
          backgroundColor: "rgba(0, 0, 0, 0.5)",
          backdropFilter: "blur(5px)",
          opacity: authModalOpen ? 1 : 0,
          transition: "opacity 0.3s ease",
          pointerEvents: authModalOpen ? "auto" : "none",
        }}
        onClick={closeAuthModal}
      />
      <Box
        sx={{
          position: "fixed",
          top: 0,
          left: 0,
          right: 0,
          bottom: 0,
          zIndex: 1300,
          pointerEvents: "none",
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
        }}
      >
        <AuthModal open={authModalOpen} onClose={closeAuthModal} initialMode={authMode} />
      </Box>
    </>
  )
}

export default AuthModalWrapper
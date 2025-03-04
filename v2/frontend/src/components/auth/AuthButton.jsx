"use client"

import { Button } from "@mui/material"
import { useAuth } from "./AuthContext"

const AuthButton = ({
  children = "Sign In",
  mode = "signin",
  variant = "contained",
  color = "primary",
  ...props
}) => {
  const { openAuthModal } = useAuth()

  const handleClick = () => {
    openAuthModal(mode)
  }

  return (
    <Button
      variant={variant}
      color={color}
      onClick={handleClick}
      sx={{
        borderRadius: '28px',
        textTransform: 'none',
        fontWeight: 500,
        ...props.sx
      }}
      {...props}
    >
      {children}
    </Button>
  )
}

export default AuthButton
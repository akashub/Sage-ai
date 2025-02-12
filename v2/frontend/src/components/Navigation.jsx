"use client"

import { useState } from "react"
import {
  AppBar,
  Toolbar,
  Button,
  Box,
  IconButton,
  Drawer,
  List,
  ListItem,
  ListItemText,
  useMediaQuery,
} from "@mui/material"
import { styled, useTheme } from "@mui/material/styles"
import { Apps, ArrowOutward, Menu as MenuIcon } from "@mui/icons-material"

const StyledAppBar = styled(AppBar)(({ theme }) => ({
  background: "transparent",
  boxShadow: "none",
  padding: theme.spacing(2, 0),
}))

const NavButton = styled(Button)(({ theme }) => ({
  color: theme.palette.text.primary,
  "&:hover": {
    backgroundColor: "rgba(255, 255, 255, 0.1)",
  },
}))

const Navigation = () => {
  const [mobileOpen, setMobileOpen] = useState(false)
  const theme = useTheme()
  const isMobile = useMediaQuery(theme.breakpoints.down("md"))

  const handleDrawerToggle = () => {
    setMobileOpen(!mobileOpen)
  }

  const navItems = [
    { text: "Features", href: "#features" },
    { text: "Docs", href: "#docs" },
    { text: "About", href: "#about" },
  ]

  const drawer = (
    <Box
      sx={{
        width: 250,
        height: "100%",
        background: "linear-gradient(135deg, #1a1c20 0%, #2c2f33 100%)",
      }}
    >
      <List>
        {navItems.map((item) => (
          <ListItem key={item.text} disablePadding>
            <Button
              fullWidth
              sx={{
                justifyContent: "flex-start",
                px: 3,
                py: 2,
                color: "white",
                "&:hover": {
                  backgroundColor: "rgba(255, 255, 255, 0.1)",
                },
              }}
              href={item.href}
            >
              <ListItemText primary={item.text} />
            </Button>
          </ListItem>
        ))}
        <ListItem disablePadding>
          <Button
            fullWidth
            variant="outlined"
            startIcon={<Apps />}
            sx={{
              m: 1,
              borderColor: "rgba(255, 255, 255, 0.8)",
              color: "white",
              "&:hover": {
                borderColor: "white",
                backgroundColor: "rgba(255, 255, 255, 0.3)",
              },
            }}
          >
            Launch App
          </Button>
        </ListItem>
        <ListItem disablePadding>
          <Button
            fullWidth
            variant="contained"
            sx={{
              m: 1,
              backgroundColor: "white",
              color: "background.default",
              display: "flex",
              alignItems: "center",
              gap: 0.2,
              "&:hover": {
                backgroundColor: "rgba(255, 255, 255, 0.9)",
              },
            }}
          >
            Join Early Access <ArrowOutward fontSize="small" />
          </Button>
        </ListItem>
      </List>
    </Box>
  )

  return (
    <StyledAppBar position="floating">
      <Toolbar sx={{ justifyContent: "space-between" }}>
        <Box component="img" src="/logo.png" alt="SQL AI" sx={{ height: 62 }} />

        {/* Desktop Navigation */}
        {!isMobile && (
          <Box sx={{ display: "flex", gap: 1 }}>
            {navItems.map((item) => (
              <NavButton key={item.text} href={item.href}>
                {item.text}
              </NavButton>
            ))}
            <Button
              variant="outlined"
              startIcon={<Apps />}
              sx={{
                borderColor: "rgba(255, 255, 255, 0.8)",
                color: "white",
                "&:hover": {
                  borderColor: "white",
                  backgroundColor: "rgba(255, 255, 255, 0.3)",
                },
              }}
            >
              Launch App
            </Button>
            <Button
              variant="contained"
              sx={{
                backgroundColor: "white",
                color: "background.default",
                display: "flex",
                alignItems: "center",
                gap: 0.2,
                "&:hover": {
                  backgroundColor: "rgba(255, 255, 255, 0.9)",
                },
              }}
            >
              Join Early Access <ArrowOutward fontSize="small" />
            </Button>
          </Box>
        )}

        {/* Mobile Navigation */}
        {isMobile && (
          <IconButton
            color="inherit"
            aria-label="open drawer"
            edge="start"
            onClick={handleDrawerToggle}
            sx={{
              ml: 2,
              "&:hover": {
                backgroundColor: "rgba(255, 255, 255, 0.1)",
              },
            }}
          >
            <MenuIcon />
          </IconButton>
        )}

        {/* Mobile Drawer */}
        <Drawer
          variant="temporary"
          anchor="right"
          open={mobileOpen}
          onClose={handleDrawerToggle}
          ModalProps={{
            keepMounted: true, // Better open performance on mobile.
          }}
          sx={{
            display: { xs: "block", md: "none" },
            "& .MuiDrawer-paper": {
              boxSizing: "border-box",
              width: 250,
              backgroundColor: "transparent",
              borderLeft: "none",
            },
          }}
        >
          {drawer}
        </Drawer>
      </Toolbar>
    </StyledAppBar>
  )
}

export default Navigation


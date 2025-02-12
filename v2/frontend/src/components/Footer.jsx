import { Box, Container, Grid, Typography, Button, Link, Stack } from "@mui/material"
import { styled } from "@mui/material/styles"
import { Apps, ArrowOutward, GitHub, Twitter, LinkedIn } from "@mui/icons-material"

const StyledLink = styled(Link)(({ theme }) => ({
  color: "rgba(255, 255, 255, 0.7)",
  textDecoration: "none",
  "&:hover": {
    color: "white",
    textDecoration: "none",
  },
}))

const Footer = () => {
  const currentYear = new Date().getFullYear()

  return (
    <Box
      component="footer"
      sx={{
        py: 8,
        px: { xs: 4, md: 8 },
        background: 'linear-gradient(180deg, rgba(64,78,237,0.2) 0%, rgba(0,0,0,0.1) 100%)',
        minHeight: '400px',
        display: 'flex',
        flexDirection: 'column',
      }}
    >
      <Grid container spacing={6} sx={{ maxWidth: '1400px', mx: 'auto', width: '100%' }}>
          {/* Logo and Description */}
          <Grid item xs={12} md={3}>
            <Box
              component="img"
              src="/logo.png"
              alt="SAGE.AI"
              sx={{
                height: 60,
                mb: 2,
              }}
            />
            <Typography variant="body2" sx={{ color: "rgba(255, 255, 255, 0.7)", mb: 2 }}>
              Transform natural language into powerful SQL queries instantly with AI.
            </Typography>
          </Grid>

          {/* Product Links */}
          <Grid item xs={6} md={1}>
            <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
              Product
            </Typography>
            <Stack spacing={1}>
              <StyledLink href="#">Features</StyledLink>
              <StyledLink href="#">Docs</StyledLink>
              <StyledLink href="#">What's New</StyledLink>
              <StyledLink href="#">Roadmap</StyledLink>
            </Stack>
          </Grid>

          {/* Company Links */}
          <Grid item xs={6} md={1}>
            <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
              Company
            </Typography>
            <Stack spacing={1}>
              <StyledLink href="#">About</StyledLink>
              <StyledLink href="#">Blog</StyledLink>
              <StyledLink href="#">Careers</StyledLink>
              <StyledLink href="#">Contact</StyledLink>
            </Stack>
          </Grid>

          {/* Resources Links */}
          <Grid item xs={6} md={1}>
            <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
              Resources
            </Typography>
            <Stack spacing={1}>
              <StyledLink href="#">Community</StyledLink>
              <StyledLink href="#">Support</StyledLink>
              <StyledLink href="#">Privacy</StyledLink>
              <StyledLink href="#">Terms</StyledLink>
            </Stack>
          </Grid>

          {/* CTA Buttons */}
          <Grid item xs={12} md={4}>
            <Stack spacing={2} direction="row" sx={{ justifyContent: 'flex-end' }}>
              <Button
                variant="outlined"
                sx={{
                  borderColor: 'rgba(255, 255, 255, 0.3)',
                  color: 'white',
                  py: 1,
                  px: 3,
                  '&:hover': {
                    borderColor: 'white',
                    backgroundColor: 'rgba(255, 255, 255, 0.1)',
                  },
                }}
              >
                Launch App
              </Button>
              <Button
                variant="contained"
                endIcon={<ArrowOutward />}
                sx={{
                  backgroundColor: 'white',
                  color: 'background.default',
                  py: 1,
                  px: 3,
                  '&:hover': {
                    backgroundColor: 'rgba(255, 255, 255, 0.9)',
                  },
                }}
              >
                Join Early Access
              </Button>
            </Stack>
          </Grid>
        </Grid>

        {/* Contact and Social Links */}
        <Box
          sx={{
            mt: 8,
            pt: 4,
            borderTop: '1px solid rgba(255, 255, 255, 0.1)',
            display: 'flex',
            flexDirection: { xs: 'column', md: 'row' },
            justifyContent: 'space-between',
            alignItems: 'center',
            maxWidth: '1400px',
            mx: 'auto',
            width: '100%',
          }}
        >
          <Typography
            variant="h6"
            sx={{
              color: "white",
              fontWeight: "normal",
              fontSize: "1.1rem",
            }}
          >
            hello@sage.ai
          </Typography>

          <Stack direction="row" spacing={2}>
            <IconLink href="https://github.com/sage-ai" icon={<GitHub />} />
            <IconLink href="https://twitter.com/sage_ai" icon={<Twitter />} />
            <IconLink href="https://linkedin.com/company/sage-ai" icon={<LinkedIn />} />
          </Stack>
        </Box>

        {/* Legal */}
        <Box
          sx={{
            mt: 4,
            display: 'flex',
            flexDirection: { xs: 'column', md: 'row' },
            justifyContent: 'space-between',
            alignItems: 'center',
            maxWidth: '1400px',
            mx: 'auto',
            width: '100%',
            color: 'rgba(255, 255, 255, 0.5)',
          }}
        >
          <Typography variant="body2" sx={{ color: "rgba(255, 255, 255, 0.7)" }}>
            © {currentYear} SAGE.AI Inc. All Rights Reserved.
          </Typography>
          <Stack
            direction="row"
            spacing={3}
            sx={{
              "& > a": {
                fontSize: "0.875rem",
              },
            }}
          >
            <StyledLink href="/terms">Terms of Service</StyledLink>
            <StyledLink href="/privacy">Privacy Policy</StyledLink>
          </Stack>
        </Box>
    </Box>
  )
}

// Helper component for social media icons
const IconLink = ({ href, icon }) => (
  <Link
    href={href}
    target="_blank"
    rel="noopener noreferrer"
    sx={{
      color: "rgba(255, 255, 255, 0.7)",
      "&:hover": {
        color: "white",
      },
    }}
  >
    {icon}
  </Link>
)

export default Footer

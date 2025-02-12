import { Box, Typography, Grid } from "@mui/material";
import { useInView } from "react-intersection-observer";

const platforms = [
  {
    name: "MySQL",
    logo: "/mysql-logo.png",
    alt: "MySQL Database",
  },
  {
    name: "PostgreSQL",
    logo: "/postgresql-logo.svg",
    alt: "PostgreSQL Database",
  },
  {
    name: "Microsoft SQL Server",
    logo: "/sqlserver-logo.png",
    alt: "Microsoft SQL Server",
  },
  {
    name: "SQLite",
    logo: "/sqlite-logo.svg",
    alt: "SQLite Database",
  },
  {
    name: "Oracle",
    logo: "/oracle-logo.png",
    alt: "Oracle Database",
  },
  {
    name: "MariaDB",
    logo: "/mariadb-logo.png",
    alt: "MariaDB Database",
  },
];

const SupportedPlatforms = () => {
  const { ref, inView } = useInView({
    triggerOnce: true,
    threshold: 0.1,
  });

  return (
    <Box
      sx={{
        py: 8,
        background: "linear-gradient(180deg, rgba(64,78,237,0.8) 0%, rgba(64,78,237,0.6) 100%)",
        borderRadius: "16px",
        mx: { xs: 2, md: 12 },
        my: 15,
        boxShadow: "0 4px 16px rgba(0, 0, 0, 0.1)",
        backdropFilter: "blur(8px)",
        textAlign: "center",
      }}
      ref={ref}
    >
      <Typography variant="h3" sx={{ fontWeight: 700, mb: 2 }}>
        Supported Platforms
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
        Generate optimized SQL queries for all major database platforms
      </Typography>
      <Grid container spacing={3} justifyContent="center">
        {platforms.map((platform, index) => (
          <Grid item xs={6} sm={4} md={4} key={platform.name} sx={{ textAlign: "center" }}>
            <Box
              component="img"
              src={platform.logo}
              alt={platform.alt}
              sx={{
                width: "auto",
                height: "120px",
                maxWidth: "140px",
                objectFit: "contain",
                filter: "brightness(1) contrast(1)",
                transition: "all 0.3s ease-in-out",
                "&:hover": {
                  transform: "scale(1.05)",
                },
              }}
            />
          </Grid>
        ))}
      </Grid>
    </Box>
  );
};

export default SupportedPlatforms;
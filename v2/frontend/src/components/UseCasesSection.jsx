import { Box, Typography, Container, Grid, Card, CardContent, CardActionArea } from "@mui/material"
import { useInView } from "react-intersection-observer"
import CodeIcon from "@mui/icons-material/Code"
import AnalyticsIcon from "@mui/icons-material/Analytics"
import TrendingUpIcon from "@mui/icons-material/TrendingUp"
import SchoolIcon from "@mui/icons-material/School"

const useCases = [
  { icon: <CodeIcon />, title: "Developers", description: "Quickly generate SQL without manual writing." },
  { icon: <AnalyticsIcon />, title: "Data Analysts", description: "Fetch insights without deep SQL knowledge." },
  { icon: <TrendingUpIcon />, title: "Product Managers", description: "Query data efficiently for reports." },
  { icon: <SchoolIcon />, title: "Students", description: "Learn SQL interactively." },
]

const UseCasesSection = () => {
  const { ref, inView } = useInView({
    triggerOnce: true,
    threshold: 0.1,
  })

  return (
    <Box sx={{ py: 8, bgcolor: "background.default" }}>
      <Container maxWidth="lg">
        <Typography variant="h3" component="h2" align="center" gutterBottom>
          Who Can Benefit?
        </Typography>
        <Grid container spacing={4} ref={ref}>
          {useCases.map((useCase, index) => (
            <Grid item xs={12} sm={6} md={3} key={index}>
              <Card
                sx={{
                  height: "100%",
                  display: "flex",
                  flexDirection: "column",
                  transition: "all 0.3s ease-in-out",
                  transform: inView ? "rotateY(0deg)" : "rotateY(90deg)",
                  transitionDelay: `${index * 100}ms`,
                }}
              >
                <CardActionArea sx={{ height: "100%" }}>
                  <CardContent
                    sx={{
                      height: "100%",
                      display: "flex",
                      flexDirection: "column",
                      alignItems: "center",
                      textAlign: "center",
                    }}
                  >
                    <Box
                      sx={{
                        mb: 2,
                        color: "secondary.main",
                        fontSize: "3rem",
                        "& > svg": {
                          fontSize: "inherit",
                        },
                      }}
                    >
                      {useCase.icon}
                    </Box>
                    <Typography variant="h6" component="h3" gutterBottom>
                      {useCase.title}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      {useCase.description}
                    </Typography>
                  </CardContent>
                </CardActionArea>
              </Card>
            </Grid>
          ))}
        </Grid>
      </Container>
    </Box>
  )
}

export default UseCasesSection


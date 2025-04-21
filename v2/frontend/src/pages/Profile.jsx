"use client"

import React from 'react';
import { 
  Container, 
  Typography, 
  Paper, 
  Box, 
  Avatar, 
  Grid, 
  Button, 
  Divider, 
  Chip,
  Card,
  CardContent,
  CardActions,
  IconButton,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  Switch,
  useTheme
} from '@mui/material';
import { 
  Settings as SettingsIcon,
  Storage as StorageIcon,
  Speed as SpeedIcon,
  Security as SecurityIcon,
  Notifications as NotificationsIcon,
  Language as LanguageIcon,
  DarkMode as DarkModeIcon,
  Edit as EditIcon,
  CheckCircle as CheckCircleIcon
} from '@mui/icons-material';
import { useAuth } from '../components/auth/AuthContext';
import Navigation from '../components/layout/Navigation';

const PlanFeature = ({ included, text }) => (
  <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
    <CheckCircleIcon
      sx={{
        mr: 1,
        color: included ? 'success.main' : 'text.disabled',
        fontSize: '1.2rem'
      }}
    />
    <Typography
      variant="body2"
      sx={{ color: included ? 'text.primary' : 'text.disabled' }}
    >
      {text}
    </Typography>
  </Box>
);

const Profile = () => {
  const { user } = useAuth();
  const theme = useTheme();

  if (!user) {
    return (
      <>
        <Navigation />
        <Container>
          <Typography variant="h4" sx={{ mt: 4, mb: 2 }}>
            Please log in to view your profile
          </Typography>
        </Container>
      </>
    );
  }

  // Mock data - replace with actual data from your backend
  const usageStats = {
    queries: 45,
    tokens: 12000,
    files: 3
  };

  const plans = [
    {
      name: 'Free',
      price: '$0',
      features: [
        { text: '100 queries per month', included: true },
        { text: 'Basic SQL generation', included: true },
        { text: 'Standard support', included: true },
        { text: 'Advanced features', included: false },
        { text: 'Priority support', included: false }
      ],
      current: user?.plan === 'Free'
    },
    {
      name: 'Pro',
      price: '$19',
      features: [
        { text: 'Unlimited queries', included: true },
        { text: 'Advanced SQL generation', included: true },
        { text: 'Priority support', included: true },
        { text: 'Custom database integration', included: true },
        { text: 'Team collaboration', included: true }
      ],
      current: user?.plan === 'Pro'
    }
  ];

  return (
    <>
      <Navigation />
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Grid container spacing={4}>
          {/* Profile Overview */}
          <Grid item xs={12} md={4}>
            <Paper elevation={3} sx={{ p: 4, height: '100%' }}>
              <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', mb: 3 }}>
                <Avatar
                  src={user?.picture}
                  alt={user?.name}
                  sx={{ 
                    width: 120, 
                    height: 120,
                    mb: 2,
                    fontSize: '3rem',
                    bgcolor: 'primary.main'
                  }}
                >
                  {user?.name?.[0]?.toUpperCase() || 'U'}
                </Avatar>
                <Typography variant="h5" gutterBottom>
                  {user?.name || 'User'}
                </Typography>
                <Typography variant="body2" color="text.secondary" paragraph>
                  {user?.email || 'No email provided'}
                </Typography>
                <Chip 
                  label={`${user?.plan || 'Free'} Plan`} 
                  color="primary" 
                  size="small"
                />
              </Box>
              <Divider sx={{ my: 2 }} />
              <Typography variant="subtitle2" color="text.secondary" gutterBottom>
                Member since
              </Typography>
              <Typography variant="body2">
                {user?.createdAt ? new Date(user.createdAt).toLocaleDateString() : 'Unknown'}
              </Typography>
            </Paper>
          </Grid>

          {/* Usage Statistics */}
          <Grid item xs={12} md={8}>
            <Paper elevation={3} sx={{ p: 4, mb: 4 }}>
              <Typography variant="h6" gutterBottom>
                Usage Statistics
              </Typography>
              <Grid container spacing={3}>
                <Grid item xs={12} sm={4}>
                  <Card>
                    <CardContent>
                      <Typography variant="h4" color="primary">
                        {usageStats.queries}
                      </Typography>
                      <Typography variant="body2" color="text.secondary">
                        Queries Used
                      </Typography>
                    </CardContent>
                  </Card>
                </Grid>
                <Grid item xs={12} sm={4}>
                  <Card>
                    <CardContent>
                      <Typography variant="h4" color="primary">
                        {usageStats.tokens.toLocaleString()}
                      </Typography>
                      <Typography variant="body2" color="text.secondary">
                        Tokens Processed
                      </Typography>
                    </CardContent>
                  </Card>
                </Grid>
                <Grid item xs={12} sm={4}>
                  <Card>
                    <CardContent>
                      <Typography variant="h4" color="primary">
                        {usageStats.files}
                      </Typography>
                      <Typography variant="body2" color="text.secondary">
                        Files Uploaded
                      </Typography>
                    </CardContent>
                  </Card>
                </Grid>
              </Grid>
            </Paper>

            {/* Subscription Plans */}
            <Paper elevation={3} sx={{ p: 4, mb: 4 }}>
              <Typography variant="h6" gutterBottom>
                Subscription Plans
              </Typography>
              <Grid container spacing={3}>
                {plans.map((plan) => (
                  <Grid item xs={12} sm={6} key={plan.name}>
                    <Card
                      sx={{
                        height: '100%',
                        display: 'flex',
                        flexDirection: 'column',
                        position: 'relative',
                        ...(plan.current && {
                          borderColor: 'primary.main',
                          borderWidth: 2,
                          borderStyle: 'solid'
                        })
                      }}
                    >
                      {plan.current && (
                        <Chip
                          label="Current Plan"
                          color="primary"
                          size="small"
                          sx={{
                            position: 'absolute',
                            top: 16,
                            right: 16
                          }}
                        />
                      )}
                      <CardContent sx={{ flexGrow: 1 }}>
                        <Typography variant="h6" gutterBottom>
                          {plan.name}
                        </Typography>
                        <Typography variant="h4" color="primary" gutterBottom>
                          {plan.price}
                          <Typography variant="caption" color="text.secondary">
                            /month
                          </Typography>
                        </Typography>
                        <Box sx={{ mt: 2 }}>
                          {plan.features.map((feature, index) => (
                            <PlanFeature
                              key={index}
                              included={feature.included}
                              text={feature.text}
                            />
                          ))}
                        </Box>
                      </CardContent>
                      <CardActions sx={{ p: 2, pt: 0 }}>
                        <Button
                          fullWidth
                          variant={plan.current ? "outlined" : "contained"}
                          disabled={plan.current}
                        >
                          {plan.current ? 'Current Plan' : 'Upgrade'}
                        </Button>
                      </CardActions>
                    </Card>
                  </Grid>
                ))}
              </Grid>
            </Paper>

            {/* Settings */}
            <Paper elevation={3} sx={{ p: 4 }}>
              <Typography variant="h6" gutterBottom>
                Settings
              </Typography>
              <List>
                <ListItem>
                  <ListItemIcon>
                    <NotificationsIcon />
                  </ListItemIcon>
                  <ListItemText primary="Email Notifications" />
                  <Switch edge="end" defaultChecked />
                </ListItem>
                <ListItem>
                  <ListItemIcon>
                    <DarkModeIcon />
                  </ListItemIcon>
                  <ListItemText primary="Dark Mode" />
                  <Switch edge="end" defaultChecked />
                </ListItem>
                <ListItem>
                  <ListItemIcon>
                    <LanguageIcon />
                  </ListItemIcon>
                  <ListItemText primary="Language" secondary="English" />
                </ListItem>
                <ListItem>
                  <ListItemIcon>
                    <SecurityIcon />
                  </ListItemIcon>
                  <ListItemText primary="Two-Factor Authentication" />
                  <Button variant="outlined" size="small">
                    Enable
                  </Button>
                </ListItem>
              </List>
            </Paper>
          </Grid>
        </Grid>
      </Container>
    </>
  );
};

export default Profile; 
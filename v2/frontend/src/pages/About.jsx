"use client"

import React from 'react';
import {
  Box,
  Container,
  Typography,
  Grid,
  Paper,
  Avatar,
  Divider,
  useTheme,
  Button
} from '@mui/material';
import {
  Code as CodeIcon,
  Security as SecurityIcon,
  Speed as SpeedIcon,
  Group as GroupIcon,
  Rocket as RocketIcon
} from '@mui/icons-material';
import Navigation from '../components/layout/Navigation';

const About = () => {
  const theme = useTheme();

  const features = [
    {
      icon: <CodeIcon sx={{ fontSize: 40 }} />,
      title: 'Advanced AI',
      description: 'Powered by state-of-the-art language models for accurate SQL generation'
    },
    {
      icon: <SecurityIcon sx={{ fontSize: 40 }} />,
      title: 'Enterprise Security',
      description: 'Bank-grade encryption and security protocols to protect your data'
    },
    {
      icon: <SpeedIcon sx={{ fontSize: 40 }} />,
      title: 'Lightning Fast',
      description: 'Optimized for performance with minimal latency'
    }
  ];

  const team = [
    {
      name: 'Akash Singh',
      role: 'Backend Developer',
      image: 'https://i.pravatar.cc/150?img=1'
    },
    {
      name: 'Nitin Reddy',
      role: 'Backend Developer',
      image: 'https://i.pravatar.cc/150?img=2'
    },
    {
      name: 'Sudiksha Rajavaram',
      role: 'Frontend Developer',
      image: 'https://i.pravatar.cc/150?img=3'
    },
    {
      name: 'Yash Kishore',
      role: 'Frontend Developer',
      image: 'https://i.pravatar.cc/150?img=4'
    }
  ];

  return (
    <>
      <Navigation />
      <Container maxWidth="lg" sx={{ py: 4 }}>
        {/* Hero Section */}
        <Box sx={{ textAlign: 'center', mb: 6 }}>
          <Typography variant="h2" gutterBottom>
            About Sage.AI
          </Typography>
          <Typography variant="h5" color="text.secondary" paragraph>
            Revolutionizing Data Analysis with Artificial Intelligence
          </Typography>
        </Box>

        {/* Mission Statement */}
        <Paper elevation={2} sx={{ p: 4, mb: 6 }}>
          <Box sx={{ display: 'flex', alignItems: 'center', mb: 3 }}>
            <RocketIcon sx={{ fontSize: 40, mr: 2, color: 'primary.main' }} />
            <Typography variant="h4">
              Our Mission
            </Typography>
          </Box>
          <Typography variant="body1" paragraph>
            At Sage AI, we're on a mission to democratize data analysis by making it accessible to everyone, 
            regardless of their technical background. We believe that the power of data should be available 
            to all, and we're building the tools to make that possible.
          </Typography>
          <Typography variant="body1">
            Our platform combines cutting-edge AI technology with intuitive design to help businesses 
            unlock the full potential of their data.
          </Typography>
        </Paper>

        {/* Features */}
        <Box sx={{ mb: 6 }}>
          <Typography variant="h4" gutterBottom sx={{ textAlign: 'center', mb: 4 }}>
            Why Choose Sage AI?
          </Typography>
          <Grid container spacing={4}>
            {features.map((feature, index) => (
              <Grid item xs={12} md={4} key={index}>
                <Paper elevation={2} sx={{ p: 4, height: '100%', textAlign: 'center' }}>
                  <Box sx={{ color: 'primary.main', mb: 2 }}>
                    {feature.icon}
                  </Box>
                  <Typography variant="h6" gutterBottom>
                    {feature.title}
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    {feature.description}
                  </Typography>
                </Paper>
              </Grid>
            ))}
          </Grid>
        </Box>

        {/* Team Section */}
        <Box sx={{ mb: 6 }}>
          <Typography variant="h4" gutterBottom sx={{ textAlign: 'center', mb: 4 }}>
            Meet Our Team
          </Typography>
          <Grid container spacing={4}>
            {team.map((member, index) => (
              <Grid item xs={12} sm={6} md={3} key={index}>
                <Paper elevation={2} sx={{ p: 3, textAlign: 'center' }}>
                  <Avatar
                    src={member.image}
                    alt={member.name}
                    sx={{ width: 120, height: 120, mx: 'auto', mb: 2 }}
                  />
                  <Typography variant="h6" gutterBottom>
                    {member.name}
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    {member.role}
                  </Typography>
                </Paper>
              </Grid>
            ))}
          </Grid>
        </Box>

        {/* CTA Section */}
        <Paper 
          elevation={2} 
          sx={{ 
            p: 6, 
            textAlign: 'center',
            background: `linear-gradient(45deg, ${theme.palette.primary.main}, ${theme.palette.primary.dark})`,
            color: 'white'
          }}
        >
          <Typography variant="h4" gutterBottom>
            Ready to Transform Your Data Analysis?
          </Typography>
          <Typography variant="body1" paragraph sx={{ mb: 4 }}>
            Join thousands of businesses already using Sage AI to make data-driven decisions.
          </Typography>
          <Button 
            variant="contained" 
            color="secondary" 
            size="large"
            sx={{ 
              px: 4, 
              py: 1.5,
              fontSize: '1.1rem'
            }}
          >
            Get Started
          </Button>
        </Paper>
      </Container>
    </>
  );
};

export default About; 
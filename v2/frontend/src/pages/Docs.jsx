"use client"

import React from 'react';
import {
  Box,
  Container,
  Typography,
  Paper,
  Grid,
  List,
  ListItem,
  ListItemText,
  ListItemIcon,
  Divider,
  useTheme,
} from '@mui/material';
import {
  School as SchoolIcon,
  Code as CodeIcon,
  Storage as StorageIcon,
  Security as SecurityIcon,
  Api as ApiIcon,
  ContactSupport as ContactSupportIcon,
} from '@mui/icons-material';
import Navigation from '../components/layout/Navigation';

const Docs = () => {
  const theme = useTheme();

  const sections = [
    {
      title: 'Getting Started',
      icon: <SchoolIcon />,
      content: 'Learn how to upload your CSV files and start analyzing your data with natural language queries.'
    },
    {
      title: 'Features',
      icon: <CodeIcon />,
      content: 'Explore our powerful features including natural language processing, SQL generation, and data visualization.'
    },
    {
      title: 'Database Integration',
      icon: <StorageIcon />,
      content: 'Understand how Sage AI integrates with your data and processes complex queries.'
    },
    {
      title: 'Security',
      icon: <SecurityIcon />,
      content: 'Learn about our security measures and how we protect your data.'
    },
    {
      title: 'API Reference',
      icon: <ApiIcon />,
      content: 'Detailed documentation of our API endpoints and integration options.'
    }
  ];

  return (
    <>
      <Navigation />
      <Container maxWidth="lg" sx={{ py: 4 }}>
        {/* Header */}
        <Box sx={{ mb: 6, textAlign: 'center' }}>
          <Typography variant="h3" gutterBottom>
            Documentation
          </Typography>
          <Typography variant="h6" color="text.secondary" sx={{ mb: 4 }}>
            Everything you need to know about using Sage AI
          </Typography>
        </Box>

        {/* Main Content */}
        <Grid container spacing={4}>
          {/* Sidebar */}
          <Grid item xs={12} md={3}>
            <Paper elevation={2} sx={{ p: 2 }}>
              <Typography variant="h6" sx={{ mb: 2 }}>
                Quick Links
              </Typography>
              <List>
                {sections.map((section) => (
                  <ListItem 
                    key={section.title} 
                    button
                    sx={{
                      borderRadius: 1,
                      mb: 1,
                      '&:hover': {
                        bgcolor: 'action.hover',
                      }
                    }}
                  >
                    <ListItemIcon sx={{ minWidth: 40 }}>
                      {section.icon}
                    </ListItemIcon>
                    <ListItemText primary={section.title} />
                  </ListItem>
                ))}
              </List>
            </Paper>
          </Grid>

          {/* Content */}
          <Grid item xs={12} md={9}>
            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 4 }}>
              {sections.map((section, index) => (
                <Paper 
                  key={section.title}
                  elevation={2}
                  sx={{ p: 4 }}
                >
                  <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                    <Box sx={{ 
                      mr: 2,
                      p: 1,
                      borderRadius: 1,
                      bgcolor: 'primary.main',
                      color: 'white',
                      display: 'flex'
                    }}>
                      {section.icon}
                    </Box>
                    <Typography variant="h5">
                      {section.title}
                    </Typography>
                  </Box>
                  <Typography variant="body1" color="text.secondary" sx={{ mb: 2 }}>
                    {section.content}
                  </Typography>
                  <Divider sx={{ my: 2 }} />
                  <Typography variant="body2" color="text.secondary">
                    This section provides detailed information about {section.title.toLowerCase()}. 
                    Check our documentation for more details.
                  </Typography>
                </Paper>
              ))}
            </Box>
          </Grid>
        </Grid>

        {/* Support Section */}
        <Paper 
          elevation={2} 
          sx={{ 
            mt: 6, 
            p: 4, 
            textAlign: 'center'
          }}
        >
          <ContactSupportIcon sx={{ fontSize: 48, color: 'primary.main', mb: 2 }} />
          <Typography variant="h5" gutterBottom>
            Need Help?
          </Typography>
          <Typography variant="body1" color="text.secondary">
            Our support team is here to help you with any questions or issues.
            Contact us at support@sageai.com
          </Typography>
        </Paper>
      </Container>
    </>
  );
};

export default Docs; 
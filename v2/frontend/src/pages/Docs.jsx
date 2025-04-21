"use client"

import React from 'react';
import {
  Container,
  Typography,
  Box,
  Grid,
  Paper,
  List,
  ListItem,
  ListItemText,
  Divider,
  Link,
  Accordion,
  AccordionSummary,
  AccordionDetails,
  useTheme
} from '@mui/material';
import {
  ExpandMore as ExpandMoreIcon,
  Code as CodeIcon,
  Storage as StorageIcon,
  Security as SecurityIcon,
  Speed as SpeedIcon,
  Help as HelpIcon
} from '@mui/icons-material';
import Navigation from '../components/layout/Navigation';

const Docs = () => {
  const theme = useTheme();

  const sections = [
    {
      title: 'Getting Started',
      icon: <HelpIcon />,
      content: [
        {
          title: 'Introduction',
          content: 'Sage AI is an advanced SQL query generation and data analysis platform. This documentation will help you get started with using our platform effectively.'
        },
        {
          title: 'Quick Start',
          content: 'To begin using Sage AI, simply sign up for an account and connect your database. You can start by asking questions in natural language, and our AI will generate the appropriate SQL queries.'
        }
      ]
    },
    {
      title: 'Features',
      icon: <SpeedIcon />,
      content: [
        {
          title: 'Natural Language Processing',
          content: 'Convert your questions into SQL queries using natural language. No SQL knowledge required!'
        },
        {
          title: 'Query Optimization',
          content: 'Our AI automatically optimizes your queries for better performance and efficiency.'
        }
      ]
    },
    {
      title: 'Database Integration',
      icon: <StorageIcon />,
      content: [
        {
          title: 'Supported Databases',
          content: 'Sage AI supports major databases including PostgreSQL, MySQL, MongoDB, and more.'
        },
        {
          title: 'Connection Setup',
          content: 'Learn how to securely connect your database to Sage AI.'
        }
      ]
    },
    {
      title: 'Security',
      icon: <SecurityIcon />,
      content: [
        {
          title: 'Data Protection',
          content: 'Your data is encrypted in transit and at rest. We never store your database credentials.'
        },
        {
          title: 'Access Control',
          content: 'Manage user permissions and access levels for your team.'
        }
      ]
    },
    {
      title: 'API Reference',
      icon: <CodeIcon />,
      content: [
        {
          title: 'Authentication',
          content: 'Learn how to authenticate with our API using OAuth 2.0.'
        },
        {
          title: 'Endpoints',
          content: 'Comprehensive documentation of all available API endpoints.'
        }
      ]
    }
  ];

  return (
    <>
      <Navigation />
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Box sx={{ mb: 4 }}>
          <Typography variant="h3" gutterBottom>
            Documentation
          </Typography>
          <Typography variant="subtitle1" color="text.secondary" paragraph>
            Everything you need to know about using Sage AI
          </Typography>
        </Box>

        <Grid container spacing={4}>
          {/* Main Content */}
          <Grid item xs={12} md={8}>
            {sections.map((section, index) => (
              <Paper
                key={index}
                elevation={2}
                sx={{ mb: 4, overflow: 'hidden' }}
              >
                <Box sx={{ p: 3, bgcolor: 'primary.main', color: 'white' }}>
                  <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                    {section.icon}
                    <Typography variant="h5" sx={{ ml: 1 }}>
                      {section.title}
                    </Typography>
                  </Box>
                </Box>
                <Box sx={{ p: 3 }}>
                  {section.content.map((item, itemIndex) => (
                    <Box key={itemIndex} sx={{ mb: 3 }}>
                      <Typography variant="h6" gutterBottom>
                        {item.title}
                      </Typography>
                      <Typography variant="body1" color="text.secondary">
                        {item.content}
                      </Typography>
                    </Box>
                  ))}
                </Box>
              </Paper>
            ))}
          </Grid>

          {/* Sidebar */}
          <Grid item xs={12} md={4}>
            <Paper elevation={2} sx={{ p: 3 }}>
              <Typography variant="h6" gutterBottom>
                Quick Links
              </Typography>
              <List>
                {sections.map((section, index) => (
                  <React.Fragment key={index}>
                    <ListItem>
                      <ListItemText
                        primary={section.title}
                        primaryTypographyProps={{ variant: 'subtitle2' }}
                      />
                    </ListItem>
                    {index < sections.length - 1 && <Divider />}
                  </React.Fragment>
                ))}
              </List>
            </Paper>

            <Paper elevation={2} sx={{ p: 3, mt: 3 }}>
              <Typography variant="h6" gutterBottom>
                Need Help?
              </Typography>
              <Typography variant="body2" color="text.secondary" paragraph>
                If you need additional assistance, please contact our support team.
              </Typography>
              <Link href="/contact" color="primary">
                Contact Support
              </Link>
            </Paper>
          </Grid>
        </Grid>
      </Container>
    </>
  );
};

export default Docs; 
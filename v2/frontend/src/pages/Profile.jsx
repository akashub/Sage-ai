// "use client"

// import React from 'react';
// import { 
//   Container, 
//   Typography, 
//   Paper, 
//   Box, 
//   Avatar, 
//   Grid, 
//   Button, 
//   Divider, 
//   Chip,
//   Card,
//   CardContent,
//   CardActions,
//   IconButton,
//   List,
//   ListItem,
//   ListItemIcon,
//   ListItemText,
//   Switch,
//   useTheme
// } from '@mui/material';
// import { 
//   Settings as SettingsIcon,
//   Storage as StorageIcon,
//   Speed as SpeedIcon,
//   Security as SecurityIcon,
//   Notifications as NotificationsIcon,
//   Language as LanguageIcon,
//   DarkMode as DarkModeIcon,
//   Edit as EditIcon,
//   CheckCircle as CheckCircleIcon
// } from '@mui/icons-material';
// import { useAuth } from '../components/auth/AuthContext';
// import Navigation from '../components/layout/Navigation';

// const PlanFeature = ({ included, text }) => (
//   <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
//     <CheckCircleIcon
//       sx={{
//         mr: 1,
//         color: included ? 'success.main' : 'text.disabled',
//         fontSize: '1.2rem'
//       }}
//     />
//     <Typography
//       variant="body2"
//       sx={{ color: included ? 'text.primary' : 'text.disabled' }}
//     >
//       {text}
//     </Typography>
//   </Box>
// );

// const Profile = () => {
//   const { user } = useAuth();
//   const theme = useTheme();

//   if (!user) {
//     return (
//       <>
//         <Navigation />
//         <Container>
//           <Typography variant="h4" sx={{ mt: 4, mb: 2 }}>
//             Please log in to view your profile
//           </Typography>
//         </Container>
//       </>
//     );
//   }

//   // Mock data - replace with actual data from your backend
//   const usageStats = {
//     queries: 45,
//     tokens: 12000,
//     files: 3
//   };

//   const plans = [
//     {
//       name: 'Free',
//       price: '$0',
//       features: [
//         { text: '100 queries per month', included: true },
//         { text: 'Basic SQL generation', included: true },
//         { text: 'Standard support', included: true },
//         { text: 'Advanced features', included: false },
//         { text: 'Priority support', included: false }
//       ],
//       current: user?.plan === 'Free'
//     },
//     {
//       name: 'Pro',
//       price: '$19',
//       features: [
//         { text: 'Unlimited queries', included: true },
//         { text: 'Advanced SQL generation', included: true },
//         { text: 'Priority support', included: true },
//         { text: 'Custom database integration', included: true },
//         { text: 'Team collaboration', included: true }
//       ],
//       current: user?.plan === 'Pro'
//     }
//   ];

//   return (
//     <>
//       <Navigation />
//       <Container maxWidth="lg" sx={{ py: 4 }}>
//         <Grid container spacing={4}>
//           {/* Profile Overview */}
//           <Grid item xs={12} md={4}>
//             <Paper elevation={3} sx={{ p: 4, height: '100%' }}>
//               <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', mb: 3 }}>
//                 <Avatar
//                   src={user?.picture}
//                   alt={user?.name}
//                   sx={{ 
//                     width: 120, 
//                     height: 120,
//                     mb: 2,
//                     fontSize: '3rem',
//                     bgcolor: 'primary.main'
//                   }}
//                 >
//                   {user?.name?.[0]?.toUpperCase() || 'U'}
//                 </Avatar>
//                 <Typography variant="h5" gutterBottom>
//                   {user?.name || 'User'}
//                 </Typography>
//                 <Typography variant="body2" color="text.secondary" paragraph>
//                   {user?.email || 'No email provided'}
//                 </Typography>
//                 <Chip 
//                   label={`${user?.plan || 'Free'} Plan`} 
//                   color="primary" 
//                   size="small"
//                 />
//               </Box>
//               <Divider sx={{ my: 2 }} />
//               <Typography variant="subtitle2" color="text.secondary" gutterBottom>
//                 Member since
//               </Typography>
//               <Typography variant="body2">
//                 {user?.createdAt ? new Date(user.createdAt).toLocaleDateString() : 'Unknown'}
//               </Typography>
//             </Paper>
//           </Grid>

//           {/* Usage Statistics */}
//           <Grid item xs={12} md={8}>
//             <Paper elevation={3} sx={{ p: 4, mb: 4 }}>
//               <Typography variant="h6" gutterBottom>
//                 Usage Statistics
//               </Typography>
//               <Grid container spacing={3}>
//                 <Grid item xs={12} sm={4}>
//                   <Card>
//                     <CardContent>
//                       <Typography variant="h4" color="primary">
//                         {usageStats.queries}
//                       </Typography>
//                       <Typography variant="body2" color="text.secondary">
//                         Queries Used
//                       </Typography>
//                     </CardContent>
//                   </Card>
//                 </Grid>
//                 <Grid item xs={12} sm={4}>
//                   <Card>
//                     <CardContent>
//                       <Typography variant="h4" color="primary">
//                         {usageStats.tokens.toLocaleString()}
//                       </Typography>
//                       <Typography variant="body2" color="text.secondary">
//                         Tokens Processed
//                       </Typography>
//                     </CardContent>
//                   </Card>
//                 </Grid>
//                 <Grid item xs={12} sm={4}>
//                   <Card>
//                     <CardContent>
//                       <Typography variant="h4" color="primary">
//                         {usageStats.files}
//                       </Typography>
//                       <Typography variant="body2" color="text.secondary">
//                         Files Uploaded
//                       </Typography>
//                     </CardContent>
//                   </Card>
//                 </Grid>
//               </Grid>
//             </Paper>

//             {/* Subscription Plans */}
//             <Paper elevation={3} sx={{ p: 4, mb: 4 }}>
//               <Typography variant="h6" gutterBottom>
//                 Subscription Plans
//               </Typography>
//               <Grid container spacing={3}>
//                 {plans.map((plan) => (
//                   <Grid item xs={12} sm={6} key={plan.name}>
//                     <Card
//                       sx={{
//                         height: '100%',
//                         display: 'flex',
//                         flexDirection: 'column',
//                         position: 'relative',
//                         ...(plan.current && {
//                           borderColor: 'primary.main',
//                           borderWidth: 2,
//                           borderStyle: 'solid'
//                         })
//                       }}
//                     >
//                       {plan.current && (
//                         <Chip
//                           label="Current Plan"
//                           color="primary"
//                           size="small"
//                           sx={{
//                             position: 'absolute',
//                             top: 16,
//                             right: 16
//                           }}
//                         />
//                       )}
//                       <CardContent sx={{ flexGrow: 1 }}>
//                         <Typography variant="h6" gutterBottom>
//                           {plan.name}
//                         </Typography>
//                         <Typography variant="h4" color="primary" gutterBottom>
//                           {plan.price}
//                           <Typography variant="caption" color="text.secondary">
//                             /month
//                           </Typography>
//                         </Typography>
//                         <Box sx={{ mt: 2 }}>
//                           {plan.features.map((feature, index) => (
//                             <PlanFeature
//                               key={index}
//                               included={feature.included}
//                               text={feature.text}
//                             />
//                           ))}
//                         </Box>
//                       </CardContent>
//                       <CardActions sx={{ p: 2, pt: 0 }}>
//                         <Button
//                           fullWidth
//                           variant={plan.current ? "outlined" : "contained"}
//                           disabled={plan.current}
//                         >
//                           {plan.current ? 'Current Plan' : 'Upgrade'}
//                         </Button>
//                       </CardActions>
//                     </Card>
//                   </Grid>
//                 ))}
//               </Grid>
//             </Paper>

//             {/* Settings */}
//             <Paper elevation={3} sx={{ p: 4 }}>
//               <Typography variant="h6" gutterBottom>
//                 Settings
//               </Typography>
//               <List>
//                 <ListItem>
//                   <ListItemIcon>
//                     <NotificationsIcon />
//                   </ListItemIcon>
//                   <ListItemText primary="Email Notifications" />
//                   <Switch edge="end" defaultChecked />
//                 </ListItem>
//                 <ListItem>
//                   <ListItemIcon>
//                     <DarkModeIcon />
//                   </ListItemIcon>
//                   <ListItemText primary="Dark Mode" />
//                   <Switch edge="end" defaultChecked />
//                 </ListItem>
//                 <ListItem>
//                   <ListItemIcon>
//                     <LanguageIcon />
//                   </ListItemIcon>
//                   <ListItemText primary="Language" secondary="English" />
//                 </ListItem>
//                 <ListItem>
//                   <ListItemIcon>
//                     <SecurityIcon />
//                   </ListItemIcon>
//                   <ListItemText primary="Two-Factor Authentication" />
//                   <Button variant="outlined" size="small">
//                     Enable
//                   </Button>
//                 </ListItem>
//               </List>
//             </Paper>
//           </Grid>
//         </Grid>
//       </Container>
//     </>
//   );
// };

// export default Profile; 

// frontend/src/pages/Profile.jsx
import React, { useState, useEffect } from "react";
import { 
  Box, 
  Container, 
  Typography, 
  Paper, 
  Avatar, 
  Button, 
  Grid, 
  Divider,
  Card,
  CardContent,
  CardHeader,
  List,
  ListItem,
  ListItemText,
  ListItemAvatar,
  ListItemSecondaryAction,
  IconButton,
  TextField,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Switch,
  FormControlLabel,
  Chip
} from "@mui/material";
import {
  Person as PersonIcon,
  Edit as EditIcon,
  Save as SaveIcon,
  Add as AddIcon,
  Delete as DeleteIcon,
  Key as KeyIcon,
  Check as CheckIcon
} from "@mui/icons-material";
import { useAuth } from "../components/auth/AuthContext";
import { getUserProfile, updateUserProfile, getAPIKeys, saveAPIKey, deleteAPIKey, setDefaultAPIKey } from "../utils/api";

const Profile = () => {
  const { user } = useAuth();
  const [profile, setProfile] = useState(null);
  const [loading, setLoading] = useState(true);
  const [editMode, setEditMode] = useState(false);
  const [name, setName] = useState("");
  const [apiKeys, setApiKeys] = useState([]);
  const [keyDialogOpen, setKeyDialogOpen] = useState(false);
  const [newKeyData, setNewKeyData] = useState({
    provider: "gemini",
    apiKey: "",
    name: "",
    isDefault: false
  });

  // Load profile data
  useEffect(() => {
    const loadProfile = async () => {
      try {
        setLoading(true);
        const profileData = await getUserProfile();
        setProfile(profileData);
        setName(profileData.name);
        
        // Load API keys
        const keys = await getAPIKeys();
        setApiKeys(keys);
      } catch (error) {
        console.error("Error loading profile:", error);
      } finally {
        setLoading(false);
      }
    };
    
    loadProfile();
  }, []);

  const handleSaveProfile = async () => {
    try {
      setLoading(true);
      const updatedProfile = await updateUserProfile({ name });
      setProfile(updatedProfile);
      setEditMode(false);
    } catch (error) {
      console.error("Error updating profile:", error);
    } finally {
      setLoading(false);
    }
  };

  const handleSaveAPIKey = async () => {
    try {
      setLoading(true);
      const newKey = await saveAPIKey(newKeyData);
      setApiKeys(prev => [...prev, newKey]);
      setKeyDialogOpen(false);
      setNewKeyData({
        provider: "gemini",
        apiKey: "",
        name: "",
        isDefault: false
      });
    } catch (error) {
      console.error("Error saving API key:", error);
    } finally {
      setLoading(false);
    }
  };

  const handleDeleteKey = async (keyId) => {
    try {
      await deleteAPIKey(keyId);
      setApiKeys(prev => prev.filter(key => key.id !== keyId));
    } catch (error) {
      console.error("Error deleting API key:", error);
    }
  };

  const handleSetDefaultKey = async (keyId) => {
    try {
      await setDefaultAPIKey(keyId);
      setApiKeys(prev => prev.map(key => ({
        ...key,
        isDefault: key.id === keyId
      })));
    } catch (error) {
      console.error("Error setting default API key:", error);
    }
  };

  if (loading && !profile) {
    return (
      <Box sx={{ display: "flex", justifyContent: "center", alignItems: "center", height: "100vh" }}>
        <Typography>Loading profile...</Typography>
      </Box>
    );
  }

  return (
    <Container maxWidth="lg" sx={{ py: 5 }}>
      <Grid container spacing={3}>
        {/* Profile Header */}
        <Grid item xs={12}>
          <Paper
            sx={{
              p: 3,
              display: "flex",
              flexDirection: { xs: "column", md: "row" },
              alignItems: { xs: "center", md: "flex-start" },
              gap: 3,
              borderRadius: 2,
              boxShadow: 2
            }}
          >
            <Avatar
              src={profile?.profilePicUrl}
              sx={{ width: 120, height: 120, bgcolor: "primary.main" }}
            >
              {profile?.name?.[0] || <PersonIcon fontSize="large" />}
            </Avatar>
            
            <Box sx={{ flex: 1 }}>
              {editMode ? (
                <Box sx={{ mb: 2 }}>
                  <TextField
                    label="Name"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    fullWidth
                    sx={{ mb: 2 }}
                  />
                  <Box sx={{ display: 'flex', gap: 1 }}>
                    <Button
                      variant="contained"
                      startIcon={<SaveIcon />}
                      onClick={handleSaveProfile}
                    >
                      Save
                    </Button>
                    <Button
                      variant="outlined"
                      onClick={() => {
                        setEditMode(false);
                        setName(profile.name);
                      }}
                    >
                      Cancel
                    </Button>
                  </Box>
                </Box>
              ) : (
                <Box sx={{ mb: 2, display: 'flex', justifyContent: 'space-between' }}>
                  <Typography variant="h4" gutterBottom>
                    {profile?.name || "User"}
                  </Typography>
                  <Button
                    variant="outlined"
                    startIcon={<EditIcon />}
                    onClick={() => setEditMode(true)}
                  >
                    Edit Profile
                  </Button>
                </Box>
              )}
              
              <Typography variant="body1" color="text.secondary" gutterBottom>
                {profile?.email || "No email"}
              </Typography>
              
              <Box sx={{ display: 'flex', alignItems: 'center', flexWrap: 'wrap', gap: 1, mt: 1 }}>
                <Chip 
                  label={`Plan: ${profile?.plan || "Free"}`} 
                  color="primary" 
                  variant="outlined" 
                />
                <Chip 
                  label={`Member since ${new Date(profile?.createdAt).toLocaleDateString()}`} 
                  variant="outlined" 
                />
              </Box>
            </Box>
          </Paper>
        </Grid>
        
        {/* Stats Cards */}
        <Grid item xs={12} md={4}>
          <Card sx={{ height: '100%' }}>
            <CardHeader title="Chats" />
            <CardContent>
              <Typography variant="h3" align="center">
                {profile?.stats?.totalChats || 0}
              </Typography>
              <Typography variant="body2" color="text.secondary" align="center">
                Total chat sessions
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        
        <Grid item xs={12} md={4}>
          <Card sx={{ height: '100%' }}>
            <CardHeader title="Queries" />
            <CardContent>
              <Typography variant="h3" align="center">
                {profile?.stats?.totalQueries || 0}
              </Typography>
              <Typography variant="body2" color="text.secondary" align="center">
                Total queries processed
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        
        <Grid item xs={12} md={4}>
          <Card sx={{ height: '100%' }}>
            <CardHeader title="API Credentials" />
            <CardContent>
              <Typography variant="h3" align="center">
                {apiKeys.length}
              </Typography>
              <Typography variant="body2" color="text.secondary" align="center">
                Saved API credentials
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        
        {/* API Keys Management */}
        <Grid item xs={12}>
          <Card>
            <CardHeader 
              title="LLM API Keys" 
              action={
                <Button 
                  variant="contained" 
                  startIcon={<AddIcon />}
                  onClick={() => setKeyDialogOpen(true)}
                >
                  Add Key
                </Button>
              }
            />
            <CardContent>
              <List>
                {apiKeys.length === 0 ? (
                  <ListItem>
                    <ListItemText 
                      primary="No API keys saved"
                      secondary="Add an API key to use with your chats"
                    />
                  </ListItem>
                ) : (
                  apiKeys.map((key) => (
                    <ListItem key={key.id} divider>
                      <ListItemAvatar>
                        <Avatar sx={{ bgcolor: 'primary.main' }}>
                          <KeyIcon />
                        </Avatar>
                      </ListItemAvatar>
                      <ListItemText
                        primary={
                          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                            {key.name}
                            {key.isDefault && (
                              <Chip
                                label="Default"
                                color="primary"
                                size="small"
                              />
                            )}
                          </Box>
                        }
                        secondary={
                          <Box>
                            <Typography variant="body2">
                              Provider: {key.provider}
                            </Typography>
                            <Typography variant="body2">
                              Key: {key.maskedKey}
                            </Typography>
                            <Typography variant="caption" display="block">
                              Last used: {new Date(key.lastUsed).toLocaleDateString()}
                            </Typography>
                          </Box>
                        }
                      />
                      <ListItemSecondaryAction>
                        {!key.isDefault && (
                          <IconButton edge="end" onClick={() => handleSetDefaultKey(key.id)}>
                            <CheckIcon />
                          </IconButton>
                        )}
                        <IconButton edge="end" onClick={() => handleDeleteKey(key.id)}>
                          <DeleteIcon />
                        </IconButton>
                      </ListItemSecondaryAction>
                    </ListItem>
                  ))
                )}
              </List>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
      
      {/* API Key Dialog */}
      <Dialog open={keyDialogOpen} onClose={() => setKeyDialogOpen(false)}>
        <DialogTitle>Add API Key</DialogTitle>
        <DialogContent>
          <Box sx={{ mt: 2 }}>
            <TextField
              label="Key Name"
              value={newKeyData.name}
              onChange={(e) => setNewKeyData({ ...newKeyData, name: e.target.value })}
              fullWidth
              sx={{ mb: 2 }}
            />
            
            <FormControl fullWidth sx={{ mb: 2 }}>
              <InputLabel>Provider</InputLabel>
              <Select
                value={newKeyData.provider}
                onChange={(e) => setNewKeyData({ ...newKeyData, provider: e.target.value })}
                label="Provider"
              >
                <MenuItem value="gemini">Google Gemini</MenuItem>
                <MenuItem value="openai">OpenAI</MenuItem>
                <MenuItem value="anthropic">Anthropic Claude</MenuItem>
                <MenuItem value="mistral">Mistral</MenuItem>
              </Select>
            </FormControl>
            
            <TextField
              label="API Key"
              value={newKeyData.apiKey}
              onChange={(e) => setNewKeyData({ ...newKeyData, apiKey: e.target.value })}
              fullWidth
              type="password"
              sx={{ mb: 2 }}
            />
            
            <FormControlLabel
              control={
                <Switch
                  checked={newKeyData.isDefault}
                  onChange={(e) => setNewKeyData({ ...newKeyData, isDefault: e.target.checked })}
                />
              }
              label="Set as default"
            />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setKeyDialogOpen(false)}>Cancel</Button>
          <Button 
            onClick={handleSaveAPIKey}
            variant="contained"
            disabled={!newKeyData.apiKey || !newKeyData.name}
          >
            Save
          </Button>
        </DialogActions>
      </Dialog>
    </Container>
  );
};

export default Profile;
import React, { useState, useEffect } from 'react';
import { 
  Box, 
  Typography, 
  List, 
  ListItem, 
  ListItemIcon, 
  ListItemText, 
  Collapse, 
  IconButton, 
  Badge,
  Button, 
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  CircularProgress,
  Alert,
  Divider,
  Card,
  CardContent,
  CardActions,
  Chip
} from '@mui/material';
import {
  ExpandLess,
  ExpandMore,
  Code as CodeIcon,
  Description as DescriptionIcon,
  DataObject as DataObjectIcon,
  Add as AddIcon,
  Refresh as RefreshIcon,
  Upload as UploadIcon,
  Delete as DeleteIcon,
  Visibility as VisibilityIcon,
  Close as CloseIcon
} from '@mui/icons-material';
import { fetchTrainingData, uploadTrainingFile, addTrainingData, deleteTrainingData, viewTrainingData } from '../../utils/api';

const TrainingDataSection = () => {
  const [trainingOpen, setTrainingOpen] = useState(false);
  const [trainingData, setTrainingData] = useState([]);
  const [loading, setLoading] = useState(false);
  const [uploadDialogOpen, setUploadDialogOpen] = useState(false);
  const [addDialogOpen, setAddDialogOpen] = useState(false);
  const [viewDialogOpen, setViewDialogOpen] = useState(false);
  const [deleteConfirmOpen, setDeleteConfirmOpen] = useState(false);
  
  // View state
  const [viewItem, setViewItem] = useState(null);
  const [deleteItem, setDeleteItem] = useState(null);
  
  // Form states
  const [formType, setFormType] = useState('ddl');
  const [formTitle, setFormTitle] = useState('');
  const [formContent, setFormContent] = useState('');
  const [formFile, setFormFile] = useState(null);
  const [formError, setFormError] = useState('');
  const [formSuccess, setFormSuccess] = useState('');

  // Tracking deleted items
  const [deletedItemIds, setDeletedItemIds] = useState(new Set());
  
  useEffect(() => {
    if (trainingOpen) {
      loadTrainingData();
    }
  }, [trainingOpen]);
  
  const loadTrainingData = async () => {
    setLoading(true);
    try {
      console.log("Fetching training data...");
      
      const data = await fetchTrainingData();
      console.log("Training data received:", data);
      
      // Filter out mock example data and previously deleted items
      const realData = Array.isArray(data) 
        ? data.filter(item => 
            !item.id?.includes('example') && 
            !deletedItemIds.has(item.id)
          ) 
        : [];
      
      setTrainingData(realData);
    } catch (error) {
      console.error("Error fetching training data:", error);
    } finally {
      setLoading(false);
    }
  };
  
  const toggleTrainingSection = () => {
    setTrainingOpen(!trainingOpen);
  };
  
  const handleOpenUploadDialog = () => {
    setUploadDialogOpen(true);
    setFormError('');
    setFormSuccess('');
    setFormFile(null);
    setFormType('ddl');
    setFormTitle('');
  };
  
  const handleCloseUploadDialog = () => {
    setUploadDialogOpen(false);
  };
  
  const handleOpenAddDialog = () => {
    setAddDialogOpen(true);
    setFormError('');
    setFormSuccess('');
    setFormContent('');
    setFormType('ddl');
    setFormTitle('');
  };
  
  const handleCloseAddDialog = () => {
    setAddDialogOpen(false);
  };
  
  const handleViewItem = async (item) => {
    setLoading(true);
    try {
      // If item already has content, use that
      if (item.content) {
        setViewItem(item);
        setViewDialogOpen(true);
      } else {
        // Otherwise fetch the full item
        const fullItem = await viewTrainingData(item.id);
        setViewItem(fullItem);
        setViewDialogOpen(true);
      }
    } catch (error) {
      console.error("Error viewing training data:", error);
    } finally {
      setLoading(false);
    }
  };
  
  const handleCloseViewDialog = () => {
    setViewDialogOpen(false);
    setViewItem(null);
  };
  
  const handleDeleteConfirmation = (item) => {
    setDeleteItem(item);
    setDeleteConfirmOpen(true);
  };
  
  const handleDeleteItem = async () => {
    if (!deleteItem) return;
    
    setLoading(true);
    
    try {
      // Track the item to be deleted
      const itemToDelete = {...deleteItem};
      
      // Add the deleted item's ID to our tracking set
      setDeletedItemIds(prev => {
        const newSet = new Set(prev);
        newSet.add(deleteItem.id);
        return newSet;
      });
      
      // Immediately update UI to remove the item
      setTrainingData(prev => prev.filter(item => item.id !== deleteItem.id));
      
      // Close the dialog immediately
      setDeleteConfirmOpen(false);
      setDeleteItem(null);
      
      // Then send the delete request to the server
      console.log(`Attempting to delete training item ${itemToDelete.id}`);
      const result = await deleteTrainingData(itemToDelete.id);
      console.log("Delete result:", result);
      
      if (!result.success) {
        console.error(`Server failed to delete training item: ${itemToDelete.id}`, result);
        // You may want to show an error message here, but keep the item removed from UI
      }
    } catch (error) {
      console.error("Error deleting training data:", error);
      // Optionally show an error message to the user
    } finally {
      setLoading(false);
    }
  };

  const handleUploadTrainingFile = async () => {
    if (!formFile) {
      setFormError('Please select a file to upload');
      return;
    }
    
    try {
      setLoading(true);
      
      const formData = new FormData();
      formData.append('file', formFile);
      formData.append('type', formType);
      formData.append('description', formTitle || formFile.name);
      
      const result = await uploadTrainingFile(formData);
      
      if (result.success) {
        setFormSuccess('File uploaded successfully!');
        
        // Reload training data
        await loadTrainingData();
        
        // Close dialog after a short delay
        setTimeout(() => {
          handleCloseUploadDialog();
        }, 1500);
      } else {
        setFormError('Upload failed: ' + (result.error || 'Unknown error'));
      }
    } catch (error) {
      setFormError(`Upload failed: ${error.message || 'Unknown error'}`);
      console.error('Error uploading training file:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleAddTrainingData = async () => {
    if (!formTitle) {
      setFormError('Please provide a title');
      return;
    }
    
    if (!formContent) {
      setFormError('Please provide content');
      return;
    }
    
    try {
      setLoading(true);
      
      await addTrainingData({
        type: formType,
        content: formContent,
        description: formTitle,
        metadata: {
          title: formTitle,
          description: formTitle
        }
      });
      
      setFormSuccess('Data added successfully!');
      
      // Reload training data
      await loadTrainingData();
      
      // Close dialog after a short delay
      setTimeout(() => {
        handleCloseAddDialog();
      }, 1500);
    } catch (error) {
      setFormError(`Failed to add data: ${error.message || 'Unknown error'}`);
      console.error('Error adding training data:', error);
    } finally {
      setLoading(false);
    }
  };
  
  // Count training data by type
  const getTrainingDataCounts = () => {
    const counts = {
      ddl: 0,
      documentation: 0,
      question_sql: 0,
      total: 0
    };
    
    if (!trainingData || !Array.isArray(trainingData)) {
      return counts;
    }
    
    // Count actual items
    trainingData.forEach(item => {
      counts.total++;
      if (item.type === 'ddl') counts.ddl++;
      else if (item.type === 'documentation') counts.documentation++;
      else if (item.type === 'question_sql') counts.question_sql++;
    });
    
    return counts;
  };
    
  const counts = getTrainingDataCounts();
  
  return (
    <>
      <ListItem button onClick={toggleTrainingSection}>
        <ListItemIcon>
          <DataObjectIcon />
        </ListItemIcon>
        <ListItemText primary="Training Data" />
        <Badge badgeContent={counts.total} color="primary" sx={{ mr: 1 }} />
        {trainingOpen ? <ExpandLess /> : <ExpandMore />}
      </ListItem>
      
      <Collapse in={trainingOpen} timeout="auto" unmountOnExit>
        <Box sx={{ bgcolor: 'rgba(0, 0, 0, 0.15)', pb: 1 }}>
          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', px: 2, pt: 1 }}>
            <Typography variant="body2" color="text.secondary">Knowledge Base</Typography>
            <Box>
              <IconButton size="small" onClick={loadTrainingData} disabled={loading}>
                <RefreshIcon fontSize="small" />
              </IconButton>
              <IconButton size="small" onClick={handleOpenUploadDialog}>
                <UploadIcon fontSize="small" />
              </IconButton>
              <IconButton size="small" onClick={handleOpenAddDialog}>
                <AddIcon fontSize="small" />
              </IconButton>
            </Box>
          </Box>
          
          {/* Training Data Items */}
          <Box sx={{ maxHeight: '200px', overflowY: 'auto', mt: 1, px: 2 }}>
            {loading ? (
              <Box sx={{ display: 'flex', justifyContent: 'center', my: 2 }}>
                <CircularProgress size={24} />
              </Box>
            ) : trainingData.length === 0 ? (
              <Typography variant="caption" sx={{ display: 'block', textAlign: 'center', py: 2 }}>
                No training data added yet
              </Typography>
            ) : (
              trainingData.map((item) => (
                <Card key={item.id} sx={{ mb: 1, bgcolor: 'rgba(0,0,0,0.2)' }}>
                  <CardContent sx={{ py: 1, px: 2, '&:last-child': { pb: 1 } }}>
                    <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                      <Box sx={{ display: 'flex', alignItems: 'center' }}>
                        {item.type === 'ddl' ? (
                          <CodeIcon fontSize="small" sx={{ mr: 1, color: 'primary.light' }} />
                        ) : item.type === 'documentation' ? (
                          <DescriptionIcon fontSize="small" sx={{ mr: 1, color: 'secondary.light' }} />
                        ) : (
                          <DataObjectIcon fontSize="small" sx={{ mr: 1, color: 'success.light' }} />
                        )}
                        <Typography variant="body2" sx={{ fontWeight: 'medium' }}>
                          {item.description || "Unnamed item"}
                        </Typography>
                      </Box>
                      <Chip 
                        label={item.type} 
                        size="small" 
                        sx={{ 
                          height: 20, 
                          fontSize: '0.625rem',
                          bgcolor: item.type === 'ddl' ? 'primary.dark' : 
                                  item.type === 'documentation' ? 'secondary.dark' : 'success.dark'
                        }} 
                      />
                    </Box>
                    <Typography variant="caption" color="text.secondary" sx={{ display: 'block', mt: 0.5 }}>
                      Added: {new Date(item.date_added).toLocaleDateString()}
                    </Typography>
                  </CardContent>
                  <CardActions sx={{ pt: 0, pb: 1, px: 1, justifyContent: 'flex-end' }}>
                    <IconButton size="small" onClick={() => handleViewItem(item)}>
                      <VisibilityIcon fontSize="small" />
                    </IconButton>
                    <IconButton size="small" onClick={() => handleDeleteConfirmation(item)}>
                      <DeleteIcon fontSize="small" />
                    </IconButton>
                  </CardActions>
                </Card>
              ))
            )}
          </Box>
          
          <Divider sx={{ my: 1 }} />
          
          <List dense sx={{ mb: 2 }}> {/* Added margin bottom to fix overlap */}
            <ListItem>
              <ListItemIcon sx={{ minWidth: 36 }}>
                <CodeIcon fontSize="small" />
              </ListItemIcon>
              <ListItemText 
                primary="DDL Schemas" 
                secondary={`${counts.ddl} items`}
                primaryTypographyProps={{ variant: 'body2' }}
                secondaryTypographyProps={{ variant: 'caption' }}
              />
            </ListItem>
            <ListItem>
              <ListItemIcon sx={{ minWidth: 36 }}>
                <DescriptionIcon fontSize="small" />
              </ListItemIcon>
              <ListItemText 
                primary="Documentation" 
                secondary={`${counts.documentation} items`}
                primaryTypographyProps={{ variant: 'body2' }}
                secondaryTypographyProps={{ variant: 'caption' }}
              />
            </ListItem>
            <ListItem>
              <ListItemIcon sx={{ minWidth: 36 }}>
                <DataObjectIcon fontSize="small" />
              </ListItemIcon>
              <ListItemText 
                primary="Question-SQL Pairs" 
                secondary={`${counts.question_sql} items`}
                primaryTypographyProps={{ variant: 'body2' }}
                secondaryTypographyProps={{ variant: 'caption' }}
              />
            </ListItem>
          </List>
        </Box>
      </Collapse>
      
      {/* Upload Training Data Dialog */}
      <Dialog open={uploadDialogOpen} onClose={handleCloseUploadDialog} maxWidth="sm" fullWidth>
        <DialogTitle>Upload Training Data</DialogTitle>
        <DialogContent>
          {formError && <Alert severity="error" sx={{ mb: 2 }}>{formError}</Alert>}
          {formSuccess && <Alert severity="success" sx={{ mb: 2 }}>{formSuccess}</Alert>}
          
          <FormControl fullWidth sx={{ mt: 2, mb: 2 }}>
            <InputLabel>Data Type</InputLabel>
            <Select
              value={formType}
              onChange={(e) => setFormType(e.target.value)}
              label="Data Type"
            >
              <MenuItem value="ddl">DDL Schema</MenuItem>
              <MenuItem value="documentation">Documentation</MenuItem>
              <MenuItem value="question_sql_json">Question-SQL Pairs (JSON)</MenuItem>
              <MenuItem value="auto">Auto-detect</MenuItem>
            </Select>
          </FormControl>
          
          <TextField
            fullWidth
            margin="dense"
            label="Title/Description"
            value={formTitle}
            onChange={(e) => setFormTitle(e.target.value)}
            sx={{ mb: 2 }}
          />
          
          <Button
            variant="outlined"
            component="label"
            startIcon={<UploadIcon />}
            fullWidth
          >
            {formFile ? formFile.name : "Select File"}
            <input
              type="file"
              hidden
              onChange={(e) => setFormFile(e.target.files[0])}
            />
          </Button>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseUploadDialog}>Cancel</Button>
          <Button 
            onClick={handleUploadTrainingFile} 
            variant="contained" 
            disabled={!formFile || loading || formSuccess}
          >
            {loading ? <CircularProgress size={24} /> : "Upload"}
          </Button>
        </DialogActions>
      </Dialog>
      
      {/* Add Training Data Dialog */}
      <Dialog open={addDialogOpen} onClose={handleCloseAddDialog} maxWidth="md" fullWidth>
        <DialogTitle>Add Training Data</DialogTitle>
        <DialogContent>
          {formError && <Alert severity="error" sx={{ mb: 2 }}>{formError}</Alert>}
          {formSuccess && <Alert severity="success" sx={{ mb: 2 }}>{formSuccess}</Alert>}
          
          <FormControl fullWidth sx={{ mt: 2, mb: 2 }}>
            <InputLabel>Data Type</InputLabel>
            <Select
              value={formType}
              onChange={(e) => setFormType(e.target.value)}
              label="Data Type"
            >
              <MenuItem value="ddl">DDL Schema</MenuItem>
              <MenuItem value="documentation">Documentation</MenuItem>
              <MenuItem value="question_sql">Question-SQL Pair</MenuItem>
            </Select>
          </FormControl>
          
          <TextField
            fullWidth
            margin="dense"
            label="Title/Description"
            value={formTitle}
            onChange={(e) => setFormTitle(e.target.value)}
            sx={{ mb: 2 }}
          />
          
          <TextField
            fullWidth
            margin="dense"
            label="Content"
            value={formContent}
            onChange={(e) => setFormContent(e.target.value)}
            multiline
            rows={10}
            sx={{ mb: 2 }}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseAddDialog}>Cancel</Button>
          <Button 
            onClick={handleAddTrainingData} 
            variant="contained" 
            disabled={!formTitle || !formContent || loading || formSuccess}
          >
            {loading ? <CircularProgress size={24} /> : "Add"}
          </Button>
        </DialogActions>
      </Dialog>
      
      {/* View Training Data Dialog */}
      <Dialog open={viewDialogOpen} onClose={handleCloseViewDialog} maxWidth="md" fullWidth>
        <DialogTitle>
          <Box sx={{ display: 'flex', alignItems: 'center' }}>
            {viewItem?.type === 'ddl' ? (
              <CodeIcon sx={{ mr: 1 }} />
            ) : viewItem?.type === 'documentation' ? (
              <DescriptionIcon sx={{ mr: 1 }} />
            ) : (
              <DataObjectIcon sx={{ mr: 1 }} />
            )}
            {viewItem?.description || "Training Data"}
          </Box>
        </DialogTitle>
        <DialogContent>
          <Box sx={{ mb: 2 }}>
            <Chip 
              label={viewItem?.type} 
              size="small" 
              sx={{ mr: 1 }} 
            />
            <Typography variant="caption" color="text.secondary">
              Added: {viewItem?.date_added ? new Date(viewItem.date_added).toLocaleString() : "Unknown"}
            </Typography>
          </Box>
          
          <TextField
            fullWidth
            multiline
            rows={15}
            value={viewItem?.content || ""}
            InputProps={{
              readOnly: true,
            }}
            sx={{ 
              fontFamily: viewItem?.type === 'ddl' ? 'monospace' : 'inherit',
            }}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseViewDialog}>Close</Button>
        </DialogActions>
      </Dialog>
      
      {/* Delete Confirmation Dialog */}
      <Dialog open={deleteConfirmOpen} onClose={() => setDeleteConfirmOpen(false)}>
        <DialogTitle>Confirm Delete</DialogTitle>
        <DialogContent>
          <Typography>
            Are you sure you want to delete "{deleteItem?.description}"?
          </Typography>
          <Typography variant="caption" color="error">
            This action cannot be undone.
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDeleteConfirmOpen(false)}>Cancel</Button>
          <Button 
            onClick={handleDeleteItem} 
            color="error" 
            variant="contained" 
            disabled={loading}
          >
            {loading ? <CircularProgress size={24} /> : "Delete"}
          </Button>
        </DialogActions>
      </Dialog>
    </>
  );
};

export default TrainingDataSection;
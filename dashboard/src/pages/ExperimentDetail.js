import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  Typography,
  Paper,
  Grid,
  Box,
  Chip,
  Button,
  CircularProgress,
  Alert,
  Card,
  CardContent,
  Divider
} from '@mui/material';
import { ArrowBack as ArrowBackIcon, Delete as DeleteIcon } from '@mui/icons-material';
import api from '../services/api';

const statusColors = {
  Pending: 'warning',
  Running: 'info',
  Completed: 'success',
  Failed: 'error'
};

const ExperimentDetail = () => {
  const { namespace, name } = useParams();
  const navigate = useNavigate();
  const [experiment, setExperiment] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const loadExperiment = async () => {
    try {
      setLoading(true);
      const experiment = await api.getExperiment(namespace, name);
      setExperiment(experiment);
      setError(null);
    } catch (err) {
      setError('Failed to load experiment details. Please try again later.');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadExperiment();
    // Set up polling to refresh data every 5 seconds
    const interval = setInterval(loadExperiment, 5000);
    return () => clearInterval(interval);
  }, [namespace, name]);

  const handleDelete = async () => {
    if (window.confirm(`Are you sure you want to delete experiment ${name}?`)) {
      try {
        await api.deleteExperiment(namespace, name);
        navigate('/');
      } catch (err) {
        setError(`Failed to delete experiment ${name}. Please try again later.`);
      }
    }
  };

  if (loading && !experiment) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', mt: 4 }}>
        <CircularProgress />
      </Box>
    );
  }

  if (error && !experiment) {
    return (
      <Alert severity="error" sx={{ mt: 2 }}>
        {error}
      </Alert>
    );
  }

  if (!experiment) {
    return (
      <Alert severity="warning" sx={{ mt: 2 }}>
        Experiment not found.
      </Alert>
    );
  }

  return (
    <div>
      <Box sx={{ display: 'flex', alignItems: 'center', mb: 3 }}>
        <Button
          startIcon={<ArrowBackIcon />}
          onClick={() => navigate('/')}
          sx={{ mr: 2 }}
        >
          Back
        </Button>
        <Typography variant="h4" component="h1">
          Experiment: {experiment.name}
        </Typography>
      </Box>

      {error && (
        <Alert severity="error" sx={{ mb: 3 }}>
          {error}
        </Alert>
      )}

      <Grid container spacing={3}>
        <Grid item xs={12} md={8}>
          <Paper sx={{ p: 3, mb: 3 }}>
            <Typography variant="h6" gutterBottom>
              Overview
            </Typography>
            <Grid container spacing={2}>
              <Grid item xs={12} sm={6}>
                <Typography variant="subtitle2" color="text.secondary">
                  Status
                </Typography>
                <Chip 
                  label={experiment.status || 'Pending'} 
                  color={statusColors[experiment.status] || 'default'} 
                  sx={{ mt: 1 }}
                />
              </Grid>
              <Grid item xs={12} sm={6}>
                <Typography variant="subtitle2" color="text.secondary">
                  Experiment Type
                </Typography>
                <Typography variant="body1">
                  {experiment.experimentType}
                </Typography>
              </Grid>
              <Grid item xs={12} sm={6}>
                <Typography variant="subtitle2" color="text.secondary">
                  Namespace
                </Typography>
                <Typography variant="body1">
                  {experiment.namespace}
                </Typography>
              </Grid>
              <Grid item xs={12} sm={6}>
                <Typography variant="subtitle2" color="text.secondary">
                  Created At
                </Typography>
                <Typography variant="body1">
                  {experiment.startTime ? new Date(experiment.startTime).toLocaleString() : 'Not started'}
                </Typography>
              </Grid>
              <Grid item xs={12} sm={6}>
                <Typography variant="subtitle2" color="text.secondary">
                  Completed At
                </Typography>
                <Typography variant="body1">
                  {experiment.endTime ? new Date(experiment.endTime).toLocaleString() : 'Not completed'}
                </Typography>
              </Grid>
              <Grid item xs={12}>
                <Typography variant="subtitle2" color="text.secondary">
                  Message
                </Typography>
                <Typography variant="body1">
                  {experiment.message || 'No message'}
                </Typography>
              </Grid>
            </Grid>
          </Paper>

          <Paper sx={{ p: 3 }}>
            <Typography variant="h6" gutterBottom>
              Timeline
            </Typography>
            <Box sx={{ position: 'relative', ml: 2, pt: 2, pb: 2 }}>
              {/* Vertical line */}
              <Box sx={{ position: 'absolute', left: 0, top: 0, bottom: 0, width: 2, bgcolor: 'divider' }} />
              
              {/* Timeline events */}
              <Box sx={{ position: 'relative', mb: 3 }}>
                <Box sx={{ position: 'absolute', left: -6, top: 0, width: 10, height: 10, borderRadius: '50%', bgcolor: 'primary.main' }} />
                <Box sx={{ ml: 3 }}>
                  <Typography variant="subtitle2">
                    Experiment Created
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    {experiment.startTime ? new Date(experiment.startTime).toLocaleString() : 'Not started'}
                  </Typography>
                </Box>
              </Box>

              {experiment.status === 'Running' && (
                <Box sx={{ position: 'relative', mb: 3 }}>
                  <Box sx={{ position: 'absolute', left: -6, top: 0, width: 10, height: 10, borderRadius: '50%', bgcolor: 'info.main' }} />
                  <Box sx={{ ml: 3 }}>
                    <Typography variant="subtitle2">
                      Experiment Running
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      {experiment.startTime ? new Date(experiment.startTime).toLocaleString() : 'Not started'}
                    </Typography>
                  </Box>
                </Box>
              )}

              {(experiment.status === 'Completed' || experiment.status === 'Failed') && (
                <Box sx={{ position: 'relative' }}>
                  <Box sx={{ 
                    position: 'absolute', 
                    left: -6, 
                    top: 0, 
                    width: 10, 
                    height: 10, 
                    borderRadius: '50%', 
                    bgcolor: experiment.status === 'Completed' ? 'success.main' : 'error.main' 
                  }} />
                  <Box sx={{ ml: 3 }}>
                    <Typography variant="subtitle2">
                      Experiment {experiment.status}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      {experiment.endTime ? new Date(experiment.endTime).toLocaleString() : 'Not completed'}
                    </Typography>
                  </Box>
                </Box>
              )}
            </Box>
          </Paper>
        </Grid>

        <Grid item xs={12} md={4}>
          <Card sx={{ mb: 3 }}>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Target
              </Typography>
              <Divider sx={{ mb: 2 }} />
              <Typography variant="subtitle2" color="text.secondary">
                Kind
              </Typography>
              <Typography variant="body1" sx={{ mb: 2 }}>
                {experiment.targetKind || 'Pod'}
              </Typography>
              <Typography variant="subtitle2" color="text.secondary">
                Name
              </Typography>
              <Typography variant="body1">
                {experiment.targetName}
              </Typography>
            </CardContent>
          </Card>

          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Actions
              </Typography>
              <Divider sx={{ mb: 2 }} />
              <Button
                variant="outlined"
                color="error"
                startIcon={<DeleteIcon />}
                onClick={handleDelete}
                fullWidth
              >
                Delete Experiment
              </Button>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </div>
  );
};

export default ExperimentDetail;

import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Typography,
  Paper,
  TextField,
  Button,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Grid,
  Box,
  Alert,
  CircularProgress
} from '@mui/material';
import api from '../services/api';

const experimentTypes = [
  { value: 'pod-failure', label: 'Pod Failure', description: 'Kills a pod to test resilience to pod failures' },
  { value: 'network-latency', label: 'Network Latency', description: 'Adds latency to network traffic' },
  { value: 'cpu-hog', label: 'CPU Hog', description: 'Consumes CPU resources' },
  { value: 'memory-hog', label: 'Memory Hog', description: 'Consumes memory resources' }
];

const targetKinds = [
  { value: 'Pod', label: 'Pod' },
  { value: 'Deployment', label: 'Deployment' },
  { value: 'StatefulSet', label: 'StatefulSet' },
  { value: 'Service', label: 'Service' }
];

const durations = [
  { value: '30s', label: '30 seconds' },
  { value: '1m', label: '1 minute' },
  { value: '5m', label: '5 minutes' },
  { value: '10m', label: '10 minutes' },
  { value: '30m', label: '30 minutes' },
  { value: '1h', label: '1 hour' }
];

const NewExperiment = () => {
  const navigate = useNavigate();
  const [formData, setFormData] = useState({
    name: '',
    namespace: 'default',
    targetName: '',
    targetKind: 'Pod',
    experimentType: 'pod-failure',
    duration: '1m',
    parameters: {}
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [parameters, setParameters] = useState({});

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData({
      ...formData,
      [name]: value
    });

    // Reset parameters when experiment type changes
    if (name === 'experimentType') {
      setParameters({});
    }
  };

  const handleParameterChange = (e) => {
    const { name, value } = e.target;
    setParameters({
      ...parameters,
      [name]: value
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      const experimentData = {
        ...formData,
        parameters
      };
      
      await api.createExperiment(experimentData);
      navigate('/');
    } catch (err) {
      setError('Failed to create experiment. Please check your inputs and try again.');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  // Render parameter inputs based on experiment type
  const renderParameterInputs = () => {
    switch (formData.experimentType) {
      case 'network-latency':
        return (
          <TextField
            fullWidth
            label="Latency (e.g., 100ms)"
            name="latency"
            value={parameters.latency || '100ms'}
            onChange={handleParameterChange}
            margin="normal"
            helperText="Amount of latency to add to network requests"
          />
        );
      case 'cpu-hog':
        return (
          <TextField
            fullWidth
            label="CPU Cores"
            name="cpuCores"
            type="number"
            value={parameters.cpuCores || '1'}
            onChange={handleParameterChange}
            margin="normal"
            helperText="Number of CPU cores to consume"
          />
        );
      case 'memory-hog':
        return (
          <TextField
            fullWidth
            label="Memory (MB)"
            name="memoryMB"
            type="number"
            value={parameters.memoryMB || '256'}
            onChange={handleParameterChange}
            margin="normal"
            helperText="Amount of memory to consume in MB"
          />
        );
      default:
        return null;
    }
  };

  return (
    <div>
      <Typography variant="h4" component="h1" gutterBottom>
        Create New Chaos Experiment
      </Typography>

      {error && (
        <Alert severity="error" sx={{ mb: 3 }}>
          {error}
        </Alert>
      )}

      <Paper sx={{ p: 3 }}>
        <form onSubmit={handleSubmit}>
          <Grid container spacing={3}>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="Experiment Name"
                name="name"
                value={formData.name}
                onChange={handleChange}
                required
                helperText="A unique name for your experiment"
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="Namespace"
                name="namespace"
                value={formData.namespace}
                onChange={handleChange}
                required
                helperText="Kubernetes namespace where the experiment will run"
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <FormControl fullWidth>
                <InputLabel>Target Kind</InputLabel>
                <Select
                  name="targetKind"
                  value={formData.targetKind}
                  onChange={handleChange}
                  label="Target Kind"
                  required
                >
                  {targetKinds.map((kind) => (
                    <MenuItem key={kind.value} value={kind.value}>
                      {kind.label}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="Target Name"
                name="targetName"
                value={formData.targetName}
                onChange={handleChange}
                required
                helperText={`Name of the ${formData.targetKind} to target`}
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <FormControl fullWidth>
                <InputLabel>Experiment Type</InputLabel>
                <Select
                  name="experimentType"
                  value={formData.experimentType}
                  onChange={handleChange}
                  label="Experiment Type"
                  required
                >
                  {experimentTypes.map((type) => (
                    <MenuItem key={type.value} value={type.value}>
                      {type.label}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
              <Typography variant="caption" color="text.secondary" sx={{ display: 'block', mt: 1 }}>
                {experimentTypes.find(t => t.value === formData.experimentType)?.description}
              </Typography>
            </Grid>
            <Grid item xs={12} md={6}>
              <FormControl fullWidth>
                <InputLabel>Duration</InputLabel>
                <Select
                  name="duration"
                  value={formData.duration}
                  onChange={handleChange}
                  label="Duration"
                  required
                >
                  {durations.map((duration) => (
                    <MenuItem key={duration.value} value={duration.value}>
                      {duration.label}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12}>
              <Typography variant="h6" gutterBottom>
                Experiment Parameters
              </Typography>
              {renderParameterInputs()}
            </Grid>
            <Grid item xs={12}>
              <Box sx={{ display: 'flex', justifyContent: 'flex-end', mt: 2 }}>
                <Button
                  variant="outlined"
                  onClick={() => navigate('/')}
                  sx={{ mr: 2 }}
                  disabled={loading}
                >
                  Cancel
                </Button>
                <Button
                  type="submit"
                  variant="contained"
                  color="primary"
                  disabled={loading}
                >
                  {loading ? <CircularProgress size={24} /> : 'Create Experiment'}
                </Button>
              </Box>
            </Grid>
          </Grid>
        </form>
      </Paper>
    </div>
  );
};

export default NewExperiment;

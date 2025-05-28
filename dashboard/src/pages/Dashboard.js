import React, { useState, useEffect } from 'react';
import { Link as RouterLink } from 'react-router-dom';
import {
  Typography,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Button,
  Chip,
  Box,
  CircularProgress,
  Alert
} from '@mui/material';
import { Add as AddIcon } from '@mui/icons-material';
import api from '../services/api';

const statusColors = {
  Pending: 'warning',
  Running: 'info',
  Completed: 'success',
  Failed: 'error'
};

const Dashboard = () => {
  const [experiments, setExperiments] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const loadExperiments = async () => {
    try {
      setLoading(true);
      const experiments = await api.getExperiments();
      setExperiments(experiments);
      setError(null);
    } catch (err) {
      setError('Failed to load experiments. Please try again later.');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadExperiments();
    // Set up polling to refresh data every 5 seconds
    const interval = setInterval(loadExperiments, 5000);
    return () => clearInterval(interval);
  }, []);

  const handleDelete = async (namespace, name) => {
    if (window.confirm(`Are you sure you want to delete experiment ${name}?`)) {
      try {
        await api.deleteExperiment(namespace, name);
        loadExperiments();
      } catch (err) {
        setError(`Failed to delete experiment ${name}. Please try again later.`);
      }
    }
  };

  return (
    <div>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          Chaos Experiments
        </Typography>
        <Button 
          variant="contained" 
          color="primary" 
          startIcon={<AddIcon />}
          component={RouterLink}
          to="/experiments/new"
        >
          New Experiment
        </Button>
      </Box>

      {error && (
        <Alert severity="error" sx={{ mb: 3 }}>
          {error}
        </Alert>
      )}

      {loading && experiments.length === 0 ? (
        <Box sx={{ display: 'flex', justifyContent: 'center', mt: 4 }}>
          <CircularProgress />
        </Box>
      ) : (
        <TableContainer component={Paper}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Name</TableCell>
                <TableCell>Namespace</TableCell>
                <TableCell>Type</TableCell>
                <TableCell>Status</TableCell>
                <TableCell>Start Time</TableCell>
                <TableCell>End Time</TableCell>
                <TableCell>Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {experiments.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={7} align="center">
                    No experiments found. Create your first experiment!
                  </TableCell>
                </TableRow>
              ) : (
                experiments.map((experiment) => (
                  <TableRow key={`${experiment.namespace}-${experiment.name}`}>
                    <TableCell>
                      <RouterLink to={`/experiments/${experiment.namespace}/${experiment.name}`} style={{ textDecoration: 'none', color: 'inherit' }}>
                        {experiment.name}
                      </RouterLink>
                    </TableCell>
                    <TableCell>{experiment.namespace}</TableCell>
                    <TableCell>{experiment.experimentType}</TableCell>
                    <TableCell>
                      <Chip 
                        label={experiment.status || 'Pending'} 
                        color={statusColors[experiment.status] || 'default'} 
                        size="small" 
                      />
                    </TableCell>
                    <TableCell>
                      {experiment.startTime ? new Date(experiment.startTime).toLocaleString() : '-'}
                    </TableCell>
                    <TableCell>
                      {experiment.endTime ? new Date(experiment.endTime).toLocaleString() : '-'}
                    </TableCell>
                    <TableCell>
                      <Button
                        variant="outlined"
                        color="error"
                        size="small"
                        onClick={() => handleDelete(experiment.namespace, experiment.name)}
                      >
                        Delete
                      </Button>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </TableContainer>
      )}
    </div>
  );
};

export default Dashboard;

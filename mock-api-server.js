const express = require('express');
const cors = require('cors');
const path = require('path');
const app = express();
const PORT = 5000;

// Sample data
const experiments = [
  {
    name: 'nginx-pod-failure',
    namespace: 'chaos-test',
    experimentType: 'pod-failure',
    status: 'Completed',
    startTime: new Date(Date.now() - 3600000).toISOString(),
    endTime: new Date(Date.now() - 3300000).toISOString(),
    message: 'Experiment completed successfully'
  },
  {
    name: 'nginx-network-latency',
    namespace: 'chaos-test',
    experimentType: 'network-latency',
    status: 'Running',
    startTime: new Date(Date.now() - 1800000).toISOString(),
    endTime: null,
    message: 'Experiment in progress'
  },
  {
    name: 'nginx-cpu-hog',
    namespace: 'chaos-test',
    experimentType: 'cpu-hog',
    status: 'Pending',
    startTime: null,
    endTime: null,
    message: 'Waiting to start'
  }
];

app.use(cors());
app.use(express.json());

// Serve static files from the React app
app.use(express.static(path.join(__dirname, 'dashboard/build')));

// API routes
app.get('/api/experiments', (req, res) => {
  res.json(experiments);
});

app.get('/api/experiments/:namespace/:name', (req, res) => {
  const { namespace, name } = req.params;
  const experiment = experiments.find(e => e.namespace === namespace && e.name === name);
  
  if (experiment) {
    res.json(experiment);
  } else {
    res.status(404).json({ error: 'Experiment not found' });
  }
});

app.post('/api/experiments', (req, res) => {
  const newExperiment = {
    ...req.body,
    status: 'Pending',
    startTime: null,
    endTime: null,
    message: 'Waiting to start'
  };
  
  experiments.push(newExperiment);
  res.status(201).json(newExperiment);
});

app.delete('/api/experiments/:namespace/:name', (req, res) => {
  const { namespace, name } = req.params;
  const index = experiments.findIndex(e => e.namespace === namespace && e.name === name);
  
  if (index !== -1) {
    experiments.splice(index, 1);
    res.status(204).send();
  } else {
    res.status(404).json({ error: 'Experiment not found' });
  }
});

// All other GET requests not handled before will return a 404
app.get('*', (req, res) => {
  res.status(404).json({ error: 'Not found' });
});

app.listen(PORT, () => {
  console.log(`Mock API server running on port ${PORT}`);
  console.log(`Access the dashboard at http://localhost:${PORT}`);
});

const http = require('http');

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

// Create HTTP server
const server = http.createServer((req, res) => {
  // Set CORS headers
  res.setHeader('Access-Control-Allow-Origin', '*');
  res.setHeader('Access-Control-Allow-Methods', 'GET, POST, PUT, DELETE, OPTIONS');
  res.setHeader('Access-Control-Allow-Headers', 'Content-Type');
  
  // Handle OPTIONS requests (for CORS preflight)
  if (req.method === 'OPTIONS') {
    res.writeHead(204);
    res.end();
    return;
  }
  
  // Parse URL
  const url = new URL(req.url, `http://${req.headers.host}`);
  const path = url.pathname;
  
  // Set content type to JSON
  res.setHeader('Content-Type', 'application/json');
  
  // Handle API routes
  if (path === '/api/experiments' && req.method === 'GET') {
    // GET all experiments
    res.writeHead(200);
    res.end(JSON.stringify(experiments));
  } 
  else if (path.match(/^\/api\/experiments\/[^\/]+\/[^\/]+$/) && req.method === 'GET') {
    // GET single experiment by namespace/name
    const parts = path.split('/');
    const namespace = parts[3];
    const name = parts[4];
    
    const experiment = experiments.find(e => e.namespace === namespace && e.name === name);
    
    if (experiment) {
      res.writeHead(200);
      res.end(JSON.stringify(experiment));
    } else {
      res.writeHead(404);
      res.end(JSON.stringify({ error: 'Experiment not found' }));
    }
  }
  else if (path === '/api/experiments' && req.method === 'POST') {
    // POST new experiment
    let body = '';
    
    req.on('data', chunk => {
      body += chunk.toString();
    });
    
    req.on('end', () => {
      try {
        const experimentData = JSON.parse(body);
        const newExperiment = {
          ...experimentData,
          status: 'Pending',
          startTime: null,
          endTime: null,
          message: 'Waiting to start'
        };
        
        experiments.push(newExperiment);
        
        res.writeHead(201);
        res.end(JSON.stringify(newExperiment));
      } catch (error) {
        res.writeHead(400);
        res.end(JSON.stringify({ error: 'Invalid JSON' }));
      }
    });
  }
  else if (path.match(/^\/api\/experiments\/[^\/]+\/[^\/]+$/) && req.method === 'DELETE') {
    // DELETE experiment
    const parts = path.split('/');
    const namespace = parts[3];
    const name = parts[4];
    
    const index = experiments.findIndex(e => e.namespace === namespace && e.name === name);
    
    if (index !== -1) {
      experiments.splice(index, 1);
      res.writeHead(204);
      res.end();
    } else {
      res.writeHead(404);
      res.end(JSON.stringify({ error: 'Experiment not found' }));
    }
  }
  else {
    // Not found
    res.writeHead(404);
    res.end(JSON.stringify({ error: 'Not found' }));
  }
});

const PORT = 5000;

server.listen(PORT, () => {
  console.log(`Mock API server running on port ${PORT}`);
  console.log(`Access the API at http://localhost:${PORT}/api/experiments`);
});

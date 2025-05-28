import React from 'react';
import { Link as RouterLink } from 'react-router-dom';
import { 
  AppBar, 
  Box, 
  Drawer, 
  Toolbar, 
  Typography, 
  List, 
  ListItem, 
  ListItemIcon, 
  ListItemText,
  Link,
  Container
} from '@mui/material';
import { 
  Dashboard as DashboardIcon, 
  Science as ScienceIcon, 
  Add as AddIcon 
} from '@mui/icons-material';

const drawerWidth = 240;

const Layout = ({ children }) => {
  return (
    <Box sx={{ display: 'flex' }}>
      <AppBar position="fixed" sx={{ zIndex: (theme) => theme.zIndex.drawer + 1 }}>
        <Toolbar>
          <ScienceIcon sx={{ mr: 2 }} />
          <Typography variant="h6" noWrap component="div">
            Chaos Engineering as a Service
          </Typography>
        </Toolbar>
      </AppBar>
      <Drawer
        variant="permanent"
        sx={{
          width: drawerWidth,
          flexShrink: 0,
          [`& .MuiDrawer-paper`]: { width: drawerWidth, boxSizing: 'border-box' },
        }}
      >
        <Toolbar />
        <Box sx={{ overflow: 'auto' }}>
          <List>
            <ListItem button component={RouterLink} to="/">
              <ListItemIcon>
                <DashboardIcon />
              </ListItemIcon>
              <ListItemText primary="Dashboard" />
            </ListItem>
            <ListItem button component={RouterLink} to="/experiments/new">
              <ListItemIcon>
                <AddIcon />
              </ListItemIcon>
              <ListItemText primary="New Experiment" />
            </ListItem>
          </List>
        </Box>
      </Drawer>
      <Box component="main" sx={{ flexGrow: 1, p: 3 }}>
        <Toolbar />
        <Container maxWidth="lg">
          {children}
        </Container>
      </Box>
    </Box>
  );
};

export default Layout;

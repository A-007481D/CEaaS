import axios from 'axios';

const API_URL = 'http://localhost:5000/api';

const api = {
  getExperiments: async () => {
    try {
      const response = await axios.get(`${API_URL}/experiments`);
      return response.data;
    } catch (error) {
      console.error('Error fetching experiments:', error);
      throw error;
    }
  },

  getExperiment: async (namespace, name) => {
    try {
      const response = await axios.get(`${API_URL}/experiments/${namespace}/${name}`);
      return response.data;
    } catch (error) {
      console.error(`Error fetching experiment ${namespace}/${name}:`, error);
      throw error;
    }
  },

  createExperiment: async (experimentData) => {
    try {
      const response = await axios.post(`${API_URL}/experiments`, experimentData);
      return response.data;
    } catch (error) {
      console.error('Error creating experiment:', error);
      throw error;
    }
  },

  deleteExperiment: async (namespace, name) => {
    try {
      await axios.delete(`${API_URL}/experiments/${namespace}/${name}`);
      return true;
    } catch (error) {
      console.error(`Error deleting experiment ${namespace}/${name}:`, error);
      throw error;
    }
  }
};

export default api;

import axios from 'axios';
import { toast } from 'sonner';
import { API_BASE_URL } from '@core/config';

const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: { 'Content-Type': 'application/json' },
});

apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    const status = error.response?.status;
    if (status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
      return Promise.reject(error);
    }

    const data = error.response?.data;
    const message = data?.error || error.message || 'Something went wrong';

    if (status === 400) {
      const fields = data?.fields;
      const desc = fields
        ? Object.entries(fields).map(([k, v]) => `${k}: ${v}`).join(', ')
        : message;
      toast.error('Validation error', { description: desc });
    } else if (status === 403) {
      toast.error('Access denied', { description: message });
    } else if (status === 404) {
      toast.error('Not found', { description: message });
    } else if (status === 409) {
      toast.error('Conflict', { description: message });
    } else if (status && status >= 500) {
      toast.error('Server error', { description: 'Please try again later.' });
    }

    return Promise.reject(error);
  }
);

export default apiClient;

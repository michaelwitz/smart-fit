import { jwtDecode } from 'jwt-decode';

const TOKEN_KEY = 'smart_fit_girl_token';

export const authUtils = {
  // Store JWT token
  setToken: (token) => {
    localStorage.setItem(TOKEN_KEY, token);
  },

  // Get JWT token
  getToken: () => {
    return localStorage.getItem(TOKEN_KEY);
  },

  // Remove JWT token
  removeToken: () => {
    localStorage.removeItem(TOKEN_KEY);
  },

  // Check if user is authenticated and token is valid
  isAuthenticated: () => {
    const token = authUtils.getToken();
    if (!token) return false;

    try {
      const decoded = jwtDecode(token);
      const currentTime = Date.now() / 1000;
      
      // Check if token is expired
      if (decoded.exp < currentTime) {
        authUtils.removeToken();
        return false;
      }
      
      return true;
    } catch (error) {
      console.error('Error decoding token:', error);
      authUtils.removeToken();
      return false;
    }
  },

  // Get user info from token
  getUserInfo: () => {
    const token = authUtils.getToken();
    if (!token) return null;

    try {
      const decoded = jwtDecode(token);
      return {
        userId: decoded.sub || decoded.user_id,
        email: decoded.email,
        fullName: decoded.fullName || decoded.name,
        exp: decoded.exp
      };
    } catch (error) {
      console.error('Error decoding token:', error);
      return null;
    }
  },

  // Check if token expires within a certain time (in hours)
  isTokenExpiringSoon: (hoursThreshold = 24) => {
    const token = authUtils.getToken();
    if (!token) return false;

    try {
      const decoded = jwtDecode(token);
      const currentTime = Date.now() / 1000;
      const timeUntilExpiry = decoded.exp - currentTime;
      const hoursUntilExpiry = timeUntilExpiry / 3600;
      
      return hoursUntilExpiry <= hoursThreshold;
    } catch (error) {
      console.error('Error checking token expiry:', error);
      return false;
    }
  }
};

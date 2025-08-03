import React, { useState, useEffect } from 'react';
import {
  AppShell,
  Burger,
  Button,
  Container,
  Title,
  Text,
  Stack,
  Group,
  ActionIcon,
  Menu,
  Avatar,
  Box
} from '@mantine/core';
import { 
  IconUser, 
  IconSettings, 
  IconLogout, 
  IconDashboard,
  IconCalendar,
  IconBarbell
} from '@tabler/icons-react';
import { notifications } from '@mantine/notifications';
import { 
  BrowserView, 
  MobileView, 
  isBrowser, 
  isMobile, 
  isTablet,
  isIOS,
  isAndroid,
  deviceType,
  browserName,
  osName
} from 'react-device-detect';
import AuthModal from './components/AuthModal';
import { authUtils } from './utils/auth';
import { authAPI } from './utils/api';

function App() {
  const [authModalOpened, setAuthModalOpened] = useState(false);
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [userInfo, setUserInfo] = useState(null);
  // Use react-device-detect for better device detection
  const isActualMobile = isMobile || isTablet;
  
  // Log device info for debugging
  useEffect(() => {
    console.log('Device Info:', {
      deviceType,
      isMobile,
      isTablet,
      isIOS,
      isAndroid,
      browserName,
      osName,
      screenWidth: window.screen.width,
      screenHeight: window.screen.height
    });
  }, []);

  // Check authentication status on component mount
  useEffect(() => {
    const checkAuth = () => {
      const authenticated = authUtils.isAuthenticated();
      setIsAuthenticated(authenticated);
      
      if (authenticated) {
        const user = authUtils.getUserInfo();
        setUserInfo(user);
      }
    };

    checkAuth();
    
    // Check token expiration every minute
    const interval = setInterval(checkAuth, 60000);
    
    return () => clearInterval(interval);
  }, []);

  // Mobile address bar hiding functionality
  useEffect(() => {
    const hideAddressBar = () => {
      // For mobile Safari - force scroll to hide address bar
      if (window.navigator.userAgent.includes('Safari') && window.navigator.userAgent.includes('Mobile')) {
        window.scrollTo(0, 1);
        setTimeout(() => window.scrollTo(0, 0), 100);
      }
    };

    // Hide address bar on load
    hideAddressBar();
    
    // Handle resize for responsive design
    const handleResize = () => {
      hideAddressBar();
    };
    
    // Also hide on orientation change
    window.addEventListener('orientationchange', handleResize);
    window.addEventListener('resize', handleResize);
    
    return () => {
      window.removeEventListener('orientationchange', handleResize);
      window.removeEventListener('resize', handleResize);
    };
  }, []);

  const handleAuthSuccess = () => {
    setIsAuthenticated(true);
    const user = authUtils.getUserInfo();
    setUserInfo(user);
  };

  const handleLogout = async () => {
    try {
      await authAPI.logout();
      setIsAuthenticated(false);
      setUserInfo(null);
      notifications.show({
        title: 'Logged out',
        message: 'See you next time!',
        color: 'blue',
      });
    } catch (error) {
      console.error('Logout error:', error);
      // Still clear local state even if API call fails
      authUtils.removeToken();
      setIsAuthenticated(false);
      setUserInfo(null);
    }
  };

  const handleContinue = () => {
    // Navigate to main app - for now just show notification
    notifications.show({
      title: 'Welcome!',
      message: 'Redirecting to your dashboard...',
      color: 'green',
    });
  };

  return (
    <>
      {/* Background Video */}
      <video
        className="video-background"
        autoPlay
        muted
        loop
        playsInline
      >
        <source src="/assets/videos/splash-background.mp4" type="video/mp4" />
        {/* Fallback for browsers that don't support video */}
        <div style={{
          position: 'fixed',
          top: 0,
          left: 0,
          width: '100vw',
          height: '100vh',
          background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
          zIndex: -1
        }} />
      </video>
      
      {/* Video Overlay */}
      <div className="video-overlay" />

      {/* Main App Layout using CSS Grid */}
      <div className="app-layout">
        {/* Header Section */}
        <div className="header-section">
          {/* Hamburger Menu */}
          <div className="hamburger-menu">
            <Menu shadow="md" width={250}>
              <Menu.Target>
                <Box
                  style={{
                    padding: '8px',
                    borderRadius: '8px',
                    backgroundColor: 'rgba(255, 255, 255, 0.2)',
                    cursor: 'pointer',
                    minWidth: '44px',
                    minHeight: '44px',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    border: '1px solid rgba(255, 255, 255, 0.3)'
                  }}
                >
                  <Burger 
                    size={isActualMobile ? 'md' : 'lg'}
                    color="white" 
                  />
                </Box>
              </Menu.Target>
              <Menu.Dropdown>
                <Menu.Label>Navigation</Menu.Label>
                <Menu.Item leftSection={<IconDashboard size={16} />}>
                  Dashboard
                </Menu.Item>
                <Menu.Item leftSection={<IconBarbell size={16} />}>
                  Workouts
                </Menu.Item>
                <Menu.Item leftSection={<IconCalendar size={16} />}>
                  Schedule
                </Menu.Item>
                <Menu.Item leftSection={<IconUser size={16} />}>
                  Profile
                </Menu.Item>
              </Menu.Dropdown>
            </Menu>
          </div>

          {/* Empty div to maintain grid structure */}
          <div></div>

          {/* User Menu */}
          <div className="user-menu">
            {isAuthenticated && userInfo && (
              <Menu shadow="md" width={200}>
                <Menu.Target>
                  <ActionIcon variant="subtle" size="lg">
                    <Avatar size="sm" color="pink">
                      {userInfo.fullName?.charAt(0) || userInfo.email?.charAt(0) || 'U'}
                    </Avatar>
                  </ActionIcon>
                </Menu.Target>
                <Menu.Dropdown>
                  <Menu.Label>Account</Menu.Label>
                  <Menu.Item leftSection={<IconUser size={14} />}>
                    Profile
                  </Menu.Item>
                  <Menu.Item leftSection={<IconSettings size={14} />}>
                    Settings
                  </Menu.Item>
                  <Menu.Divider />
                  <Menu.Item 
                    leftSection={<IconLogout size={14} />} 
                    color="red"
                    onClick={handleLogout}
                  >
                    Logout
                  </Menu.Item>
                </Menu.Dropdown>
              </Menu>
            )}
          </div>
        </div>

        {/* Main Content Section */}
        <div className="main-content">
          {/* This area can be used for main app content later */}
        </div>

        {/* Footer Section */}
        <div className="footer-section">
          {isAuthenticated ? (
            <Button
              size="xl"
              radius="xl"
              style={{
                backgroundColor: 'rgba(255, 255, 255, 0.9)',
                color: '#333',
                border: 'none',
                minWidth: '200px',
                height: '60px',
                fontSize: '1.2rem',
                fontWeight: 600,
                boxShadow: '0 4px 15px rgba(0,0,0,0.2)',
              }}
              onClick={handleContinue}
            >
              Continue
            </Button>
          ) : (
            <Button
              size="xl"
              radius="xl"
              style={{
                backgroundColor: 'rgba(255, 255, 255, 0.9)',
                color: '#333',
                border: 'none',
                minWidth: '200px',
                height: '60px',
                fontSize: '1.2rem',
                fontWeight: 600,
                boxShadow: '0 4px 15px rgba(0,0,0,0.2)',
              }}
              onClick={() => setAuthModalOpened(true)}
            >
              Sign In / Sign Up
            </Button>
          )}
        </div>
      </div>

      {/* Title - horizontally centered, mobile safe */}
      <div 
        style={{
          position: 'absolute',
          top: isActualMobile ? (isIOS ? 'calc(env(safe-area-inset-top) + 80px)' : '80px') : '10%',
          left: '50%',
          transform: 'translate(-50%, -50%)',
          textAlign: 'center',
          zIndex: 1,
          pointerEvents: 'none'
        }}
      >
        <Title 
          order={1} 
          size={isActualMobile ? (isTablet ? "2.8rem" : "2.3rem") : "3rem"}
          c="white" 
          style={{ 
            fontWeight: 900,
            textShadow: '2px 2px 4px rgba(0,0,0,0.7)',
            marginBottom: '0.5rem',
            whiteSpace: 'nowrap'
          }}
        >
          Smart Fit Girl
        </Title>
        
        <Text 
          size={isActualMobile ? (isTablet ? "xl" : "lg") : "xl"}
          c="white" 
          style={{ 
            textShadow: '1px 1px 2px rgba(0,0,0,0.7)'
          }}
        >
          Your personal fitness journey starts here
        </Text>
      </div>

      <AuthModal
        opened={authModalOpened}
        onClose={() => setAuthModalOpened(false)}
        onAuthSuccess={handleAuthSuccess}
      />
    </>
  );
}

export default App;

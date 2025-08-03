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
import AuthModal from './components/AuthModal';
import { authUtils } from './utils/auth';
import { authAPI } from './utils/api';

function App() {
  const [authModalOpened, setAuthModalOpened] = useState(false);
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [userInfo, setUserInfo] = useState(null);
  const [isMobile, setIsMobile] = useState(window.innerWidth <= 768);

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
      setIsMobile(window.innerWidth <= 768);
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

      {/* Hamburger Menu - Inline with Title */}
      <Box className="hamburger-menu">
        <Menu shadow="md" width={250}>
          <Menu.Target>
            <Box
              style={{
                padding: '8px',
                borderRadius: '8px',
                backgroundColor: 'rgba(255, 255, 255, 0.1)',
                cursor: 'pointer',
                minWidth: '44px', // iOS touch target minimum
                minHeight: '44px',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center'
              }}
            >
              <Burger 
              size={isMobile ? 'md' : 'lg'}
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
      </Box>

      {/* User Avatar Menu - Aligned with Title */}
      {isAuthenticated && userInfo && (
        <Box
          style={{
            position: 'fixed',
            top: 'calc(20vh - 20px)', // Align with hamburger menu
            right: '20px',
            zIndex: 10
          }}
        >
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
        </Box>
      )}
          {/* Title Section - Top 20% */}
          <Box
            className="title-section"
            style={{
              position: 'fixed',
              top: 0,
              left: 0,
              right: 0,
              height: '20vh',
              display: 'flex',
              flexDirection: 'column',
              justifyContent: 'center',
              alignItems: 'center',
              textAlign: 'center',
              zIndex: 1
            }}
          >
            <Title 
              order={1} 
              size="3rem"
              c="white" 
              style={{ 
                fontWeight: 900,
                textShadow: '2px 2px 4px rgba(0,0,0,0.7)',
                marginBottom: '1rem'
              }}
            >
              Smart Fit Girl
            </Title>
            
            <Text 
              size="xl" 
              c="white" 
              style={{ 
                textShadow: '1px 1px 2px rgba(0,0,0,0.7)'
              }}
            >
              Your personal fitness journey starts here
            </Text>
          </Box>

          {/* Button Section - Bottom 20% */}
          <Box
            style={{
              position: 'fixed',
              bottom: 0,
              left: 0,
              right: 0,
              height: '20vh',
              display: 'flex',
              justifyContent: 'center',
              alignItems: 'center',
              zIndex: 1
            }}
          >
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
          </Box>

      <AuthModal
        opened={authModalOpened}
        onClose={() => setAuthModalOpened(false)}
        onAuthSuccess={handleAuthSuccess}
      />
    </>
  );
}

export default App;

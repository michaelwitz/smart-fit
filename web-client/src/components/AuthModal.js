import React, { useState } from 'react';
import {
  Modal,
  Paper,
  TextInput,
  PasswordInput,
  Button,
  Title,
  Text,
  Anchor,
  Stack,
  Group,
  Alert,
  LoadingOverlay,
  Tabs
} from '@mantine/core';
import { useForm } from '@mantine/form';
import { notifications } from '@mantine/notifications';
import { IconAlertCircle, IconCheck } from '@tabler/icons-react';
import { authAPI } from '../utils/api';
import { authUtils } from '../utils/auth';

const AuthModal = ({ opened, onClose, onAuthSuccess }) => {
  const [activeTab, setActiveTab] = useState('login');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [resetEmailSent, setResetEmailSent] = useState(false);

  const loginForm = useForm({
    initialValues: {
      email: '',
      password: '',
    },
    validate: {
      email: (value) => (/^\S+@\S+$/.test(value) ? null : 'Invalid email'),
      password: (value) => (value.length < 6 ? 'Password must be at least 6 characters' : null),
    },
  });

  const registerForm = useForm({
    initialValues: {
      email: '',
      password: '',
      confirmPassword: '',
      fullName: '',
      phoneNumber: '',
    },
    validate: {
      email: (value) => (/^\S+@\S+$/.test(value) ? null : 'Invalid email'),
      password: (value) => (value.length < 6 ? 'Password must be at least 6 characters' : null),
      confirmPassword: (value, values) =>
        value !== values.password ? 'Passwords do not match' : null,
      fullName: (value) => (value.trim().length < 2 ? 'Full name is required' : null),
      phoneNumber: (value) => {
        const phoneRegex = /^[\+]?[1-9][\d]{0,15}$/;
        return phoneRegex.test(value.replace(/[\s\-\(\)]/g, '')) ? null : 'Valid phone number is required';
      },
    },
  });

  const forgotPasswordForm = useForm({
    initialValues: {
      email: '',
    },
    validate: {
      email: (value) => (/^\S+@\S+$/.test(value) ? null : 'Invalid email'),
    },
  });

  const handleLogin = async (values) => {
    setLoading(true);
    setError('');
    
    try {
      const response = await authAPI.login(values.email, values.password);
      
      if (response.token) {
        authUtils.setToken(response.token);
        notifications.show({
          title: 'Success!',
          message: 'Welcome back!',
          color: 'green',
          icon: <IconCheck size={16} />,
        });
        
        onAuthSuccess();
        onClose();
      } else {
        setError('Login failed. Please try again.');
      }
    } catch (err) {
      const errorMessage = err.response?.data?.message || 'Login failed. Please check your credentials.';
      setError(errorMessage);
      notifications.show({
        title: 'Login Failed',
        message: errorMessage,
        color: 'red',
        icon: <IconAlertCircle size={16} />,
      });
    } finally {
      setLoading(false);
    }
  };

  const handleRegister = async (values) => {
    setLoading(true);
    setError('');
    
    try {
      const userData = {
        email: values.email,
        password: values.password,
        fullName: values.fullName,
        phoneNumber: values.phoneNumber,
      };
      
      const response = await authAPI.register(userData);
      
      if (response.token) {
        authUtils.setToken(response.token);
        notifications.show({
          title: 'Welcome!',
          message: 'Account created successfully!',
          color: 'green',
          icon: <IconCheck size={16} />,
        });
        
        onAuthSuccess();
        onClose();
      } else {
        notifications.show({
          title: 'Success!',
          message: 'Account created! Please check your email to verify your account.',
          color: 'blue',
          icon: <IconCheck size={16} />,
        });
        
        setActiveTab('login');
        registerForm.reset();
      }
    } catch (err) {
      const errorMessage = err.response?.data?.message || 'Registration failed. Please try again.';
      setError(errorMessage);
      notifications.show({
        title: 'Registration Failed',
        message: errorMessage,
        color: 'red',
        icon: <IconAlertCircle size={16} />,
      });
    } finally {
      setLoading(false);
    }
  };

  const handleForgotPassword = async (values) => {
    setLoading(true);
    setError('');
    
    try {
      await authAPI.requestPasswordReset(values.email);
      
      setResetEmailSent(true);
      notifications.show({
        title: 'Reset Email Sent!',
        message: 'Check your email for password reset instructions.',
        color: 'green',
        icon: <IconCheck size={16} />,
      });
    } catch (err) {
      const errorMessage = err.response?.data?.message || 'Failed to send reset email. Please try again.';
      setError(errorMessage);
      notifications.show({
        title: 'Error',
        message: errorMessage,
        color: 'red',
        icon: <IconAlertCircle size={16} />,
      });
    } finally {
      setLoading(false);
    }
  };

  const handleClose = () => {
    setError('');
    setResetEmailSent(false);
    loginForm.reset();
    registerForm.reset();
    forgotPasswordForm.reset();
    setActiveTab('login');
    onClose();
  };

  return (
    <Modal
      opened={opened}
      onClose={handleClose}
      title=""
      centered
      size="sm"
      padding="xl"
      radius="lg"
      styles={{
        modal: {
          backgroundColor: 'rgba(26, 27, 30, 0.95)',
          backdropFilter: 'blur(15px)',
          border: '1px solid rgba(255, 255, 255, 0.1)',
          borderRadius: '16px',
        },
        overlay: {
          backgroundColor: 'rgba(0, 0, 0, 0.6)',
        }
      }}
    >
      <Paper 
        radius="lg" 
        style={{ 
          position: 'relative',
          backgroundColor: 'transparent',
          border: 'none'
        }}
      >
        <LoadingOverlay visible={loading} />
        
        <Tabs value={activeTab} onChange={setActiveTab}>
          <Tabs.List grow>
            <Tabs.Tab value="login">Sign In</Tabs.Tab>
            <Tabs.Tab value="register">Sign Up</Tabs.Tab>
            <Tabs.Tab value="forgot" style={{ display: activeTab === 'forgot' ? 'block' : 'none' }}>Reset Password</Tabs.Tab>
          </Tabs.List>

          <Tabs.Panel value="login" pt="md">
            <form onSubmit={loginForm.onSubmit(handleLogin)}>
              <Stack spacing="md">
                <Title order={3} ta="center" c="white">
                  Welcome Back
                </Title>
                
                {error && (
                  <Alert icon={<IconAlertCircle size={16} />} color="red">
                    {error}
                  </Alert>
                )}

                <TextInput
                  label="Email"
                  placeholder="your@email.com"
                  required
                  {...loginForm.getInputProps('email')}
                />

                <PasswordInput
                  label="Password"
                  placeholder="Your password"
                  required
                  {...loginForm.getInputProps('password')}
                />

                <Button type="submit" fullWidth mt="md" size="md" disabled={loading}>
                  Sign In
                </Button>

                <Text ta="center" size="sm" c="dimmed">
                  Forgot your password?{' '}
                  <Anchor 
                    size="sm" 
                    component="button" 
                    type="button"
                    onClick={() => setActiveTab('forgot')}
                  >
                    Reset it here
                  </Anchor>
                </Text>
              </Stack>
            </form>
          </Tabs.Panel>

          <Tabs.Panel value="register" pt="md">
            <form onSubmit={registerForm.onSubmit(handleRegister)}>
              <Stack spacing="md">
                <Title order={3} ta="center" c="white">
                  Join Smart Fit Girl
                </Title>
                
                {error && (
                  <Alert icon={<IconAlertCircle size={16} />} color="red">
                    {error}
                  </Alert>
                )}

                <TextInput
                  label="Full Name"
                  placeholder="Jane Doe"
                  required
                  {...registerForm.getInputProps('fullName')}
                />

                <TextInput
                  label="Phone Number"
                  placeholder="+1 555-123-4567"
                  required
                  {...registerForm.getInputProps('phoneNumber')}
                />

                <TextInput
                  label="Email"
                  placeholder="your@email.com"
                  required
                  {...registerForm.getInputProps('email')}
                />

                <PasswordInput
                  label="Password"
                  placeholder="At least 6 characters"
                  required
                  {...registerForm.getInputProps('password')}
                />

                <PasswordInput
                  label="Confirm Password"
                  placeholder="Confirm your password"
                  required
                  {...registerForm.getInputProps('confirmPassword')}
                />

                <Button type="submit" fullWidth mt="md" size="md" disabled={loading}>
                  Create Account
                </Button>

                <Text ta="center" size="xs" c="dimmed">
                  By signing up, you agree to our Terms of Service and Privacy Policy
                </Text>
              </Stack>
            </form>
          </Tabs.Panel>

          <Tabs.Panel value="forgot" pt="md">
            {resetEmailSent ? (
              <Stack spacing="md">
                <Title order={3} ta="center" c="white">
                  Check Your Email
                </Title>
                
                <Alert icon={<IconCheck size={16} />} color="green">
                  We've sent password reset instructions to your email address.
                </Alert>

                <Text ta="center" size="sm" c="dimmed">
                  Didn't receive the email? Check your spam folder or{' '}
                  <Anchor 
                    size="sm" 
                    component="button" 
                    type="button"
                    onClick={() => {
                      setResetEmailSent(false);
                      setError('');
                    }}
                  >
                    try again
                  </Anchor>
                </Text>

                <Button 
                  variant="outline" 
                  fullWidth 
                  mt="md" 
                  onClick={() => setActiveTab('login')}
                >
                  Back to Sign In
                </Button>
              </Stack>
            ) : (
              <form onSubmit={forgotPasswordForm.onSubmit(handleForgotPassword)}>
                <Stack spacing="md">
                  <Title order={3} ta="center" c="white">
                    Reset Your Password
                  </Title>
                  
                  <Text ta="center" size="sm" c="dimmed">
                    Enter your email address and we'll send you a link to reset your password.
                  </Text>
                  
                  {error && (
                    <Alert icon={<IconAlertCircle size={16} />} color="red">
                      {error}
                    </Alert>
                  )}

                  <TextInput
                    label="Email"
                    placeholder="your@email.com"
                    required
                    {...forgotPasswordForm.getInputProps('email')}
                  />

                  <Button type="submit" fullWidth mt="md" size="md" disabled={loading}>
                    Send Reset Email
                  </Button>

                  <Text ta="center" size="sm" c="dimmed">
                    Remember your password?{' '}
                    <Anchor 
                      size="sm" 
                      component="button" 
                      type="button"
                      onClick={() => setActiveTab('login')}
                    >
                      Back to Sign In
                    </Anchor>
                  </Text>
                </Stack>
              </form>
            )}
          </Tabs.Panel>
        </Tabs>
      </Paper>
    </Modal>
  );
};

export default AuthModal;

import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import { MantineProvider } from '@mantine/core';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { Provider } from 'react-redux';
import { configureStore } from '@reduxjs/toolkit';
import authReducer from '@/redux/slices/authSlice';
import { RegisterPage } from './RegisterPage';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => {
      const translations: Record<string, string> = {
        'auth.signUp': 'Create an account',
        'auth.name': 'Name',
        'auth.namePlaceholder': 'Your name',
        'auth.email': 'Email',
        'auth.emailPlaceholder': 'your@email.com',
        'auth.password': 'Password',
        'auth.passwordPlaceholder': 'Your password',
        'auth.confirmPassword': 'Confirm Password',
        'auth.confirmPasswordPlaceholder': 'Repeat your password',
        'auth.register': 'Sign Up',
        'auth.alreadyHaveAccount': 'Already have an account? Sign in',
        'auth.nameMinLength': 'Name must be at least 2 characters',
        'auth.invalidEmail': 'Invalid email address',
        'auth.passwordMinLength': 'Password must be at least 8 characters',
        'auth.passwordsMismatch': 'Passwords do not match',
        'auth.registerError': 'Registration failed',
      };
      return translations[key] || key;
    },
    i18n: { language: 'en' },
  }),
}));

vi.mock('@tanstack/react-router', () => ({
  Link: ({ children, ...props }: any) => <a {...props}>{children}</a>,
  useNavigate: () => vi.fn(),
}));

function renderWithProviders(ui: React.ReactElement) {
  const store = configureStore({ reducer: { auth: authReducer } });
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });

  return render(
    <Provider store={store}>
      <QueryClientProvider client={queryClient}>
        <MantineProvider>{ui}</MantineProvider>
      </QueryClientProvider>
    </Provider>,
  );
}

describe('RegisterPage', () => {
  it('renders create account title', () => {
    renderWithProviders(<RegisterPage />);
    expect(screen.getByText('Create an account')).toBeInTheDocument();
  });

  it('renders all required form fields', () => {
    renderWithProviders(<RegisterPage />);
    expect(screen.getByLabelText('Name')).toBeInTheDocument();
    expect(screen.getByLabelText('Email')).toBeInTheDocument();
    expect(screen.getByLabelText('Password')).toBeInTheDocument();
    expect(screen.getByLabelText('Confirm Password')).toBeInTheDocument();
  });

  it('renders sign up button', () => {
    renderWithProviders(<RegisterPage />);
    const buttons = screen.getAllByRole('button');
    const signUpButton = buttons.find((btn) => btn.textContent?.includes('Sign Up'));
    expect(signUpButton).toBeInTheDocument();
  });

  it('renders link to login page', () => {
    renderWithProviders(<RegisterPage />);
    expect(screen.getByText('Already have an account? Sign in')).toBeInTheDocument();
  });
});

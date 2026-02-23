import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import { MantineProvider } from '@mantine/core';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { Provider } from 'react-redux';
import { configureStore } from '@reduxjs/toolkit';
import authReducer from '@/redux/slices/authSlice';
import { LoginPage } from './LoginPage';

// Mock i18next
vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => {
      const translations: Record<string, string> = {
        'auth.signIn': 'Sign in to FeedbackBot',
        'auth.email': 'Email',
        'auth.emailPlaceholder': 'your@email.com',
        'auth.password': 'Password',
        'auth.passwordPlaceholder': 'Your password',
        'auth.noAccount': "Don't have an account? Sign up",
        'auth.invalidEmail': 'Invalid email address',
        'auth.passwordMinLength': 'Password must be at least 8 characters',
        'auth.invalidCredentials': 'Invalid email or password',
      };
      return translations[key] || key;
    },
    i18n: { language: 'en' },
  }),
}));

// Mock router
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

describe('LoginPage', () => {
  it('renders sign in title', () => {
    renderWithProviders(<LoginPage />);
    expect(screen.getByRole('heading', { name: 'Sign in to FeedbackBot' })).toBeInTheDocument();
  });

  it('renders email and password inputs', () => {
    renderWithProviders(<LoginPage />);
    expect(screen.getByLabelText('Email')).toBeInTheDocument();
    expect(screen.getByLabelText('Password')).toBeInTheDocument();
  });

  it('renders sign in button', () => {
    renderWithProviders(<LoginPage />);
    const buttons = screen.getAllByRole('button');
    const signInButton = buttons.find((btn) => btn.textContent?.includes('Sign in'));
    expect(signInButton).toBeInTheDocument();
  });

  it('renders link to register page', () => {
    renderWithProviders(<LoginPage />);
    expect(screen.getByText("Don't have an account? Sign up")).toBeInTheDocument();
  });
});

import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import { MantineProvider } from '@mantine/core';
import { Provider } from 'react-redux';
import { configureStore } from '@reduxjs/toolkit';
import authReducer from '@/redux/slices/authSlice';
import { Navbar } from './Navbar';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => {
      const translations: Record<string, string> = {
        'nav.setup': 'Setup',
        'nav.bots': 'Bots',
        'nav.groups': 'Groups',
        'nav.feedbacks': 'Feedbacks',
        'auth.logout': 'Logout',
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

function renderNavbar() {
  const store = configureStore({ reducer: { auth: authReducer } });
  return render(
    <Provider store={store}>
      <MantineProvider>
        <Navbar />
      </MantineProvider>
    </Provider>,
  );
}

describe('Navbar', () => {
  it('renders the FeedbackBot title', () => {
    renderNavbar();
    expect(screen.getByText('FeedbackBot')).toBeInTheDocument();
  });

  it('renders navigation links', () => {
    renderNavbar();
    expect(screen.getByText('Setup')).toBeInTheDocument();
    expect(screen.getByText('Bots')).toBeInTheDocument();
    expect(screen.getByText('Groups')).toBeInTheDocument();
    expect(screen.getByText('Feedbacks')).toBeInTheDocument();
  });

  it('renders logout button', () => {
    renderNavbar();
    expect(screen.getByText('Logout')).toBeInTheDocument();
  });
});

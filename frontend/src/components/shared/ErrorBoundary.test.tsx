import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import { MantineProvider } from '@mantine/core';
import { ErrorBoundary } from './ErrorBoundary';

function ThrowingComponent({ message }: { message: string }): never {
  throw new Error(message);
}

describe('ErrorBoundary', () => {
  it('renders children when no error', () => {
    render(
      <MantineProvider>
        <ErrorBoundary>
          <div>Working content</div>
        </ErrorBoundary>
      </MantineProvider>,
    );
    expect(screen.getByText('Working content')).toBeInTheDocument();
  });

  it('renders error UI when child throws', () => {
    // Suppress console.error for expected error
    const spy = vi.spyOn(console, 'error').mockImplementation(() => {});
    render(
      <MantineProvider>
        <ErrorBoundary>
          <ThrowingComponent message="Test error message" />
        </ErrorBoundary>
      </MantineProvider>,
    );
    expect(screen.getByText('Something went wrong')).toBeInTheDocument();
    expect(screen.getByText('Test error message')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Reload page' })).toBeInTheDocument();
    spy.mockRestore();
  });
});

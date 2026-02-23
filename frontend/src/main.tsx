import React from 'react';
import ReactDOM from 'react-dom/client';
import { MantineProvider } from '@mantine/core';
import { Notifications } from '@mantine/notifications';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { Provider } from 'react-redux';
import { RouterProvider, createRouter } from '@tanstack/react-router';
import { store } from '@/redux/store';
import { routeTree } from '@/routes/routeTree';
import { ErrorBoundary } from '@/components/shared';
import '@mantine/core/styles.css';
import '@mantine/notifications/styles.css';
import '@/i18n';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      refetchOnWindowFocus: false,
    },
  },
});

const router = createRouter({ routeTree });

declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router;
  }
}

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <Provider store={store}>
      <QueryClientProvider client={queryClient}>
        <MantineProvider>
          <ErrorBoundary>
            <Notifications />
            <RouterProvider router={router} />
          </ErrorBoundary>
        </MantineProvider>
      </QueryClientProvider>
    </Provider>
  </React.StrictMode>,
);

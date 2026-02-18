import { createRootRoute, createRoute, Outlet } from '@tanstack/react-router';
import { Navbar } from '@/components/shared';
import { SetupPage, BotConfigPage } from '@/components/pageComponents';

const rootRoute = createRootRoute({
  component: () => (
    <>
      <Navbar />
      <Outlet />
    </>
  ),
});

const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/',
  component: SetupPage,
});

const setupRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/setup',
  component: SetupPage,
});

const botSettingsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/settings/bot',
  component: BotConfigPage,
});

export const routeTree = rootRoute.addChildren([indexRoute, setupRoute, botSettingsRoute]);

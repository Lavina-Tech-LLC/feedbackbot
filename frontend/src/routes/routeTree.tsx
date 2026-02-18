import { createRootRoute, createRoute, Outlet } from '@tanstack/react-router';
import { Navbar } from '@/components/shared';
import { SetupPage, BotConfigPage, GroupsPage, FeedbacksPage } from '@/components/pageComponents';

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

const groupsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/groups',
  component: GroupsPage,
});

const feedbacksRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/feedbacks',
  component: FeedbacksPage,
});

export const routeTree = rootRoute.addChildren([indexRoute, setupRoute, botSettingsRoute, groupsRoute, feedbacksRoute]);

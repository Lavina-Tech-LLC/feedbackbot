import { createRootRoute, createRoute, Outlet } from '@tanstack/react-router';
import { Navbar } from '@/components/shared';
import { SetupPage, BotConfigPage, GroupsPage, FeedbacksPage, LoginPage, CallbackPage } from '@/components/pageComponents';

const rootRoute = createRootRoute({
  component: () => <Outlet />,
});

const authLayout = createRoute({
  getParentRoute: () => rootRoute,
  id: 'auth-layout',
  component: () => <Outlet />,
});

const appLayout = createRoute({
  getParentRoute: () => rootRoute,
  id: 'app-layout',
  component: () => (
    <>
      <Navbar />
      <Outlet />
    </>
  ),
});

const loginRoute = createRoute({
  getParentRoute: () => authLayout,
  path: '/login',
  component: LoginPage,
});

const callbackRoute = createRoute({
  getParentRoute: () => authLayout,
  path: '/auth/callback',
  component: CallbackPage,
});

const indexRoute = createRoute({
  getParentRoute: () => appLayout,
  path: '/',
  component: SetupPage,
});

const setupRoute = createRoute({
  getParentRoute: () => appLayout,
  path: '/setup',
  component: SetupPage,
});

const botSettingsRoute = createRoute({
  getParentRoute: () => appLayout,
  path: '/settings/bot',
  component: BotConfigPage,
});

const groupsRoute = createRoute({
  getParentRoute: () => appLayout,
  path: '/groups',
  component: GroupsPage,
});

const feedbacksRoute = createRoute({
  getParentRoute: () => appLayout,
  path: '/feedbacks',
  component: FeedbacksPage,
});

export const routeTree = rootRoute.addChildren([
  authLayout.addChildren([loginRoute, callbackRoute]),
  appLayout.addChildren([indexRoute, setupRoute, botSettingsRoute, groupsRoute, feedbacksRoute]),
]);

import { createRootRoute, createRoute, Outlet, redirect } from '@tanstack/react-router';
import { Navbar } from '@/components/shared';
import { SetupPage, BotConfigPage, GroupsPage, FeedbacksPage, LoginPage, RegisterPage, RedirectPage } from '@/components/pageComponents';
import { store } from '@/redux/store';

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
  beforeLoad: () => {
    const token = store.getState().auth.token;
    if (!token) {
      throw redirect({ to: '/login' });
    }
  },
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

const registerRoute = createRoute({
  getParentRoute: () => authLayout,
  path: '/register',
  component: RegisterPage,
});

const indexRoute = createRoute({
  getParentRoute: () => appLayout,
  path: '/',
  component: RedirectPage,
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
  authLayout.addChildren([loginRoute, registerRoute]),
  appLayout.addChildren([indexRoute, setupRoute, botSettingsRoute, groupsRoute, feedbacksRoute]),
]);

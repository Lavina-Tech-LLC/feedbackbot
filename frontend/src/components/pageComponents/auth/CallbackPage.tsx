import { useEffect, useRef } from 'react';
import { Center, Loader, Stack, Text } from '@mantine/core';
import { notifications } from '@mantine/notifications';
import { useTranslation } from 'react-i18next';
import { useNavigate } from '@tanstack/react-router';
import { useExchangeToken } from '@/service';
import { useAppDispatch } from '@/redux/store';
import { setUser, setToken } from '@/redux/slices';

export function CallbackPage() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const dispatch = useAppDispatch();
  const exchangeToken = useExchangeToken();
  const calledRef = useRef(false);

  useEffect(() => {
    if (calledRef.current) return;
    calledRef.current = true;

    const params = new URLSearchParams(window.location.search);
    const code = params.get('code');

    if (!code) {
      notifications.show({ title: t('common.error'), message: t('auth.loginError'), color: 'red' });
      setTimeout(() => navigate({ to: '/login' }), 2000);
      return;
    }

    exchangeToken.mutate(
      { code, redirect_uri: window.location.origin + '/auth/callback' },
      {
        onSuccess: (res) => {
          dispatch(setToken(res.data.access_token));
          dispatch(setUser(res.data));
          navigate({ to: '/' });
        },
        onError: () => {
          notifications.show({ title: t('common.error'), message: t('auth.loginError'), color: 'red' });
          setTimeout(() => navigate({ to: '/login' }), 2000);
        },
      },
    );
  }, []);

  return (
    <Center h="100vh">
      <Stack align="center" gap="md">
        <Loader size="lg" />
        <Text c="dimmed">{t('auth.redirecting')}</Text>
      </Stack>
    </Center>
  );
}

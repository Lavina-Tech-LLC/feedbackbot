import { Alert, Anchor, Button, Center, Divider, Paper, PasswordInput, Stack, TextInput, Title } from '@mantine/core';
import { useForm } from '@mantine/form';
import { IconAlertCircle, IconLogin } from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';
import { Link, useNavigate } from '@tanstack/react-router';
import { useAuthConfig, useForgotPassword, useLogin } from '@/service';
import { useAppDispatch } from '@/redux/store';
import { setToken, setUser } from '@/redux/slices';
import { api } from '@/api';
import { useState } from 'react';
import type { User } from '@/types';

export function LoginPage() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const dispatch = useAppDispatch();
  const login = useLogin();
  const forgotPassword = useForgotPassword();
  const { data: configRes, isLoading: isConfigLoading } = useAuthConfig();
  const [error, setError] = useState('');

  const form = useForm({
    initialValues: {
      email: '',
      password: '',
    },
    validate: {
      email: (v) => (/^\S+@\S+\.\S+$/.test(v) ? null : t('auth.invalidEmail')),
      password: (v) => (v.length < 1 ? t('auth.passwordMinLength') : null),
    },
  });

  const handleSubmit = form.onSubmit((values) => {
    setError('');
    login.mutate(
      { email: values.email, password: values.password },
      {
        onSuccess: async (res) => {
          dispatch(setToken(res.data.access_token));
          try {
            const meRes = await api.get<never, { data: User }>('/auth/me');
            dispatch(setUser(meRes.data));
          } catch {
            // token set, user will be fetched on next load
          }
          navigate({ to: '/' });
        },
        onError: (err) => {
          const message =
            (err as any)?.response?.data?.message || t('auth.invalidCredentials');
          setError(message);
        },
      },
    );
  });

  const handleOAuthLogin = async () => {
    const config = configRes?.data;
    if (!config?.authorize_url) return;

    try {
      const res = await fetch(config.authorize_url);
      const json = await res.json();
      const googleUrl = json.data?.url ?? json.url;
      if (googleUrl) {
        window.location.href = googleUrl;
      }
    } catch {
      // ignore fetch errors; user can retry
    }
  };

  const handleForgotPassword = () => {
    forgotPassword.mutate(undefined, {
      onSuccess: (res) => {
        window.open(res.data.url, '_blank');
      },
    });
  };

  return (
    <Center h="100vh">
      <Paper shadow="md" p="xl" radius="md" w={400}>
        <form onSubmit={handleSubmit}>
          <Stack gap="md">
            <Title order={2} ta="center">
              {t('auth.signIn')}
            </Title>

            {error && (
              <Alert icon={<IconAlertCircle size={16} />} color="red" variant="light">
                {error}
              </Alert>
            )}

            <TextInput
              label={t('auth.email')}
              placeholder={t('auth.emailPlaceholder')}
              type="email"
              {...form.getInputProps('email')}
            />

            <PasswordInput
              label={t('auth.password')}
              placeholder={t('auth.passwordPlaceholder')}
              {...form.getInputProps('password')}
            />

            <Anchor component="button" type="button" size="sm" onClick={handleForgotPassword}>
              {t('auth.forgotPassword')}
            </Anchor>

            <Button
              type="submit"
              fullWidth
              size="md"
              leftSection={<IconLogin size={20} />}
              loading={login.isPending}
            >
              {t('auth.signIn')}
            </Button>

            <Divider label={t('auth.or')} labelPosition="center" />

            <Button
              fullWidth
              size="md"
              variant="outline"
              leftSection={<IconLogin size={20} />}
              onClick={handleOAuthLogin}
              loading={isConfigLoading}
            >
              {t('auth.loginWithLavina')}
            </Button>

            <Anchor component={Link} to="/register" ta="center" size="sm">
              {t('auth.noAccount')}
            </Anchor>
          </Stack>
        </form>
      </Paper>
    </Center>
  );
}

import { Alert, Anchor, Button, Center, Paper, PasswordInput, Stack, TextInput, Title } from '@mantine/core';
import { useForm } from '@mantine/form';
import { IconAlertCircle, IconLogin } from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';
import { Link, useNavigate } from '@tanstack/react-router';
import { useLogin } from '@/service';
import { useAppDispatch } from '@/redux/store';
import { setToken, setRefreshToken, setUser } from '@/redux/slices';
import { useState } from 'react';

export function LoginPage() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const dispatch = useAppDispatch();
  const login = useLogin();
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
        onSuccess: (res) => {
          dispatch(setToken(res.data.access_token));
          dispatch(setRefreshToken(res.data.refresh_token));
          dispatch(setUser({
            id: res.data.user_id,
            email: res.data.email,
            name: res.data.name,
            role: res.data.role,
            tenant_id: res.data.tenant_id,
          }));
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

            <Button
              type="submit"
              fullWidth
              size="md"
              leftSection={<IconLogin size={20} />}
              loading={login.isPending}
            >
              {t('auth.signIn')}
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

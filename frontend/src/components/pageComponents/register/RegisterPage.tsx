import { Alert, Anchor, Button, Center, Paper, PasswordInput, Stack, TextInput, Title } from '@mantine/core';
import { useForm } from '@mantine/form';
import { IconAlertCircle, IconUserPlus } from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';
import { Link, useNavigate } from '@tanstack/react-router';
import { useRegister } from '@/service';
import { useAppDispatch } from '@/redux/store';
import { setToken, setUser } from '@/redux/slices';
import { api } from '@/api';
import { useState } from 'react';
import type { User } from '@/types';

export function RegisterPage() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const dispatch = useAppDispatch();
  const register = useRegister();
  const [error, setError] = useState('');

  const form = useForm({
    initialValues: {
      name: '',
      email: '',
      password: '',
      confirmPassword: '',
    },
    validate: {
      name: (v) => (v.trim().length < 2 ? t('auth.nameMinLength') : null),
      email: (v) => (/^\S+@\S+\.\S+$/.test(v) ? null : t('auth.invalidEmail')),
      password: (v) => (v.length < 8 ? t('auth.passwordMinLength') : null),
      confirmPassword: (v, values) => (v !== values.password ? t('auth.passwordsMismatch') : null),
    },
  });

  const handleSubmit = form.onSubmit((values) => {
    setError('');
    register.mutate(
      { name: values.name, email: values.email, password: values.password },
      {
        onSuccess: async (res) => {
          dispatch(setToken(res.data.access_token));
          try {
            const meRes = await api.get<never, { data: User }>('/auth/me');
            dispatch(setUser(meRes.data));
          } catch {
            dispatch(setUser(res.data));
          }
          navigate({ to: '/' });
        },
        onError: (err) => {
          const message =
            (err as any)?.response?.data?.message || t('auth.registerError');
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
              {t('auth.signUp')}
            </Title>

            {error && (
              <Alert icon={<IconAlertCircle size={16} />} color="red" variant="light">
                {error}
              </Alert>
            )}

            <TextInput
              label={t('auth.name')}
              placeholder={t('auth.namePlaceholder')}
              {...form.getInputProps('name')}
            />

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

            <PasswordInput
              label={t('auth.confirmPassword')}
              placeholder={t('auth.confirmPasswordPlaceholder')}
              {...form.getInputProps('confirmPassword')}
            />

            <Button
              type="submit"
              fullWidth
              size="md"
              leftSection={<IconUserPlus size={20} />}
              loading={register.isPending}
            >
              {t('auth.register')}
            </Button>

            <Anchor component={Link} to="/login" ta="center" size="sm">
              {t('auth.alreadyHaveAccount')}
            </Anchor>
          </Stack>
        </form>
      </Paper>
    </Center>
  );
}

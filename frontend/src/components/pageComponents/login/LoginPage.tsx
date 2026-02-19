import { Button, Center, Paper, Stack, Title } from '@mantine/core';
import { IconLogin } from '@tabler/icons-react';
import { useTranslation } from 'react-i18next';
import { useAuthConfig } from '@/service';

export function LoginPage() {
  const { t } = useTranslation();
  const { data: configRes, isLoading } = useAuthConfig();

  const handleLogin = () => {
    const config = configRes?.data;
    if (!config) return;

    const params = new URLSearchParams({
      client_id: config.client_id,
      redirect_uri: config.redirect_uri,
      response_type: 'code',
      scope: config.scope,
    });

    window.location.href = `${config.authorize_url}?${params.toString()}`;
  };

  return (
    <Center h="100vh">
      <Paper shadow="md" p="xl" radius="md" w={400}>
        <Stack align="center" gap="lg">
          <Title order={2}>{t('auth.signIn')}</Title>
          <Button
            fullWidth
            size="lg"
            leftSection={<IconLogin size={20} />}
            onClick={handleLogin}
            loading={isLoading}
          >
            {t('auth.loginWithLavina')}
          </Button>
        </Stack>
      </Paper>
    </Center>
  );
}

import { Group, Title, Anchor, Button } from '@mantine/core';
import { useTranslation } from 'react-i18next';
import { Link, useNavigate } from '@tanstack/react-router';
import { IconMessageReport, IconLogout } from '@tabler/icons-react';
import { useAppDispatch } from '@/redux/store';
import { clearAuth } from '@/redux/slices';

export function Navbar() {
  const { t } = useTranslation();
  const dispatch = useAppDispatch();
  const navigate = useNavigate();

  const handleLogout = () => {
    dispatch(clearAuth());
    navigate({ to: '/login' });
  };

  return (
    <Group h={60} px="md" justify="space-between" style={{ borderBottom: '1px solid #e9ecef' }}>
      <Group>
        <IconMessageReport size={28} />
        <Title order={4}>FeedbackBot</Title>
      </Group>
      <Group>
        <Anchor component={Link} to="/setup">
          {t('nav.setup')}
        </Anchor>
        <Anchor component={Link} to="/settings/bot">
          {t('nav.bots')}
        </Anchor>
        <Anchor component={Link} to="/groups">
          {t('nav.groups')}
        </Anchor>
        <Anchor component={Link} to="/feedbacks">
          {t('nav.feedbacks')}
        </Anchor>
        <Button
          variant="subtle"
          color="red"
          leftSection={<IconLogout size={16} />}
          onClick={handleLogout}
        >
          {t('auth.logout')}
        </Button>
      </Group>
    </Group>
  );
}

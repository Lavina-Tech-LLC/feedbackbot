import { Group, Title, Anchor, Button, Burger, Drawer, Stack } from '@mantine/core';
import { useDisclosure, useMediaQuery } from '@mantine/hooks';
import { useTranslation } from 'react-i18next';
import { Link, useNavigate } from '@tanstack/react-router';
import { IconMessageReport, IconLogout } from '@tabler/icons-react';
import { useAppDispatch } from '@/redux/store';
import { clearAuth } from '@/redux/slices';

export function Navbar() {
  const { t } = useTranslation();
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const [opened, { toggle, close }] = useDisclosure(false);
  const isMobile = useMediaQuery('(max-width: 768px)');

  const handleLogout = () => {
    dispatch(clearAuth());
    navigate({ to: '/login' });
    close();
  };

  const navLinks = (
    <>
      <Anchor component={Link} to="/setup" onClick={close}>
        {t('nav.setup')}
      </Anchor>
      <Anchor component={Link} to="/settings/bot" onClick={close}>
        {t('nav.bots')}
      </Anchor>
      <Anchor component={Link} to="/groups" onClick={close}>
        {t('nav.groups')}
      </Anchor>
      <Anchor component={Link} to="/feedbacks" onClick={close}>
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
    </>
  );

  return (
    <Group h={60} px="md" justify="space-between" style={{ borderBottom: '1px solid #e9ecef' }}>
      <Group>
        <IconMessageReport size={28} />
        <Title order={4}>FeedbackBot</Title>
      </Group>

      {isMobile ? (
        <>
          <Burger opened={opened} onClick={toggle} size="sm" />
          <Drawer opened={opened} onClose={close} title="Menu" position="right" size="xs">
            <Stack gap="md" p="md">
              {navLinks}
            </Stack>
          </Drawer>
        </>
      ) : (
        <Group>{navLinks}</Group>
      )}
    </Group>
  );
}

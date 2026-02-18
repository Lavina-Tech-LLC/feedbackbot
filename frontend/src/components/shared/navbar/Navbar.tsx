import { Group, Title, Anchor } from '@mantine/core';
import { useTranslation } from 'react-i18next';
import { Link } from '@tanstack/react-router';
import { IconMessageReport } from '@tabler/icons-react';

export function Navbar() {
  const { t } = useTranslation();

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
      </Group>
    </Group>
  );
}
